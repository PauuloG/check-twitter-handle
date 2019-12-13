package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli"
)

type Environment struct {
	TwitterHandles	[]string `required:"true" envconfig:"TWITTER_HANDLES"`
}

type Notifier interface {
	Notify(twitterHandle string, statusCode int) bool
}

func main() {
	var notifierParameter string

	err := godotenv.Load("/go/bin/.env")
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	var e Environment
	err = envconfig.Process("checktwitterhandle", &e)
	if err != nil {
		log.Fatalf("envconfig.Process: %w", err)
	}

	app := cli.NewApp()
	app.Name = "check-twitter-handle"
	app.Usage = "Check if a specific twitter handle is available"

	app.Flags = []cli.Flag {
		&cli.StringFlag{
		  Name:        "notifier, n",
		  Value:       "slack",
		  Usage:       "Chose a notifier. Supported values : slack",
		  Destination: &notifierParameter,
		},
	  }

	app.Commands = []*cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "check if twitter handle is available",
			Action: func(c *cli.Context) error {
				var statusCode int

				notifier := getNotifier(notifierParameter)

				for _, twitterHandle := range e.TwitterHandles {
					statusCode = getTwitterHandleStatusCode(twitterHandle)
					if statusCode != 200 {
						log.Println(fmt.Sprintf("Twitter Handle %s looks available, notifying", twitterHandle))

						notifier.Notify(twitterHandle, statusCode)
					} else {
						log.Println(fmt.Sprintf("Twitter Handle %s is not available, nothing to do.", twitterHandle))
					}
				}

				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Test notifier settings",
			Action: func(c *cli.Context) error {
				log.Println("Sending test message")

				notifier := getNotifier(notifierParameter)
				notifier.Notify("test", 404)
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getTwitterHandleStatusCode(twitterHandle string) int {
	twitterUrl := fmt.Sprintf("https://twitter.com/%s", twitterHandle)

	log.Println(fmt.Sprintf("Visiting %v", twitterUrl))
	resp, err := http.Get(twitterUrl)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Response for url %s : %d", twitterUrl, resp.StatusCode))
	return resp.StatusCode
}

func getNotifier(notifierParameter string) (Notifier) {

  switch notifierParameter {
  case "slack":
    return makeSlackNotifier()
  }

  log.Println("No notifier found, or notifier is not valid. Falling back to slack notifier")
  return makeSlackNotifier()
}
