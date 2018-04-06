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

type ServiceInfo struct {
	ServiceName string `json:"ServiceName"`
	UidOrEmail  string `json:"UidOrEmail"`
	Password    string `json:"Password"`
}

const (
	exitCodeOK = iota
	exitCodeErr
)

type ServicesInfo []ServiceInfo

func chHomeDir() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	err2 := os.Chdir(usr.HomeDir)
	if err2 != nil {
		return err2
	}
	return nil
}

func hCtrlC(ch chan os.Signal) error {
	<-ch
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil,
	}
	var ws syscall.WaitStatus
	pid, err := syscall.ForkExec("/bin/stty", []string{"stty", "echo"}, &attrs)
	if err != nil {
		return err
	}
	_, err2 := syscall.Wait4(pid, &ws, 0, nil)
	if err2 != nil {
		return err2
	}
	os.Exit(exitCodeOK)
	return nil
}

func encodingJson(serviceInfo ServiceInfo) (data []byte, err error) {
	data, err = json.Marshal(serviceInfo)
	if err != nil {
		return nil, err
	}
	return
}

func getServicesInfo(dir string) (servicesInfo []ServiceInfo, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, serviceName := range files {
		var serviceInfo ServiceInfo
		var servicesInfo ServicesInfo
		err := os.Chdir(".armadillo")
		if err != nil {
			return nil, err
		}
		file, err2 := ioutil.ReadFile(string(serviceName.Name()))
		if err2 != nil {
			return nil, err2
		}
		err3 := json.Unmarshal(file, &serviceInfo)
		if err3 != nil {
			return nil, err3
		}
		servicesInfo = append(servicesInfo, serviceInfo)
	}
	return servicesInfo, nil
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
					err := os.Mkdir(".armadillo", 0777)
					if err != nil {
						return err
					}
					fmt.Printf("Successful initialization.\n")
				} else {
					fmt.Printf("Already initialized.\n")
				}
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "armadillo create [service_name] <- setting password for service.",
			Action: func(c *cli.Context) error {
				serviceInfo := ServiceInfo{}

				for {
					fmt.Printf("Enter service name: ")
					stdIn1 := bufio.NewScanner(os.Stdin)
					stdIn1.Scan()
					serviceInfo.ServiceName = stdIn1.Text()

					if len(serviceInfo.ServiceName) != 0 {
						break
					} else {
						fmt.Printf("Input is empty! Cancel with Ctrl + C\n")
					}
				}

				for {
					fmt.Printf("Enter UserID or Email used for login: ")
					stdIn2 := bufio.NewScanner(os.Stdin)
					stdIn2.Scan()
					serviceInfo.UidOrEmail = stdIn2.Text()

					if len(serviceInfo.UidOrEmail) != 0 {
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
					servicePass, err := terminal.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return err
					}

					fmt.Printf("\nRetype password: ")
					retypePass, err2 := terminal.ReadPassword(int(syscall.Stdin))
					if err2 != nil {
						return err2
					}

					serviceInfo.Password = string(servicePass)
					retypePassStr := string(retypePass)

					if len(serviceInfo.Password) != 0 {
						if retypePassStr == serviceInfo.Password {
							chHomeDir()
							err := os.Chdir(".armadillo")
							if err != nil {
								return err
							}
							bdata, err := encodingJson(serviceInfo)
							if err != nil {
								return err
							}
							content := []byte(bdata)
							err2 := ioutil.WriteFile(serviceInfo.ServiceName+".json", content, os.ModePerm)
							if err2 != nil {
								return err2
							}
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
				serviceInfo := ServiceInfo{}
				chHomeDir()

				var items []string
				servicesInfo, err := getServicesInfo(".armadillo")
				if err != nil {
					return err
				}
				for _, serviceInfo := range servicesInfo {
					items = append(items, serviceInfo.ServiceName)
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
						return err
					}
					serviceInfo.ServiceName = result

					for {
						fmt.Printf("Enter UserID or Email used for login: ")
						stdIn2 := bufio.NewScanner(os.Stdin)
						stdIn2.Scan()
						serviceInfo.UidOrEmail = stdIn2.Text()

						if len(serviceInfo.UidOrEmail) != 0 {
							break
						} else {
							fmt.Printf("Input is empty! Cancel with Ctrl + C\n")
						}
					}

					for {
						fmt.Printf("Enter service password: ")
						servicePass, err := terminal.ReadPassword(int(syscall.Stdin))
						if err != nil {
							return err
						}

						fmt.Printf("\nRetype password: ")
						retypePass, err2 := terminal.ReadPassword(int(syscall.Stdin))
						if err2 != nil {
							return err2
						}

						serviceInfo.Password = string(servicePass)
						retypePassStr := string(retypePass)

						if len(serviceInfo.Password) != 0 {
							if retypePassStr == serviceInfo.Password {
								chHomeDir()
								err := os.Chdir(".armadillo")
								if err != nil {
									return err
								}
								bdata, err := encodingJson(serviceInfo)
								if err != nil {
									return err
								}
								content := []byte(bdata)
								err2 := ioutil.WriteFile(serviceInfo.ServiceName+".json", content, os.ModePerm)
								if err2 != nil {
									return err2
								}
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
				chHomeDir()

				var items []string
				servicesInfo, err := getServicesInfo(".armadillo")
				if err != nil {
					return err
				}
				for _, serviceInfo := range servicesInfo {
					items = append(items, serviceInfo.ServiceName)
				}

				if len(items) != 0 {
					ch := make(chan os.Signal)
					signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
					go hCtrlC(ch)
					prompt := promptui.Select{
						Label: "Delete the information. Please select a service.",
						Items: items,
					}
					_, result, err := prompt.Run()
					if err != nil {
						return err
					}

					fileName := result + ".json"

					err2 := os.Remove(fileName)
					if err2 != nil {
						return err2
					}
					fmt.Printf("Information on the service has been deleted.\n")

				} else {
					fmt.Printf("Information on the service is not registered.\n")
				}
				return nil
			},
		},
		{
			Name:  "show",
			Usage: "armadillo show <- show password.",
			Action: func(c *cli.Context) error {
				serviceInfo := ServiceInfo{}
				chHomeDir()

				var items []string
				servicesInfo, err := getServicesInfo(".armadillo")
				if err != nil {
					return err
				}
				for _, serviceInfo := range servicesInfo {
					items = append(items, serviceInfo.ServiceName)
				}

				if len(items) != 0 {
					ch := make(chan os.Signal)
					signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
					go hCtrlC(ch)
					prompt := promptui.Select{
						Label: "Show the information. Please select a service.",
						Items: items,
					}
					_, result, err := prompt.Run()
					if err != nil {
						return err
					}

					fileName := result + ".json"

					file, err := ioutil.ReadFile(fileName)
					if err != nil {
						return err
					}

					err2 := json.Unmarshal(file, &serviceInfo)
					if err2 != nil {
						return err2
					}

					fmt.Printf("Service name -> %s\nUserID or Email -> %s\nPassword -> %s\n", serviceInfo.ServiceName, serviceInfo.UidOrEmail, serviceInfo.Password)

				} else {
					fmt.Printf("Information on the service is not registered.\n")
				}
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
