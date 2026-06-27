package cmd

import (
	"fmt"
	"strconv"

	"github.com/GreenStarMatter/zenzore/internal/navigator"
	"github.com/spf13/cobra"
)

var navServerURL string

var navCmd = &cobra.Command{
	Use:   "nav",
	Short: "Navigate a running zenzore server's zyztems, devices, and sensors",
	RunE: func(cmd *cobra.Command, args []string) error {
		if navServerURL == "" {
			return fmt.Errorf("--server is required, e.g. --server http://localhost:8080")
		}
		baseURL := navServerURL

		nav, _, err := navigator.LoadFromServer(baseURL)
		if err != nil {
			return err
		}

		for {
			fmt.Printf("[%s] %s\n", nav.CurrentNode.Level, nav.CurrentNode.ID)
			fmt.Println(nav.CurrentNode.List())
			fmt.Print("> (number to select, 'u' up, 'r' refresh, 'q' quit) ")

			var input string
			fmt.Scanln(&input)

			switch input {
			case "q":
				return nil
			case "u":
				if err := nav.Up(); err != nil {
					fmt.Println(err)
				}
				continue
			case "r":
				if err := nav.CurrentNode.Populate(baseURL); err != nil {
					fmt.Println("refresh failed:", err)
				}
				continue
			}

			choice, err := strconv.Atoi(input)
			if err != nil || choice < 0 || choice >= len(nav.CurrentNode.Children) {
				fmt.Println("invalid selection")
				continue
			}

			child := nav.CurrentNode.Children[choice]
			if len(child.Children) == 0 {
				if err := child.Populate(baseURL); err != nil {
					fmt.Println("failed to load children:", err)
					continue
				}
			}
			if err := nav.Down(child); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	navCmd.Flags().StringVar(&navServerURL, "server", "", "base URL of the running zenzore server, e.g. http://localhost:8080")
	rootCmd.AddCommand(navCmd)
}
