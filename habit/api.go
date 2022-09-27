package habit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Headers = map[string]string

type API struct {
	baseURL string
	token   string
}

func NewAPI(token string) *API {
	return &API{
		baseURL: "https://gastro.fly.dev",
		token:   token,
	}
}

func (a API) List() (*http.Response, error) {
	return get(a.baseURL+"/habits", Headers{"Authorization": a.token})
}

func (a API) Create(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits?name=%s", a.baseURL, name)
	return post(url, Headers{"Authorization": a.token}, &bytes.Buffer{})
}

func (a API) Get(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits/%s", a.baseURL, name)
	return get(url, Headers{"Authorization": a.token})
}

func (a API) Delete(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits/%s", a.baseURL, name)
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", a.token)
	return http.DefaultClient.Do(req)
}

func (a API) AddActivity(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits/%s", a.baseURL, name)
	return post(url, Headers{"Authorization": a.token}, &bytes.Buffer{})
}

func get(url string, headers Headers) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, new(bytes.Buffer))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}

func post(url string, headers Headers, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}
