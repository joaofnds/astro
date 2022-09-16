package habit

import (
	"fmt"
	"net/http"
	"strings"
)

type API struct {
	baseURL string
}

func NewAPI() *API {
	return &API{"https://gastro.fly.dev"}
}

func (a API) List() (*http.Response, error) {
	return http.Get(a.baseURL + "/habits")
}

func (a API) Create(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits?name=%s", a.baseURL, name)
	return http.Post(url, "application/text", strings.NewReader(""))
}

func (a API) Get(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits/%s", a.baseURL, name)
	return http.Get(url)
}

func (a API) Delete(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits/%s", a.baseURL, name)
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func (a API) AddActivity(name string) (*http.Response, error) {
	url := fmt.Sprintf("%s/habits/%s", a.baseURL, name)
	return http.Post(url, "application/json", strings.NewReader(""))
}
