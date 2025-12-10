package file

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
	"github.com/cthulhu-platform/gateway/internal/repository/local"
	"github.com/google/uuid"
)

type localFileService struct {
	conns    *microservices.ServiceConnectionContainer
	fileRepo local.FileRepository
}

func NewLocalFileService(conns *microservices.ServiceConnectionContainer, fileRepo local.FileRepository) *localFileService {
	return &localFileService{
		conns:    conns,
		fileRepo: fileRepo,
	}
}

func (s *localFileService) UploadFiles(ctx context.Context, files []*multipart.FileHeader, userID *string) (*filemanager.UploadResult, error) {
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	fm := s.filemanager()
	if fm == nil {
		return nil, errors.New("filemanager connection not configured")
	}

	res := &filemanager.UploadResult{
		TransactionID: uuid.New().String(),
		Success:       false,
	}

	// Generate storage_id (bucket_id) - 10 char alphanumeric
	storageID := s.generateStorageID()

	// Check if bucket already exists
	existingBucket, err := s.fileRepo.GetBucketByID(storageID)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}
	if existingBucket != nil {
		// Retry with new storage_id
		storageID = s.generateStorageID()
		existingBucket, err = s.fileRepo.GetBucketByID(storageID)
		if err != nil {
			res.Error = err.Error()
			return res, err
		}
		if existingBucket != nil {
			res.Error = "failed to generate unique storage id"
			return res, errors.New(res.Error)
		}
	}

	// Create bucket in DB
	now := time.Now().Unix()
	bucket := &local.Bucket{
		ID:        storageID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.fileRepo.CreateBucket(bucket); err != nil {
		res.Error = err.Error()
		return res, err
	}

	// Add bucket admin if user is logged in (validate user exists first)
	if userID != nil && *userID != "" {
		// Validate user exists through authentication microservice
		authConn := s.conns.Authentication
		if authConn != nil {
			valid, err := authConn.ValidateUserID(*userID)
			if err == nil && valid {
				admin := &local.BucketAdmin{
					UserID:    *userID,
					BucketID:  storageID,
					CreatedAt: now,
				}
				if err := s.fileRepo.AddBucketAdmin(admin); err != nil {
					// Log error but don't fail upload
					_ = err
				}
			}
		}
	}

	// Open all files
	objects := make([]filemanager.UploadObject, 0, len(files))
	closers := make([]io.Closer, 0, len(files))
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			res.Error = err.Error()
			return res, err
		}
		closers = append(closers, file)
		objects = append(objects, filemanager.UploadObject{
			Name:        fh.Filename,
			Size:        fh.Size,
			ContentType: fh.Header.Get("Content-Type"),
			Body:        file,
		})
	}
	defer func() {
		for _, f := range closers {
			f.Close()
		}
	}()

	// Process each file: generate string_id, upload to S3, save to DB
	totalSize := int64(0)
	fileInfos := make([]filemanager.FileInfo, 0, len(objects))

	for _, obj := range objects {
		// Generate unique string_id (UUID)
		stringID := s.generateUniqueStringID(ctx)
		if stringID == "" {
			res.Error = "failed to generate unique string_id"
			return res, errors.New(res.Error)
		}

		// Upload to S3 using bucket_id/string_id as key
		s3Key := storageID + "/" + stringID
		uploadObj := filemanager.UploadObject{
			Name:        stringID, // Use string_id as the S3 object name
			Size:        obj.Size,
			ContentType: obj.ContentType,
			Body:        obj.Body,
		}

		// Upload to S3 using filemanager
		// We need to access the local implementation
		// For now, we'll use a type assertion - this works because we know we're using local
		type uploader interface {
			UploadSingleObject(ctx context.Context, storageID, stringID string, obj filemanager.UploadObject) error
		}

		uploaderFM, ok := fm.(uploader)
		if !ok {
			res.Error = "filemanager does not support single object upload"
			return res, errors.New(res.Error)
		}

		// Upload to S3
		if err := uploaderFM.UploadSingleObject(ctx, storageID, stringID, uploadObj); err != nil {
			res.Error = err.Error()
			return res, err
		}

		// Validate user_id if provided before storing as owner
		var ownerID *string
		if userID != nil && *userID != "" {
			authConn := s.conns.Authentication
			if authConn != nil {
				valid, err := authConn.ValidateUserID(*userID)
				if err == nil && valid {
					ownerID = userID
				}
			}
		}

		// Save file metadata to DB
		dbFile := &local.File{
			StringID:     stringID,
			BucketID:     storageID,
			OriginalName: obj.Name,
			OwnerID:      ownerID,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
			S3Key:        s3Key,
			CreatedAt:    now,
		}
		if err := s.fileRepo.CreateFile(dbFile); err != nil {
			res.Error = err.Error()
			return res, err
		}

		fileInfos = append(fileInfos, filemanager.FileInfo{
			OriginalName: obj.Name,
			StringID:     stringID,
			Key:          s3Key,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
		})
		totalSize += obj.Size
	}

	res.StorageID = storageID
	res.Files = fileInfos
	res.TotalSize = totalSize
	res.Success = true
	return res, nil
}

