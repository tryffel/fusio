package dtos

import (
	"encoding/json"
	"github.com/thedevsaddam/govalidator"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
	"time"
)

const (
	NotificationDefaultTemplateFire  = "Alarm '{{.AlarmName}}' fired, value: {{.Value}}"
	NotificationDefaultTemplateClear = "Alarm '{{.AlarmName}}' cleared, value: {{.Value}}"
	NotificationDefaultTemplateError = "Alarm '{{.AlarmName}}' has error, {{.Error}} "
)

type NewOutput struct {
	Name          string        `json:"name"`
	AlarmId       string        `json:"alarm_id"`
	OutputChannel string        `json:"output_channel_id"`
	FireTemplate  string        `json:"template_fire"`
	ClearTemplate string        `json:"template_clear"`
	ErrorTemplate string        `json:"template_error"`
	Repeat        util.Interval `json:"repeat"`
}

func (o *NewOutput) ValidationMap() *govalidator.MapData {
	// TODO: add custom template validator
	return &govalidator.MapData{
		"name":              []string{},
		"alarm_id":          []string{"required", "uuid"},
		"output_channel_id": []string{"required", "uuid"},
		"template_fire":     []string{""},
		"template_clear":    []string{""},
		"template_error":    []string{""},
		"repeat":            []string{},
	}
}

func (o *NewOutput) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"name":              []string{"Optional name for output"},
		"alarm_id":          []string{"Alarm that fires output"},
		"output_channel_id": []string{"Output channel id. This is the channel that output is sent to."},
		"template_fire":     []string{"Template'd string to send when fired with filled data. See docs."},
		"template_clear":    []string{"Template'd string to send when cleared with filled data. See docs."},
		"template_error":    []string{"Template'd string to send on error with filled data. See docs."},
		"repeat":            []string{"Repeat interval. Can be e.g. '1s', '1m', '1d'. Leave empty to disable."},
	}
}

//ToOutput return Dto as models.Output. If any of templates is empty, it is filled with
// default template
func (o *NewOutput) ToOutput() *models.Output {
	output := &models.Output{
		Name:            o.Name,
		AlarmId:         o.AlarmId,
		OutputChannelId: o.OutputChannel,
		FireTemplate:    o.FireTemplate,
		ClearTemplate:   o.ClearTemplate,
		ErrorTemplate:   o.ErrorTemplate,
		Repeat:          time.Duration(o.Repeat),
	}

	if output.FireTemplate == "" {
		output.FireTemplate = NotificationDefaultTemplateFire
	}
	if output.ClearTemplate == "" {
		output.ClearTemplate = NotificationDefaultTemplateClear
	}
	if output.ErrorTemplate == "" {
		output.ErrorTemplate = NotificationDefaultTemplateError
	}
	return output
}

type NewOutputChannel struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (o *NewOutputChannel) ValidationMap() *govalidator.MapData {
	// TODO: add custom validator for data
	return &govalidator.MapData{
		"name": []string{},
		"type": []string{"required"},
		"data": []string{"required"},
	}
}

func (o *NewOutputChannel) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"name": []string{"Optional name for output channel"},
		"type": []string{"Output type. See docs for valid types."},
		"data": []string{"Custom data depending on output type. See docs."},
	}
}

func (o *NewOutputChannel) ToOutputChannel() (*models.OutputChannel, error) {

	data, err := json.Marshal(o.Data)
	if err != nil {
		return &models.OutputChannel{}, err
	}

	str := string(data)

	return &models.OutputChannel{
		Name:       o.Name,
		Data:       str,
		OutputType: o.Type,
	}, nil
}

type WebHook struct {
	Url     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
	Method  string                 `json:"method"`
}

func (h *WebHook) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"url": []string{"required", "url"},
		// TODO: validate headers
		"headers": []string{""},
		"method":  []string{"required", "in:GET,POST,PUT,DELETE"},
	}
}

func (h *WebHook) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"url":     []string{"Valid url to call. Required."},
		"headers": []string{"Key-value headers. Optional."},
		"method":  []string{"HTTP-method, either: get|post|put|delete. Required."},
	}
}

type OutputChannel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Matrix struct {
	Host   string `json:"host"`
	RoomId string `json:"room_id"`
	Token  string `json:"token"`
}

func (m *Matrix) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"url":     []string{"required", "url"},
		"room_id": []string{"required", `regex:^(![a-zA-Z0-9]+:[\w\-\_\.]+)$`},
		"token":   []string{"required"},
	}
}

func (m *Matrix) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"url":     []string{"Server host with protocol, e.g. https://matrix.org"},
		"room_id": []string{"Valid room id, in the form !<id>:<host>"},
		"token":   []string{"User/bot token to send message from"},
	}
}
