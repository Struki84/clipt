package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/struki84/clipt/config"
	"github.com/struki84/clipt/internal"
	"github.com/struki84/clipt/ui"
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
			if err := RunNode(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running node: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cliptCmd.AddCommand(nodeCmd)

	uiCmd := &cobra.Command{
		Use:   "ui",
		Short: "Run UI",
		Run: func(cmd *cobra.Command, args []string) {
			ui.ShowUI(internal.NewAgent(appConfig))
		},
	}

	cliptCmd.AddCommand(uiCmd)

	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Run test",
		Run: func(cmd *cobra.Command, args []string) {
			ShowTestUI(internal.NewAgent(appConfig))
		},
	}

	cliptCmd.AddCommand(testCmd)
}

func main() {
	if err := cliptCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
