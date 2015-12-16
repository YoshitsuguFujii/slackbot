package main

import (
	log "./lib/logger"
	"./qiitabot"
	"./slackbot_responder"
	"./todo_list"
	"./twitterbot"
	"./util"
	"fmt"
	"github.com/fukata/golang-stats-api-handler"
	"net/http"
	"os"
	"syscall"
)

//var p = fmt.Println

const PidFilePath = "tmp.pid"

func slackBotResponder(w http.ResponseWriter, r *http.Request) {
	checkUser(w, r, func(text string, channel_name string) {
		return_text := slackbot_responder.DetectWord(text)
		fmt.Fprintf(w, "{\"text\": \"%s\"}", return_text)
	})
}

func qiitaBotResponder(w http.ResponseWriter, r *http.Request) {
	checkUser(w, r, func(text string, channel_name string) {
		return_text := qiitabot.UserStockSample(text)
		fmt.Fprintf(w, "{\"text\": \"%s\"}", return_text)
	})
}

func todoListBot(w http.ResponseWriter, r *http.Request) {
	checkUser(w, r, func(text string, channel_name string) {
		return_text := todo.Accept(text, channel_name)
		fmt.Fprintf(w, "{\"text\": \"%s\"}", return_text)
	})
}

func checkUser(w http.ResponseWriter, r *http.Request, proc func(text string, channel_name string)) {
	if r.Method == "POST" {
		text := r.FormValue("text")
		user_name := r.FormValue("user_name")
		channel_name := r.FormValue("channel_name")

		if user_name != "slackbot" {
			log.Info("user_name: " + user_name)
			log.Info("channel_name: " + channel_name)
			log.Info("text: " + text)
			proc(text, channel_name)
		}
	}
}

func postTwitterMessage() {
	twitterbot.WatchUser()
}

func watchWord() {
	twitterbot.WatchWord()
}

func prepare() {
	if ferr := os.Remove(PidFilePath); ferr != nil {
		if !os.IsNotExist(ferr) {
			panic(ferr.Error())
		}
	}
	pidf, perr := os.OpenFile(PidFilePath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)

	if perr != nil {
		panic(perr.Error())
	}
	if _, err := fmt.Fprint(pidf, syscall.Getpid()); err != nil {
		panic(err.Error())
	}
	pidf.Close()

	log.InitLog()
}

func main() {
	prepare()
	log.Info("START => " + util.JpCurrentTIme())
	go postTwitterMessage()
	go watchWord()
	http.HandleFunc("/", slackBotResponder)
	http.HandleFunc("/qiita", qiitaBotResponder)
	http.HandleFunc("/todo", todoListBot)
	http.HandleFunc("/stats", stats_api.Handler)
	http.ListenAndServe(":8888", nil)
}
