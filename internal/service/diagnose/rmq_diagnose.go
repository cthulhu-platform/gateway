package diagnose

import (
	"encoding/json"
	"fmt"

	"github.com/cthulhu-platform/common/pkg/messages"
	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/google/uuid"
)

type rmqDiagnoseService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewRMQDiagnoseService(conns *microservices.ServiceConnectionContainer) *rmqDiagnoseService {
	return &rmqDiagnoseService{
		conns: conns,
	}
}

func (s *rmqDiagnoseService) ServiceFanoutTest() (string, error) {
	// Generate transaction ID
	transactionID := uuid.New().String()

	// Create diagnostic message
	diagnoseMsg := messages.DiagnoseMessage{
		TransactionID: transactionID,
		Operation:     "all", // Operation type: "all" means check if service is up
		Message:       "Health check - are you up?",
	}

	// Marshal to JSON
	messageBody, err := json.Marshal(diagnoseMsg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish message to the exchange through the diagnose connection
	if err := s.conns.Diagnose.PublishDiagnose(messageBody); err != nil {
		return "", fmt.Errorf("failed to publish message: %w", err)
	}

	return transactionID, nil
}
