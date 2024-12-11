package config

import (
	"encoding/json"
	"log"
	"os"
)

var stateFilePath = "./config/state.json"

type State struct {
	CurrentSession string
}

func SetState(state State) {
	encodedState, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error encoding state: %v", err)
		return
	}

	err = os.WriteFile(stateFilePath, encodedState, 0644)
	if err != nil {
		log.Printf("Error writing state file: %v", err)
		return
	}
}

func LoadState() State {
	stateFile, err := os.Open(stateFilePath)
	if err != nil {
		log.Printf("Error opening state file: %v", err)
		return State{}
	}

	defer stateFile.Close()

	var state State

	decoder := json.NewDecoder(stateFile)
	err = decoder.Decode(&state)
	if err != nil {
		log.Printf("Error decoding state file: %v", err)
		return State{}
	}

	return state
}
