package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
)

func main() {
	// ワーク構造体を正常時設定で初期化
	p := new(Param)
	p.httpStatus = 200
	p.contentType = "application/pdf"

	// WAF初期化
	m := martini.Classic()

	// ダウンロードAPI [ POST /rfm/logic/api/pdfdownload ]
	m.Post("/rfm/logic/api/pdfdownload", func(w http.ResponseWriter, r *http.Request) string {

		// パラメータを読み取りエラー処理
		errorCode := p.ParseParam(r)
		if errorCode != 0 {
			p.SetErrorCode(errorCode)
		}

		// 応答電文を編集
		w.Header().Set("Content-Type", p.contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=''%s", p.fileName))
		w.WriteHeader(p.httpStatus)

		return p.message
	})

	// サーバを開始
	m.RunOnAddr("0.0.0.0:80")
}
