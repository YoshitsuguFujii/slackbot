package qiitabot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var p = fmt.Println

const PER_PAGE = 100

type UserResponce struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type StockResponce struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	RenderedBody string `json:"rendered_body"`
	Body         string `json:"body"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func getUser(text string) string {
	return strings.Split(text, " ")[1]
}

func isUserExist(user string) bool {
	var user_responce UserResponce
	p("request: -> " + "https://qiita.com/api/v2/users/" + user)
	request, _ := http.NewRequest("GET", "https://qiita.com/api/v2/users/"+user, nil)
	request.Header.Set("Authorization", "Bearer "+getToken())

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	perror(err)

	p(string(body))
	parse_err := json.Unmarshal(body, &user_responce)
	perror(parse_err)
	//fmt.Printf("id: %s", user_responce.Id)
	if user_responce.Id != "" {
		return true
	} else {
		return false
	}
}

func getStock(user string) string {
	var total_count int
	var total_page int
	var stock_responces []StockResponce

	total_count = getTotalCount(user)
	total_page = total_count / PER_PAGE
	rand.Seed(time.Now().UnixNano())
	page := rand.Intn(total_page)
	if page == 0 {
		page = 1
	}
	get_url := "https://qiita.com/" + "/api/v2/users/" + user + "/stocks" + "?page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(PER_PAGE)
	p("request: -> " + get_url)
	request, _ := http.NewRequest("GET", get_url, nil)
	request.Header.Set("Authorization", "Bearer "+getToken())

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	perror(err)

	p("body:" + response.Header.Get("Total-Count"))
	parse_err := json.Unmarshal(body, &stock_responces)
	perror(parse_err)
	//fmt.Printf("id: %s", len(stock_responces))
	var url string
	p(len(stock_responces))
	if len(stock_responces) > 0 {
		rand.Seed(time.Now().UnixNano())
		stock_responce := stock_responces[rand.Intn(len(stock_responces))]
		url = stock_responce.Title + "\n" + stock_responce.Url
	}
	return url
}

func UserStockSample(text string) string {
	//currrentUser()
	var stock_url string
	user := getUser(text)
	if isUserExist(user) {
		stock_url = getStock(user)
	} else {
		stock_url = "そんな人いませんでした"
	}
	return stock_url
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	return out.Bytes(), err
}

func currrentUser() {
	request, _ := http.NewRequest("GET", "https://qiita.com/api/v2/authenticated_user", nil)
	request.Header.Set("Authorization", "Bearer "+getToken())

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	data, err := prettyPrint(contents)

	fmt.Println(string(data))
}

func getTotalCount(user string) int {
	get_url := "https://qiita.com/" + "/api/v2/users/" + user + "/stocks" + "?per_page=1"
	p("count request: -> " + get_url)
	request, _ := http.NewRequest("GET", get_url, nil)
	request.Header.Set("Authorization", "Bearer "+getToken())

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	total_count := response.Header.Get("Total-Count")
	p("total count :" + total_count)
	i, _ := strconv.Atoi(total_count)
	return i
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func getToken() string {
	return os.Getenv("qiita_token")
}
