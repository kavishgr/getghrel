package main

import (
	"fmt"
	"github.com/kavishgr/getghrel/github"
	"github.com/kavishgr/getghrel/options"
	"github.com/kavishgr/getghrel/utils"
	"os"
	"sync"
)

func main() {

	var (
		opts           = options.ParseFlags()
		skipextraction = opts.SkipExtraction
		token          = opts.GHToken
		tempdir        = opts.TempDir
		ost, arch      = utils.OsInfo()
		regex          = utils.SetRegex(ost, arch)
		stdInUrls      = make(chan string)
		jobs           sync.WaitGroup
	)

	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable is not found.")
		fmt.Println("Nor is -ghtoken provided on the command line.")
		fmt.Println("")
		fmt.Println("Run 'getghrel -h'")
		fmt.Println("Or browse to: 'https://github.com/kavishgr/getghrel'")
		os.Exit(1)
	}

	go utils.ScanStdIn(stdInUrls)

	if opts.List {
		for c := 0; c < opts.Concurrency; c++ {
			jobs.Add(1)
			go github.FetchGithubReleaseUrl(stdInUrls, &jobs, regex, token)
		}
	}

	if opts.Download {
		_, err := os.Stat(tempdir)

		if os.IsNotExist(err) {
			err = os.Mkdir(tempdir, 0755)
			if err != nil {
				panic(err)
			}
		}

		for c := 0; c < opts.Concurrency; c++ {
			jobs.Add(1)
			go github.DownloadRelease(stdInUrls, &jobs, token, tempdir, skipextraction)
		}
	}

	jobs.Wait()

	switch {
	case opts.Version:
		fmt.Println("getghrel version: 0.1.0")

	case opts.List:
		return

	case skipextraction:
		fmt.Println("Archives are inside: ", tempdir)

	default:
		cleanup(tempdir)
		fmt.Println("")
		fmt.Println("All Binaries are inside: ", tempdir)
	}
}
