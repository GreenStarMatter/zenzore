package navigator

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
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

func TestFileFindInOperatingEnvironment(t *testing.T) {
	//attempt to find files in srv path (do happy and unhappy path)
	//Place a file manually

	//TODO: Make test which finds files in server location (Should be in PWD or Stored in environment variable to point at data location)
	//verify folder location exists
	//verify folder empty when starts (terminate test if not empty because this means that the system is already operating)
	//Add dummy files with mocked PID
	//Read all dummy files and instance navigator
	//Verify navigator structure
	//verify PID read in from each file

	//Create a rootNode (this logic should implicitly locate or create a /srv dir)
	//fileNameArray := os.ListFiles("./srv") != 0 then fail and stop test
	//write 5 files with names 1..5

	//iterate through fileNameArray and add a NavigatorNode for each PID
	//assert that the rootNode has 5 children
	//iterate through Children and assert that each ID is in 1..5 and none repeat
	dir := t.TempDir() // fresh, guaranteed-empty dir; auto-cleaned after test

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(entries), "expected fresh dir to start empty")

	expectedIDs := []string{"1", "2", "3", "4", "5"}
	for _, name := range expectedIDs {
		f, err := os.Create(filepath.Join(dir, name))
		assert.NoError(t, err)
		assert.NoError(t, f.Close())
	}

	_, rootNode, err := LoadFromOperatingEnvironment(dir)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedIDs), len(rootNode.Children))

	seen := make(map[string]bool)
	for _, child := range rootNode.Children {
		id := string(child.ID)
		assert.Contains(t, expectedIDs, id, "unexpected ID %q found", id)
		assert.False(t, seen[id], "duplicate ID %q found", id)
		seen[id] = true
	}

	missingDir := filepath.Join(t.TempDir(), "does-not-exist")
	_, _, err = LoadFromOperatingEnvironment(missingDir)
	assert.Error(t, err)

}

//TODO: Add an integration test in the actual command test which verifies that the PID is running when a real item is created
