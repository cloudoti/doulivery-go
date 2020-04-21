package doulivery

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime/multipart"
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

func request(client *http.Client, method string, url, secret string, body []byte, m *Mailer) ([]byte, error) {

	var b *bytes.Buffer
	var writer *multipart.Writer

	if m == nil {
		b = bytes.NewBuffer(body)
	} else {
		b = &bytes.Buffer{}
		writer = multipart.NewWriter(b)

		_ = writer.WriteField("service", m.Service)
		_ = writer.WriteField("from", m.From)
		for _, to := range m.To {
			_ = writer.WriteField("to", to)
		}
		_ = writer.WriteField("subject", m.Subject)
		_ = writer.WriteField("is_html", strconv.FormatBool(m.Html))
		_ = writer.WriteField("body", m.Body)

		for _, fil := range m.Files {
			part, err := writer.CreateFormFile("files", fil.Name)
			if err != nil {
				panic(err)
			}
			part.Write(fil.Content)
		}

		err := writer.Close()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, b)
	req.Header.Set("Authorization", secret)

	if m == nil {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

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
