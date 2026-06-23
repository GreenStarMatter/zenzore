package message

import (
	"cloud.google.com/go/pubsub/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOutputMessage(t *testing.T) {
	//TODO: verify that a message is properly formed
	message := New()
	msg := map[string]any{
		"SN":        "ABC123",
		"PN":        "PART-456",
		"Reading_1": 3,
		"Reading_2": 2,
	}
	assert.NoError(t, message.FormatMessage(msg))
	assert.Equal(t, []byte("Fail For Sure"), message.Message)
	message.AcceptGenericJson(message.Message)
	assert.Equal(t, []byte("Fail For Sure"), message.Message)

}

func TestGCPConnectionFailureHandling(t *testing.T) {

	//TODO: verify that app can handle a failed connection to GCP
	//I'm thinking that this sends a "dormant" command to rest of zyztems
	//This just stops them from advancing their clocks and keeps them in a waiting state
	//Will have to fill out result (maybe with temp interface that mocks different errors)
	message := New()
	var result *pubsub.PublishResult
	err := message.HandlePubSubResults(result)
	assert.NoError(t, err)
}
