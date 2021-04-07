package wordpress

import (
	"bytes"
	"github.com/divan/gorilla-xmlrpc/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpRT struct {
	http.RoundTripper
}

// RoundTrip implement a http.RoundTripper can print the request body
func (t HttpRT) RoundTrip(req *http.Request) (*http.Response, error) {
	// you can customize  to get more control over connection options
	// example: print the request body

	b, err := req.GetBody()
	if err != nil {
		log.Println(err)
	}
	r, err := ioutil.ReadAll(b)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(r))
	return t.RoundTripper.RoundTrip(req)
}

// Client Packaging the xmlrpc client
type Client struct {
	Url string
	UserInfo
}

// UserInfo wordpress's username and password
type UserInfo struct {
	Username string
	Password string
}

// NewClient without  http.RoundTripper
func NewClient(url string, info UserInfo) *Client {
	return &Client{Url: url, UserInfo: info}
}

// Call abstract to proxy xmlrpc call
func (c *Client) Call(method string, args interface{}) (reply struct{Message string}, err error) {
	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := http.Post(c.Url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = xml.DecodeClientResponse(resp.Body, &reply)
	return reply, err
}
