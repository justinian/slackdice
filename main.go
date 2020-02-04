package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/justinian/dice"
	"github.com/kelseyhightower/envconfig"
)

type SlackMessage struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Icon     string `json:"icon_emoji"`
}

func rollHandler(c Config, private bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		desc := r.PostFormValue("text")
		user := r.PostFormValue("user_name")
		channel := r.PostFormValue("channel_name")
		channelId := r.PostFormValue("channel_id")

		result, reason, err := dice.Roll(desc)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		if private {
			fmt.Fprintf(w, "You rolled %s: %v", reason, result)
			log.Printf("Private roll for %s (%s): %s = %v", user, reason, result.Description(), result)
		} else {
			if reason != "" {
				reason = fmt.Sprintf(" *%s*", reason)
			}

			resultStrs := strings.Split(result.String(), "\n")
			for i, s := range resultStrs {
				if i == 0 {
					resultStrs[i] = fmt.Sprintf("*%s*", s)
				} else {
					resultStrs[i] = fmt.Sprintf("_%s_", s)
				}
			}

			m := SlackMessage{
				Text:     fmt.Sprintf("*<@%s>* rolled `%s`:%s\n%s", user, result.Description(), reason, strings.Join(resultStrs, "\n")),
				Username: "rollbot",
				Channel:  channelId,
				Icon:     ":d20:",
			}

			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			err := enc.Encode(m)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err.Error())
				return
			}

			_, err = http.Post(c.SlackUrl, "text/json", &buf)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err.Error())
			}

			log.Printf("Roll for %s in %s(%s) [%s]: %s = %v", user, channel, channelId, reason, result.Description(), result)
		}
	}
}

type Config struct {
	Port     int `envconfig:"port"`
	SlackUrl string `envconfig:"slack_url"`
}

func main() {
	c := Config{Port: 8000}
	rand.Seed(time.Now().UnixNano())

	err := envconfig.Process("slackdice", &c)
	if err != nil {
		log.Fatal("Getting config: " + err.Error())
	}

	http.HandleFunc("/roll", rollHandler(c, false))
	http.HandleFunc("/roll/private", rollHandler(c, true))

	listen := fmt.Sprintf(":%d", c.Port)
	log.Printf("Starting up, listening on %s", listen)

	log.Fatal(http.ListenAndServe(listen, nil))
}
