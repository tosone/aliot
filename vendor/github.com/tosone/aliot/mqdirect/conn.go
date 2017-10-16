package mqdirect

import (
	"fmt"
	"strconv"
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
	aliot "github.com/tosone/aliot"
)

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

func mqttConnDirect(productKey, deviceName, deviceSecret string) (conn, error) {
	clientID := uuid.NewV4().String()
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := aliot.Sign(productKey, deviceName, deviceSecret, clientID, timestamp)
	opts := mq.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s.iot-as-mqtt.cn-shanghai.aliyuncs.com:1883", productKey))
	opts.SetClientID(fmt.Sprintf("%s|securemode=3,signmethod=hmacsha1,timestamp=%s|", clientID, timestamp))
	opts.SetKeepAlive(time.Second * 60)
	opts.SetUsername(fmt.Sprintf("%s&%s", deviceName, productKey))
	opts.SetPassword(sign)
	client := mq.NewClient(opts)
	var c conn
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return c, token.Error()
	}
	c.client = client
	return c, nil
}
