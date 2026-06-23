//go:build integration

package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGCPScout(t *testing.T) {
	//TODO: verify that a message can be sent to GCP
	message := New()
	message.CreatePubSubClient()
	defer message.Client.Close()

	msg := map[string]any{
		"SN":        "ABC123",
		"PN":        "PART-456",
		"Reading_1": 3,
		"Reading_2": 2,
	}
	assert.NoError(t, message.FormatMessage(msg))

	topicName := os.Getenv(TOPIC_ID_ENV_VAR)
	assert.NotEqual(t, "", topicName)
	err := message.SendMessageToPubSub(topicName)
	assert.NoError(t, err)
}
