package zyztem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRemoveZyztem(t *testing.T) {
	//TODO: test create logic for zyztem
	//test that new process can be created
	//test that a new process cannot be created on top of an existing one
	//test that trying to remove a non-existent zyztem fails
	//test remove logic works and fails appropriately if zyztem is unavailable to be removed
	//make sure that PID ./srv reflects state of zyztems
	//New zyztem, Get Nav, Remove zyztem
	exists, err := CheckForExistingZyztem()
	assert.NoError(t, err)
	assert.True(t, exists)
	zyztem := New()
	assert.NoError(t, zyztem.Remove())
	assert.Error(t, zyztem.Remove())
	//test that navigator is properly initialized

}

func TestSensor(t *testing.T) {
	sensor := NewSensor("sn 1", "pn 1")
	assert.Equal(t, 1, sensor.Signal.GenerateSignalSample())
	sensor.UpdateSignal("constant", 2, 0)
	assert.Equal(t, 2, sensor.Signal.GenerateSignalSample())
}

func TestDevice(t *testing.T) {

	unassignedSensor := NewSensor("sn 1", "pn 1")
	device := NewDevice("sn 1", "pn 1")
	err := device.RemoveSensor(unassignedSensor)
	assert.Error(t, err)
	sensor := device.AddSensor("dummy sensor sn", "dummy part pn")
	err = device.RemoveSensor(unassignedSensor)
	assert.Error(t, err)
	err = device.RemoveSensor(sensor)
	assert.NoError(t, err)
	err = device.RemoveSensor(sensor)
	assert.Error(t, err)
	err = device.RemoveSensor(unassignedSensor)
	assert.Error(t, err)

	_ = device.AddSensor("sn 2", "pn 1")
	sensor3 := device.AddSensor("sn 2", "pn 3")

	sensor3.UpdateSignal("updated", 4, 3)
	jsonByteOut := device.ExportDeviceData()
	assert.Equal(t, []byte("bad match"), jsonByteOut)

}

func TestSubZyztemAddRemove(t *testing.T) {
	//TODO: test all subcomponents can be added to a zyztem
	//Make sure that devices and sensors can be appropriately added
	//New zyztem, add subcomponents, remove subcomponents
	//Need to rethink this.  Now that everything is not attached in a graph structure directly, I cannot globably access.  Perhaps I could do this if I moved the nav structure in too
	//This makes the atomic removals a bit more esoteric as they can be broken apart from each other by just disconnecting the chain.  This does make modular testing a bit easier, but makes the ful zyztem test more difficult

	//TODO: Break this test into a bunch of smaller tests which test each component (Make sure each component change works)
	//TODO: Make an integrated test which tests individual changes against the global structure (break the chain and verify the global state isn't changed)

	exists, err := CheckForExistingZyztem()
	assert.NoError(t, err)
	assert.True(t, exists)
	zyztem := New()
	device := zyztem.AddDevice("dummy device serial", "dummy device part")
	sensor := device.AddSensor("dummy sensor sn", "dummy part pn")
	sensor2 := device.AddSensor("dummy sensor sn 2", "dummy part pn 2")
	device2 := zyztem.AddDevice("dummy device serial 2", "dummy device part 2")
	sensor3 := device2.AddSensor("dummy sensor sn", "dummy part pn 2")

	assert.Equal(t, 2, len(device.Sensors))
	assert.NoError(t, device.RemoveSensor(sensor))
	assert.Equal(t, 1, len(device.Sensors))
	assert.Equal(t, sensor2, device.Sensors[0])

	assert.Error(t, device.RemoveSensor(sensor))

	assert.NoError(t, zyztem.RemoveDevice(device2))
	assert.Error(t, device2.RemoveSensor(sensor3))

	assert.NoError(t, zyztem.Remove())
	assert.Error(t, zyztem.Remove())
}

func TestZyztemPrint(t *testing.T) {
	//TODO: test print logic coming from zyztem can be properly retrieved
	//New zyztem, add subcomponents, print
	ExpectedZyztemOutput := "Auto Fail"
	exists, err := CheckForExistingZyztem()
	assert.NoError(t, err)
	assert.True(t, exists)
	zyztem := New()
	device := zyztem.AddDevice("dummy device serial", "dummy device part")
	_ = device.AddSensor("sn 1", "pn 1")
	_ = device.AddSensor("sn 2", "pn 2")
	device2 := zyztem.AddDevice("dummy device serial", "dummy device part 2")
	_ = device2.AddSensor("sn 3", "pn 1")
	assert.Equal(t, ExpectedZyztemOutput, zyztem.Print())

}
