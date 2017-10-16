package aliot

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func Sign(productKey, deviceName, deviceSecret, clientID, timestamp string) string {
	mac := hmac.New(sha1.New, []byte(deviceSecret))
	mac.Write([]byte(fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%s", clientID, deviceName, productKey, timestamp)))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}
