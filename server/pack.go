package server

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func packHandler(baseDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsed, _ := url.Parse(r.RequestURI)

		z := zip.NewWriter(w)

		err := packDir(baseDir+parsed.Path, "", z)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "pack directory failed: %s", err)
			return
		}

		z.Close()
	})
}

func packDir(baseDir string, dir string, z *zip.Writer) error {
	dirItems, err := os.ReadDir(baseDir + dir)
	if err != nil {
		return err
	}

	for _, item := range dirItems {
		if item.IsDir() {
			packDir(baseDir, dir+"/"+item.Name(), z)
			continue
		}

		zipFile, err := z.Create(dir + "/" + item.Name())
		if err != nil {
			return err
		}

		path := baseDir + dir + item.Name()

		localFile, err := os.Open(path)
		if err != nil {
			log.Printf("file %s open failed: %s", path, err)
			continue
		}

		_, err = io.Copy(zipFile, localFile)

		localFile.Close()

		if err != nil {
			log.Printf("file %s copy failed: %s", path, err)
			continue
		}
	}

	return nil
}
