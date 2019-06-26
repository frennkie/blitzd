package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
	"path/filepath"
)

var AssetsFs http.FileSystem = http.Dir(filepath.FromSlash("assets/src"))

func main() {
	err := vfsgen.Generate(
		AssetsFs,
		vfsgen.Options{
			Filename:     "./assets/assets_vfsdata.go",
			PackageName:  "assets",
			VariableName: "AssetsFs",
		})
	if err != nil {
		log.Fatalln(err)
	}
}
