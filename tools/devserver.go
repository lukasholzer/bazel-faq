package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func listDirectory(folder string) {
	var files []string

	err := filepath.Walk(folder, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}

func main() {
	entry := flag.String("entry", "", "The entry file. (Required)")
	root := flag.String("root", "./", "The root to serve the files from")
	flag.Parse()

	port := 3000

	// listDirectory(a)
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
				}
			},
		},
	})

	if len(result.Errors) > 0 {
		fmt.Printf("Error occured", result.Errors)
		os.Exit(1)
	}

	fmt.Printf("👌 build success ")
	fileServer := http.FileServer(http.Dir(*root))
	http.Handle("/", fileServer)

	fmt.Printf("Starting server at http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}

}
