package telegram

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	sendPhotoMethod   = "sendPhoto"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// Updates sends HTTPs request to Telegram API to get updates.
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.sendRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, errors.New("send request: " + err.Error())
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, errors.New("unmarshal json: " + err.Error())
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.sendRequest(sendMessageMethod, q)
	if err != nil {
		return errors.New("send request: " + err.Error())
	}

	return nil
}

func (c *Client) SendPhoto(chatID int, photoURL, caption string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", photoURL)
	q.Add("caption", caption)

	_, err := c.sendRequest(sendPhotoMethod, q)
	if err != nil {
		return errors.New("send request: " + err.Error())
	}

	return nil
}

func (c *Client) sendRequest(method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.New("new http request: " + err.Error())
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("failed to do http request: " + err.Error())
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body: " + err.Error())
	}

	return body, nil
}
