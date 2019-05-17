package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nulla-go/tester/reciever"
	"github.com/nulla-go/tester/requester"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Nulla RTMP tester"
	app.Usage = "Push and recive RTMP streams"
	var fileWithAddresses string
	var serverPort string

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
				file, err := os.Open(fileWithAddresses)
				if err != nil {
					fmt.Println("File didn't find")
					return nil
				}
				body := make([]byte, 256)
				l, err := file.Read(body)
				if err != nil {
					fmt.Println("Failed read file, ", fileWithAddresses, err)
					return err
				}
				addresses := strings.Split(string(body[:l]), "\n")
				//strings := mainStrings.Split()
				//fmt.Println(addresses)
				requesterServer := &requester.Requester{}
				requesterServer.Push(addresses)

				//requestServer := &requester.Requester{}
				//requestServer.Start()
				return nil
			},
		},
		{
			Name:    "recieve",
			Aliases: []string{"r"},
			Usage:   "Recieve streams",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "port, p",
					Usage:       "Port for server",
					Value:       ":1935",
					Destination: &serverPort,
				},
			},
			Action: func(c *cli.Context) error {
				reciverServer := reciever.NewReciever()
				reciverServer.Start(serverPort)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
