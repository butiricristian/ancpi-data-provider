package main

import (
	"fmt"

	"golang.org/x/net/html"
)

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func checkId(n *html.Node, elem string, key string, attr string) bool {
	if n.Type == html.ElementNode && n.Data == elem {
		s, ok := getAttribute(n, key)
		if ok && s == attr {
			return true
		}
	}

	return false
}

func traverse(n *html.Node, elem string, key string, attr string) *html.Node {
	if checkId(n, elem, key, attr) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := traverse(c, elem, key, attr)
		if res != nil {
			return res
		}
	}

	return nil
}

// func renderNode(n *html.Node) string {

// 	var buf bytes.Buffer
// 	w := io.Writer(&buf)

// 	err := html.Render(w, n)

// 	if err != nil {
// 		return ""
// 	}

// 	return buf.String()
// }

func findByElemAttr(elem string, key string, attr string) string {
	html_resp := requestPage("https://www.ancpi.ro/statistica-decembrie-2022/")
	defer html_resp.Close()

	doc, err := html.Parse(html_resp)
	if err != nil {
		fmt.Printf("An error occurred while parsing the html page: %v", err)
	}

	node := traverse(doc, elem, key, attr)
	url, ok := getAttribute(node, "href")
	if ok {
		return url
	}

	return ""
}
