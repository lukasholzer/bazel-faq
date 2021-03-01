package bundle

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type Asset struct {
	Location string
	FilePath string
}

type BundleArgs struct {
	Root      string
	Entry     string
	IndexHTML string
	Outdir    string
	Assets    []Asset
	WatchMode api.WatchMode
}

func injectAssets(indexHTML string, assets []Asset) string {
	// Read the indexHTML file
	byteArray, err := ioutil.ReadFile(indexHTML)

	if err != nil {
		panic(err)
	}

	var file = string(byteArray)

	for _, a := range assets {
		var location string
		var tag string

		switch filepath.Ext(a.FilePath) {
		case ".css":
			location = "head"
			tag = fmt.Sprintf(`<link href="%s" rel="stylesheet" />`, a.FilePath)
		case ".js":
			location = "body"
			tag = fmt.Sprintf(`<script src="%s"></script>`, a.FilePath)
		}

		if len(a.Location) > 0 {
			location = a.Location
		}

		regex := regexp.MustCompile(fmt.Sprintf("</%s>", location))
		file = regex.ReplaceAllString(file, fmt.Sprintf("%s</%s>", tag, location))
	}

	return file
}

func relative(path string, relativeTo string) string {
	return strings.TrimPrefix(filepath.ToSlash(path), relativeTo)
}

func Bundle(bundleArgs BundleArgs) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bundleArgs)
	fmt.Printf("CWD: %s\n", cwd)

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{bundleArgs.Entry},
		Bundle:      true,
		Write:       true,
		Target:      api.Target(api.ESNext),
		Format:      api.Format(api.FormatESModule),
		Outdir:      bundleArgs.Outdir,
		Color:       api.ColorAlways,
		Watch:       &bundleArgs.WatchMode,
	})

	if len(result.Errors) > 0 {
		fmt.Println(result.Errors)
	}

	assets := bundleArgs.Assets

	for _, out := range result.OutputFiles {
		fmt.Println(out.Path)
		assets = append(assets, Asset{FilePath: relative(out.Path, filepath.Join(cwd, bundleArgs.Root))})
	}

	updatedIndex := injectAssets(bundleArgs.IndexHTML, assets)

	ioutil.WriteFile(filepath.Join(bundleArgs.Root, "index.html"), []byte(updatedIndex), 0644)
}
