package public

import (
	"embed"
	"io/fs"
	"net/http"
)

var (

	//go:embed assets/*
	assetsFs embed.FS
	//go:embed templates/*
	Templates embed.FS
)

// Assets 返回以Assets作为根目录的对象
func Assets() http.FileSystem {
	// even assets is empty, fs.Sub won't fail
	stripped, err := fs.Sub(assetsFs, "assets")
	if err != nil {
		panic(err)
	}
	return http.FS(stripped)
}
