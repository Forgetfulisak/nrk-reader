package nrk

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var (
	titleColor      = color.New(color.FgWhite, color.Bold)
	smallTitleColor = color.New(color.FgWhite)
	linkColor       = color.New(color.FgBlue)
)

type Article struct {
	SmallTitle string `json:"smalltitle"`
	Title      string `json:"title"`
	PageLink   string `json:"pagelink"`
	LeadText   string `json:"leadtext"`
}

func (a *Article) Print() {
	linkColor.Println(a.PageLink)

	if a.SmallTitle != "" {
		smallTitleColor.Println(a.SmallTitle)
	}

	titleColor.Println(a.Title)

	if a.LeadText != "" {
		smallTitleColor.Println(a.LeadText)
	}
}

func (a *Article) Equal(other *Article) bool {
	return *a == *other 
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

func PrintNodeTree(root *html.Node, indent int) {
	escapedData := strings.ReplaceAll(root.Data, "\n", "\\n")

	fmt.Printf("%s%v, \"%v\" %v\n", strings.Repeat("\t", indent), nodeTypeStr(root.Type), escapedData, root.Attr)

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		PrintNodeTree(child, indent+1)
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

			if node.Type == html.ElementNode && strings.Contains(node.Data, "small") {
				article.SmallTitle = text
			} else {
				article.Title += text
			}
		}
		parseTitleText(child, article)
	}
}

func parseLeadText(node *html.Node, article *Article) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {

		if child.Type == html.TextNode {
			text := cleanData(child.Data)
			article.LeadText += text
		}
		parseLeadText(child, article)
	}
}

func parseArticle(root *html.Node, article *Article) bool {

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		for _, attr := range child.Attr {
			if attr.Key == "href" {
				article.PageLink = attr.Val
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
					SmallTitle: "",
					Title:      "",
					PageLink:   "",
				}

				parseArticle(child, &article)
				if article.Title != "" {
					*articles = append(*articles, article)
				}
			}
		}
		_parseArticles(child, articles)
	}
}

func ParseArticles(root *html.Node) []Article {
	out := make([]Article, 0)
	_parseArticles(root, &out)
	return out
}

func FetchArticles() ([]Article, error) {
	resp, err := http.Get("https://www.nrk.no")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	articles := ParseArticles(root)
	return articles, nil
}
