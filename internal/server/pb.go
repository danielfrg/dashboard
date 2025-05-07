package server

import (
	"os"
	"log"

	"dashboard/internal/server/views"

	"github.com/a-h/templ"
	"github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
)



// runServer runs a new HTTP server with the loaded environment variables.
func RunServer() error {
	// port, err := strconv.Atoi(os.Getenv("BACKEND_PORT", "3000"))
	// if err != nil {
	// 	return err
	// }

	// Create a new PocketBase server.
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Route for the index page using templ
		se.Router.GET("/", func(e *core.RequestEvent) error {
			return templ.Handler(views.IndexPage()).ServeHTTP(e	.Response(), e.Request())
		})

		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})

    if err := app.Start(); err != nil {
        log.Fatal(err)
		return err
    }

	return nil
}
