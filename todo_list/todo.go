package todo

import (
	log "../lib/logger"
	"../models"
	"../util"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const ADD = "add"
const DEL = "del"
const LIST = "list"
const CLEAR = "clear"
const STORE_DIR = "todo_list/stored_files"

var p = fmt.Println

func Accept(text string, channel_name string) interface{} {
	var command string
	var message string
	var rtn_text interface{}

	if validateParams(text) {
		command, message = parseText(text)
	} else {
		return "入力されたパラメータが不正です→ " + text
	}

	if command == ADD {
		rtn_text = add(channel_name, message)
	} else if command == DEL {
		rtn_text = del(channel_name, message)
	} else if command == LIST {
		rtn_text = list(channel_name)
	} else if command == CLEAR {
		rtn_text = clear(channel_name)
	}

	return rtn_text
}

func validateParams(text string) bool {
	var command string

	if len(strings.Split(text, " ")) < 2 {
		return false
	}

	command = getCommand(text)
	correct_commands := []string{ADD, LIST, DEL, CLEAR}
	if !util.Contains(correct_commands, command) {
		return false
	}

	return true
}

func parseText(text string) (command string, post_text string) {
	command = getCommand(text)
	post_text = getMessage(text)
	return
}

func getTriggerWord(text string) string {
	return strings.Split(text, " ")[0]
}

func getCommand(text string) string {
	return strings.Split(text, " ")[1]
}

func getMessage(text string) string {
	if len(strings.Split(text, " ")) > 2 {
		return strings.Split(text, " ")[2]
	} else {
		return ""
	}
}

func add(channel_name string, message string) string {
	// ディレクトリの作成
	err := os.MkdirAll(STORE_DIR, 0777)
	if err != nil {
		util.Perror(err)
	}

	file_path := getStoredPath(channel_name)
	// ファイルがなかったら作る
	if !util.FileExists(file_path) {
		util.CreateFile(file_path)
	}
	f, err := os.OpenFile(file_path, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return "ファイルのオープンに失敗しました"
	}
	defer f.Close()

	message = strings.Replace(message, "\n", " ", -1)
	if _, err = f.WriteString(message + "\n"); err != nil {
		return "書き込みに失敗しました"
	}

	lines, _ := getList(file_path)
	return "登録に成功しました :wink: \n 現在のタスク \n " + strings.Join(lines, "\n")
}

func list(channel_name string) string {
	file_path := getStoredPath(channel_name)

	if !util.FileExists(file_path) {
		return "まだ何も書き込まれていません"
	}

	lines, _ := getList(file_path)

	// TODO 文言をランダムで変えたい
	return "現在のタスクです。\n気張っていきましょー :kissing_heart: \n\n " + strings.Join(lines, "\n")
}

func del(channel_name string, message string) string {
	var del_flg bool
	var new_lines []string
	message = strings.Replace(message, "\n", " ", -1)
	file_path := getStoredPath(channel_name)

	if !util.FileExists(file_path) {
		return "まだ何も書き込まれていません"
	}

	lines, err := util.ReadLines(file_path)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for i := 0; i < len(lines); i++ {
		if lines[i] == message || strconv.Itoa(i+1) == message {
			del_flg = true
		} else {
			new_lines = append(new_lines, lines[i])
		}
	}

	if del_flg {
		if len(new_lines) == 0 {
			if err := os.Remove(file_path); err != nil {
				return "クリアに失敗しました"
			}
		} else {
			content := []byte(strings.Join(new_lines, "\n") + "\n")
			ioutil.WriteFile(file_path, content, 0600)
		}

		if len(new_lines) != 0 {
			lines, _ = getList(file_path)
		}
		return "削除に成功しました :neutral_face: \n 残りのタスク \n " + strings.Join(lines, "\n")
	} else {
		return "一致するものが見つかりませんでした"
	}
}

func clear(channel_name string) models.SlackMessage {
	file_path := getStoredPath(channel_name)
	if err := os.Remove(file_path); err != nil {
		return models.SlackMessage{Text: "クリアに失敗しました"}
	}

	return models.SlackMessage{Text: "クリアしました"}
}

func getStoredPath(channel_name string) string {
	return STORE_DIR + "/" + channel_name
}

func getList(file_path string) ([]string, error) {
	lines, err := util.ReadLines(file_path)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	var items []string
	for i, line := range lines {
		items = append(items, strconv.Itoa(i+1)+". "+line)
	}

	items = append([]string{"```"}, items...)
	items = append(items, "```")

	return items, err
}
