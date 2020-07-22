package main

import "gopattern/app/controllers"

func main() {
	app := controllers.App{}

	app.RunServer()
}
