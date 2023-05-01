package main

import (
	"fmt"
	"os"

	"github.com/stuton/xm-golang-exercise/internal/application"
	"github.com/stuton/xm-golang-exercise/utils/logger"
)

func main() {
	log, err := logger.New("xm-api")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := application.Run(log); err != nil {
		log.Errorw("unable to start application", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}
