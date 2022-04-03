package client

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
func TestPlaceholderClient_GetPostDetail(t *testing.T) {
	client := NewClient("https://jsonplaceholder.typicode.com", "", time.Second*10)
	post, err := client.GetPostDetail("posts", "1")

	assert.Equal(t, (*ApiError)(nil), err, "Error should be nil")
	assert.Equal(t, 1, post.Id, "PostId should be 1")
}

func TestPlaceholderClient_GetPostDetail_NotFound(t *testing.T) {
	client := NewClient("https://jsonplaceholder.typicode.com", "", time.Second*10)

	post, err := client.GetPostDetail("posts", "1123123")

	assert.Equal(t, (*Post)(nil), post, "Post should be nil")
	assert.Equal(t, http.StatusNotFound, err.StatusCode, "Repsonse status code should be 404")
	assert.Equal(t, "Post not found", err.Message, "Error message should be : 'Post not found' ")
}

func TestPlaceholderClient_GetPostDetail_Error(t *testing.T) {
	client := NewClient("https://restcountries.com/v3.1", "", time.Second*10)
	_, err := client.GetPostDetail("alpha", "pe")

	assert.Equal(t, http.StatusInternalServerError, err.StatusCode, "Repsonse status code should be 404")
	assert.Equal(t, "Cannot unmarshal response to Post struct", err.Message, "Error message should be : 'Post not found' ")
}
*/

func TestPlaceholderClient_GetAllPosts(t *testing.T) {
	resp := `{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}`
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/posts/1", req.URL.String())
		// Send response to be tested
		rw.Write([]byte(resp))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := PlaceholderClient{client: server.Client(), baseURL: server.URL, apiKey: "dfsfsfd", timeout: time.Second * 10}
	post, err := api.GetPostDetail("posts", "1")

	assert.Equal(t, ((*ApiError)(nil)), err, "Error should be nil")
	assert.Equal(t, 1, post.Id, "Post Id wrong")
	assert.Equal(t, 1, post.UserId, "User Id wrong")
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestPlaceholderClient_GetAllPostsWithRoundTrip(t *testing.T) {
	resp := `{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}`
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "http://jsonplaceholder.typicode.com/posts/1")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(resp)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	api := PlaceholderClient{client: client, baseURL: "http://jsonplaceholder.typicode.com", timeout: time.Second * 10}
	post, err := api.GetPostDetail("posts", "1")
	assert.Equal(t, 1, post.Id, "Post Id wrong")
	assert.Equal(t, 1, post.UserId, "User Id wrong")

	assert.Equal(t, ((*ApiError)(nil)), err, "Error should be nil")

}
