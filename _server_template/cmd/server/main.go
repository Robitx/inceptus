package main

import (
	do "github.com/robitx/inceptus/_server_template/internal/do"
	env "github.com/robitx/inceptus/_server_template/internal/env"
)

func main() {
	// bundle for application state
	// (config, logger, app wide context and so on..)
	appEnv := env.New()

	// For app specific logic
	go do.It(appEnv)

	// Keep main running until panic or ending syscall
	appEnv.RunForever()
}
