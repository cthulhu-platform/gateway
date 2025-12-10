package authentication

import "github.com/wagslane/go-rabbitmq"

type AuthenticationConnection interface{}

type rmqAuthenticationConnection struct{}

func NewRMQAuthenticationConn(conn *rabbitmq.Conn) *rmqAuthenticationConnection {
	return &rmqAuthenticationConnection{}
}
