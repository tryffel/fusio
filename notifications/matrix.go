package notifications

import (
	"fmt"
	"github.com/tryffel/fusio/err"
	"time"
)

type Matrix struct {
	Host   string `json:"host"`
	RoomId string `json:"room_id"`
	Token  string `json:"token"`
}

// Notify pushes data to matrix chat as text message
func (m *Matrix) Notify(data string) error {
	txid := time.Now().Nanosecond() % 1000
	url := fmt.Sprintf("%s/_matrix/client/r0/rooms/%s/send/m.room.message/%d", m.Host, m.RoomId, txid)
	auth := fmt.Sprintf("Bearer %s", m.Token)

	hook := WebHook{
		Url:        url,
		Method:     "PUT",
		ExpectCode: 200,
	}
	hook.Headers.Authorization = auth
	hook.Headers.ContentType = "application/json"

	body := fmt.Sprintf(`{"msgtype": "m.text", "body": "%s"}`, data)
	err := hook.Notify(body)
	if err == nil {
		return nil
	}

	err.(*Err.Error).Wrap("Failed to push matrix message")
	return err
}
