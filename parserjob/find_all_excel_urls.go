package parserjob

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

type ExcelUrl struct {
	url   string
	name  string
	month string
	year  string
}

func compareAttributeValues(val1 string, val2 string) bool {
	val1 = strings.ToLower(val1)
	val1 = strings.ReplaceAll(val1, "_", " ")

	val2 = strings.ToLower(val2)
	val2 = strings.ReplaceAll(val2, "_", " ")

	return val1 == val2
}

func getAttributeValue(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func attrMatches(n *html.Node, elem string, key string, attr string) bool {
	if n.Type == html.ElementNode && n.Data == elem {
		s, ok := getAttributeValue(n, key)
		if ok {
			return compareAttributeValues(s, attr)
		}
	}

	return false
}

func traverse(n *html.Node, elem string, key string, attr string, collector []*html.Node) []*html.Node {
	if attrMatches(n, elem, key, attr) {
		collector = append(collector, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collector = traverse(c, elem, key, attr, collector)
	}

	return collector
}

func formatExcelUrls(validNodes []*html.Node) []*ExcelUrl {
	var urls []*ExcelUrl

	for _, node := range validNodes {
		url, urlOk := getAttributeValue(node, "href")
		title, titleOk := getAttributeValue(node, "title")
		if urlOk && titleOk {
			title, err := getDocumentName(title)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
			urls = append(urls, &ExcelUrl{url: url, name: title})
		}
	}

	return urls
}

func findAllNodesHavingAttribute(html_resp io.Reader, elem string, key string, attr string) []*html.Node {
	doc, err := html.Parse(html_resp)
	if err != nil {
		fmt.Printf("An error occurred while parsing the html page: %v", err)
		return make([]*html.Node, 0)
	}

	return traverse(doc, elem, key, attr, make([]*html.Node, 0))
}

func getMonths() []string {
	return []string{"ianuarie", "februarie", "martie", "aprilie", "mai", "iunie", "iulie",
		"august", "septembrie", "octombrie", "noiembrie", "decembrie"}
	// return []string{"noiembrie", "decembrie"}
}

func getYears() []string {
	return []string{"2022", "2021", "2020", "2019", "2018", "2017"}
	// return []string{"2022"}
}

func getDocumentName(name string) (string, error) {
	name = strings.ToLower(name)
	if strings.Contains(name, "cer") || strings.Contains(name, "cereri") {
		return "CERERI", nil
	}
	if strings.Contains(name, "vanzari") || strings.Contains(name, "tranzactii") || strings.Contains(name, "vânzări") {
		return "VANZARI", nil
	}
	if strings.Contains(name, "ipoteci") {
		return "IPOTECI", nil
	}

	return name, fmt.Errorf("could not find a match for %s", name)
}

func getUrlsFromLink(url string) ([]*ExcelUrl, bool) {
	// fmt.Printf("Requesting page %s...\n", url)

	html_resp, ok := requestPage(url)
	if !ok {
		return make([]*ExcelUrl, 0), false
	}
	defer html_resp.Close()

	validNodes := findAllNodesHavingAttribute(html_resp, "a", "class", "attachment-link")
	if len(validNodes) <= 0 {
		// fmt.Println("Node could not be found")
		return make([]*ExcelUrl, 0), false
	}

	urls := formatExcelUrls(validNodes)
	if len(urls) <= 0 {
		return make([]*ExcelUrl, 0), false
	}

	return urls, true
}

func addMonthAndYearToUrls(urls []*ExcelUrl, month string, year string) {
	for _, excelUrl := range urls {
		excelUrl.month = month
		excelUrl.year = year
	}
}

func sendUrls(urls []*ExcelUrl, urlsChannel chan<- *ExcelUrl, month string, year string) {
	addMonthAndYearToUrls(urls, month, year)
	for _, url := range urls {
		urlsChannel <- url
	}
}

func getUrls(month string, year string, urlsChannel chan<- *ExcelUrl, wg *sync.WaitGroup) {
	defer wg.Done()
	urls, ok := getUrlsFromLink(fmt.Sprintf("https://www.ancpi.ro/statistica-%s-%s/", month, year))
	if ok {
		sendUrls(urls, urlsChannel, month, year)
		return
	}

	// fmt.Println("Retrying with second URL version")
	urls, ok = getUrlsFromLink(fmt.Sprintf("https://www.ancpi.ro/statistica-%s-%s/", month[0:3], year))
	if ok {
		sendUrls(urls, urlsChannel, month, year)
		return
	}

	// fmt.Println("Retrying with third URL version")
	urls, ok = getUrlsFromLink(fmt.Sprintf("https://www.ancpi.ro/statistici-%s-%s/", month, year))
	if ok {
		sendUrls(urls, urlsChannel, month, year)
		return
	}
}

func FindAllExcelUrls() []*ExcelUrl {
	var excelUrls []*ExcelUrl

	urlsChannel := make(chan *ExcelUrl)
	var wg sync.WaitGroup

	for _, year := range getYears() {
		for _, month := range getMonths() {
			wg.Add(1)
			go getUrls(month, year, urlsChannel, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(urlsChannel)
	}()

	totalMonths := len(getYears()) * len(getMonths())
	bar := progressbar.NewOptions(
		totalMonths,
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetDescription("[1/2] Retrieving excel URLS: "),
	)
	for url := range urlsChannel {
		excelUrls = append(excelUrls, url)
		bar.Add(1)
	}
	fmt.Println()

	return excelUrls
}
