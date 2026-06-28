package zyztem

import (
	"fmt"
	"time"
)

const BAD_SAMPLE = Sample(-9999)

type Sample int

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
	ID      string
	Devices []*Device
}

func FormatTimeStamp(tstamp time.Time) string {
	return tstamp.UTC().Format(time.RFC3339Nano)
}

func New() *Zyztem {
	return &Zyztem{Devices: make([]*Device, 0)}
}

func NewDevice(sn, pn string) *Device {
	sensors := []*Sensor{}
	return &Device{SN: sn, PN: pn, Sensors: sensors}
}

func (zyztem *Zyztem) AddDevice(sn, pn string) (*Device, error) {
	for _, d := range zyztem.Devices {
		if d.SN == sn && d.PN == pn {
			return nil, fmt.Errorf("device with SN %q and PN %q already exists in zyztem", sn, pn)
		}
	}
	device := NewDevice(sn, pn)
	zyztem.Devices = append(zyztem.Devices, device)
	return device, nil
}

func (z *Zyztem) RemoveDevice(sn, pn string) error {
	for i, d := range z.Devices {
		if d.SN == sn && d.PN == pn {
			z.Devices = append(z.Devices[:i], z.Devices[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("device not found")
}

func (z *Zyztem) FindDevice(sn, pn string) (*Device, bool) {
	for _, d := range z.Devices {
		if d.SN == sn && d.PN == pn {
			return d, true
		}
	}
	return nil, false
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

func (device *Device) AddSensor(sn, pn string) (*Sensor, error) {

	for _, d := range device.Sensors {
		if d.SN == sn && d.PN == pn {
			return nil, fmt.Errorf("sensor with SN %q and PN %q already exists in device", sn, pn)
		}
	}

	sensor := NewSensor(sn, pn)
	device.Sensors = append(device.Sensors, sensor)
	return sensor, nil
}

func (d *Device) FindSensor(sn, pn string) (*Sensor, bool) {
	for _, s := range d.Sensors {
		if s.SN == sn && s.PN == pn {
			return s, true
		}
	}
	return nil, false
}

func (d *Device) RemoveSensor(sn, pn string) error {
	for i, s := range d.Sensors {
		if s.SN == sn && s.PN == pn {
			d.Sensors = append(d.Sensors[:i], d.Sensors[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("sensor not found")
}

func NewSignal(typeSignal string, expected int, random int) *Signal {
	return &Signal{Type: typeSignal, ExpectedValue: expected, RandomValue: random}
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

func GenerateConstantSignal(expected int) Sample {
	return Sample(expected)
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
