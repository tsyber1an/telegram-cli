package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

func NewBot(apiToken string) (*tb.Bot, error) {
	if apiToken == "" {
		return nil, errors.New("token can not be empty")
	}

	var bot *tb.Bot
	tbSettings := tb.Settings{
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	tbSettings.Token = apiToken
	var err error
	bot, err = tb.NewBot(tbSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to init a bot, got the error: %s", err)
	}

	return bot, nil
}

func main() {

	app := &cli.App{
		Name: "Telegram Bot REPL",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api_token",
				Value: "",
				Usage: "api_token",
			},
		},
		Usage: "Telegram Bot REPL",
		Commands: []*cli.Command{
			{
				Name:    "tell",
				Aliases: []string{"install"},
				Usage:   "handle a message and reply",
				Subcommands: []*cli.Command{
					{
						Name:  "sendMessage",
						Usage: "tell which command and which reply text should",
						Action: func(c *cli.Context) error {
							bot, err := NewBot(c.String("api_token"))
							if err != nil {
								return err
							}

							bot.Handle(os.Args[5], func(m *tb.Message) {
								log.Printf("command %s received. Repling with %s", os.Args[5], os.Args[6])
								bot.Send(m.Sender, os.Args[6])
								log.Println("a message was send. Stopping a bot")
								bot.Stop()
								return
							})

							log.Println("Waiting a text messge")
							bot.Start()

							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					_, err := NewBot(c.String("api_token"))
					if err != nil {
						return err
					}

					log.Println("The Bot is on. Ready to receive commands")

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
