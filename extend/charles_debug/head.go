package charles_debug

import (
	"bytes"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CharlesDebug 线上需要关闭这个

const CharlesHttpProxyAddr = "http://127.0.0.1:8888"

func SendPbToCharlesIfDefine(pb proto.Message, tag string) {

	if Flag {
		var data, err = proto.Marshal(pb)
		if err == nil {
			SendDataToCharlesIfDefine(data, tag)
		}
	}

}

// PbDataSendWithTag send protobuf log to charles
func SendDataToCharlesIfDefine(data []byte, tag string) {

	if Flag {

		var proxy, _ = url.Parse(CharlesHttpProxyAddr)
		var transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}

		var client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 15,
		}

		var webUrl = "http://127.0.0.1:8880/" + tag
		var body io.Reader = bytes.NewReader(data)

		var contentType = "application/x-google-protobuf"

		if strings.HasSuffix(tag, "png") {
			contentType = "image/png"
		} else if strings.HasSuffix(tag, "jpg") {
			contentType = "image/jpg"
		} else if strings.HasSuffix(tag, "json") {
			contentType = "application/json"
		}

		var _, err = client.Post(webUrl, contentType, body)
		if err != nil {

		}

	}
}
