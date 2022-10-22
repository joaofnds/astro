package habit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Headers = map[string]string

type API struct {
	baseURL string
}

func NewAPI() *API {
	return &API{
		baseURL: "https://astro.joaofnds.com",
	}
}

func (a API) List(token string) (*http.Response, error) {
	return req(
		http.MethodGet,
		a.baseURL+"/habits",
		map[string]string{"Authorization": token},
		new(bytes.Buffer),
	)
}

func (a API) Create(token, name string) (*http.Response, error) {
	return req(
		http.MethodPost,
		a.baseURL+"/habits?name="+url.QueryEscape(name),
		Headers{"Content-Type": "application/json", "Authorization": token},
		new(bytes.Buffer),
	)
}

func (a API) Get(token, id string) (*http.Response, error) {
	return req(
		http.MethodGet,
		a.baseURL+"/habits/"+id,
		map[string]string{"Authorization": token},
		new(bytes.Buffer),
	)
}

func (a API) Delete(token, id string) (*http.Response, error) {
	return req(
		http.MethodDelete,
		a.baseURL+"/habits/"+id,
		map[string]string{"Authorization": token},
		new(bytes.Buffer),
	)
}

func (a API) AddActivity(token, id, desc string) (*http.Response, error) {
	return req(
		http.MethodPost,
		a.baseURL+"/habits/"+id,
		map[string]string{"Authorization": token, "Content-Type": "application/json"},
		strings.NewReader(fmt.Sprintf(`{"description":%q}`, desc)),
	)
}

func (a API) DeleteActivity(token, habitID, activityID string) (*http.Response, error) {
	return req(
		http.MethodDelete,
		a.baseURL+"/habits/"+habitID+"/"+activityID,
		map[string]string{"Authorization": token},
		new(bytes.Buffer),
	)
}

func (a API) CreateToken() (*http.Response, error) {
	return req(
		http.MethodPost,
		a.baseURL+"/token",
		map[string]string{"Content-Type": "application/text"},
		new(bytes.Buffer),
	)
}

func (a API) TestToken(token string) (*http.Response, error) {
	return req(
		http.MethodGet,
		a.baseURL+"/tokentest",
		map[string]string{"Authorization": token},
		new(bytes.Buffer),
	)
}

func req(method string, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}
