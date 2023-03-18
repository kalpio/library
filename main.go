package main

import "library/application"

func main() {
	app := &application.App{}

	app.Host("127.0.0.1")
	app.Port("8089")
	app.Initialize()
	app.Run()
}
