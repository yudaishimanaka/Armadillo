package main

import (
	"os"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Armadillo"
	app.Usage = "Password management CLI tool"
	app.Version = "1.0.0"

	app.Run(os.Args)
}
