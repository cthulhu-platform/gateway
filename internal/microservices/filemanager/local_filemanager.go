package filemanager

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/google/uuid"
)

// Local implementation uses its own S3-backed connection (e.g., LocalStack).
type localFilemanagerConnection struct {
	client   *s3.Client
	bucket   string
	idLength int
}

func NewLocalFilemanagerConnection() (*localFilemanagerConnection, error) {
	return &localFilemanagerConnection{
		client:   newS3Client(),
		bucket:   pkg.S3_BUCKET,
		idLength: parseLength(pkg.S3_STORAGE_ID_LENGTH, 10),
	}, nil
}

func (c *localFilemanagerConnection) Upload(ctx context.Context, storageID string, objects []UploadObject) (*UploadResult, error) {
	if len(objects) == 0 {
		return nil, errors.New("no files provided")
	}

	res := &UploadResult{
		TransactionID: uuid.New().String(),
		Success:       false,
	}

	if storageID == "" {
		storageID = c.generateStorageID()
	}

	exists, err := c.prefixExists(ctx, storageID)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}
	if exists {
		err := fmt.Errorf("storage id %s already exists", storageID)
		res.Error = err.Error()
		return res, err
	}

	totalSize := int64(0)
	fileInfos := make([]FileInfo, 0, len(objects))

	for _, obj := range objects {
		key := fmt.Sprintf("%s/%s", storageID, obj.Name)
		_, putErr := c.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(c.bucket),
			Key:         aws.String(key),
			Body:        obj.Body,
			ContentType: aws.String(obj.ContentType),
		})
		if putErr != nil {
			res.Error = putErr.Error()
			return res, putErr
		}

		fileInfos = append(fileInfos, FileInfo{
			FileName:    obj.Name,
			Key:         key,
			Size:        obj.Size,
			ContentType: obj.ContentType,
		})
		totalSize += obj.Size
	}

	res.StorageID = storageID
	res.Files = fileInfos
	res.TotalSize = totalSize
	res.Success = true
	return res, nil
}

func (c *localFilemanagerConnection) List(ctx context.Context, storageID string) (*BucketMetadata, error) {
	if storageID == "" {
		return nil, errors.New("storage id is required")
	}

	prefix := storageID + "/"
	resp, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Contents) == 0 {
		return nil, fmt.Errorf("storage id %s not found", storageID)
	}

	files := make([]FileInfo, 0, len(resp.Contents))
	var totalSize int64

	for _, obj := range resp.Contents {
		name := strings.TrimPrefix(aws.ToString(obj.Key), prefix)
		if name == "" {
			continue
		}
		size := aws.ToInt64(obj.Size)
		files = append(files, FileInfo{
			FileName: name,
			Key:      aws.ToString(obj.Key),
			Size:     size,
		})
		totalSize += size
	}

	return &BucketMetadata{
		StorageID: storageID,
		Files:     files,
		TotalSize: totalSize,
	}, nil
}

func (c *localFilemanagerConnection) Download(ctx context.Context, storageID, filename string) (*DownloadResult, error) {
	if storageID == "" || filename == "" {
		return nil, errors.New("storage id and filename are required")
	}

	key := fmt.Sprintf("%s/%s", storageID, filename)
	obj, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &DownloadResult{
		Body:           obj.Body,
		ContentType:    aws.ToString(obj.ContentType),
		ContentLength:  aws.ToInt64(obj.ContentLength),
		DownloadedFile: filename,
	}, nil
}

func (c *localFilemanagerConnection) prefixExists(ctx context.Context, storageID string) (bool, error) {
	prefix := storageID + "/"
	out, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(c.bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int32(1),
	})
	if err != nil {
		var nf *types.NoSuchBucket
		if errors.As(err, &nf) {
			return false, fmt.Errorf("bucket %s not found", c.bucket)
		}
		return false, err
	}
	return len(out.Contents) > 0, nil
}

func (c *localFilemanagerConnection) generateStorageID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := c.idLength
	if length <= 0 {
		length = 10
	}
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		idx := rand.IntN(len(letters))
		sb.WriteByte(letters[idx])
	}
	return sb.String()
}

func newS3Client() *s3.Client {
	creds := credentials.NewStaticCredentialsProvider(pkg.S3_ACCESS_KEY_ID, pkg.S3_SECRET_ACCESS_KEY, "")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if pkg.S3_ENDPOINT != "" {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               pkg.S3_ENDPOINT,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(pkg.S3_REGION),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	usePathStyle := strings.ToLower(pkg.S3_FORCE_PATH_STYLE) == "true"

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = usePathStyle
	})
}

func parseLength(val string, def int) int {
	i, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil || i <= 0 {
		return def
	}
	return i
}
