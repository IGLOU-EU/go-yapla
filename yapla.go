// Copyright 2021 Iglou.eu
// license that can be found in the LICENSE file

package yapla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Api struct {
	key    string
	token  token
	Config Config
}

type Config struct {
	URL     string
	Timeout time.Duration
}

type token struct {
	Key    string
	Expire time.Time
}

type Reply struct {
	Result bool                   `json:"result"`
	Data   map[string]interface{} `json:"data"`
}

func (api *Api) LoginMember(login, password string) (Reply, error) {
	return api.login("/member/login", login, password)
}

func (api *Api) LoginContact(login, password string) (Reply, error) {
	return api.login("/contact/login", login, password)
}

func (api *Api) login(path, login, password string) (Reply, error) {
	var rep Reply

	if err := api.renewToken(); err != nil {
		return rep, err
	}

	rep, err := requestPost(
		"/authentication",
		api.token.Key,
		map[string]string{"login": login, "password": password},
		api.Config,
	)
	if err != nil {
		return rep, err
	}

	return rep, nil
}

func (api *Api) renewToken() error {
	if api.token.Expire.After(time.Now().Add(time.Minute * 2)) {
		return nil
	}

	rep, err := requestPost(
		"/authentication",
		api.token.Key,
		map[string]string{"api_key": api.key},
		api.Config,
	)
	if err != nil {
		return err
	}

	if rep.Data["session_token"] == nil || rep.Data["expire_date"] == nil {
		return fmt.Errorf("session token or expire date is missing:\n%v", rep)
	}

	api.token.Key = rep.Data["session_token"].(string)
	api.token.Expire, err = expireToTime(rep.Data["expire_date"].(string))
	if err != nil {
		return err
	}

	return nil
}

func NewSession(apiKey string, config ...Config) (*Api, error) {
	api := &Api{
		key: apiKey,
		token: token{
			Key:    "",
			Expire: time.Now(),
		},
		Config: Config{
			URL:     "https://s1.yapla.com/api/2",
			Timeout: time.Second * 10,
		},
	}

	if len(config) > 0 {
		api.Config = config[0]
	}

	err := api.renewToken()

	return api, err
}

func expireToTime(e string) (time.Time, error) {
	e = strings.TrimSpace(e)

	if t, err := time.Parse(time.RFC3339, e); err == nil {
		return t, nil
	}

	var s []string
	if strings.ContainsRune(e, 'T') {
		s = strings.Split(e, "T")
	} else {
		s = strings.Split(e, " ")
	}

	if len(s) != 2 {
		return time.Time{}, fmt.Errorf("json `expire_date` time format unexpected `%s`", e)
	}

	e = fmt.Sprintf("%sT%sZ", s[0], s[1])

	return time.Parse(time.RFC3339, e)
}

func requestPost(path, token string, content map[string]string, config Config) (Reply, error) {
	var rep Reply

	body, err := json.Marshal(content)
	if err != nil {
		return rep, fmt.Errorf("%s:\n%s", err, body)
	}

	res, err := request(
		http.MethodPost,
		fmt.Sprintf("%s%s", config.URL, path),
		token,
		body,
		config.Timeout,
	)
	if err != nil {
		return rep, err
	}

	if err := json.Unmarshal(res, &rep); err != nil {
		return rep, fmt.Errorf("%s:\n%s", err, res)
	}

	return rep, nil
}

func request(method, url, token string, body []byte, tOut time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: tOut,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("token", token)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("POST %s: %s", url, res.Status)
	}

	p, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return p, nil
}
