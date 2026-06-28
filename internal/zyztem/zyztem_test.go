package zyztem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSensor(t *testing.T) {
	sensor := NewSensor("sn 1", "pn 1")
	assert.Equal(t, Sample(1), sensor.Signal.GenerateSignalSample())
	sensor.UpdateSignal("constant", 2, 0)
	assert.Equal(t, Sample(2), sensor.Signal.GenerateSignalSample())
}

func TestDevice(t *testing.T) {

	unassignedSensor := NewSensor("sn 1", "pn 1")
	device := NewDevice("sn 1", "pn 1")
	err := device.RemoveSensor(unassignedSensor.SN, unassignedSensor.PN)
	assert.Error(t, err)
	sensor, err := device.AddSensor("dummy sensor sn", "dummy part pn")
	assert.NoError(t, err)
	_, err = device.AddSensor("dummy sensor sn", "dummy part pn")
	assert.Error(t, err)
	err = device.RemoveSensor(unassignedSensor.SN, unassignedSensor.PN)
	assert.Error(t, err)
	err = device.RemoveSensor(sensor.SN, sensor.PN)
	assert.NoError(t, err)
	err = device.RemoveSensor(sensor.SN, sensor.PN)
	assert.Error(t, err)
	err = device.RemoveSensor(unassignedSensor.SN, unassignedSensor.PN)
	assert.Error(t, err)

	_, _ = device.AddSensor("sn 2", "pn 1")
	sensor3, _ := device.AddSensor("sn 2", "pn 3")

	sensor3.UpdateSignal("updated", 4, 3)
}

func TestSubZyztemAddRemove(t *testing.T) {

	zyztem := New()
	device, err := zyztem.AddDevice("dummy device serial", "dummy device part")
	assert.NoError(t, err)
	_, err = zyztem.AddDevice("dummy device serial", "dummy device part")
	assert.Error(t, err)
	sensor, _ := device.AddSensor("dummy sensor sn", "dummy part pn")
	sensor2, _ := device.AddSensor("dummy sensor sn 2", "dummy part pn 2")
	device2, err := zyztem.AddDevice("dummy device serial 2", "dummy device part 2")
	assert.NoError(t, err)
	_, _ = device2.AddSensor("dummy sensor sn", "dummy part pn 2")

	assert.Equal(t, 2, len(device.Sensors))
	assert.NoError(t, device.RemoveSensor(sensor.SN, sensor.PN))
	assert.Equal(t, 1, len(device.Sensors))
	assert.Equal(t, sensor2, device.Sensors[0])

	assert.Error(t, device.RemoveSensor(sensor.SN, sensor.PN))

	assert.NoError(t, zyztem.RemoveDevice(device2.SN, device2.PN))
	assert.Error(t, zyztem.RemoveDevice(device2.SN, device2.PN))

}
