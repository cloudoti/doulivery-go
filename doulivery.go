package doulivery

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

const VERSION = 0.1

var SETTINGS = Settings{
	Timeout: time.Second * 10,
	Debug:   false,
	Secure:  true,
	Host:    "api.doulivery.io",
}

type Trigger struct {
	EventName string   `json:"eventName"`
	Data      string   `json:"data"`
	Channels  []string `json:"channels"`
	BodyMd5   string   `json:"-"`
}

type Client struct {
	Settings    Settings
	AppId       string
	Key         string
	Secret      string
	Environment string
	HTTPClient  *http.Client
}

type Settings struct {
	Debug   bool
	Secure  bool
	Timeout time.Duration
	Host    string
}

func CreateClient(appId, key, secret string) *Client {
	d := Client{
		Settings:    SETTINGS,
		AppId:       appId,
		Key:         key,
		Secret:      secret,
		Environment: "production",
		HTTPClient:  nil,
	}

	return &d
}

func (d *Client) Trigger(channels []string, event string, data interface{}, encoded bool) (err error) {
	if err = d.validateChannels(channels); err != nil {
		return
	}

	path := fmt.Sprintf("/api/app/%s/publish", d.AppId)

	if d.Environment != "production" {
		path = fmt.Sprintf("/api/app/%s/environment/%s/publish", d.AppId, d.Environment)
	}

	var b []byte
	if !encoded {
		if b, err = json.Marshal(data); err != nil {
			return
		}
	} else {
		b = []byte(data.(string))
	}

	postTrigger := Trigger{
		Channels:  channels,
		EventName: event,
		Data:      string(b),
	}

	if b, err = json.Marshal(postTrigger); err != nil {
		return
	}

	hasher := md5.New()
	hasher.Write(b)

	postTrigger.BodyMd5 = hex.EncodeToString(hasher.Sum(nil))

	triggerURL, err := createRequestURL(d.Settings.Host, path, d.Key, d.Secret, d.Settings.Secure, postTrigger.BodyMd5)
	if err != nil {
		return err
	}

	if b, err = json.Marshal(postTrigger); err != nil {
		return
	}

	_, err = d.request("POST", triggerURL, d.Secret, b, nil)

	return err
}

func (d *Client) requestDoulivery() *http.Client {
	if d.HTTPClient == nil {
		d.HTTPClient = &http.Client{Timeout: d.Settings.Timeout}
	}

	return d.HTTPClient
}

func (d *Client) request(method, url, secret string, body []byte, m *Mailer) ([]byte, error) {
	return request(d.requestDoulivery(), method, url, secret, body, m)
}

func (d *Client) validateChannels(channels []string) error {
	if len(channels) > 100 {
		return fmt.Errorf("an event can be triggered on a maximum of 100 channels in a single call")
	}

	for _, ch := range channels {
		if !regexp.MustCompile("\\A[-a-zA-Z0-9_=@,.;]+\\z").MatchString(ch) {
			return fmt.Errorf("Invalid channel name " + ch)
		}
	}

	return nil
}
