package main

import (
	"os"

	"github.com/urfave/cli"
)

var Version string = "0.1.0"

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ec2sshtb"
	app.Version = Version
	app.Author = "reizist"
	app.Email = "reizist@gmail.com"
	app.Commands = Commands
	return app
}

func main() {
	newApp().Run(os.Args)
}
