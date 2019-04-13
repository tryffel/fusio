package notifications

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tryffel/fusio/err"
	"net/http"
)

const urlFmt = "https://api.telegram.org/bot%s/sendMessage"

type Telegram struct {
	BotKey string `json:"bot_key"`
	ChatId string `json:"chat_id"`
}

func (t *Telegram) Notify(data string) error {
	body := bytes.NewBufferString(data)
	req, err := http.NewRequest("post", fmt.Sprintf(urlFmt, t.BotKey), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &Err.Error{Code: Err.Econflict, Err: errors.Wrap(err, "failed to create request for telegram message")}
	}

	if resp.StatusCode == 200 {
		return nil
	}
	return &Err.Error{Code: Err.Econflict, Err: errors.Wrapf(err, "failed to push telegram message, statuscode: %d", resp.StatusCode)}
}
