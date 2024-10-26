package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/struki84/clipt/config"
	"github.com/struki84/clipt/internal"
)

var cliptCmd = &cobra.Command{
	Use:   "clipt",
	Short: "CLI RAG tool for interacting with LLMs.",
	Long:  "",
}

func init() {
	appConfig := config.NewConfig()

	promptCmd := &cobra.Command{
		Use:   "prompt",
		Short: "Run a prompt.",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			agent := internal.NewAgent(appConfig)

			agent.Stream(context.Background(), func(ctx context.Context, chunk []byte) {
				fmt.Print(string(chunk))
			})

			input := args[0]

			agent.Run(context.Background(), input)
		},
	}

	cliptCmd.AddCommand(promptCmd)

	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Run chat UI",
		Run: func(cmd *cobra.Command, args []string) {
			agent := internal.NewAgent(appConfig)
			ShowUI(agent)
		},
	}

	cliptCmd.AddCommand(chatCmd)

	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Run node",
		Run: func(cmd *cobra.Command, args []string) {
			RunNode()
		},
	}

	cliptCmd.AddCommand(nodeCmd)
}
func main() {
	if err := cliptCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
