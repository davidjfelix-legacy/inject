package main

import (
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"syscall"
)

var Version string

//FIXME: remove this
func init() {
	Version = "test"
}

func promptAffirm(prompt string) bool {
	affirmative := []string{"y", "yes"}
	negative := []string{"n", "no"}

	var response string
	fmt.Print(prompt + " [y/N]: ")
	for {
		if len, err := fmt.Scanln(&response); err == nil && len == 1 {
			ans := strings.ToLower(response)
			for _, k := range affirmative {
				if ans == k {
					return true
				}
			}
			for _, k := range negative {
				if ans == k {
					return false
				}
			}
			fmt.Println("Response not recognized. Ctrl+C to quit or y/N...")
		} else {
			fmt.Println("Problem getting response. Ctrl+C to quit.")
		}
		fmt.Print(prompt + " [y/N]: ")
	}
}

func proxyPrompt() (*http.Transport, error) {
	if promptAffirm("Do you have an HTTP proxy?") {
		rawUrl := "http://"
		// FIXME: get protocol
		if promptAffirm("Does your proxy require username/pass?") {
			var username string
			fmt.Print("Username: ")
			fmt.Scanln(&username)
			// FIXME: check for bad scans and empty responses
			fmt.Print("Password (hidden, just press enter if no password): ")
			rawPassword, _ := terminal.ReadPassword(int(syscall.Stdin))
			// FIXME: check for bad scans
			password := string(rawPassword)
			if password != "" {
				rawUrl = rawUrl + username + ":" + password + "@"
			} else {
				rawUrl = rawUrl + username + "@"
			}
		}
		// FIXME: get proxy host and port
		rawUrl += "example.com:3128"
		if proxyUrl, err := url.Parse(""); err == nil {
			return &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}, nil
		} else {
			return nil, err
		}
	} else {
		return &http.Transport{}, nil
	}
}

func run(ctx *cli.Context) {
	fmt.Println("Injecting...")
	proxyPrompt()
	getDependencies()
}

func main() {
	inject := cli.NewApp()
	inject.EnableBashCompletion = true
	inject.Name = "inject"
	inject.Usage = "A standalone binary for rapidly injecting genes into a developer workstation"
	inject.Version = Version
	inject.Author = "David J Felix <davidjfelix@davidjfelix.com>"
	inject.Action = run
	inject.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "don't print any messages",
		},
	}
	inject.Run(os.Args)
}
