package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func listingHandler(baseDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsed, _ := url.Parse(r.RequestURI)

		uri := parsed.Path

		dir := baseDir + "/" + uri

		dirItems, err := os.ReadDir(dir)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "listing failed: %s", err)
			return
		}

		trimmed := make([]fs.DirEntry, 0)

		for _, item := range dirItems {
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

		prevDir := strings.TrimRight(uri, "/")
		if prevDir != "" {
			slashPos := strings.LastIndex(prevDir, "/")
			prevDir = prevDir[:slashPos]
		}

		if prevDir == "" {
			prevDir = "/"
		}

		buf := bytes.NewBuffer(nil)

		err = template.Must(template.New("listing").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Dir}} - File listing</title>
			<script src="https://cdn.tailwindcss.com"></script>
			<script src="https://cdn.jsdelivr.net/npm/@mdi/font@7.3.67/scripts/verify.min.js"></script>
			<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@mdi/font@7.3.67/css/materialdesignicons.min.css" integrity="sha256-akFxqbgnSEftsMESNX9beHAwLq+cU+tEQPGC8Ft9U2Y=" crossorigin="anonymous">
		</head>
		<body>
			<div class="max-w-lg mx-auto">
				<div class="flex items-center justify-between py-2 mb-2 border-b-2">
					{{if ne .Dir "/"}}
						<a href="{{.PrevDir}}">
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
						<li class="hover:bg-gray-200 hover:cursor-pointer p-2 flex items-center justify-content gap-2">
							{{if .IsDir}}
								<span class="mdi mdi-folder-outline"></span>
							{{else}}
								<span class="mdi mdi-file-document-outline"></span>
							{{end}}

							<span class="grow">
								<a href="{{$dir}}{{.Name}}">
									{{.Name}}
								</a>
							</span>

							<div class="flex gap-3">
								{{if .IsDir}}
									<a class="p-1 hover:bg-gray-400 rounded-md" href="{{$dir}}{{.Name}}?pack" title="Download ZIP archive">
										<span class="mdi mdi-folder-zip"></span>
									</a>
								{{else}}
									<a class="p-1 hover:bg-gray-400 rounded-md" href="{{$dir}}{{.Name}}" title="Download file">
										<span class="mdi mdi-file-download-outline"></span>
									</a>
								{{end}}
							</div>
						</li>
					{{end}}
				</ul>
			</div>
		</body>
		</html>
	`)).Execute(buf, struct {
			Dir     string
			PrevDir string
			Items   []fs.DirEntry
		}{
			Dir:     uri,
			PrevDir: prevDir,
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
	})
}
