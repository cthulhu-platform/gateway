package file

type FileService interface {
	UploadFile() error
	DownloadFile() error
	RetrieveFileBucket() error
}
