package main

func main() {
	var app App
	app.CreateConnection()
	app.Migrate()
	app.CreateRoutes()
	app.Run()
}
