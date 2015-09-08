package slackbot_responder

import (
	"../util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

const WordTable = "slackbot_responder/word.yml"

func getWordTable() map[string][]string {
	buf, err := ioutil.ReadFile(WordTable)
	if err != nil {
		panic(err)
	}

	m := make(map[string][]string)
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}

	return m
}

func DetectWord(post_text string) (bot_message string) {
	dict := getWordTable()       // 設定を取得
	word_keys := util.Keys(dict) // キーのみ抽出
	var suggestion []string

	for i := 0; i < len(word_keys); i++ {
		if strings.Contains(post_text, word_keys[i]) {
			suggestion = append(suggestion, dict[word_keys[i]]...)
		}
	}

	if len(suggestion) > 0 {
		rand.Seed(time.Now().UnixNano())
		bot_message = suggestion[rand.Intn(len(suggestion))]
	}
	return bot_message
}
