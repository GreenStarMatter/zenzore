package signal

import (
	"encoding/json"
	"log"
	"time"
)

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

func GenerateConstantSignal(expected int) Sample {
	return Sample(expected)
}

func FormatTimeStamp(tstamp time.Time) string {
	return tstamp.UTC().Format(time.RFC3339Nano)
}

func (incomingSignal *Signal) GenerateSignalSample() {
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
}

func (device *Device) AddSensor(sensor *Sensor) {
	device.Sensors = append(device.Sensors, sensor)
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

func NewSensor(sn, pn string, signal *Signal) *Sensor {
	sensor := &Sensor{
		SN:     sn,
		PN:     pn,
		Signal: signal,
	}
	signal.ParentSensor = sensor
	return sensor
}

func NewSignal(typeSignal string, expected int, random int) *Signal {
	return &Signal{Type: typeSignal, ExpectedValue: expected, RandomValue: random}
}
