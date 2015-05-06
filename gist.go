package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	fileName string
	public   bool
	token    string
)

func init() {
	flag.StringVar(&fileName, "f", "", "gist file name")
	flag.BoolVar(&public, "p", false, "make gist public")
	flag.StringVar(&token, "t", os.Getenv("GITHUB_TOKEN"), "github token")
}

type TokenSource oauth2.Token

func (t *TokenSource) Token() (*oauth2.Token, error) {
	return (*oauth2.Token)(t), nil
}

type GistFiles map[github.GistFilename]github.GistFile

func getFilesFromStdin() (GistFiles, error) {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	content := string(data)
	return GistFiles{
		github.GistFilename(fileName): github.GistFile{
			Content: &content,
		},
	}, nil
}

func readFile(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getFilesFromArgs() (GistFiles, error) {
	files := make(GistFiles)
	for _, arg := range flag.Args() {
		content, err := readFile(arg)
		if err != nil {
			return nil, err
		}
		files[github.GistFilename(arg)] = github.GistFile{
			Content: &content,
		}
	}
	return files, nil
}

func getFiles() (GistFiles, error) {
	if flag.NArg() > 0 {
		return getFilesFromArgs()
	} else {
		return getFilesFromStdin()
	}
}

func main() {
	flag.Parse()
	files, err := getFiles()
	if err != nil {
		log.Fatal(err)
	}
	ts := TokenSource{AccessToken: token}
	client := github.NewClient(
		oauth2.NewClient(oauth2.NoContext, &ts),
	)
	gist, _, err := client.Gists.Create(&github.Gist{
		Files:  files,
		Public: &public,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*gist.HTMLURL)
}
