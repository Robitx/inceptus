// Package env deals with application state
// bundle here stuff you need to pass throughout your app
// config, logger, app context, DB, message broker, and such..
package env

import (
	"fmt"
	"os"
	"syscall"
	"time"

	conf "github.com/robitx/inceptus/conf"
	helpers "github.com/robitx/inceptus/helpers"
	life "github.com/robitx/inceptus/life"
	log "github.com/robitx/inceptus/log"
)

// App holds state of the application..
// config, logger, context and so on
type App struct {
	// App wide context with cancel function
	// and some helper functions to register syscalls for smooth app termination
	// and reloading (conf, logger and such)
	life.AppContext

	// App specific config struct populated from conf file or envs
	Config

	// Logger wrapping zerolog
	*log.Logger
}

// New prepares application environment
func New() *App {
	// Preparing config
	var config Config
	confFile, envPrefix := conf.ReadFlagsHelper()
	// preferably don't use both at once
	// (overriting previously loaded slices, can get ugly)
	if confFile != "" {
		conf.LoadYAML(confFile, &config)
	}
	if envPrefix != "" {
		conf.LoadENV(envPrefix, &config)
	}

	// Preparing logger
	logger := log.New(config.Log.File, config.Log.Mask,
		log.TimestampHook(),
		log.StaticHook("pid", os.Getpid()),
	)

	logger.Info().
		Interface("config", config).
		Msg("showing config")

	logger.Info().
		Interface("envs", helpers.GetEnvs()).
		Msg("showing envs")

	// Preparing app wide context that supports smooth app termination
	appCtx := life.New()
	appCtx.RegisterStopSignals(config.Control.DieTimeout,
		syscall.SIGINT, syscall.SIGTERM)

	appEnv := &App{
		AppContext: appCtx,
		Logger:     logger,
		Config:     config,
	}

	// Cleanup starts after app context finishes
	go appEnv.Cleanup()

	return appEnv
}

// Cleanup after application environment
func (app *App) Cleanup() {
	<-app.Context.Done()

	time.Sleep(app.Config.Control.DieTimeout)

	app.Logger.Info().Msg("Bye")
	app.Logger.Close()

	fmt.Fprintf(os.Stderr, "Bye\n")
}
