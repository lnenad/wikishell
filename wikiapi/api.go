package wikiapi

import (

	//"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// cookiejar for persistent connection
type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

func FetchResponseBody(resp *http.Response) []byte {
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return rbody
}

//Struct for decoding json during login
type InnerLogin struct {
	Result, Token, Cookieprefix, Sessionid string
}

type OuterLogin struct {
	Login InnerLogin
}

type BaseToken struct {
	Tokens *Token
}

type Token struct {
	Edittoken  string
	Watchtoken string
}

type WikiAPI struct {
	//struct contains URL for making requests
	Url *url.URL
	/* All the request will be sent using client, using client makes
	   easier to maintain session
	*/
	jar    *Jar
	client *http.Client
	format string
	Tokens *Token
}

func NewWikiAPI(mUrl *url.URL) *WikiAPI {
	jar := new(Jar)
	client := &http.Client{Transport: nil, CheckRedirect: nil, Jar: jar, Timeout: time.Duration(time.Second * 60)}
	tokens := new(Token)
	return &WikiAPI{mUrl, jar, client, "json", tokens}
}

func (m WikiAPI) Get(params url.Values) *http.Response {
	params.Add("format", m.format)
	m.Url.RawQuery = params.Encode()
	resp, err := m.client.Get(m.Url.String())
	if err != nil {
		panic(err)
	}
	return resp
}
