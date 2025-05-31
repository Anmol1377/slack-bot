package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
	"github.com/shomali11/slacker"
)

func PrintCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {

		fmt.Println("command events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

type NewsHandler struct{}

func (h *NewsHandler) FetchNews() (string, error) {
	feedURL := "https://inc42.com/feed/"

	// Initialize gofeed parser
	fp := gofeed.NewParser()

	// Fetch and parse the RSS feed
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return "", fmt.Errorf("could not fetch or parse feed: %v", err)
	}

	// Prepare a message with the latest news
	var newsMessage string
	if len(feed.Items) == 0 {
		newsMessage = "No recent news available at the moment."
	} else {
		newsMessage = "Here are the latest headlines:\n"
		for i, item := range feed.Items {
			if i >= 5 {
				break
			}
			newsMessage += fmt.Sprintf("%d. *%s*\n%s\n\n", i+1, item.Title, item.Link)
		}
	}

	return newsMessage, nil
}

func main() {

	os.Setenv("SLACK_BOT_TOKEN", "")
	os.Setenv("SLACK_APP_TOKEN", "")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go PrintCommandEvents(bot.CommandEvents())
	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCTX slacker.BotContext, Request slacker.Request, Response slacker.ResponseWriter) {
			Response.Reply("pong")
		},
	})
	bot.Command("hello", &slacker.CommandDefinition{
		Handler: func(botCTX slacker.BotContext, Request slacker.Request, Response slacker.ResponseWriter) {
			Response.Reply("hello, how are you?")
		},
	})

	newsHandlerGo := &NewsHandler{}
	bot.Command("news from go", &slacker.CommandDefinition{
		Handler: func(botCTX slacker.BotContext, Request slacker.Request, Response slacker.ResponseWriter) {
			// Fetch the latest news
			news, err := newsHandlerGo.FetchNews()
			if err != nil {
				// If there's an error, reply with an error message
				Response.Reply(fmt.Sprintf("Oops! Something went wrong: %v", err))
				return
			}

			// Send the latest news to the user
			Response.Reply(news)
		},
	})

	newsHandler := &NewsHandler{}
	bot.Command("news from  make", &slacker.CommandDefinition{
		Handler: func(botCTX slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			// Fetch the latest news
			news, err := newsHandler.FetchNews()
			if err != nil {
				response.Reply(fmt.Sprintf("Oops! Something went wrong: %v", err))
				return
			}

			//  Trigger the Make scenario webhook
			webhookURL := "https://hook.eu2.make.com/kev3jirwkvz4attya8y36oav619du9e5"
			payload := map[string]string{
				"event": "news_requested",
				"text":  news,
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Failed to marshal JSON for webhook:", err)
			} else {
				resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error triggering Make webhook:", err)
				} else {
					defer resp.Body.Close()
					fmt.Println("Make webhook triggered. Status:", resp.Status)
				}
			}

			// Reply with the news to the Slack user
			// response.Reply(news)

		},
	})

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(context)

	if err != nil {
		log.Fatal(err, "error here")
	}
}
