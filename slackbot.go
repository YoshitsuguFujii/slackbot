package main

import (
	"./qiitabot"
	"./slackbot_responder"
	"./twitterbot"
	"fmt"
	"net/http"
	"time"
)

var p = fmt.Println

func slackBotResponder(w http.ResponseWriter, r *http.Request) {
	checkUser(w, r, func(text string) {
		return_text := slackbot_responder.DetectWord(text)
		fmt.Fprintf(w, "{\"text\": \"%s\"}", return_text)
	})
}

func qiitaBotResponder(w http.ResponseWriter, r *http.Request) {
	checkUser(w, r, func(text string) {
		return_text := qiitabot.UserStockSample(text)
		fmt.Fprintf(w, "{\"text\": \"%s\"}", return_text)
	})
}

func checkUser(w http.ResponseWriter, r *http.Request, proc func(text string)) {
	if r.Method == "POST" {
		text := r.FormValue("text")
		user_name := r.FormValue("user_name")

		if user_name != "slackbot" {
			p("user_name:", user_name)
			proc(text)
		}
	}
}

func PostTwitterMessage() {
	for {
		tweet_text, tweet_url := twitterbot.GetTweet()
		if tweet_text != "" {
			twitterbot.PostTweet(tweet_text, tweet_url)
		}
		time.Sleep(60 * time.Second)
	}
}

func main() {
	go PostTwitterMessage()
	http.HandleFunc("/", slackBotResponder)
	http.HandleFunc("/qiita", qiitaBotResponder)
	http.ListenAndServe(":8888", nil)
}
