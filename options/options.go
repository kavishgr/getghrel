package options

import (
    "flag"
    "os"
    "strings"
    "fmt"
)

type options struct {
    List       bool
    Download    bool
    NoExtract bool
    Concurrency int
    GHToken     string
    TempDir     string
}

func ParseFlags() options {
    
flag.Usage = func() {
        h := []string{
            "Download releases from github and retain only the binaries",
            "",
            "Usage:",
            "  echo 'https://github.com/sharkdp/bat' | getghrel -list | sort",
            "  cat urls.txt | getghrel -list | sort | tee releases.txt",
            "  cat releases.txt | getghrel -download",
            "  cat releases.txt | getghrel -download -tempdir '/tmp/bin'",
            "  cat urls.txt | getghrel -list | grep -vi '^n/a'",
            " ",
            "If a repository doesn't have a latest release tag, the tool will print for the latest tag available.",
            "This can include nightly/unstable releases.\n",
            "Checksums, SBOMs, and files related to your OS and Arch will also be included.",
            "You can use `grep` to filter out these results.\n",
            "Repos that do not have a release section, \n",
            "The url format for -list: \n",
            "  A github url -> 'https://github.com/owner/repo'",
            "  For e.g 'https://github.com/sharkdp/bat'\n",
            "  Or only owner and repository -> 'owner/repo'",
            "  For e.g -> 'sharkdp/bat'\n",
            "For more examples, browse to: https://github.com/kavishgr/getghrel",
            "",
            "Options:",
            "  -list list all the releases found",
            "\tWill list all the release/releases found for your OS and Architecture.\n",
            "\tExample: cat urls.txt | getghrel -list | sort",
            "\tExample: echo 'https://github.com/sharkdp/bat' | getghrel -list | sort",
            "",
            "  -con <int> set the concurrency level (default: 2)\n",
            "\t Example: cat urls.txt | getghrel -list -con 3",
            "\t Example: cat releases.txt | getghrel -download -con 3",
            "",
            "  -ghtoken <string> specify you GITHUB TOKEN",
            "\t Default is the GIHUB_TOKEN environment variable.\n",
            "\t Example: cat urls.txt | getghrel -list -ghtoken 'YOUR TOKEN'",
            "",
            "  -download download the releases",
            "\tDefault directory in which the release will be downloaded is '/tmp/getghrel'",
            "\tIf the release is compressed or in an archive format, the tool will automatically extract and unpack it",
            "\tno matter how it's compressed, and keep only the binary.\n",
            "\tExample: cat urls_from_list_results.txt | getghrel -download",
            "\tExample: cat urls_from_list_results.txt | getghrel -download -tempdir '/tmp/test'",
            "",
            "  -tempdir <string> specify a temporary directory to download/extract the binaries\n",
            "\tExample: cat urls_from_list_results.txt | getghrel -download -tempdir '/tmp/test'",
            "",

        }

        // fmt.Fprint(os.Stderr, strings.Join(h, "\n"))
        fmt.Fprint(os.Stdout, strings.Join(h, "\n"))

    }

    opts := options{}
    flag.BoolVar(&opts.Download, "download", false, "")
    flag.BoolVar(&opts.List, "list", false, "")
    flag.BoolVar(&opts.NoExtract, "noextract", false, "")
    flag.IntVar(&opts.Concurrency, "con", 2, "")
    default_ghtoken := os.Getenv("GITHUB_TOKEN")
    flag.StringVar(&opts.GHToken, "ghtoken", default_ghtoken, "")
    flag.StringVar(&opts.TempDir, "tempdir", "/tmp/getghrel", "")

    flag.Parse()

    return opts
}
