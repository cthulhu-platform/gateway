package diagnose

import (
	"github.com/cthulhu-platform/common/pkg/messages"
	"github.com/wagslane/go-rabbitmq"
)

type DiagnoseConnection interface {
	PublishDiagnose(msg []byte) error
	Close()
}

type rmqDiagnoseConnection struct {
	publisher *rabbitmq.Publisher
}

func NewRMQDiagnoseConn(conn *rabbitmq.Conn) (*rmqDiagnoseConnection, error) {
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

	return &rmqDiagnoseConnection{
		publisher: publisher,
	}, nil
}

func (c *rmqDiagnoseConnection) PublishDiagnose(msg []byte) error {
	return c.publisher.Publish(
		msg,
		[]string{messages.TopicDiagnoseServicesAll},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}

func (c *rmqDiagnoseConnection) Close() {
	c.publisher.Close()
}
