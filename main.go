package main

import (
	"fmt"
	"os"
	"os/user"
	"bufio"
	"sort"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func main() {
	app := cli.NewApp()

	app.Name = "Armadillo"
	app.Usage = "Password management CLI tool"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "armadillo init <- Initialization processing, done only once.",
			Action: func(c *cli.Context) error {
				usr, err := user.Current()
				if err != nil {
					fmt.Println(err)
				}
				os.Chdir(usr.HomeDir)
				if _, err := os.Stat(".armadillo"); os.IsNotExist(err) {
					os.Mkdir(".armadillo", 0777)
					fmt.Printf("Successful initialization.\n")
				} else {
					fmt.Printf("Already initialized.\n")
				}
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "armadillo create [site_name] <- setting password for site.",
			Action: func(c *cli.Context) error {
				fmt.Printf("Enter site name.: ")
				stdIn1 := bufio.NewScanner(os.Stdin)
				stdIn1.Scan()
				siteName := stdIn1.Text()

				fmt.Printf("Enter UserID or Email.: ")
				stdIn2 := bufio.NewScanner(os.Stdin)
				stdIn2.Scan()
				idOrEmail := stdIn2.Text()

				fmt.Printf("Enter site password.: ")
				sitePass, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("\n" + siteName, idOrEmail, string(sitePass))

				return nil
			},
		},
		{
			Name:  "update",
			Usage: "armadillo update <- update password.",
			Action: func(c *cli.Context) error {
				fmt.Printf("Update password.")
				return nil
			},
		},
		{
			Name:  "show",
			Usage: "armadillo show <- show password.",
			Action: func(c *cli.Context) error {
				fmt.Printf("Show password.")
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
