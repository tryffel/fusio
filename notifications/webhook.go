package notifications

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/err"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Headers struct {
	Authorization string `json:"Authorization"`
	ContentType   string `json:"Content-Type"`
}

type WebHook struct {
	Url        string `json:"url"`
	Headers    `json:"headers"`
	Method     string `json:"method"`
	ExpectCode uint   `json:"expect_status_code"`
}

func (h *WebHook) Notify(data string) error {
	body := bytes.NewBufferString(data)
	uri, err := url.Parse(h.Url)

	encoded := fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.EscapedPath())
	req, err := http.NewRequest(strings.ToUpper(h.Method), encoded, body)
	if err != nil {
		return &Err.Error{Code: Err.Econflict, Err: errors.Wrap(err, "failed to build request")}
	}

	if h.Headers.Authorization != "" {
		req.Header.Add("Content-Type", h.Headers.ContentType)
	}
	if h.Headers.ContentType != "" {
		req.Header.Add("Authorization", h.Headers.Authorization)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err == nil && uint(resp.StatusCode) == h.ExpectCode {
		return nil
	}

	if err != nil {
		return &Err.Error{Code: Err.Econflict, Err: errors.Wrapf(err, "webhook failed. Statuscode: %d,", resp.StatusCode)}
	}

	out, _ := ioutil.ReadAll(resp.Body)
	logrus.Error(string(out))
	return &Err.Error{Code: Err.Econflict, Err: errors.Wrapf(err, "webhook failed. Statuscode: %d, body:'%s'", resp.StatusCode, out)}
}
