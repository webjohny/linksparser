package wordpress

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Client Packaging the xmlrpc client
type Client struct {
	Url *url.URL
	Credentials
}

// UserInfo wordpress's username and password
type Credentials struct {
	Username string
	Password string
}

type Response struct {
	Body string
	Status string
	StatusCode int
}

var (
	baseURL = &url.URL{Host: "testwp.com", Scheme: "http", Path: "/wp-json/wp/v2/"}
)

// NewClient without  http.RoundTripper
func NewClient(url *url.URL, info Credentials) (*Client, error) {
	return &Client{Url: url, Credentials: info}, nil
}

// Call abstract to proxy xmlrpc call
func (c *Client) Get(schema string, params *interface{}) (result *Response, err error) {
	var resp *http.Response
	var path string
	path += schema

	if !isNil(params) {
		v, _ := query.Values(params)
		if v != nil {
			path += "?" + v.Encode()
		}
	}

	// Create a new request using http
	req, err := http.NewRequest("GET", baseURL.String() + path, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", "Basic cHJvZ2VyOnF3ZXJ0eTEyMzQ1")
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	return &Response{
		Status: resp.Status,
		StatusCode: resp.StatusCode,
		Body: string(body),
	}, err
}

// Call abstract to proxy xmlrpc call
func (c *Client) Post(schema string, params interface{}) (result *Response, err error) {
	var resp *http.Response
	var path string
	path += schema

	var requestBody []byte
	if !isNil(params) {
		requestBody, _ = json.Marshal(params)
	}
	// Create a new request using http
	req, err := http.NewRequest("POST", baseURL.String() + path, bytes.NewBuffer(requestBody))

	// add authorization header to the req
	req.Header.Add("Authorization", "Basic cHJvZ2VyOnF3ZXJ0eTEyMzQ1")
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	return &Response{
		Status: resp.Status,
		StatusCode: resp.StatusCode,
		Body: string(body),
	}, err
}
