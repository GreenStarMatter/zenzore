package zyztem

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//create a New zyztem which spawns a process
//zyztem print state
//zyztem create slots for all device and sensor data

const BAD_SAMPLE = Sample(-9999)

type Sample int //`json:"SampleValue"`

type Signal struct {
	Type          string  `json:"SignalType"`
	ParentSensor  *Sensor `json:"-"`
	ExpectedValue int     `json:"-"`
	RandomValue   int     `json:"-"`
}

type Sensor struct {
	SN             string  `json:"SensorSN"`
	PN             string  `json:"SensorPN"`
	LastSample     Sample  `json:"SampleValue"`
	LastSampleTime string  `json:"SampleValueTime"`
	Signal         *Signal `json:"-"`
}

type Device struct {
	SN               string    `json:"DeviceSN"`
	PN               string    `json:"DevicePN"`
	DeviceExportTime string    `json:"DeviceExportTime"`
	Sensors          []*Sensor `json:"Sensors"`
}

type Zyztem struct {
	//Root     *ZyztemNode
	ID      string
	Devices []*Device
}

//type ZyztemNode struct {
//	Children []*ZyztemNode
//}

func CheckForExistingZyztem() (bool, error) {
	return true, fmt.Errorf("unable to access zyztem location")
}

func New() *Zyztem {
	//NewNode := &ZyztemNode{Children: make([]*ZyztemNode, 0)}
	//return &Zyztem{Root: NewNode}

	return &Zyztem{Devices: make([]*Device, 0)}

}

func (zyztem *Zyztem) AddDevice(sn, pn string) *Device {
	device := NewDevice(sn, pn)
	zyztem.Devices = append(zyztem.Devices, device)
	return device
}

func (zyztem *Zyztem) RemoveDevice(device *Device) error {
	return fmt.Errorf("failed to remove device from zyztem")
}

//func (*ZyztemNode) Remove(child *ZyztemNode) error {
//	return fmt.Errorf("failed to delete node")
//}

func (*Zyztem) Print() string {
	return "No Zyztem Logic Yet"
}

func (device *Device) AddSensor(sn, pn string) *Sensor {
	sensor := NewSensor(sn, pn)
	device.Sensors = append(device.Sensors, sensor)
	return sensor
}

func (device *Device) RemoveSensor(sensor *Sensor) error {
	return fmt.Errorf("failed to remove sensor from device")
}

func (device *Device) ExportDeviceData() []byte {
	device.DeviceExportTime = FormatTimeStamp(time.Now())
	jsonData, err := json.Marshal(device)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func NewDevice(sn, pn string) *Device {
	sensors := []*Sensor{}
	return &Device{SN: sn, PN: pn, Sensors: sensors}
}

func NewSensor(sn, pn string) *Sensor {
	sensor := &Sensor{
		SN: sn,
		PN: pn,
	}
	sensor.Signal = sensor.AddSignal("constant", 1, 0)
	sensor.Signal.ParentSensor = sensor
	return sensor
}

func (sensor *Sensor) AddSignal(typeSignal string, expected int, random int) *Signal {
	signal := NewSignal(typeSignal, expected, random)
	return signal
}

func (sensor *Sensor) UpdateSignal(typeSignal string, expected int, random int) *Signal {
	signal := NewSignal(typeSignal, expected, random)
	sensor.Signal = signal
	signal.ParentSensor = sensor
	return signal
}

func NewSignal(typeSignal string, expected int, random int) *Signal {
	return &Signal{Type: typeSignal, ExpectedValue: expected, RandomValue: random}
}

func GenerateConstantSignal(expected int) Sample {
	return Sample(expected)
}

func FormatTimeStamp(tstamp time.Time) string {
	return tstamp.UTC().Format(time.RFC3339Nano)
}

func (incomingSignal *Signal) GenerateSignalSample() Sample {
	incomingSignal.ParentSensor.LastSampleTime = FormatTimeStamp(time.Now())
	if incomingSignal.ParentSensor == nil {
		incomingSignal.ParentSensor.LastSample = Sample(BAD_SAMPLE)
	}
	switch incomingSignal.Type {
	case "constant":
		incomingSignal.ParentSensor.LastSample = GenerateConstantSignal(incomingSignal.ExpectedValue)
	default:
		incomingSignal.ParentSensor.LastSample = Sample(BAD_SAMPLE)
	}
	return incomingSignal.ParentSensor.LastSample
}
