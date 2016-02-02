package weather

import (
	log "../lib/logger"
	"../models"
	"../util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Response struct {
	Main     Main      `json:"main"`
	Weathers []Weather `json:"weather"`
	Wind     Wind      `json:"wind"`
	Name     string    `json:"name"`
}

type Main struct {
	Temp     float32 `json:"temp"`     // 気温
	TempMin  float32 `json:"temp_min"` // 最低気温
	TempMax  float32 `json:"temp_max"` // 最高気温
	Humidity int     `json:"humidity"` // 湿度
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Wind struct {
	Speed float32 `json:"speed"` // 風速
	Deg   float32 `json:"deg"`   // 方角
}

type Counds struct {
	All int `json:"all"` // 雲の割合
}

const URL = "http://api.openweathermap.org/data/2.5/weather"

// 指定できる都市一覧 http://openweathermap.org/help/city_list.txt
var JP_CITYS = [2]string{"Tokyo", "Fukuoka-shi"}

var WEATHER_NAME = map[int]string{
	200: "小雨と雷雨",
	201: "雨と雷雨",
	202: "大雨と雷雨",
	210: "光雷雨",
	211: "雷雨",
	212: "重い雷雨",
	221: "ぼろぼろの雷雨",
	230: "小雨と雷雨",
	231: "霧雨と雷雨",
	232: "重い霧雨と雷雨",
	300: "光強度霧雨",
	301: "霧雨",
	302: "重い強度霧雨",
	310: "光強度霧雨の雨",
	311: "霧雨の雨",
	312: "重い強度霧雨の雨",
	313: "にわかの雨と霧雨",
	314: "重いにわかの雨と霧雨",
	321: "にわか霧雨",
	500: "小雨",
	501: "適度な雨",
	502: "重い強度の雨",
	503: "非常に激しい雨",
	504: "極端な雨",
	511: "雨氷",
	520: "光強度のにわかの雨",
	521: "にわかの雨",
	522: "重い強度にわかの雨",
	531: "不規則なにわかの雨",
	600: "小雪",
	601: "雪",
	602: "大雪",
	611: "みぞれ",
	612: "にわかみぞれ",
	615: "光雨と雪",
	616: "雨や雪",
	620: "光のにわか雪",
	621: "にわか雪",
	622: "重いにわか雪",
	701: "ミスト",
	711: "煙",
	721: "ヘイズ",
	731: "砂、ほこり旋回する",
	741: "霧",
	751: "砂",
	761: "ほこり",
	762: "火山灰",
	771: "スコール",
	781: "竜巻",
	800: "晴天",
	801: "薄い雲",
	802: "雲",
	803: "曇りがち",
	804: "厚い雲",
}

func Post(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(JP_CITYS); i++ {
		api_url := buildUrl(JP_CITYS[i])
		println(api_url)
		result := getResult(api_url)
		message := format(result)
		slack_message := models.SlackMessageNew(message, convertCityName(result.Name)+" お天気bot", "", iconUrl(result.Weathers[0].Icon), "#general")
		models.SlackMessagePost(slack_message)

		fmt.Fprintf(w, message+"\n")
		fmt.Printf("%v")
	}
	fmt.Fprintf(w, "")
}

func buildUrl(city_name string) string {
	secret := url.QueryEscape(os.Getenv("openweathermap_api_key"))
	return URL + "?q=" + city_name + ",jp&APPID=" + secret
}

func getResult(url string) Response {
	var weather_response Response

	request, err := http.NewRequest("GET", url, nil)
	client := new(http.Client)
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	util.Perror(err)

	println(string(body))

	parse_err := json.Unmarshal(body, &weather_response)
	util.Perror(parse_err)
	return weather_response
}

// 摂氏に変換
func toCelsius(kelvin float32) int {
	return int(kelvin - 273.15)
}

func iconUrl(icon string) string {
	return "http://openweathermap.org/img/w/" + strings.Replace(icon, "n", "d", 1) + ".png"
}

func windPowerAndDegree(deg float32) string {

	wk_deg := deg
	if wk_deg < 22.5 {
		wk_deg = wk_deg + 360
	}

	if 337.5 <= wk_deg && 382.5 >= wk_deg {
		return "北"
	} else if 22.5 <= wk_deg && 67.5 >= wk_deg {
		return "北東"
	} else if 67.5 <= wk_deg && 112.5 >= wk_deg {
		return "東"
	} else if 112.5 <= wk_deg && 157.5 >= wk_deg {
		return "南東"
	} else if 157.5 <= wk_deg && 202.5 >= wk_deg {
		return "南"
	} else if 202.5 <= wk_deg && 247.5 >= wk_deg {
		return "南西"
	} else if 247.5 <= wk_deg && 292.5 >= wk_deg {
		return "西"
	} else if 292.5 <= wk_deg && 337.5 >= wk_deg {
		return "北西"
	}

	return ""
}

func format(result Response) string {
	temp := strconv.Itoa(toCelsius(result.Main.Temp))
	temp_max := strconv.Itoa(toCelsius(result.Main.TempMax))
	temp_min := strconv.Itoa(toCelsius(result.Main.TempMin))
	humidity := strconv.Itoa(result.Main.Humidity)
	tenki := WEATHER_NAME[result.Weathers[0].Id]
	wind := windPowerAndDegree(result.Wind.Deg) + "の風" + fmt.Sprintf("%.2f", result.Wind.Speed) + "m/s."
	message := "-----------------------------\n" + util.JpCurrentDate() + "\n-----------------------------\n\n" + "天気:" + tenki + "\n" + "気温:" + temp + "\n" + "最高気温:" + temp_max + "\n" + "最低気温:" + temp_min + "\n" + "湿度:" + humidity + "\n" + wind
	return message
}

func convertCityName(name string) string {
	if name == "Tokyo" {
		return "東京"
	} else if name == "Fukuoka-shi" {
		return "福岡"
	}
	return ""
}
