package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PlaceholderClient struct {
	client  *http.Client
	baseURL string
	timeout time.Duration
	apiKey  string
}

func NewClient(baseUrl, apiKey string, timeout time.Duration) *PlaceholderClient {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  true,
			MaxIdleConnsPerHost: 10, // istek atılan adres başına boşta bekleyebilecek connection sayısı
		},
	}
	return &PlaceholderClient{client: client, apiKey: apiKey, baseURL: baseUrl}
}

func (p PlaceholderClient) GetPostDetail(path, id string) (*Post, *ApiError) {
	resp, err := p.client.Get(fmt.Sprintf("%s/%s/%s", p.baseURL, path, id))
	if err != nil {
		return nil, &ApiError{
			StatusCode: resp.StatusCode,
			Message:    "An error occured while getting post detail",
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)

		var post *Post
		err = json.Unmarshal(body, &post)

		if err != nil {
			return nil, &ApiError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Cannot unmarshal response to Post struct",
			}
		}
		return post, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, &ApiError{
			StatusCode: resp.StatusCode,
			Message:    "Post not found",
		}
	}

	return nil, &ApiError{
		StatusCode: resp.StatusCode,
		Message:    "An error occured while getting post detail",
	}
}
