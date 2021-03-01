package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
	"github.com/lukasholzer/bazel-faq/bundle"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// The Websocket connection
var conn *websocket.Conn = nil

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

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func onReload() {
}

var watchMode = api.WatchMode{
	SpinnerBusy: "Rebuilding...",
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
}

func writeLiveReload(port int, root string) {
	liveReload := fmt.Sprintf(`
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
	ioutil.WriteFile(filepath.Join(root, "_reload.js"), []byte(liveReload), 0644)
}

func main() {
	root := flag.String("root", "", "The root path that should be served (Required)")
	port := flag.Int("port", 3000, "The port where it should be served")
	entry := flag.String("entry", "", "The entry file. (Required)")
	index := flag.String("index", "", "The index HTML file that should be served! (Required)")
	openBrowser := flag.Bool("open", true, "If it should open the url in the borwser")
	flag.Parse()

	bundle.Bundle(bundle.BundleArgs{
		Entry:     *entry,
		Root:      *root,
		IndexHTML: *index,
		WatchMode: watchMode,
		Outdir:    filepath.Join(*root, "dist"),

		Assets: []bundle.Asset{{FilePath: "/_reload.js", Location: "head"}},
		// Assets:
	})

	writeLiveReload(*port, *root)

	url := fmt.Sprintf("http://localhost:%d/", *port)
	fileServer := http.FileServer(http.Dir(*root))
	http.Handle("/", fileServer)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ = upgrader.Upgrade(w, r, nil)
	})

	if *openBrowser {
		go func() {
			<-time.After(100 * time.Millisecond)
			open(url)
		}()
	}

	fmt.Printf("\n🌍 Starting server at %s\n", url)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
