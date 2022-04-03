package client

/*
import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
		assert.Equal(t, req.URL.String(), "/posts")
		// Send response to be tested
		rw.Write([]byte(resp))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := PlaceholderClient{client: server.Client(), baseURL: server.URL, apiKey: "dfsfsfd", timeout: time.Second * 10}
	post, err := api.GetPostDetail("posts", "1")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, post.Id, 1, "Post Id wrong")
	assert.Equal(t, post.UserId, 1, "User Id wrong")
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
	assert.Equal(t, post.Id, 1, "Post Id wrong")
	assert.Equal(t, post.UserId, 1, "User Id wrong")

	assert.Equal(t, err, nil, "Error should be nil")

}
*/
