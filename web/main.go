package main

import (
	"fmt"
	"lottery/bootstrap"
	"lottery/web/routes"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("lottery", "YuBo")

	app.Bootstrap()

	app.Configure(routes.Configure)

	return app
}

func main() {
	app := newApp()

	app.Listen(fmt.Sprintf(":%d", port))
}
