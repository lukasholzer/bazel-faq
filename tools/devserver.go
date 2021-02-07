package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"runtime"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type asset struct {
	location string
	filePath string
}

func injectScripts(indexHTML string, assets []asset) string {
	// Read the indexHTML file
	byteArray, err := ioutil.ReadFile(indexHTML)

	if err != nil {
		panic(err)
	}

	var file = string(byteArray)

	for _, a := range assets {
		var location string
		var tag string

		switch filepath.Ext(a.filePath) {
		case ".css":
			location = "head"
			tag = fmt.Sprintf(`<link href="%s" rel="stylesheet" />`, a.filePath)
		case ".js":
			location = "body"
			tag = fmt.Sprintf(`<script src="%s"></script>`, a.filePath)
		}

		if len(a.location) > 0 {
			location = a.location
		}

		regex := regexp.MustCompile(fmt.Sprintf("</%s>", location))
		file = regex.ReplaceAllString(file, fmt.Sprintf("%s</%s>", tag, location))
	}

	return file
}

func relative(path string, relativeTo string) string {
	return strings.TrimPrefix(filepath.ToSlash(path), relativeTo)
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	cwd, _ := os.Getwd()
	entry := flag.String("entry", "", "The entry file. (Required)")
	index := flag.String("index", "", "The index file.")
	root := flag.String("root", "./", "The root to serve the files from")
	flag.Parse()

	var conn *websocket.Conn = nil
	port := 3000

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{*entry},
		Bundle:      true,
		Write:       true,
		Target:      api.Target(api.ESNext),
		Format:      api.Format(api.FormatESModule),
		Outdir:      filepath.Join(*root, "dist"),
		Watch: &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				if len(result.Errors) > 0 {
					fmt.Printf("watch build failed: %d errors\n", len(result.Errors))
				} else {
					fmt.Printf("watch build succeeded: %d warnings\n", len(result.Warnings))
					if conn != nil {
						conn.WriteMessage(1, []byte(`{"type": "reload"}`))
						fmt.Println("Reloading Browser...")
					}

				}
			},
		},
	})

	reload := fmt.Sprintf(`
var socket = new WebSocket("ws://localhost:%d/ws");

socket.onopen = () => {
		console.log("Connected to Live reload!\n");
};

socket.onmessage = (e) => {
	const {type} = JSON.parse(e.data);
	if (type === 'reload') {
		window.location.reload()
	}
};
`, port)

	ioutil.WriteFile(filepath.Join(*root, "_reload.js"), []byte(reload), 0644)

	assets := []asset{{filePath: "/_reload.js", location: "head"}}

	for _, out := range result.OutputFiles {
		assets = append(assets, asset{filePath: relative(out.Path, filepath.Join(cwd, *root))})
	}

	updated := injectScripts(*index, assets)
	ioutil.WriteFile(filepath.Join(*root, "index.html"), []byte(updated), 0644)

	if len(result.Errors) > 0 {
		fmt.Println("Error occured", result.Errors)
		os.Exit(1)
	}

	fmt.Println("👌 build success ")
	fileServer := http.FileServer(http.Dir(*root))
	http.Handle("/", fileServer)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ = upgrader.Upgrade(w, r, nil)
	})

	url := fmt.Sprintf("http://localhost:%d/", port)

	go func() {
		<-time.After(100 * time.Millisecond)
		open(url)
	}()

	fmt.Printf("\n🌍 Starting server at %s\n", url)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
