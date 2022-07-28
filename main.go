package main

func main() {
	app := &App{}

	app.Host("127.0.0.1")
	app.Port("8089")
	app.Initialize("test.db")
	app.Run()
}
