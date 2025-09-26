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
		model := NewTestModel()
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
