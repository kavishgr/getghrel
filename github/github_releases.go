package github

import (
	"context"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/kavishgr/getghrel/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/shurcooL/githubv4"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

/*
- does the url contains https://
  - returns true or false
  - got the regex from chatgpt
  - works fine
*/

func isValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

/*
- if githubUrl is valid -> https://github.com/sharkdp/bat

  - take only the username and repo -> e.g sharkdp/bat

  - prepend with apiDomain and append apiDomainSuffix

  - to make it a valid api url

  - if not valid (not a url) -> only username/repo

  - prepend/append the same thing and make it a valid api url

  - return the api url
*/

func fixUrl(githubUrl string) (string, string) {
	apiDomain := "https://api.github.com/repos"
	apiDomainSuffix := "/releases/latest"
	u, _ := url.Parse(githubUrl)
	fortag := fmt.Sprintf("%s", u.Path)

	if isValidURL(githubUrl) {
		result := fmt.Sprintf("%s%s%s", apiDomain, u.Path, apiDomainSuffix)
		return result, fortag
	}

	result := fmt.Sprintf("%s/%s%s", apiDomain, githubUrl, apiDomainSuffix)
	return result, fortag
}

/*
- takes a GITHUB_TOKEN and a valid api url
  - craft the request with "Authorization header: GITHUB_TOKEN"
  - return the request
*/
func craftGithubReq(ghtoken, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", ghtoken))
	// req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("User-Agent", "getghrel-cli")
	return req
}

/* - takes stdinurls channel from main
   - waitgroup for job
   - token to parse to craftGithubReq()
   - and tempdir to store downloaded binaries
   - download the asset with a download progressbar in the output
   - the asset is downloaded/extracted in the tempdir/hashdir directory
   - then find exectutable/binary and move it ../ to tempdir
*/

func DownloadRelease(urlsChan chan string, job *sync.WaitGroup, ghtoken, tempdir string, skipextraction bool) {

	defer job.Done()

	// anonymous func() to handle file download and processing
	// so that defer() gets called upon each iteration
	// instead of waiting for DownloadRelease() to return
	downloadAndProcessFile := func(u string) {
		// get the assetname of each url -> e.g bat.tar.gz
		file := path.Base(u)
		src := filepath.Join(tempdir, file)

		req := craftGithubReq(ghtoken, u)
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		f, _ := os.OpenFile(src, os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()

		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			file,
		)

		io.Copy(io.MultiWriter(f, bar), resp.Body)

		if skipextraction {
			return
		}

		utils.Extractor(src, tempdir)
	}

	// iterate over urls sent by stdin
	for u := range urlsChan {
		downloadAndProcessFile(u)
	}
}

func split(ownerNrepo string) (string, string) {
	var str string
	if strings.HasPrefix(ownerNrepo, "/") {
		str = strings.TrimPrefix(ownerNrepo, "/")
	}
	parts := strings.Split(str, "/")
	return parts[0], parts[1]
}

func getTagByName(ghtoken, ownerNrepo string) []byte {
	// GitHub API token

	var tagname string
	owner, name := split(ownerNrepo) // owner and name of the repo

	// Create an OAuth2 token source
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghtoken},
	)

	// Create an HTTP client with the token source
	oauthClient := oauth2.NewClient(context.Background(), src)

	// Create a new GitHub GraphQL client
	gqlClient := githubv4.NewClient(oauthClient)

	// Define the GraphQL query
	var query struct {
		Repository struct {
			Refs struct {
				Edges []struct {
					Node struct {
						Name string
					}
				}
			} `graphql:"refs(refPrefix: $refPrefix, first: $first, orderBy: $orderBy)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	// Set the query variables
	variables := map[string]interface{}{
		"owner":     githubv4.String(owner),
		"name":      githubv4.String(name),
		"refPrefix": githubv4.String("refs/tags/"),
		"first":     githubv4.Int(1),
		"orderBy": githubv4.RefOrder{
			Field:     githubv4.RefOrderFieldTagCommitDate,
			Direction: githubv4.OrderDirectionDesc,
		},
	}

	// Execute the GraphQL query
	err := gqlClient.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	// Access the query result
	for _, edge := range query.Repository.Refs.Edges {
		// fmt.Println("Tag Name:", edge.Node.Name)
		tagname = edge.Node.Name
		// fmt.Println(tagname)
	}

	tagUrl := fmt.Sprintf("https://api.github.com/repos%s/releases/tags/%s", ownerNrepo, tagname)
	// fmt.Println(tagUrl)
	// fmt.Println("TAGURL:", tagUrl)


	req := craftGithubReq(ghtoken, tagUrl)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func mapIsEmpty(m map[string]int) bool {
	return len(m) == 0 // returns true if map is empty
}

/* - fetch the asset urls from latest release for each url or username/repo (used by -list)
   - the regex is used to find the required asset url for your os/arch
   - return/print found urls for each asset that matched
   - repos that do not have a release for the os and arch will be printed like so:
   - N/A: https://github.com/user/repo
*/

//TODEL rescue
func FetchGithubReleaseUrl(urlsChan chan string, job *sync.WaitGroup, regex, ghtoken string) {

	defer job.Done()
	// var github_release []string
	// github_release := make(map[string]int)

	fetch := func(u string) {
		github_release := make(map[string]int)
		// fmt.Println(u)
		// map to keep assets
		// sometimes there are multiple assets for same os/architecture
		// for e.g gnu and musl for linux
		re2 := regexp2.MustCompile(regex, 0) // regex for os/arch
		// fmt.Println("Regex: ", re)
		githubUrl, ownerNrepo := fixUrl(u) // fix url and return valid api url
		// fmt.Println(githubUrl)

		// HTTP client starts
		req := craftGithubReq(ghtoken, githubUrl) // craft request with token and valid api url
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		// HTTP client ends

		message := gjson.Get(fmt.Sprintf("%s", body), "message")
		// if the message is "Not Found"
		// release/asset section is EMPTY or is using tags instead of latest release
		if message.Str == "Not Found" {
			// fetch assets for most recent tag
			body = getTagByName(ghtoken, ownerNrepo)
		}

		// fetch all the browser_download_url keys which contains the asset urls
		// results := gjson.Get(fmt.Sprintf("%s", body), "assets.#.browser_download_url")
		results := gjson.Get(fmt.Sprintf("%s", body), "assets.#.browser_download_url")

		// fmt.Println("Matching values for URL:", u)
		results.ForEach(func(key, value gjson.Result) bool {
			// fmt.Println(key.String())
			asset_url := value.String()
			// fmt.Println(asset_url)
			isMatch, _ := re2.MatchString(asset_url)

			if isMatch == true {
				github_release[asset_url] = 1
			}

			return true // keep iterating in case there are multiple urls that match
		})

		if mapIsEmpty(github_release) {
			fmt.Println("N/A:", u)
		} else {
			for k, _ := range github_release {
				// fmt.Println("URL found:", k)
				fmt.Println(k)
			}
		}

	}

	for u := range urlsChan {
		fetch(u)
	}
}
