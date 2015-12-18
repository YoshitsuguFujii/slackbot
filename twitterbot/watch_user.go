package twitterbot

import (
	log "../lib/logger"
	"../util"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

//var p = fmt.Println
var prev_id []int
var image_url string

var user_settings WatchUserSettings

type WatchUserSettings struct {
	Groups []Channel `groups`
}

type Channel struct {
	Channel string   `channel`
	Users   []string `users`
}

var users []User

type User struct {
	Channel    string
	Identifier string `users`
	PrevID     int
}

const Watchuser = "twitterbot/watch_user.yml"

func GetTweet(target *User) (string, string) {
	var tweet_text, tweet_url string
	bearer := getToken()
	url := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=" + target.Identifier
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	util.Perror(err)

	var tweet []Tweet
	parse_err := json.Unmarshal(body, &tweet)
	util.Perror(parse_err)

	if target.PrevID != tweet[0].Id {
		tweet_url = buildUrl(tweet[0].Id, target.Identifier)
		tweet_text = tweet[0].Text
	} else {
		tweet_url = ""
		tweet_text = ""
	}
	target.PrevID = tweet[0].Id
	return tweet_text, tweet_url
}

func buildUrl(id int, target_name string) string {
	return "https://twitter.com/" + target_name + "/status/" + strconv.Itoa(id)
}

func PostTweet(tweet_text string, tweet_url string, target *User) {
	user_info := getUserInfo(target.Identifier)
	params, _ := json.Marshal(Slack{
		tweet_text,
		user_info.Name + "Bot",
		"",
		user_info.ProfileImageUrl,
		target.Channel})

	resp, _ := http.PostForm(
		slackUrl(),
		url.Values{"payload": {string(params)}},
	)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	println(string(body))
}

func slackUrl() string {
	return os.Getenv("slack_incoming_url")
}

func initialize_watch_user() {
	buf, err := ioutil.ReadFile(Watchuser)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &user_settings)

	for i := 0; i < len(user_settings.Groups); i++ {
		for j := 0; j < len(user_settings.Groups[i].Users); j++ {
			users = append(users,
				User{Channel: user_settings.Groups[i].Channel,
					Identifier: user_settings.Groups[i].Users[j],
					PrevID:     0})
		}
	}
}

func WatchUser() {
	initialize_watch_user()

	is_first := true

	for {
		log.Info(users)
		for i := 0; i < len(users); i++ {
			// 指定された回数分ループ
			tweet_text, tweet_url := GetTweet(&users[i])
			if !is_first && tweet_text != "" {
				PostTweet(tweet_text, tweet_url, &users[i])
			}
		}
		is_first = false
		time.Sleep(5 * 60 * time.Second)
	}
}
