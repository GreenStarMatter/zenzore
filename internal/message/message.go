package message

import (
	"cloud.google.com/go/pubsub/v2"
	"context"
	"encoding/json"
	"fmt"
	//	"google.golang.org/api/option"
	"os"
)

const PROJECT_ID_ENV_VAR = "ZENZOREPROJECTID"
const TOPIC_ID_ENV_VAR = "ZENZORETOPICID"

type PubSubMessage struct {
	Ctx     context.Context
	Message []byte
	Client  *pubsub.Client
}

func New() *PubSubMessage {
	ctx := context.Background()
	return &PubSubMessage{Ctx: ctx}
}

func (psm *PubSubMessage) CreatePubSubClient() error {
	projectId := os.Getenv(PROJECT_ID_ENV_VAR)
	if projectId == "" {
		return fmt.Errorf("failed to find %s env var", PROJECT_ID_ENV_VAR)
	}
	client, err := pubsub.NewClient(psm.Ctx, projectId)
	if err != nil {
		return fmt.Errorf("creating pubsub client: %w", err)
	}
	psm.Client = client
	return nil
}

func (psm *PubSubMessage) FormatMessage(msg map[string]any) error {

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	psm.Message = data
	return nil
}

func (psm *PubSubMessage) SendMessageToPubSub(topicName string) error {
	publisher := psm.Client.Publisher(topicName)
	result := publisher.Publish(psm.Ctx, &pubsub.Message{Data: psm.Message})
	publisher.Stop()
	err := psm.HandlePubSubResults(result)
	if err != nil {
		return err
	}
	return nil
}

func (psm *PubSubMessage) HandlePubSubResults(result *pubsub.PublishResult) error {
	//TODO: Handle for different types of errors
	//retry on connection failures
	//pass error on credential failures
	id, err := result.Get(psm.Ctx)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}
	fmt.Printf("Sent Message: %s\n", id)
	return nil
}

func (psm *PubSubMessage) AcceptGenericJson(incomingJson []byte) {
	psm.Message = incomingJson
}
