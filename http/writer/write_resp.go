package writer

import (
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
