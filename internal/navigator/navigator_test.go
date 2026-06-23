package navigator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddRemoveNode(t *testing.T) {
	_, rootNode := New()
	systemNode := rootNode.Add()
	deviceNode := systemNode.Add()
	sensorNode := deviceNode.Add()
	sensorNode2 := deviceNode.Add()

	assert.Equal(t, 2, len(deviceNode.Children))
	deviceNode.Remove(sensorNode)
	assert.Equal(t, 1, len(deviceNode.Children))
	assert.Equal(t, sensorNode2, deviceNode.Children[0])

	//assert.Error(t, rootNode.Remove(sensorNode))
	//assert.NoError(t, rootNode.Remove(systemNode))
}

func TestListSystems(t *testing.T) {
	/*attempt to list all items in the navigation
	test all items are stored in correct data structure
	test string print is correct*/
	_, rootNode := New()
	//for i := 0; i < 5; i++ {
	//	_ = rootNode.Add()
	//}

	assert.Equal(t, "always fail", rootNode.List())
	//TODO: add sprint to match system list
}

func TestNavigation(t *testing.T) {
	/*attempt to navigate between different levels from system to sensor*/
	navigator, rootNode := New()
	navigator.Set(rootNode)
	systemNode := rootNode.Add()
	_ = rootNode.Add()
	deviceNode := systemNode.Add()
	_ = systemNode.Add()
	sensorNode := deviceNode.Add()
	_ = deviceNode.Add()

	assert.Error(t, navigator.Up())

	assert.Error(t, navigator.Down(sensorNode.ID))
	assert.NoError(t, navigator.Down(systemNode.ID))

	navigator.Down(deviceNode.ID)
	navigator.Down(sensorNode.ID)
	assert.Error(t, navigator.Down(sensorNode.ID))

	navigator.Up()
	navigator.Up()
	navigator.Up()
	assert.Equal(t, rootNode, navigator.CurrentNode)

}

func TestFileFindInOperatingEnvironment(t *testing.T) {
	//attempt to find files in srv path (do happy and unhappy path)
	//Place a file manually

	//TODO: Make test which finds files in server location
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
	//TODO: Add an integration test in the actual command test which verifies that the PID is running when a real item is created

}
