package filemanager

import "github.com/wagslane/go-rabbitmq"

type FilemanagerConnection interface{}

type rmqFilemanagerConnection struct{}

func NewRMQFilemanagerConn(conn *rabbitmq.Conn) *rmqFilemanagerConnection {
	return &rmqFilemanagerConnection{}
}
