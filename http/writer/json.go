package writer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
