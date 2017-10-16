package mq

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
	gorequest "github.com/parnurzeal/gorequest"
	uuid "github.com/satori/go.uuid"
	aliot "github.com/tosone/aliot"
)

const authURL = "https://iot-auth.cn-shanghai.aliyuncs.com/auth/devicename"

type mqttParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	Sign       string `json:"sign"`
	ClientID   string `json:"clientId"`
	Signmethod string `json:"signmethod"`
	Timestamp  string `json:"timestamp"`
	Resources  string `json:"resources"`
}

type mqttResp struct {
	Code int `json:"code"`
	Data struct {
		IOTID     string `json:"iotId"`
		IOToken   string `json:"iotToken"`
		Resources struct {
			Mqtt struct {
				Host string `json:"host"`
				Port int    `json:"port"`
			}
		} `json:"resources"`
	}
	Msg string `json:"message"`
}

// Conn connect object
type Conn interface {
	Sub()
	UnSub()
	Pub()
}

type conn struct {
	client mq.Client
}

// Callback callback
type Callback = mq.MessageHandler

// Sub 订阅
func (c conn) Sub(topic string, qos byte, callback Callback) error {
	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// UnSub 取消订阅
func (c conn) UnSub(topic string) error {
	if token := c.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Pub 发布消息
func (c conn) Pub(topic string, qos byte, retain bool, payload []byte) error {
	if token := c.client.Publish(topic, qos, retain, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func mqttConn(productKey, deviceName, deviceSecret string) (conn, error) {
	var err error
	var c conn
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	clientID := uuid.NewV4().String()
	sign := aliot.Sign(productKey, deviceName, deviceSecret, clientID, timestamp)
	params := mqttParams{
		ProductKey: productKey,
		DeviceName: deviceName,
		Sign:       sign,
		Signmethod: "hmacsha1",
		ClientID:   clientID,
		Timestamp:  timestamp,
		Resources:  "mqtt",
	}
	request := gorequest.New()
	resp, body, _ := request.Post(authURL).Type("form").Send(params).End()
	if err != nil {
		return c, err
	}
	if resp.StatusCode != 200 {
		return c, fmt.Errorf("network error")
	}
	var respStruct mqttResp
	err = json.Unmarshal([]byte(body), &respStruct)
	if err != nil {
		return c, err
	}
	opts := mq.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", respStruct.Data.Resources.Mqtt.Host, respStruct.Data.Resources.Mqtt.Port))
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Second * 60)
	opts.SetUsername(respStruct.Data.IOTID)
	opts.SetPassword(respStruct.Data.IOToken)
	client := mq.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return c, token.Error()
	}
	c.client = client
	return c, nil
}
