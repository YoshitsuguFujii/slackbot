package api_mock

import (
	log "../lib/logger"
	"../util"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"net/http"
	"os"
)

const API_MOCK_STORED_DIR = "api-mock"

func Show(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) > 0 {
		filename := r.URL.Query().Get("file_name")
		log.Info("ファイル名: " + filename)

		// ディレクトリの作成
		err := os.MkdirAll(API_MOCK_STORED_DIR, 0777)
		if err != nil {
			util.Perror(err)
		}

		file_path := API_MOCK_STORED_DIR + "/" + filename
		if !util.FileExists(file_path) {
			log.Fatal("指定されたファイルが見つかりません -> [%s]", filename)
			fmt.Fprintf(w, "[%s] 指定されたファイルが見つかりません ", filename)
		}

		fileBytes, err := ioutil.ReadFile(file_path)

		buffer, err := yaml.YAMLToJSON(fileBytes)

		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "[%s] のパース処理でエラーが発生しました ", filename)
		}

		// 正常に返す
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(buffer)))
	} else {
		// エラー
		log.Fatal("urlが不正です")
		fmt.Fprintf(w, "urlが不正です")
	}
}
