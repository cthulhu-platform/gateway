package service

type rmqFileManagerService struct{}

func NewRMQFileManagerService() FileManagerService {
	return &rmqFileManagerService{}
}
