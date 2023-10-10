package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"
)

type listingHandler struct {
	Dir string
}

func (h listingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := strings.TrimPrefix(r.RequestURI, "/files")

	dir := h.Dir + "/" + uri

	list, err := os.ReadDir(dir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "listing failed: %s", err)
		return
	}

	trimmed := make([]fs.DirEntry, 0)

	for _, item := range list {
		if !strings.HasPrefix(item.Name(), ".") {
			trimmed = append(trimmed, item)
		}
	}

	sort.Slice(trimmed, func(i, j int) bool {
		if trimmed[i].IsDir() && !trimmed[j].IsDir() {
			return true
		} else if !trimmed[i].IsDir() && trimmed[j].IsDir() {
			return false
		}

		return strings.Compare(trimmed[i].Name(), trimmed[j].Name()) < 0
	})

	buf := bytes.NewBuffer(nil)

	err = template.Must(template.New("listing").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Dir}} - File listing</title>
			<script src="https://cdn.tailwindcss.com"></script>
		</head>
		<body>
			<div class="max-w-lg mx-auto">
				<div class="text-lg text-center">{{.Dir}}</div>
				<ul>
					{{range .Items}}
						<li class="py-2 border-t-2 my-2">
							{{.Name}}
						</li>
					{{end}}
				</ul>
			</div>
		</body>
		</html>
	`)).Execute(buf, struct {
		Dir   string
		Items []fs.DirEntry
	}{
		Dir:   uri,
		Items: trimmed,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "render failed: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprint(w, buf.String())
}
