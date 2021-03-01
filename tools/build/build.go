package main

import (
	"flag"
	"path/filepath"

	"github.com/lukasholzer/bazel-faq/bundle"
)

func main() {
	root := flag.String("root", "", "The root path that should be served (Required)")
	entry := flag.String("entry", "", "The entry file. (Required)")
	index := flag.String("index", "", "The index HTML file that should be served! (Required)")
	flag.Parse()

	bundle.Bundle(bundle.BundleArgs{
		Entry:     *entry,
		Root:      *root,
		IndexHTML: *index,
		Outdir:    filepath.Join(*root, "out"),
	})

}
