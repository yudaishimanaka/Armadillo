package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
	"sort"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type SiteInfo struct {
	SiteName   string `json:"SiteName"`
	UidOrEmail string `json:"UidOrEmail"`
	Password   string `json:"Password"`
}

type SitesInfo []SiteInfo

func chHomeDir() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	os.Chdir(usr.HomeDir)
}

func hCtrlC(ch chan os.Signal) {
	<-ch
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil,
	}
	var ws syscall.WaitStatus
	pid, _ := syscall.ForkExec("/bin/stty", []string{"stty", "echo"}, &attrs)
	syscall.Wait4(pid, &ws, 0, nil)
	os.Exit(0)
}

func encodingJson(siteInfo SiteInfo) []byte {
	data, _ := json.Marshal(siteInfo)
	return data
}

func getSiteInfo(dir string) []SiteInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	var siteInfo SiteInfo
	var sitesInfo SitesInfo
	for _, siteName := range files {
		os.Chdir(".armadillo")
		file, err := ioutil.ReadFile(string(siteName.Name()))
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(file, &siteInfo)
		sitesInfo = append(sitesInfo, siteInfo)
	}
	return sitesInfo
}

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
				chHomeDir()

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
			Usage: "armadillo create [site_name] <- setting password for service.",
			Action: func(c *cli.Context) error {
				siteInfo := SiteInfo{}

				for {
					fmt.Printf("Enter service name: ")
					stdIn1 := bufio.NewScanner(os.Stdin)
					stdIn1.Scan()
					siteInfo.SiteName = stdIn1.Text()

					if len(siteInfo.SiteName) != 0 {
						break
					} else {
						fmt.Printf("Input is empty! Cancel with Ctrl + C\n")
					}
				}

				for {
					fmt.Printf("Enter UserID or Email used for login: ")
					stdIn2 := bufio.NewScanner(os.Stdin)
					stdIn2.Scan()
					siteInfo.UidOrEmail = stdIn2.Text()

					if len(siteInfo.UidOrEmail) != 0 {
						break
					} else {
						fmt.Printf("Input is empty! Cancel with Ctrl + C\n")
					}
				}

				ch := make(chan os.Signal)
				signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
				go hCtrlC(ch)
				for {
					fmt.Printf("Enter service password: ")
					sitePass, _ := terminal.ReadPassword(int(syscall.Stdin))

					fmt.Printf("\nRetype password: ")
					retypePass, _ := terminal.ReadPassword(int(syscall.Stdin))

					siteInfo.Password = string(sitePass)
					retypePassStr := string(retypePass)

					if len(siteInfo.Password) != 0 {
						if retypePassStr == siteInfo.Password {
							chHomeDir()
							os.Chdir(".armadillo")
							bdata := encodingJson(siteInfo)
							content := []byte(bdata)
							ioutil.WriteFile(siteInfo.SiteName+".json", content, os.ModePerm)
							fmt.Printf("\nCreate succeeded!!!\n")
							break
						} else {
							fmt.Printf("\nPasswords do not match\n")
						}
					} else {
						fmt.Printf("\nInput is empty! Cancel with Ctrl + C\n")
					}
				}
				return nil
			},
		},
		{
			Name:  "update",
			Usage: "armadillo update <- update password.",
			Action: func(c *cli.Context) error {
				siteInfo := SiteInfo{}
				chHomeDir()

				var items []string
				for _, siteInfo := range getSiteInfo(".armadillo") {
					items = append(items, siteInfo.SiteName)
				}

				if len(items) != 0 {
					ch := make(chan os.Signal)
					signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
					go hCtrlC(ch)
					prompt := promptui.Select{
						Label: "Update the information. Please select a service.",
						Items: items,
					}
					_, result, err := prompt.Run()
					if err != nil {
						fmt.Println(err)
					}
					siteInfo.SiteName = result

					for {
						fmt.Printf("Enter UserID or Email used for login: ")
						stdIn2 := bufio.NewScanner(os.Stdin)
						stdIn2.Scan()
						siteInfo.UidOrEmail = stdIn2.Text()

						if len(siteInfo.UidOrEmail) != 0 {
							break
						} else {
							fmt.Printf("Input is empty! Cancel with Ctrl + C\n")
						}
					}

					for {
						fmt.Printf("Enter service password: ")
						sitePass, _ := terminal.ReadPassword(int(syscall.Stdin))

						fmt.Printf("\nRetype password: ")
						retypePass, _ := terminal.ReadPassword(int(syscall.Stdin))

						siteInfo.Password = string(sitePass)
						retypePassStr := string(retypePass)

						if len(siteInfo.Password) != 0 {
							if retypePassStr == siteInfo.Password {
								chHomeDir()
								os.Chdir(".armadillo")
								bdata := encodingJson(siteInfo)
								content := []byte(bdata)
								ioutil.WriteFile(siteInfo.SiteName+".json", content, os.ModePerm)
								fmt.Printf("\nUpdate succeeded!!!\n")
								break
							} else {
								fmt.Printf("\nPasswords do not match\n")
							}
						} else {
							fmt.Printf("\nInput is empty! Cancel with Ctrl + C\n")
						}
					}
				} else {
					fmt.Printf("Information on the service is not registered.\n")
				}

				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "armadillo delete <- Delete service information.",
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
