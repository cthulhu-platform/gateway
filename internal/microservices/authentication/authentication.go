package authentication

import (
	"github.com/cthulhu-platform/common/pkg/messages"
	"github.com/wagslane/go-rabbitmq"
)

type AuthenticationConnection interface {
	PublishDiagnose(msg []byte) error
	Close()
}

type rmqAuthenticationConnection struct {
	publisher *rabbitmq.Publisher
}

func NewRMQAuthenticationConn(conn *rabbitmq.Conn) (*rmqAuthenticationConnection, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsExchangeName(messages.DiagnoseExchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		return nil, err
	}

	return &rmqAuthenticationConnection{
		publisher: publisher,
	}, nil
}

func (c *rmqAuthenticationConnection) PublishDiagnose(msg []byte) error {
	return c.publisher.Publish(
		msg,
		[]string{messages.TopicDiagnoseServicesAll},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}

func (c *rmqAuthenticationConnection) Close() {
	c.publisher.Close()
}
