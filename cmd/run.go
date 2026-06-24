/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/GreenStarMatter/zenzore/internal/server"
	"github.com/spf13/cobra"
)

var runPort string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the zenzore root server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if runPort != "" {
			if err := os.Setenv(server.PortEnvVar, runPort); err != nil {
				return fmt.Errorf("setting port env var: %w", err)
			}
		}

		s := server.NewServer()
		return s.Run()

	},
}

func init() {
	runCmd.Flags().StringVarP(&runPort, "port", "p", "", "port for the root server to listen on (overrides "+server.PortEnvVar+")")
	rootCmd.AddCommand(runCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
