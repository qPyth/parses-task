package parsers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/qPyth/parses-task/internal/types"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	host    = "hypeauditor.com"
	attrKey = "data-v-b11c405a"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) ParseTopInstagram(category, country string) ([]types.Influencer, error) {

	data, err := doRequest(category, country)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("doc parsing error: %w", err)
	}

	attrs := make(map[string]string)
	attrs["class"] = "row"
	attrs["data-v-b11c405a"] = ""

	nodes, _ := p.findElementsByAttr(doc, attrs)

	var persons []types.Influencer

	for _, node := range nodes {
		person, err := p.parsePersonNode(node.FirstChild.FirstChild)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (p Parser) parsePersonNode(n *html.Node) (t types.Influencer, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v; current node value: %v\n", r, n)
			n = n.NextSibling
		}
	}()

	for n.NextSibling != nil {
		switch n.Attr[0].Val {
		case "row-cell rank":
			t.Rank, err = strconv.Atoi(getValueFromNode(n.FirstChild))
			if err != nil {
				return t, fmt.Errorf("parsing rank error: %w", err)
			}
		case "row-cell contributor":
			info := n.FirstChild.NextSibling.FirstChild
			t.Info.IGUsername = getValueFromNode(info.FirstChild.FirstChild)
			t.Info.Name = getValueFromNode(info.FirstChild.NextSibling)
		case "row-cell category":
			head := n
			attr := map[string]string{
				"class": "tag__content ellipsis",
			}
			catNodes, _ := p.findElementsByAttr(head, attr)
			for _, node := range catNodes {
				t.Category = append(t.Category, getValueFromNode(node))
			}
		case "row-cell subscribers":
			t.Followers = getValueFromNode(n)
		case "row-cell audience":
			t.Country = getValueFromNode(n)
		case "row-cell authentic":
			t.EngAuth = getValueFromNode(n)
		case "row-cell engagement":
			t.EngAvg = getValueFromNode(n)
		}
		n = n.NextSibling
	}
	switch n.Attr[0].Val {
	case "row-cell rank":

	}

	return t, err
}

func (p Parser) findElementsByAttr(n *html.Node, attrs map[string]string) ([]*html.Node, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v; current node value: %v\n", r, n)
			n = n.NextSibling
		}
	}()
	var nodes []*html.Node
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			matches := true
			for key, val := range attrs {
				attrVal, err := getAttr(n, key)
				if err != nil || attrVal != val {
					matches = false
					break
				}
			}
			if matches {
				nodes = append(nodes, n)
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(n)
	return nodes, nil
}

func getAttr(n *html.Node, key string) (string, error) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, nil
		}
	}
	return "", errors.New("not found")
}

func doRequest(category, country string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   fmt.Sprintf("top-instagram-%s-%s", category, country),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("request error:%w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 code from %s, code: %d", u.String(), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body reading error: %w", err)
	}

	return body, nil
}
func getValueFromNode(n *html.Node) string {
	return n.FirstChild.Data
}
