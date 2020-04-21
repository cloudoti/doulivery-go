package doulivery

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Mailer struct {
	Client  *Client
	Service string
	From    string
	To      []string
	Subject string
	Body    string
	Html    bool
	Files   []File
}

type File struct {
	Name    string
	Content []byte
}

func CreateMailer(client *Client, service string) *Mailer {
	m := Mailer{
		Client:  client,
		Service: service,
	}

	return &m
}

func (m *Mailer) AddFrom(from string) *Mailer {
	m.From = from

	return m
}

func (m *Mailer) AddTo(to ...string) *Mailer {
	if m.To == nil {
		m.To = make([]string, 0)
	}

	m.To = append(m.To, to...)

	return m
}

func (m *Mailer) AddSubject(subject string) *Mailer {
	m.Subject = subject

	return m
}

func (m *Mailer) AddBody(body string) *Mailer {
	m.Body = body

	return m
}

func (m *Mailer) IsHtml(band bool) *Mailer {
	m.Html = band

	return m
}

func (m *Mailer) AddFile(file ...File) *Mailer {
	if m.Files == nil {
		m.Files = make([]File, 0)
	}

	m.Files = append(m.Files, file...)

	return m
}

func (m *Mailer) SendEmail() (err error) {

	path := fmt.Sprintf("/api/app/%s/email/send", m.Client.AppId)

	var b []byte
	if b, err = json.Marshal(m); err != nil {
		return
	}

	hasher := md5.New()
	hasher.Write(b)

	bodyMd5 := hex.EncodeToString(hasher.Sum(nil))

	mailerURL, err := createRequestURL(m.Client.Settings.Host, path, m.Client.Key, m.Client.Secret, m.Client.Settings.Secure, bodyMd5)
	//mailerURL, err := createRequestURL("localhost:8080", path, m.Client.Key, m.Client.Secret, false, bodyMd5)
	if err != nil {
		return err
	}

	_, err = m.Client.request("POST", mailerURL, m.Client.Secret, nil, m)

	return err
}
