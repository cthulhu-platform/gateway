package authentication

import "github.com/wagslane/go-rabbitmq"

type AuthenticationConnection interface {
	Close()
}

type rmqAuthenticationConnection struct{}

func NewRMQAuthenticationConn(conn *rabbitmq.Conn) (*rmqAuthenticationConnection, error) {
	// TODO: Add authentication-specific publishers/consumers here
	return &rmqAuthenticationConnection{}, nil
}

func (c *rmqAuthenticationConnection) Close() {
	// TODO: Close any authentication-specific connections
}
