package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Reads all .json files in the specified srcPath folder
// and encodes them as strings literals in dstPath
func main() {
	srcPath := filepath.FromSlash("../api/swagger/v1/")
	dstPath := filepath.FromSlash("../pkg/api/v1/swagger.pb.go")
	fs, _ := ioutil.ReadDir(srcPath)
	out, _ := os.Create(dstPath)

	fmt.Println("writing", dstPath)

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
