package main

type PdfFile struct {
	currentVersion int
	baseName       string
}

var PDF_FILES = map[string]PdfFile{
	"a001": {currentVersion: 1, baseName: "テスト用日本語ファイル名称①"},
	"b002": {currentVersion: 2, baseName: "テスト用日本語ファイル名称❷"},
	"c003": {currentVersion: 3, baseName: "テスト用日本語ファイル名称Ⅲ"},
}
