package tdx

import (
	"BusServer/config"
	"BusServer/constant"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	TDXTokenAPI    string = "https://tdx.transportdata.tw/auth/realms/TDXConnect/protocol/openid-connect/token"
	TDXAccessToken string
)

func GetTDXToken(c *gin.Context) (string, error) {
	uri, err := url.Parse(TDXTokenAPI)
	if err != nil {
		return "", err
	}

	postData := url.Values{}
	postData.Add("grant_type", "client_credentials")
	postData.Add("client_id", config.TDXClientID)
	postData.Add("client_secret", config.TDXClientSecret)

	var respBody []byte
	if respBody, _, err = constant.Request(c, constant.RequestParam{
		Method: http.MethodPost,
		URL:    uri.String(),
		Body:   strings.NewReader(postData.Encode()),
		Header: http.Header{
			"content-type": []string{"application/x-www-form-urlencoded"},
		},
	}); err != nil {
		return "", err
	}

	// resp := map[string]interface{}{}
	// if err = json.Unmarshal(respBody, &resp); err != nil {
	// 	fmt.Println(err.Error(), string(respBody))
	// 	return
	// }
	// fmt.Println("resp:", resp)

	var respData TDXToken
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return "", err
	}

	// respBytes, _ := json.MarshalIndent(respData, "", "    ")
	// fmt.Println(string(respBytes))
	TDXAccessToken = respData.AccessToken
	return TDXAccessToken, nil
}
