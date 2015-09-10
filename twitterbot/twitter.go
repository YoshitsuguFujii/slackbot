// 参考
// http://venkat.io/posts/twitter-api-auth-golang/

package twitterbot

import (
	"../util"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	_ "strings"
	"time"
)

var p = fmt.Println
var prev_id []int
var image_url string

var TARGET_NAMES = [2]string{"@YoshitsuguFujii", "@sasata299"}

type BearerToken struct {
	AccessToken string `json:"access_token"`
}

type UserInfo struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	CreatedAt       string `json:"created_at"`
	ProfileImageUrl string `json:"profile_image_url""`
}

type Tweet struct {
	Id        int    `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type Slack struct {
	Text       string `json:"text"`       //投稿内容
	Username   string `json:"username"`   //投稿者名 or Bot名（存在しなくてOK）
	Icon_emoji string `json:"icon_emoji"` //アイコン絵文字
	Icon_url   string `json:"icon_url"`   //アイコンURL（icon_emojiが存在する場合は、適応されない）
	Channel    string `json:"channel"`    //#部屋名
}

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

func getBearer() string {
	key := url.QueryEscape(os.Getenv("twitter_consumer_key"))
	secret := url.QueryEscape(os.Getenv("twitter_consumer_secret"))
	encodedKeySecret := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", key, secret)))
	return encodedKeySecret
}

func getToken() string {
	url := "https://api.twitter.com/oauth2/token"
	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	request, _ := http.NewRequest("POST", url, reqBody)
	request.Header.Set("Authorization", "Basic "+getBearer())
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	util.Perror(err)

	var bearer_token BearerToken
	json.Unmarshal(body, &bearer_token)
	return bearer_token.AccessToken
}

func buildUrl(id int, target_name string) string {
	return "https://twitter.com/" + target_name + "/status/" + strconv.Itoa(id)
}

func getUserInfo(target_name string) UserInfo {
	bearer := getToken()
	url := "https://api.twitter.com/1.1/users/show.json?screen_name=" + target_name
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

	var user_info UserInfo
	parse_err := json.Unmarshal(body, &user_info)
	util.Perror(parse_err)

	return user_info
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

func initialize() {
	for i := 0; i < len(TARGET_NAMES); i++ {
		prev_id = append(prev_id, 0)
	}
}

func Watch() {
	initialize()

	for {
		p(TARGET_NAMES)
		for i := 0; i < len(TARGET_NAMES); i++ {
			// 指定された回数分ループ
			tweet_text, tweet_url := GetTweet(i, TARGET_NAMES[i])
			if tweet_text != "" {
				PostTweet(tweet_text, tweet_url, TARGET_NAMES[i])
			}
		}
		time.Sleep(10 * time.Second)
	}
}
