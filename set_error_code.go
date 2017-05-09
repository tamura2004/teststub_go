package main

import (
	"fmt"
	"net/url"
)

func (p *Param) SetErrorCode(errorCode int) {
	switch errorCode {

	case 1001:
		p.httpStatus = 422
		tmpl := "必須パラメータが指定されていません。%s"
		p.message = fmt.Sprintf(tmpl, "id")

	case 1002:
		p.httpStatus = 404
		tmpl := "ファイル(ID:%s)を取得できませんでした。"
		p.message = fmt.Sprintf(tmpl, p.id)

	case 1003:
		p.httpStatus = 422
		tmpl := "バージョン不正です(最新バージョン:%d、パラメータ指定バージョン:%s)。"
		p.message = fmt.Sprintf(tmpl, p.id)

	case 1007:
		p.httpStatus = 422
		tmpl := "バージョン不正です(パラメータ指定バージョン:%d)。"
		p.message = fmt.Sprintf(tmpl, p.id)

	case 1004:
		p.httpStatus = 500
		//		tmpl := "キャビネット「%s」のセッションが生成できませんでした。"

	case 1005:
		p.httpStatus = 500
		//		tmpl := "ファイル(ファイルID:%s,ver.%s)を取得できませんでした。"

	case 1006:
		p.httpStatus = 500
		//		tmpl := "システムエラーが発生しました。%s"

	default:
		panic("不正なエラーコードです")
	}

	p.contentType = "text/plain"
	p.baseName = url.QueryEscape("ダウンロードエラー")
	p.fileName = fmt.Sprintf("%s_%s.txt", p.baseName, TimeStamp())
}
