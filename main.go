package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Article struct {
	smallTitle string
	title      string
	pageLink   string
	leadText   string
}

var (
	titleColor      = color.New(color.FgWhite, color.Bold)
	smallTitleColor = color.New(color.FgWhite)
	linkColor       = color.New(color.FgBlue)
)

func (a *Article) Print() {
	linkColor.Println(a.pageLink)
	if a.smallTitle != "" {
		smallTitleColor.Println(a.smallTitle)

	}
	titleColor.Println(a.title)
	if a.leadText != "" {
		smallTitleColor.Println(a.leadText)
	}
}

func nodeTypeStr(t html.NodeType) string {
	switch t {
	case html.ErrorNode:
		return "ErrorNode"
	case html.TextNode:
		return "TextNode"
	case html.DocumentNode:
		return "DocumentNode"
	case html.ElementNode:
		return "ElementNode"
	case html.CommentNode:
		return "CommentNode"
	case html.DoctypeNode:
		return "DoctypeNode"
	case html.RawNode:
		return "RawNode"
	default:
		return "Undefined"
	}
}

func printNodeTree(root *html.Node, indent int) {
	escapedData := strings.ReplaceAll(root.Data, "\n", "\\n")

	fmt.Printf("%s%v, \"%v\" %v\n", strings.Repeat("\t", indent), nodeTypeStr(root.Type), escapedData, root.Attr)

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		printNodeTree(child, indent+1)
	}
}

func cleanData(data string) string {
	cleanedData := strings.TrimSpace(data)
	cleanedData = strings.ReplaceAll(cleanedData, "\n", "")
	cleanedData = strings.ReplaceAll(cleanedData, "\u2013", " \u2013 ")
	cleanedData = strings.ReplaceAll(cleanedData, "\u00ad", "")
	return cleanedData
}

func parseTitleText(node *html.Node, article *Article) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {

		if child.Type == html.TextNode {
			text := cleanData(child.Data)

			if text == "" {
				continue
			}

			if node.Type == html.ElementNode && strings.Contains(node.Data, "small") {
				article.smallTitle = text
			} else {
				article.title += text
			}
		}
		parseTitleText(child, article)
	}
}

func parseLeadText(node *html.Node, article *Article) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {

		if child.Type == html.TextNode {
			text := cleanData(child.Data)

			if text == "" {
				continue
			}

			article.leadText += text
		}
		parseLeadText(child, article)
	}
}

func parseArticle(root *html.Node, article *Article) bool {

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		for _, attr := range child.Attr {
			if attr.Key == "href" {
				article.pageLink = attr.Val
			}
			if attr.Key == "class" && strings.Contains(attr.Val, "kur-room__title") {
				parseTitleText(child, article)
			}
			if attr.Key == "class" && strings.Contains(attr.Val, "kur-room__leadtext") {
				parseLeadText(child, article)
			}

		}
		parseArticle(child, article)
	}
	return true
}

func _parseArticles(root *html.Node, articles *[]Article) {
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		for _, attr := range child.Attr {

			if attr.Key == "id" && strings.Contains(attr.Val, "kur-room-id") {
				article := Article{
					smallTitle: "",
					title:      "",
					pageLink:   "",
				}

				parseArticle(child, &article)
				if article.title != "" {
					*articles = append(*articles, article)
				}
			}
		}
		_parseArticles(child, articles)
	}
}

func parseArticles(root *html.Node) []Article {
	out := make([]Article, 0)
	_parseArticles(root, &out)
	return out
}

func main() {
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

	// printNodeTree(root, 0)
	articles := parseArticles(root)

	fmt.Printf("%v articles found\n\n", len(articles))
	for _, article := range articles {
		article.Print()
		fmt.Scanln()
	}
}
