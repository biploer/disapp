package conf

import (
	"embed"
	"io/fs"
)

//go:embed *.yaml
var config embed.FS

func ConfigFS() fs.FS {
	return config
}
