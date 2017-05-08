package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-martini/martini"
)

func main() {
	messages := map[string]string{
		"パラメータチェック": "必須パラメータが指定されていません。%s",
		"ID不正":      "ファイル(ID:%s)を取得できませんでした。",
		"バージョン不正":   "バージョン不正です(最新バージョン:%d、パラメータ指定バージョン:%d)。",
		"キャビネット不正":  "キャビネット「%s」のセッションが生成できませんでした。",
		"ファイル不正":    "ファイル(ファイルID:%s,ver.%s)を取得できませんでした。",
		"システムエラー":   "システムエラーが発生しました。%s",
	}

	files := map[string]string{
		"1": "テスト用日本語ファイル名称①",
		"2": "テスト用日本語ファイル名称❷",
		"3": "テスト用日本語ファイル名称Ⅲ",
	}

	m := martini.Classic()
	m.Post("/rfm/logic/api/pdfdownload", func(w http.ResponseWriter, r *http.Request) string {
		r.ParseForm()
		id := r.PostFormValue("id")
		version := r.PostFormValue("version")
		version_num, version_err := strconv.Atoi(version)
		basename, file_err := files[id]
		var message string

		switch {

		case id == "":
			w.WriteHeader(422)
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", "filename*=ダウンロードエラー_yyyyMMddHHmmss.txt")
			message = fmt.Sprintf(messages["パラメータチェック"], "ID")

		case file_err:
			w.WriteHeader(404)
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", "filename*=ダウンロードエラー_yyyyMMddHHmmss.txt")
			message = fmt.Sprintf(messages["ID不正"], id)

		case version_num > 3 || version_num < 1 || version_err != nil:
			w.WriteHeader(422)
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", "filename*=ダウンロードエラー_yyyyMMddHHmmss.txt")
			message = fmt.Sprintf(messages["バージョン不正"], 3, version_num)

		default:
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/pdf")
			filename := fmt.Sprintf("%s_%s.pdf", url.QueryEscape(basename), version)
			value := fmt.Sprintf("attachment; filename*=''%s", filename)
			w.Header().Set("Content-Disposition", value)
		}

		return message
	})
	m.RunOnAddr("0.0.0.0:9999")
}
