package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/struki84/clipt/config"
	"github.com/struki84/clipt/files"
	"github.com/struki84/clipt/internal"
	"github.com/struki84/clipt/internal/callbacks"
	"github.com/struki84/clipt/internal/graphs"
	"github.com/struki84/clipt/internal/tools/library"
	"github.com/struki84/clipt/network"
	"github.com/struki84/clipt/ui"
)

var cliptCmd = &cobra.Command{
	Use:   "clipt",
	Short: "CLI RAG tool for interacting with LLMs.",
	Long:  "",
}

var appConfig config.AppConfig

func init() {
	appConfig = config.NewConfig()

	promptCmd := &cobra.Command{
		Use:   "prompt",
		Short: "Run a prompt.",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			agent := internal.NewAgent(appConfig)

			ctx := context.Background()

			log.Printf("Running prompt")

			agent.Stream(ctx, func(ctx context.Context, chunk []byte) {
				fmt.Print(string(chunk))
			})

			input := args[0]

			err := agent.Run(ctx, input)

			if err != nil {
				log.Println("Error running prompt:", err)
			}
		},
	}

	cliptCmd.AddCommand(promptCmd)

	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Run chat UI",
		Run: func(cmd *cobra.Command, args []string) {
			agent := internal.NewAgent(appConfig)
			ui.ShowUI(agent)
		},
	}

	cliptCmd.AddCommand(chatCmd)

	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Run node",
		Run: func(cmd *cobra.Command, args []string) {
			if err := network.RunNode(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running node: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cliptCmd.AddCommand(nodeCmd)

	graphCmd := &cobra.Command{
		Use:   "graph",
		Short: "Run graph",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			callback := callbacks.NewReActCallbackHandler()
			graphs.ReactGraph(ctx, args[0], callback)
		},
	}

	cliptCmd.AddCommand(graphCmd)

	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Run test",
		Run: func(cmd *cobra.Command, args []string) {

			// _ = internal.NewAgent(appConfig)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client, err := library.NewChromaClient()
			if err != nil {
				log.Println("Error creating chroma client:", err)
				return
			}

			sentry := files.NewFileSentry("./files", client)

			err = sentry.ScanFiles(ctx)
			if err != nil {
				log.Println("Error scanning for files:", err)
				return
			}

			go func() {
				err = sentry.WatchFiles(ctx)
				if err != nil {
					log.Println("Error watching files:", err)
				}
			}()
		},
	}

	cliptCmd.AddCommand(testCmd)

}

func main() {
	appConfig.InitDB()

	if err := cliptCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
