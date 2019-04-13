package notifications

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebHook_NotifyPost(t *testing.T) {

	data := `{"data": "test"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		n, err := r.Body.Read(body)
		if err != nil && err.Error() != "EOF" {
			t.Error(err)
		}
		if n != len(data) {
			t.Error("http body length doesn't match")
		}

		if string(body) != data {
			t.Errorf("Http body doesn't match expected: got '%s', expected '%s'", string(body), data)
		}

		if r.Header.Get("Authorization") != "Bearer abcd1" {
			t.Error("Http authorization header doesn't match")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Http content-type header doesn't match")
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	hook := WebHook{
		Url:        server.URL,
		Method:     "post",
		ExpectCode: http.StatusOK,
	}
	hook.Headers.Authorization = "Bearer abcd1"
	hook.Headers.ContentType = "application/json"

	err := hook.Notify(data)
	if err != nil {
		t.Error(err)
	}
}
