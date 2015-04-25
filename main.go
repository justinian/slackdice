package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		channel := r.PostFormValue("channel")

		result, err := dice.Roll(desc)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		if private {
			fmt.Fprintf(w, "You rolled: %v", result)
		} else {
			m := SlackMessage{
				Text:     fmt.Sprintf("*%s* rolled `%s` and got:\n_%v_", user, desc, result),
				Username: "rollbot",
				Channel:  channel,
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
		}
	}
}

type Config struct {
	Listen   string
	SlackUrl string `envconfig:"slack_url"`
}

func main() {
	c := Config{Listen: ":8000"}

	err := envconfig.Process("slackdice", &c)
	if err != nil {
		log.Fatal("Getting config: " + err.Error())
	}

	http.HandleFunc("/roll", rollHandler(c, false))
	http.HandleFunc("/roll/private", rollHandler(c, true))

	log.Fatal(http.ListenAndServe(c.Listen, nil))
}
