package main

import (
	"context"
	"fmt"
)

func main() {

	input := "What is the meaning of life?"

	config := NewConfig()

	agent := NewAgent(config)

	agent.Read(context.Background(), func(ctx context.Context, chunk []byte) {
		fmt.Print(string(chunk))
	})

	agent.Run(context.Background(), input)
}
