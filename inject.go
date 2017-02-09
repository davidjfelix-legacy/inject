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
	"strconv"
)

var Version string

//FIXME: remove this
func init() {
	Version = "test"
}

func promptForString(prompt string) string {
	var response string
	for {
		fmt.Print(prompt)
		if i, err := fmt.Scanln(&response); err == nil && i == 1 {
			return response
		}
		fmt.Println("Problem getting response. Ctrl+C to quit.")
	}
}

func promptForPassword(prompt string) string {
	for {
		fmt.Print(prompt)
		if rawPassword, err := terminal.ReadPassword(int(syscall.Stdin)); err == nil {
			fmt.Println("") // Need to add newline
			return string(rawPassword)
		}
		fmt.Println("\nSomething went wrong. Ctrl+C or try again...")
	}
}

func promptAffirm(question string) bool {
	affirmative := []string{"y", "yes"}
	negative := []string{"n", "no"}

	for {
		response := promptForString(question + " [y/N]: ")
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
	}
}

func promptProxyProto() string {
	for {
		response := promptForString("Proxy Protocol? [HTTP/HTTPS]: ")
		ans := strings.ToLower(response)
		if ans == "http" {
			return "http://"
		} else if ans == "https" {
			return "https://"
		}
		fmt.Println("Protocol not recognized. Ctrl+C to quit or choose \"https\" or \"http\"...")
	}
}

func promptForPort(prompt string) string {
	for {
		response := promptForString(prompt)
		if i, err := strconv.Atoi(response); err == nil && i > 0 && i < 65536 {
			return response
		}
		fmt.Println("Response was not an integer between 0-65535. Ctrl+C to quit or try again...")
	}
}

func proxyPrompt() string {
	if promptAffirm("Do you have an HTTP proxy?") {
		rawURL := promptProxyProto()
		diagRawURL := rawURL
		if promptAffirm("Does your proxy require username/pass?") {
			var username string
			for {
				username = promptForString("Username: ")
				if username != "" {
					break
				}
				fmt.Println("You can't have a blank username. Ctrl+C to quit or...")
			}
			password := promptForPassword("Password (hidden, just press enter if no password): ")
			if password != "" {
				diagRawURL = rawURL + username + ":" + "********" + "@"
				rawURL = rawURL + username + ":" + password + "@"
			} else {
				rawURL = rawURL + username + "@"
				diagRawURL = rawURL
			}
		}
		port := promptForPort("Proxy port: ")
		host := promptForString("Proxy hostname: ")
		rawURL = rawURL + host + ":" + port
		diagRawURL = diagRawURL + host + ":" + port
		fmt.Printf("Using proxy: %s\n", diagRawURL)
		return rawURL
	}
	return ""
}

func run(ctx *cli.Context) {
	fmt.Println("Injecting...")
	rawURL := proxyPrompt()
	var transport *http.Transport
	if proxyUrl, err := url.Parse(rawURL); err == nil {
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	} else {
		transport = &http.Transport{}
	}
	getDependencies(transport)
}

func
main() {
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
