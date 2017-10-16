package mqhttp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	gorequest "github.com/parnurzeal/gorequest"
	uuid "github.com/satori/go.uuid"
	aliot "github.com/tosone/aliot"
)

const httpAuthBaseURL = "https://iot-as-http.cn-shanghai.aliyuncs.com"

// Conn connect object
type Conn interface {
	Pub()
}

type conn struct {
	token string
}

type httpResp struct {
	Code int `json:"code"`
	Info struct {
		Token string `json:"token"`
	} `json:"info"`
	Message string `json:"message"`
}

type mqttParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	Sign       string `json:"sign"`
	ClientID   string `json:"clientId"`
	Signmethod string `json:"signmethod"`
	Timestamp  string `json:"timestamp"`
	Resources  string `json:"resources"`
}

func httpConnect(productKey, deviceName, deviceSecret string) (client conn, err error) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	clientID := uuid.NewV4().String()
	var params = mqttParams{
		Signmethod: "hmacsha1",
		ProductKey: productKey,
		DeviceName: deviceName,
		Timestamp:  timestamp,
		ClientID:   clientID,
		Sign:       aliot.Sign(productKey, deviceName, deviceSecret, clientID, timestamp),
	}
	request := gorequest.New()
	resp, body, _ := request.Post(fmt.Sprintf("%s/auth", httpAuthBaseURL)).Type("json").Send(params).End()
	if err != nil {
		return client, err
	}
	if resp.StatusCode != 200 {
		return client, fmt.Errorf("network error")
	}
	var respParams httpResp
	json.Unmarshal([]byte(body), &respParams)
	client.token = respParams.Info.Token
	return client, nil
}

func (c conn) Pub(payload interface{}) []error {
	request := gorequest.New()
	request.Post(fmt.Sprintf("")).Type("json").Send("").End(func(resp gorequest.Response, body string, errs []error) {
		if len(errs) == 0 && resp != nil && resp.StatusCode == 200 {

		} else {
		}
	})
	return nil
}
