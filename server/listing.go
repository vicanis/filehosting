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
	uri := r.RequestURI

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

	backdir := strings.TrimRight(uri, "/")
	if backdir != "" {
		slashPos := strings.LastIndex(backdir, "/")
		backdir = backdir[:slashPos]
	}

	if backdir == "" {
		backdir = "/"
	}

	buf := bytes.NewBuffer(nil)

	err = template.Must(template.New("listing").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Dir}} - File listing</title>
			<script src="https://cdn.tailwindcss.com"></script>
			<link rel="stylesheet" href="https://maxst.icons8.com/vue-static/landings/line-awesome/line-awesome/1.3.0/css/line-awesome.min.css">
		</head>
		<body>
			<div class="max-w-lg mx-auto">
				<div class="flex items-center justify-between py-2 mb-2 border-b-2">
					{{if ne .Dir "/"}}
						<a href="{{.Backdir}}">
							<span class="las la-angle-left"></span>
							<span>Back</span>
						</a>
					{{else}}
						<span>Root</span>
					{{end}}
					<div class="text-lg text-center">{{.Dir}}</div>
				</div>

				{{$dir := .Dir}}
				<ul>
					{{range .Items}}
						<a href="{{$dir}}{{.Name}}">
							<li class="hover:bg-gray-200 hover:cursor-pointer p-2 flex items-center justify-content gap-2">
								{{if .IsDir}}
									<div class="las la-folder"></div>
								{{end}}

								<span class="grow">{{.Name}}</span>
							</li>
						</a>
					{{end}}
				</ul>
			</div>
		</body>
		</html>
	`)).Execute(buf, struct {
		Dir     string
		Backdir string
		Items   []fs.DirEntry
	}{
		Dir:     uri,
		Backdir: backdir,
		Items:   trimmed,
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
