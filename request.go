package doulivery

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func unsignedParams(key string, md5Body string) url.Values {
	params := url.Values{
		"auth_key": {key},
	}

	if md5Body != "" {
		params.Add("body_md5", md5Body)
	}

	return params

}

func unescapeURL(_url url.Values) string {
	unesc, _ := url.QueryUnescape(_url.Encode())
	return unesc
}

func createRequestURL(host, path, key, secret string, secure bool, md5Body string) (string, error) {
	params := unsignedParams(key, md5Body)

	var base string
	if secure {
		base = "https://"
	} else {
		base = "http://"
	}
	base += host

	endpoint, err := url.ParseRequestURI(base + path)
	if err != nil {
		return "", err
	}
	endpoint.RawQuery = unescapeURL(params)

	return endpoint.String(), nil
}

func request(client *http.Client, method string, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return processResponse(resp)
}

func processResponse(response *http.Response) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}
	message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(response.StatusCode), string(responseBody))
	err = errors.New(message)
	return nil, err
}
