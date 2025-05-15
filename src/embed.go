package main

import (
	"embed"
	_ "image/jpeg" // jpeg画像を扱うために必要
)

//go:embed static/*
var staticFiles embed.FS

// GetEmbeddedFile は埋め込まれたファイルを返します
func GetEmbeddedFile() embed.FS {
	return staticFiles
}
