package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Reads all .json files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	srcPath := filepath.FromSlash("../api/swagger/v1/")
	fs, _ := ioutil.ReadDir(srcPath)
	out, _ := os.Create(filepath.FromSlash("../pkg/api/v1/swagger.pb.go"))
	_, _ = out.Write([]byte("package v1 \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			name := strings.TrimPrefix(f.Name(), "blitzd.")
			_, _ = out.Write([]byte(strings.TrimSuffix(name, ".json") + " = `"))
			f, _ := os.Open(path.Join(srcPath, f.Name()))
			_, _ = io.Copy(out, f)
			_, _ = out.Write([]byte("`\n"))
		}
	}
	_, _ = out.Write([]byte(")\n"))
}
