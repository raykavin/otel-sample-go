package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Setup() {
	loggingPath := "logs"
	if lgEnvPath := os.Getenv("LOGGING_PATH"); lgEnvPath != "" {
		loggingPath = lgEnvPath
	}

	if err := os.MkdirAll(loggingPath, 0o755); err != nil {
		log.Fatal(err)
	}

	fullPath := filepath.Join(loggingPath, "app.log")

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
}
