package twitterbot

import (
	"../util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var p = fmt.Println
var prev_id []int
var image_url string

var TARGET_NAMES = [1]string{"@YoshitsuguFujii"}

func GetTweet(counter int, target_name string) (string, string) {
	var tweet_text, tweet_url string
	bearer := getToken()
	url := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=" + target_name
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

	p(prev_id[counter])
	if prev_id[counter] != tweet[0].Id {
		tweet_url = buildUrl(tweet[0].Id, target_name)
		tweet_text = tweet[0].Text
	} else {
		tweet_url = ""
		tweet_text = ""
	}
	prev_id[counter] = tweet[0].Id
	return tweet_text, tweet_url
}

func buildUrl(id int, target_name string) string {
	return "https://twitter.com/" + target_name + "/status/" + strconv.Itoa(id)
}

func PostTweet(tweet_text string, tweet_url string, target_name string) {
	user_info := getUserInfo(target_name)
	params, _ := json.Marshal(Slack{
		tweet_text,
		user_info.Name + "Bot",
		"",
		user_info.ProfileImageUrl,
		"#test"})

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
	for i := 0; i < len(TARGET_NAMES); i++ {
		prev_id = append(prev_id, 0)
	}
}

func WatchUser() {
	initialize_watch_user()

	is_first := true

	for {
		p(TARGET_NAMES)
		for i := 0; i < len(TARGET_NAMES); i++ {
			// 指定された回数分ループ
			tweet_text, tweet_url := GetTweet(i, TARGET_NAMES[i])
			if !is_first && tweet_text != "" {
				PostTweet(tweet_text, tweet_url, TARGET_NAMES[i])
			}
		}
		is_first = false
		time.Sleep(5 * 60 * time.Second)
	}
}
