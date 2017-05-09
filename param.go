package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Param struct {
	errorCode   int
	id          string
	version     string
	vnum        int
	baseName    string
	fileName    string
	message     string
	httpStatus  int
	contentType string
	pdfFile     PdfFile
}

func (p *Param) ParseParam(r *http.Request) int {
	r.ParseForm()

	// idなし
	p.id = r.PostFormValue("id")
	if p.id == "" {
		return 1001
	}

	// id該当ファイルなし
	pdfFile, ok := PDF_FILES[p.id]
	if !ok {
		return 1002
	}
	p.pdfFile = pdfFile

	// versionパラメータがなければ最新のバージョンとして扱う
	p.version = r.PostFormValue("version")
	if p.version == "" {
		p.vnum = p.pdfFile.currentVersion

	} else {

		// 数字として不正なversionパラメータ
		vnum, err := strconv.Atoi(p.version)
		if err != nil {
			return 1007
		}
		p.vnum = vnum
	}

	// version範囲が不正
	if p.vnum < 1 || p.pdfFile.currentVersion < p.vnum {
		return 1003
	}

	// エラー無し
	p.baseName = url.QueryEscape(p.pdfFile.baseName)
	p.fileName = fmt.Sprintf("%s_%d.pdf", p.baseName, p.vnum)
	return 0
}
