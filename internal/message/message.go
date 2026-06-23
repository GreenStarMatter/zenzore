package message

import (
	"cloud.google.com/go/pubsub/v2"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
)

//Reformat needs to essentially be
//Make message
//Connect to sender
//Send message

const ZENSOREKEY_ENV_VAR = "ZENZOREKEY"
const PROJECT_ID_ENV_VAR = "ZENZOREPROJECTID"
const TOPIC_ID_ENV_VAR = "ZENZORETOPICID"

type PubSubMessage struct {
	Ctx     context.Context
	Message []byte
	Client  *pubsub.Client
}

func (psm *PubSubMessage) CreatePubSubClient() {
	keyPath := os.Getenv(ZENSOREKEY_ENV_VAR)
	if keyPath == "" {
		fmt.Printf("Failed to find zenzore key env var\n")
		os.Exit(1)
	}

	projectId := os.Getenv(PROJECT_ID_ENV_VAR)
	if projectId == "" {
		fmt.Printf("Failed to find projectId env var\n")
		os.Exit(1)
	}
	client, err := pubsub.NewClient(psm.Ctx, projectId,
		option.WithAuthCredentialsFile(option.ServiceAccount, keyPath))
	if err != nil {
		log.Fatal(err)
	}
	psm.Client = client
}

func (psm *PubSubMessage) FormatMessage(msg map[string]any) error {

	data, err := json.Marshal(msg)
	if err != nil {
		return err
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
	/*
		id, err := result.Get(psm.Ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Sent Message: %s\n", id)
		return nil
	*/
	return fmt.Errorf("not implemented yet")

}

func (psm *PubSubMessage) AcceptGenericJson(incomingJson []byte) {
	psm.Message = incomingJson
}

func New() *PubSubMessage {
	ctx := context.Background()
	return &PubSubMessage{Ctx: ctx}
}

func main() {
	//TODO: Move this into application logic
	message := New()
	message.CreatePubSubClient()
	defer message.Client.Close()

	msg := map[string]any{
		"SN":        "ABC123",
		"PN":        "PART-456",
		"Reading_1": 3,
		"Reading_2": 2,
	}
	message.FormatMessage(msg)

	topicName := os.Getenv(TOPIC_ID_ENV_VAR)
	if topicName == "" {
		fmt.Printf("Failed to find topic env var\n")
		os.Exit(1)
	}
	message.SendMessageToPubSub(topicName)
}
