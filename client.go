package main

import (
	"net/http"
	"time"

	"github.com/ddliu/go-httpclient"
)

const sessionKey = "YTSESSIONID"
const principalKey = "jetbrains.charisma.main.security.PRINCIPAL"

type YouTrackClient struct {
	BaseUrl   string
	Session   string
	Principal string
	Expires   time.Time
}

func NewYouTrackClient(host, login, password string) YouTrackClient {
	client := YouTrackClient{BaseUrl: host + "/rest/"}

	client.login(login, password)

	return client
}

func (c *YouTrackClient) GetIssue(id string) (string, error) {
	res, err := httpclient.WithCookie(&http.Cookie{
		Name:  sessionKey,
		Value: c.Session,
	}).WithCookie(&http.Cookie{
		Name:  principalKey,
		Value: c.Principal,
	}).Get(c.BaseUrl+"issue/"+id, nil)

	if err != nil {
		return "", err
	}

	c.setCredsFromCookies(res.Cookies())

	body, err := res.ToString()
	if err != nil {
		return "", err
	}

	return body, nil
}

func (client *YouTrackClient) CommandIssue(id, command, comment string) (string, error) {
	url := client.BaseUrl + "issue/" + id + "/execute"

	var params = make(map[string]string)

	params["command"] = command
	if comment != "" {
		params["comment"] = comment
	}

	res, err := httpclient.WithCookie(&http.Cookie{
		Name:  sessionKey,
		Value: client.Session,
	}).WithCookie(&http.Cookie{
		Name:  principalKey,
		Value: client.Principal,
	}).Post(url, params)

	if err != nil {
		return "", err
	}

	client.setCredsFromCookies(res.Cookies())

	body, err := res.ToString()
	if err != nil {
		return "", err
	}

	return body, nil

}

func (client *YouTrackClient) login(login, password string) error {
	res, err := httpclient.Post(client.BaseUrl+"user/login", map[string]string{
		"login":    login,
		"password": password,
	})

	if err != nil {
		return err
	}

	client.setCredsFromCookies(res.Cookies())

	return nil
}

func (client *YouTrackClient) setCredsFromCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		switch cookie.Name {
		case sessionKey:
			client.Session = cookie.Value
		case principalKey:
			client.Principal = cookie.Value
			client.Expires = cookie.Expires
		}
	}
}
