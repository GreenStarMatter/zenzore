/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/GreenStarMatter/zenzore/cmd"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) //TODO: handle the logs with an env variable (iterate between where ever zenzore is installed or just the std.out)
	slog.SetDefault(logger)
	cmd.Execute()
}
