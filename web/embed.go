package web

import (
	"embed"
	"io/fs"
)

//go:embed public/assets/*
var assets embed.FS

func AssetsFS() (fs.FS, error) {
	return fs.Sub(assets, "public/assets")
}
