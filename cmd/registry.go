package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/GreenStarMatter/zenzore/internal/registry"
	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Push entries to the zenzore registry",
}

// --- zyztem ---

var zyztemRandomFlag bool
var zyztemNameFlag string

var zyztemCmd = &cobra.Command{
	Use:   "zyztem",
	Short: "Push a zyztem entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !zyztemRandomFlag {
			return fmt.Errorf("nothing to do: pass --random to generate and push an entry")
		}
		ctx := context.Background()
		db, cleanup, err := registry.Connect(ctx)
		if err != nil {
			return fmt.Errorf("connecting to registry db: %w", err)
		}
		defer cleanup()

		z, err := db.CreateRandomZyztem(ctx, zyztemNameFlag)
		if err != nil {
			return fmt.Errorf("creating zyztem: %w", err)
		}
		return json.NewEncoder(os.Stdout).Encode(z)
	},
}

// --- device ---

var deviceRandomFlag bool
var devicePNFlag string
var deviceSNFlag string

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Push a device entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !deviceRandomFlag {
			return fmt.Errorf("nothing to do: pass --random to generate and push an entry")
		}
		ctx := context.Background()
		db, cleanup, err := registry.Connect(ctx)
		if err != nil {
			return fmt.Errorf("connecting to registry db: %w", err)
		}
		defer cleanup()

		d, err := db.CreateRandomDevice(ctx, devicePNFlag, deviceSNFlag)
		if err != nil {
			return fmt.Errorf("creating device: %w", err)
		}
		return json.NewEncoder(os.Stdout).Encode(d)
	},
}

// --- sensor ---

var sensorRandomFlag bool
var sensorNameFlag string
var sensorPNFlag string
var sensorSNFlag string

var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "Push a sensor entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !sensorRandomFlag {
			return fmt.Errorf("nothing to do: pass --random to generate and push an entry")
		}
		ctx := context.Background()
		db, cleanup, err := registry.Connect(ctx)
		if err != nil {
			return fmt.Errorf("connecting to registry db: %w", err)
		}
		defer cleanup()

		s, err := db.CreateRandomSensor(ctx, sensorNameFlag, sensorPNFlag, sensorSNFlag)
		if err != nil {
			return fmt.Errorf("creating sensor: %w", err)
		}
		return json.NewEncoder(os.Stdout).Encode(s)
	},
}

func init() {
	zyztemCmd.Flags().BoolVar(&zyztemRandomFlag, "random", false, "generate and push a zyztem with random attributes")
	zyztemCmd.Flags().StringVar(&zyztemNameFlag, "name", "", "override the zyztem_id (random if omitted)")

	deviceCmd.Flags().BoolVar(&deviceRandomFlag, "random", false, "generate and push a device with random attributes")
	deviceCmd.Flags().StringVar(&devicePNFlag, "pn", "", "override the device_pn (random if omitted)")
	deviceCmd.Flags().StringVar(&deviceSNFlag, "sn", "", "override the device_sn (random if omitted)")

	sensorCmd.Flags().BoolVar(&sensorRandomFlag, "random", false, "generate and push a sensor with random attributes")
	sensorCmd.Flags().StringVar(&sensorNameFlag, "name", "", "override the sensor_id (random if omitted)")
	sensorCmd.Flags().StringVar(&sensorPNFlag, "pn", "", "override the sensor_pn (random if omitted)")
	sensorCmd.Flags().StringVar(&sensorSNFlag, "sn", "", "override the sensor_sn (random if omitted)")

	registryCmd.AddCommand(zyztemCmd, deviceCmd, sensorCmd)
	rootCmd.AddCommand(registryCmd)
}
