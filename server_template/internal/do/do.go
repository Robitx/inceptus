package do


import (
	env "github.com/robitx/inceptus/server_template/internal/env"
)

// It starts doing the real application specific work
func It(app *env.App) {
	app.Logger.Info().Msg("starting the real work..")
}