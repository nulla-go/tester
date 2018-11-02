package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Strengine RTMP tester"
	app.Usage = "Push and recive RTMP streams"
	var fileWithAddresses string

	app.Commands = []cli.Command{
		{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Push stream to specify addrs",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "addresses, a",
					Usage:       "Addresses for push stream",
					Value:       "./addresses",
					Destination: &fileWithAddresses,
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("Push")
				return nil
			},
		},
		{
			Name:    "recieve",
			Aliases: []string{"r"},
			Usage:   "Recieve streams",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port, p",
					Usage: "Port for server",
					Value: "1935",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("Recieve")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
