package writer

import (
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

/*
后缀为0, 无返回
1, 返回error
*/

func Err0(w http.ResponseWriter, err error) {
	w.Header().Set("X-Error", "1")
	_, _ = w.Write([]byte(err.Error()))
}

func Err1(w http.ResponseWriter, err error) error {
	Err0(w, err)
	return nil
}

func JsonData0(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func JsonData1(w http.ResponseWriter, data []byte) error {
	JsonData0(w, data)
	return nil
}

func JsonDataOrErr0(w http.ResponseWriter, data []byte, err error) {
	if err != nil {
		Err0(w, err)
	} else {
		JsonData0(w, data)
	}
}

func JsonDataOrErr1(w http.ResponseWriter, data []byte, err error) error {
	JsonDataOrErr0(w, data, err)
	return nil
}

func Text1(w http.ResponseWriter, text string) error {
	w.Header().Set("Content-Type", "plain/txt")
	_, _ = w.Write([]byte(text))
	return nil
}

func HtmlData0(w http.ResponseWriter, text []byte) {
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	_, _ = w.Write(text)
}

func HtmlData1(w http.ResponseWriter, text []byte) error {
	HtmlData0(w, text)
	return nil
}

func HtmlDataOrErr1(w http.ResponseWriter, text []byte, err error) error {
	if err != nil {
		return Err1(w, err)
	} else {
		return HtmlData1(w, text)
	}
}

func HtmlDataOrErr0(w http.ResponseWriter, text []byte, err error) {
	if err != nil {
		Err0(w, err)
	} else {
		HtmlData0(w, text)
	}
}

func Html1(w http.ResponseWriter, text string) error {
	return HtmlData1(w, []byte(text))
}

func Json0(w http.ResponseWriter, msg interface{}) {
	data, err := json.Marshal(msg)

	if err != nil {
		Err0(w, fmt.Errorf("json marshal error: %v", err))
	} else {
		JsonData0(w, data)
	}

}

func Json1(w http.ResponseWriter, msg interface{}) error {
	Json0(w, msg)
	return nil
}

func JsonOrErr1(w http.ResponseWriter, msg interface{}, err error) error {
	if err != nil {
		return Err1(w, err)
	} else {
		return Json1(w, msg)
	}
}

//func WriteProtoData0(w http.ResponseWriter, data []byte) {
//	w.Header().Set("Content-Type", "application/x-google-protobuf")
//	_, _ = w.Write(data)
//}
//
//func WriteProtoMsg0(w http.ResponseWriter, msg proto.Message) {
//	data, err := proto.Marshal(msg)
//	if err != nil {
//		_ = Err1(w, fmt.Errorf("proto marshal error: %v", err))
//		return
//	}
//
//	WriteProtoData0(w, data)
//}
