package navigator

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRemoveNode(t *testing.T) {
	_, rootNode := New()
	systemNode, err := rootNode.Add()
	assert.NoError(t, err)
	deviceNode, err := systemNode.Add()
	assert.NoError(t, err)
	sensorNode, err := deviceNode.Add()
	assert.NoError(t, err)
	sensorNode2, err := deviceNode.Add()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(deviceNode.Children))
	err = deviceNode.Remove(sensorNode)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(deviceNode.Children))
	assert.Equal(t, sensorNode2, deviceNode.Children[0])

	//assert.Error(t, rootNode.Remove(sensorNode))
	//assert.NoError(t, rootNode.Remove(systemNode))
}

func TestListSystems(t *testing.T) {
	/*attempt to list all items in the navigation
	test all items are stored in correct data structure
	test string print is correct*/

	expectedString := "0: 0\n1: 0\n2: 0\n3: 0\n4: 0\n"
	_, rootNode := New()
	for range 5 {
		_, _ = rootNode.Add()
	}

	assert.Equal(t, expectedString, rootNode.List())
}

func TestNavigation(t *testing.T) {
	/*attempt to navigate between different levels from system to sensor*/
	navigator, rootNode := New()
	navigator.Set(rootNode)
	systemNode, _ := rootNode.Add()
	_, _ = rootNode.Add()
	deviceNode, _ := systemNode.Add()
	_, _ = systemNode.Add()
	sensorNode, _ := deviceNode.Add()
	_, _ = deviceNode.Add()

	assert.Error(t, navigator.Up())

	assert.Error(t, navigator.Down(sensorNode))
	assert.NoError(t, navigator.Down(systemNode))

	navigator.Down(deviceNode)
	navigator.Down(sensorNode)
	assert.Error(t, navigator.Down(sensorNode))

	navigator.Up()
	navigator.Up()
	navigator.Up()
	assert.Equal(t, rootNode, navigator.CurrentNode)

}

func TestLoadFromServer(t *testing.T) {
	expectedIDs := []string{"1", "2", "3", "4", "5"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var zyztems []zyztemSummary
		for _, id := range expectedIDs {
			zyztems = append(zyztems, zyztemSummary{ID: id})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(zyztems)
	}))
	defer mockServer.Close()

	_, rootNode, err := LoadFromServer(mockServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedIDs), len(rootNode.Children))

	seen := make(map[string]bool)
	for _, child := range rootNode.Children {
		id := string(child.ID)
		assert.Contains(t, expectedIDs, id, "unexpected ID %q found", id)
		assert.False(t, seen[id], "duplicate ID %q found", id)
		seen[id] = true
	}

	_, _, err = LoadFromServer("http://localhost:1")
	assert.Error(t, err)
}

//TODO: Add an integration test in the actual command test which verifies that the PID is running when a real item is created
