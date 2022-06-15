package main

import (
	"flag"
	"fmt"
	"net/http"
	"nrk-reader"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var (
	debug = flag.Bool("d", false, "prints the parsed html from nrk")
)

func main() {
	flag.Parse()

	fmt.Println("fetching...")
	resp, err := http.Get("https://www.nrk.no")
	if err != nil {
		logrus.Error("Could not fetch nrk.no:", err)
		return
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		logrus.Error("error parsing response:", err)
		return
	}

	if *debug {
		nrk.PrintNodeTree(root, 0)
		return
	}

	articles := nrk.ParseArticles(root)

	fmt.Printf("%v articles found\n\n", len(articles))
	fmt.Println("press enter to read the next article")
	for _, article := range articles {
		fmt.Scanln()
		article.Print()
	}
}
