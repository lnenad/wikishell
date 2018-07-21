package wikiapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Redirect struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type Query struct {
	Pages     map[string]Page `json:"pages,omitempty"`
	Redirects []Redirect      `json:"redirects,omitempty"`
}

type Page struct {
	PageID    int        `json:"pageid,omitempty"`
	NS        int        `json:"ns,omitempty"`
	Title     string     `json:"title,omitempty"`
	Revisions []Revision `json:"revisions,omitempty"`
}

type Revision struct {
	ContentFormat string  `json:"contentformat,omitempty"`
	ContentModel  string  `json:"contentmodel,omitempty"`
	Content       Content `json:"*"`
}

type WikiResponse struct {
	BatchComplete bool  `json:"batch_complete"`
	Query         Query `json:"query,omitempty"`
}

type ParsedWikiResponse struct {
	Parsed Parsed `json:"parse"`
}

type Section struct {
	Title      string `json:"line,omitempty"`
	ByteOffset int    `json:"byte_offset,omitempty"`
}

type Parsed struct {
	Title    string     `json:"title,omitempty"`
	Text     ParsedText `json:"text,omitempty"`
	Sections []Section  `json:"sections,omitempty"`
}

type ParsedText struct {
	Content Content `json:"*"`
}

func FetchLatestRevision(page string, wapi *WikiAPI) *Content {
	params := url.Values{
		"action":    {"query"},
		"prop":      {"revisions"},
		"titles":    {page},
		"format":    {"json"},
		"redirects": {""},
		"rvprop":    {"content"},
	}

	resp := wapi.Get(params)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		rbody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rbody))
		panic("Failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var jsonBody WikiResponse
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		panic(err)
	}

	for _, page := range jsonBody.Query.Pages {
		if len(page.Revisions) > 0 {
			return &page.Revisions[0].Content
		}
	}

	return nil
}

func FetchParsedText(page string, wapi *WikiAPI) *ParsedWikiResponse {
	params := url.Values{
		"action": {"parse"},
		"prop":   {"text|sections"},
		"page":   {page},
		"format": {"json"},
	}

	resp := wapi.Get(params)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		rbody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rbody))
		panic("Failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var jsonBody ParsedWikiResponse
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		panic(err)
	}

	return &jsonBody
}
