package main

import "github.com/artumont/DotSlashStream/backend/cmd/bootstrap"

func main() {
	app := bootstrap.NewApplication()
	defer app.Shutdown()
	app.Start()
}
