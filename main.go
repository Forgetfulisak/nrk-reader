package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type nrkArticle struct {
	smallTitle string
	title      string
	pageLink   string
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

func parseTitleNode(node *html.Node, sb *strings.Builder) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {

		if child.Type == html.TextNode {
			cleanedData := strings.TrimSpace(child.Data)
			cleanedData = strings.ReplaceAll(cleanedData, "\n", "")
			cleanedData = strings.ReplaceAll(cleanedData, "\u2013", " \u2013 ")
			cleanedData = strings.ReplaceAll(cleanedData, "\u00ad", "")
			sb.WriteString(cleanedData)

		}
		parseTitleNode(child, sb)
	}
	if node.Type == html.ElementNode && strings.Contains(node.Data, "small") {
		sb.WriteString("\n")
	}

}

func _getTitleNodes(root *html.Node, titles *[]string) {

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		for _, attr := range child.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "kur-room__title") {
				sb := strings.Builder{}
				parseTitleNode(child, &sb)
				*titles = append(*titles, sb.String())
			}
		}
		_getTitleNodes(child, titles)
	}
}
func getTitleNodes(root *html.Node) []string {
	out := make([]string, 0)
	_getTitleNodes(root, &out)
	return out
}

func main() {
	// logrus.Info("hei!")
	resp, err := http.Get("https://nrk.no")
	if err != nil {
		logrus.Error("Could not get nrk.no. err: ", err)
	}
	defer resp.Body.Close()

	// logrus.Info("parsing...")

	root, err := html.Parse(resp.Body)
	if err != nil {
		logrus.Error("htmlparse: ", err)
	}

	whiteBold := color.New(color.FgWhite, color.Bold)
	// printNodeTree(root, 0)
	titles := getTitleNodes(root)
	for _, title := range titles {
		// fmt.Println(title)
		whiteBold.Println(title + "\n")
		fmt.Scanln()

	}
}
