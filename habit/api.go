package habit

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type Headers = map[string]string

type API struct {
	baseURL string
}

func NewAPI() *API {
	return &API{
		baseURL: "https://gastro.fly.dev",
	}
}

func (a API) List(token string) (*http.Response, error) {
	return get(a.baseURL+"/habits", map[string]string{"Authorization": token})
}

func (a API) Create(token, name string) (*http.Response, error) {
	return post(
		a.baseURL+"/habits?name="+name,
		Headers{"Content-Type": "application/json", "Authorization": token},
		&bytes.Buffer{},
	)
}

func (a API) Get(token, id string) (*http.Response, error) {
	return get(
		a.baseURL+"/habits/"+id,
		map[string]string{"Authorization": token},
	)
}

func (a API) Delete(token, id string) (*http.Response, error) {
	return delete(
		a.baseURL+"/habits/"+id,
		map[string]string{"Authorization": token},
	)
}

func (a API) AddActivity(token, id string) (*http.Response, error) {
	return post(
		a.baseURL+"/habits/"+id,
		map[string]string{"Authorization": token},
		&bytes.Buffer{},
	)
}

func (a API) CreateToken() (*http.Response, error) {
	return http.Post(a.baseURL+"/token", "application/text", strings.NewReader(""))
}

func (a API) TestToken(token []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, a.baseURL+"/tokentest", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(token))
	return http.DefaultClient.Do(req)
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

func delete(url string, headers Headers) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}
