package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cliptCmd = &cobra.Command{
	Use:   "clipt",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		providers := []ChatProvider{}

		model := NewTestModel()

		providers = append(providers, model)

		ShowChatViewLight(model)
	},
}

func init() {

}

func main() {
	if err := cliptCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
