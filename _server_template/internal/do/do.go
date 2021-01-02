package do

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	env "github.com/robitx/inceptus/_server_template/internal/env"
	rest "github.com/robitx/inceptus/_server_template/internal/rest"

	route "github.com/robitx/inceptus/route"
	middleware "github.com/robitx/inceptus/route/middleware"
)

// It starts doing the real application specific work
func It(app *env.App) {
	app.Logger.Info().Msg("starting the real work..")

	// main router with some middleware
	router := route.New()
	if app.Auth.FirebaseAccountKeyFile != "" {
		router.Use(middleware.Auth("api", "accessToken",
			app.Auth.FirebaseAccountKeyFile))
	}
	router.Use(middleware.Base(app.Logger, "x-request-id"))

	// FileSserver (dirs are unaccessible) with custom html error pages
	workDir, _ := os.Getwd()
	r404, _ := ioutil.ReadFile(
		filepath.Join(workDir, "static/errors/404.html"))

	contentType := "text/html; charset=utf-8"
	customErrors := make(map[int][]byte)
	customErrors[http.StatusNotFound] = r404
	htmlErrors := middleware.Error(contentType, customErrors)
	dir := filepath.Join(workDir, "static")
	router.Mount("/static", htmlErrors(route.FilesServer("/static", dir)))

	// Dummy rest api
	apiv1 := route.New()
	apiv1.Get("/", rest.Echo)
	apiv1.Get("/ping", rest.Ping)
	apiv1.Get("/panic", rest.Panic)
	apiv1.Get("/echo", rest.Echo)
	apiv1.Get("/hi-user", rest.HiUser)
	router.Mount("/api/v1", apiv1)

	// Write routes schema to the beginning of the log
	if app.Router.GenerateDoc {
		go func() {
			app.Logger.Info().
				RawJSON("doc", []byte(route.GenerateDocs(router))).
				Msg("")
		}()
	}

	go func() {
		_, err := app.ExecContext(context.Background(), "insert into dummy(name) values($1)", "app item")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
			os.Exit(1)
		}

		rows, _ := app.QueryContext(context.Background(), "select * from dummy")

		for rows.Next() {
			var name string
			err := rows.Scan(&name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			}
			fmt.Printf("%s\n", name)
		}
	}()

	server := &http.Server{
		Addr:    app.Router.Addr,
		Handler: router,
		// Handler: http.HandlerFunc(rest.Ping),
	}
	err := server.ListenAndServe()
	app.Logger.Error().
		Err(err).
		Msg("server failed")
}
