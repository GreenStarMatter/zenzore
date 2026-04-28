/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/GreenStarMatter/zenzore/internal/signal"
	"github.com/spf13/cobra"
)

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Creates a sample from parameters",
	Long:  `Creates a sample from parameters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		typ, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}
		fmt.Printf("type of %s\n", typ)
		exp, err := cmd.Flags().GetInt("exp")
		if err != nil {
			return err
		}
		fmt.Printf("expected val of %d\n", exp)
		rand, err := cmd.Flags().GetInt("rand")
		if err != nil {
			return err
		}
		fmt.Printf("rand val of %d\n", rand)
		sig := &signal.Signal{Type: typ, ExpectedValue: exp, RandomValue: rand}
		sample := signal.GenerateSignalSample(sig)
		fmt.Printf("Returned sample of: %d\n", sample)
		return nil
	},
}

func init() {
	sampleCmd.Flags().String("type", "", "type of signal to create")
	sampleCmd.MarkFlagRequired("type")

	sampleCmd.Flags().Int("exp", 0, "expected value of signal created")
	sampleCmd.MarkFlagRequired("exp")

	sampleCmd.Flags().Int("rand", 0, "noise of signal created")
	sampleCmd.MarkFlagRequired("rand")
	rootCmd.AddCommand(sampleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
