package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// PostForm post a request with form
func PostForm(url string, data url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
