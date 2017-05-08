package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-martini/martini"
	"github.com/k0kubun/pp"
)

type PdfFile struct {
	currentVersion int
	baseName       string
}

type ErrorMessage struct {
	status  int
	label   string
	message string
}

type Param struct {
	errorCode     int
	id            string
	version       string
	vnum          int
	baseName      string
	message       string
	pdfFile       PdfFile
	errorMessage  ErrorMessage
	pdfFiles      map[string]PdfFile
	errorMessages map[int]ErrorMessage
}

func main() {
	p := new(Param)
	p.GetErrorMessages()
	p.GetPdfFiles()

	m := martini.Classic()

	m.Post("/rfm/logic/api/pdfdownload", func(w http.ResponseWriter, r *http.Request) string {
		p.ParseParam(r)

		switch p.errorCode {
		case 1001, 1002, 1003:
			w.Header().Set("Content-Type", "text/plain")
			baseName := url.QueryEscape("ダウンロードエラー")
			timestamp := time.Now().Format("20060102150405")
			filename := fmt.Sprintf("%s_%s.txt", baseName, timestamp)
			w.Header().Set("Content-Disposition", fmt.Sprintf("filename*=''%s", filename))
			w.WriteHeader(p.Status())

		default:
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=''%s", p.FileName()))
			w.WriteHeader(200)
		}

		return p.message
	})
	m.RunOnAddr("0.0.0.0:8080")
}

func (p *Param) GetErrorMessages() {
	p.errorMessages = map[int]ErrorMessage{
		1001: {status: 422, label: "パラメータチェック", message: "必須パラメータが指定されていません。%s"},
		1002: {status: 404, label: "ID不正", message: "ファイル(ID:%s)を取得できませんでした。"},
		1003: {status: 422, label: "バージョン不正", message: "バージョン不正です(最新バージョン:%d、パラメータ指定バージョン:%s)。"},
		1007: {status: 422, label: "バージョン不正", message: "バージョン不正です(パラメータ指定バージョン:%s)。"},
		1004: {status: 422, label: "キャビネット不正", message: "キャビネット「%s」のセッションが生成できませんでした。"},
		1005: {status: 422, label: "ファイル不正", message: "ファイル(ファイルID:%s,ver.%s)を取得できませんでした。"},
		1006: {status: 500, label: "システムエラー", message: "システムエラーが発生しました。%s"},
	}
}

func (p *Param) GetPdfFiles() {
	p.pdfFiles = map[string]PdfFile{
		"a001": {currentVersion: 1, baseName: "テスト用日本語ファイル名称①"},
		"b002": {currentVersion: 2, baseName: "テスト用日本語ファイル名称❷"},
		"c003": {currentVersion: 3, baseName: "テスト用日本語ファイル名称Ⅲ"},
	}
}

func (p *Param) FileName() string {
	return fmt.Sprintf("%s_%d.pdf", url.QueryEscape(p.pdfFile.baseName), p.vnum)
}

func (p *Param) Tmpl() string {
	return p._err().message
}

func (p *Param) Status() int {
	return p._err().status
}

func (p *Param) _err() ErrorMessage {
	e, ok := p.errorMessages[p.errorCode]
	if !ok {
		panic(p)
	}
	return e
}

func (p *Param) ParseParam(r *http.Request) {
	r.ParseForm()

	p.id = r.PostFormValue("id")
	if p.id == "" {
		p.errorCode = 1001
		p.message = fmt.Sprintf(p.Tmpl(), "ID")
		return
	}

	pdfFile, ok := p.pdfFiles[p.id]
	if !ok {
		p.errorCode = 1002
		p.message = fmt.Sprintf(p.Tmpl(), p.id)
		return
	}
	p.pdfFile = pdfFile

	p.version = r.PostFormValue("version")
	if p.version == "" {
		p.vnum = p.pdfFile.currentVersion

	} else {

		vnum, err := strconv.Atoi(p.version)
		if err != nil {
			p.errorCode = 1007
			p.message = fmt.Sprintf(p.Tmpl(), p.version)
			pp.Print(p)
			return
		}
		p.vnum = vnum
	}

	if p.vnum < 1 || p.pdfFile.currentVersion < p.vnum {
		p.errorCode = 1003
		p.message = fmt.Sprintf(p.Tmpl(), p.pdfFile.currentVersion, p.vnum)
		return
	}
}