func (s *localFileService) DownloadFile(ctx context.Context, storageID, stringID string) (*filemanager.DownloadResult, error) {
	if storageID == "" || stringID == "" {
		return nil, errors.New("storage id and string id are required")
	}

	// Look up file by bucket_id and string_id from DB
	file, err := s.fileRepo.GetFileByStringID(stringID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file not found")
	}

	// Verify the file belongs to the specified bucket
	if file.BucketID != storageID {
		return nil, errors.New("file does not belong to the specified bucket")
	}

	// Use stored s3_key to download from S3
	fm := s.filemanager()
	if fm == nil {
		return nil, errors.New("filemanager connection not configured")
	}

	// Extract just the filename part from s3_key for the download
	// s3_key format: "bucket_id/string_id"
	parts := strings.Split(file.S3Key, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid s3_key format")
	}

	downloadResult, err := fm.Download(ctx, storageID, parts[1])
	if err != nil {
		return nil, err
	}

	// Override the downloaded filename with the original name from DB
	downloadResult.DownloadedFile = file.OriginalName
	return downloadResult, nil
}

func (s *localFileService) RetrieveFileBucket(ctx context.Context, storageID string) (*filemanager.BucketMetadata, error) {
	if storageID == "" {
		return nil, errors.New("storage id is required")
	}

	// Query DB for bucket
	bucket, err := s.fileRepo.GetBucketByID(storageID)
	if err != nil {
		return nil, err
	}
	if bucket == nil {
		return nil, errors.New("bucket not found")
	}

	// Query DB for files in bucket
	dbFiles, err := s.fileRepo.GetFilesByBucketID(storageID)
	if err != nil {
		return nil, err
	}

	// Convert DB files to FileInfo
	files := make([]filemanager.FileInfo, 0, len(dbFiles))
	var totalSize int64

	for _, dbFile := range dbFiles {
		files = append(files, filemanager.FileInfo{
			OriginalName: dbFile.OriginalName,
			StringID:     dbFile.StringID,
			Key:          dbFile.S3Key,
			Size:         dbFile.Size,
			ContentType:  dbFile.ContentType,
		})
		totalSize += dbFile.Size
	}

	return &filemanager.BucketMetadata{
		StorageID: storageID,
		Files:     files,
		TotalSize: totalSize,
	}, nil
}

func (s *localFileService) GetBucketAdmins(ctx context.Context, bucketID string) (*BucketAdminsResponse, error) {
	if bucketID == "" {
		return nil, errors.New("bucket id is required")
	}

	// Verify bucket exists
	bucket, err := s.fileRepo.GetBucketByID(bucketID)
	if err != nil {
		return nil, err
	}
	if bucket == nil {
		return nil, errors.New("bucket not found")
	}

	// Get all admins for the bucket
	admins, err := s.fileRepo.GetBucketAdminsByBucketID(bucketID)
	if err != nil {
		return nil, err
	}

	if len(admins) == 0 {
		return &BucketAdminsResponse{
			BucketID: bucketID,
			Owner:    nil,
			Admins:   []AdminInfo{},
		}, nil
	}

	// Determine owner: first admin by created_at (oldest)
	var ownerAdmin *local.BucketAdmin
	if len(admins) > 0 {
		ownerAdmin = admins[0]
		for _, admin := range admins[1:] {
			if admin.CreatedAt < ownerAdmin.CreatedAt {
				ownerAdmin = admin
			}
		}
	}

	// Access auth repository to get user details
	authConn := s.conns.Authentication
	if authConn == nil {
		return nil, errors.New("authentication connection not configured")
	}

	// Type assert to local connection to access GetRepo
	type repoGetter interface {
		GetRepo() local.AuthRepository
	}

	localAuthConn, ok := authConn.(repoGetter)
	if !ok {
		return nil, errors.New("authentication connection does not support repository access")
	}

	authRepo := localAuthConn.GetRepo()

	// Build admin info list with user details
	adminInfos := make([]AdminInfo, 0, len(admins))
	var ownerInfo *AdminInfo

	for _, admin := range admins {
		// Fetch user details
		user, err := authRepo.GetUserByID(admin.UserID)
		if err != nil {
			// If user not found, skip or use minimal info
			continue
		}
		if user == nil {
			continue
		}

		isOwner := ownerAdmin != nil && admin.UserID == ownerAdmin.UserID && admin.CreatedAt == ownerAdmin.CreatedAt

		adminInfo := AdminInfo{
			UserID:    admin.UserID,
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
			IsOwner:   isOwner,
			CreatedAt: admin.CreatedAt,
		}

		if isOwner {
			ownerInfo = &adminInfo
		}

		adminInfos = append(adminInfos, adminInfo)
	}

	return &BucketAdminsResponse{
		BucketID: bucketID,
		Owner:    ownerInfo,
		Admins:   adminInfos,
	}, nil
}

func (s *localFileService) filemanager() filemanager.FilemanagerConnection {
	if s == nil || s.conns == nil {
		return nil
	}
	return s.conns.Filemanager
}

// generateStorageID generates a 10-character alphanumeric storage ID
func (s *localFileService) generateStorageID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		idx := rand.IntN(len(letters))
		sb.WriteByte(letters[idx])
	}
	return sb.String()
}

// generateUniqueStringID generates a UUID and checks for uniqueness in the database
func (s *localFileService) generateUniqueStringID(ctx context.Context) string {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		stringID := uuid.New().String()
		// Check if string_id already exists in DB
		existing, err := s.fileRepo.GetFileByStringID(stringID)
		if err != nil {
			// If error checking, try again
			continue
		}
		if existing == nil {
			return stringID
		}
	}
	return "" // Failed to generate unique ID after retries
}
