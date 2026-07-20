package registry

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	cloudsqlconn "cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
)

const (
	ProjectIDEnvVar  = "TF_VAR_project_id"
	RegionEnvVar     = "TF_VAR_region"
	InstanceIDEnvVar = "ZENZOREDBINSTANCE"
)

var adjectives = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel"}
var manufacturers = []string{"Initech", "Umbrella", "Cyberdyne", "Weyland", "Palantir"}
var models = []string{"MK1", "MK2", "MK3", "X100", "X200", "Z9"}
var revisions = []string{"A", "B", "C", "D"}

func randomPick(opts []string) string { return opts[rand.Intn(len(opts))] }
func randomSN() string                { return fmt.Sprintf("%08d", rand.Intn(100000000)) }
func randomPN(prefix string) string   { return fmt.Sprintf("%s-%04d", prefix, rand.Intn(10000)) }
func randomID() string                { return fmt.Sprintf("%s-%04d", randomPick(adjectives), rand.Intn(10000)) }

type DB struct {
	conn *sql.DB
}

func Connect(ctx context.Context) (*DB, func(), error) {
	projectID := os.Getenv(ProjectIDEnvVar)
	region := os.Getenv(RegionEnvVar)
	instanceID := os.Getenv(InstanceIDEnvVar)
	instanceName := fmt.Sprintf("%s:%s:%s", projectID, region, instanceID)
	dbUser := "registry_admin"
	dbPass := os.Getenv("TF_VAR_cloudsql_password")

	cleanupDriver, err := pgxv5.RegisterDriver("cloudsql-postgres",
		cloudsqlconn.WithDefaultDialOptions(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("registering cloudsql driver: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=zenzore_registry sslmode=disable",
		instanceName, dbUser, dbPass,
	)

	db, err := sql.Open("cloudsql-postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("opening db: %w", err)
	}

	cleanup := func() {
		db.Close()
		cleanupDriver()
	}
	return &DB{conn: db}, cleanup, nil
}

type Zyztem struct {
	ZyztemKey    int       `json:"zyztem_key"`
	ZyztemID     string    `json:"zyztem_id"`
	ZyztemPN     string    `json:"zyztem_pn"`
	ZyztemSN     string    `json:"zyztem_sn"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	Revision     string    `json:"revision"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateRandomZyztem builds a Zyztem with fully randomized attributes
// (including the ID) and inserts it.
func (db *DB) CreateRandomZyztem(ctx context.Context, id string) (*Zyztem, error) {
	if id == "" {
		id = randomID()
	}
	z := &Zyztem{
		ZyztemID:     id,
		ZyztemPN:     randomPN("ZYZ"),
		ZyztemSN:     randomSN(),
		Manufacturer: randomPick(manufacturers),
		Model:        randomPick(models),
		Revision:     randomPick(revisions),
	}
	row := db.conn.QueryRowContext(ctx, `
		INSERT INTO dim_zyztem (zyztem_id, zyztem_pn, zyztem_sn, manufacturer, model, revision)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING zyztem_key, created_at
	`, z.ZyztemID, z.ZyztemPN, z.ZyztemSN, z.Manufacturer, z.Model, z.Revision)
	if err := row.Scan(&z.ZyztemKey, &z.CreatedAt); err != nil {
		return nil, fmt.Errorf("inserting zyztem: %w", err)
	}
	return z, nil
}

var firmwareVersions = []string{"1.0.0", "1.2.0", "2.0.0", "2.1.3", "3.0.0"}
var hardwareRevisions = []string{"RevA", "RevB", "RevC", "RevD"}
var sensorTypes = []string{"temperature", "pressure", "humidity", "vibration", "flow"}
var engineeringUnits = []string{"C", "F", "psi", "%RH", "Hz", "m3/h"}

func randomOEMPN() string { return randomPN("OEM") }
func randomOEMSN() string { return randomSN() }

type Device struct {
	DeviceKey        int       `json:"device_key"`
	DevicePN         string    `json:"device_pn"`
	DeviceSN         string    `json:"device_sn"`
	Manufacturer     string    `json:"manufacturer"`
	OEMDevicePN      string    `json:"oem_device_pn"`
	OEMDeviceSN      string    `json:"oem_device_sn"`
	FirmwareVersion  string    `json:"firmware_version"`
	HardwareRevision string    `json:"hardware_revision"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateRandomDevice builds a Device with randomized attributes and inserts
// it. If pn or sn is non-empty, it's used instead of a random one. OEM
// fields are always randomized.
func (db *DB) CreateRandomDevice(ctx context.Context, pn, sn string) (*Device, error) {
	if pn == "" {
		pn = randomPN("DEV")
	}
	if sn == "" {
		sn = randomSN()
	}
	d := &Device{
		DevicePN:         pn,
		DeviceSN:         sn,
		Manufacturer:     randomPick(manufacturers),
		OEMDevicePN:      randomOEMPN(),
		OEMDeviceSN:      randomOEMSN(),
		FirmwareVersion:  randomPick(firmwareVersions),
		HardwareRevision: randomPick(hardwareRevisions),
	}
	row := db.conn.QueryRowContext(ctx, `
		INSERT INTO dim_device (device_pn, device_sn, manufacturer, oem_device_pn, oem_device_sn, firmware_version, hardware_revision)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING device_key, created_at
	`, d.DevicePN, d.DeviceSN, d.Manufacturer, d.OEMDevicePN, d.OEMDeviceSN, d.FirmwareVersion, d.HardwareRevision)
	if err := row.Scan(&d.DeviceKey, &d.CreatedAt); err != nil {
		return nil, fmt.Errorf("inserting device: %w", err)
	}
	return d, nil
}

type Sensor struct {
	SensorKey        int       `json:"sensor_key"`
	SensorID         string    `json:"sensor_id"`
	SensorPN         string    `json:"sensor_pn"`
	SensorSN         string    `json:"sensor_sn"`
	Manufacturer     string    `json:"manufacturer"`
	OEMSensorPN      string    `json:"oem_sensor_pn"`
	OEMSensorSN      string    `json:"oem_sensor_sn"`
	SensorType       string    `json:"sensor_type"`
	EngineeringUnits string    `json:"engineering_units"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateRandomSensor builds a Sensor with randomized attributes and inserts
// it. If id, pn, or sn is non-empty, it's used instead of a random one. OEM
// fields are always randomized.
func (db *DB) CreateRandomSensor(ctx context.Context, id, pn, sn string) (*Sensor, error) {
	if id == "" {
		id = randomID()
	}
	if pn == "" {
		pn = randomPN("SEN")
	}
	if sn == "" {
		sn = randomSN()
	}
	s := &Sensor{
		SensorID:         id,
		SensorPN:         pn,
		SensorSN:         sn,
		Manufacturer:     randomPick(manufacturers),
		OEMSensorPN:      randomOEMPN(),
		OEMSensorSN:      randomOEMSN(),
		SensorType:       randomPick(sensorTypes),
		EngineeringUnits: randomPick(engineeringUnits),
	}
	row := db.conn.QueryRowContext(ctx, `
		INSERT INTO dim_sensor (sensor_id, sensor_pn, sensor_sn, manufacturer, oem_sensor_pn, oem_sensor_sn, sensor_type, engineering_units)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING sensor_key, created_at
	`, s.SensorID, s.SensorPN, s.SensorSN, s.Manufacturer, s.OEMSensorPN, s.OEMSensorSN, s.SensorType, s.EngineeringUnits)
	if err := row.Scan(&s.SensorKey, &s.CreatedAt); err != nil {
		return nil, fmt.Errorf("inserting sensor: %w", err)
	}
	return s, nil
}
