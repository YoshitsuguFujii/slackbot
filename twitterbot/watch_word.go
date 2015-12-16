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
	"strconv"
	"time"
)

const WordTable = "twitterbot/watch_word.yml"

var word_settings WatchWordSettings

type WatchWordSetting struct {
	Word        string `word`
	Channel     string `channel`
	BeforeId    int    `before_id`
	GetCountPer int    `get_count_per`
	ExceptWord  string `except_word`
}

type WatchWordSettings struct {
	Word []WatchWordSetting `watch_words`
}

type SearchResponce struct {
	Statuses       []SearchResult `json:"statuses"`
	SearchMetadata struct {
		MaxId       int     `json:"max_id"`
		SinceId     int     `json:"since_id"`
		RefreshUrl  string  `json:"refresh_url"`
		NextResults string  `json:"next_results"`
		Count       int     `json:"count"`
		CompletedIn float32 `json:"completed_in"`
		SinceIdStr  string  `json:"since_id_str"`
		Query       string  `json:"query"`
		MaxIdStr    string  `json:"max_id_str"`
	} `json:"search_metadata"`
}

type SearchResult struct {
	Id         int    `json:"id"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
	FromUser   string `json:"from_user"`
	FromUserId int    `json:"from_user_id"`
	User       struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

func initialize_watch_word_bot() {
	buf, err := ioutil.ReadFile(WordTable)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &word_settings)

	util.Perror(err)
}

func WatchWord() {
	initialize_watch_word_bot()

	for {
		for i := 0; i < len(word_settings.Word); i++ {
			go func(search_info *WatchWordSetting) {
				text, from_user, last_id := search(search_info)
				if last_id > 0 {
					for i := 0; i < len(text); i++ {
						postTweetToSlack(text[i], from_user[i], search_info.Channel)
					}
				}
			}(&word_settings.Word[i])
		}
		time.Sleep(300 * time.Second)
	}
}

func search(search_info *WatchWordSetting) (search_text []string, user_name []string, last_id int) {
	bearer := getToken()
	values := url.Values{}
	log.Info("検索ワード:" + search_info.Word)
	values.Add("q", "\""+search_info.Word+"\""+search_info.ExceptWord)
	values.Add("count", strconv.Itoa(search_info.GetCountPer))
	values.Add("lang", "ja")
	values.Add("locale", "ja")
	values.Add("result_type", "recent")
	values.Add("since_id", strconv.Itoa(search_info.BeforeId))
	url := "https://api.twitter.com/1.1/search/tweets.json"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.URL.RawQuery = values.Encode()
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	util.Perror(err)

	var search_responce SearchResponce
	parse_err := json.Unmarshal(body, &search_responce)
	util.Perror(parse_err)

	if len(search_responce.Statuses) > 0 {
		log.Info(search_info)
		log.Info("取得id:" + strconv.Itoa(search_responce.Statuses[0].Id))

		total_count := len(search_responce.Statuses)
		for i := 0; i < total_count; i++ {
			//if search_info.BeforeId != search_responce.Statuses[i].Id {
			search_text = append(search_text, search_responce.Statuses[i].Text)
			user_name = append(user_name, "@"+search_responce.Statuses[i].User.ScreenName)
			//}
		}
		last_id = search_responce.Statuses[0].Id
		search_info.BeforeId = last_id
	} else {
		last_id = 0
	}

	return
}
