package library

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/tools"
)

var _ tools.Tool = &FileListTool{}

type FileListTool struct {
}

func NewFileListTool() *FileListTool {
	return &FileListTool{}
}

func (search *FileListTool) Name() string {
	return "ListFiles"
}

func (search *FileListTool) Description() string {
	return "Lists all files you can read and search."
}

func (search *FileListTool) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Listing files...")

	dirPath := "./files"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "No files found", nil
	}

	var str string
	for _, file := range files {
		str += "- " + file.Name() + "\n"
	}

	return str, nil
}
