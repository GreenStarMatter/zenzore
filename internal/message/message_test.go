package message

import (
	"encoding/json"
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

	err := message.FormatMessage(msg)
	assert.NoError(t, err)

	var got map[string]any
	err = json.Unmarshal(message.Message, &got)
	assert.NoError(t, err)

	assert.Equal(t, "ABC123", got["SN"])
	assert.Equal(t, "PART-456", got["PN"])
	assert.Equal(t, float64(3), got["Reading_1"])
	assert.Equal(t, float64(2), got["Reading_2"])

	newJSON := []byte(`{"SN":"XYZ789","PN":"PART-999","Reading_1":10,"Reading_2":20}`)
	message.AcceptGenericJson(newJSON)
	assert.Equal(t, newJSON, message.Message)

	var got2 map[string]any
	err = json.Unmarshal(message.Message, &got2)
	assert.NoError(t, err)
	assert.Equal(t, "XYZ789", got2["SN"])
}
