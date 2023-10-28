package httpx

import (
	"errors"
	"net/http"
	"strings"
)

const (
	ContentLocation = "Content-Location"
	ContentTypeText = "text/plain" // charset=utf-8
	ContentTypeJson = "application/json"
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
)

func GetContentLocation(req *http.Request) string {
	if req != nil && req.Header != nil {
		return req.Header.Get(ContentLocation)
	}
	return ""
}

func GetContentType(attrs []Attr) string {
	for _, attr := range attrs {
		if attr.Key == ContentType {
			return attr.Val
		}
	}
	return ""
}

func CreateHeaders(h http.Header, resp *http.Response, keys ...string) {
	if resp == nil || len(keys) == 0 {
		return
	}
	if keys[0] == "*" {
		keys = []string{}
		for k := range resp.Header {
			keys = append(keys, k)
		}
	}
	if len(keys) > 0 {
		for _, k := range keys {
			if k != "" {
				h.Add(k, resp.Header.Get(k))
			}
		}
	}
}

func SetHeaders(w http.ResponseWriter, kv []Attr) {
	//err := ValidateKVHeaders(kv...)
	//if err != nil {
	//	return err
	//}
	for _, pair := range kv {
		//key :=
		w.Header().Set(strings.ToLower(pair.Key), pair.Val)
		//if i+1 >= len(kv) {
		//	w.Header().Set(key, "")
		//} else {
		//	w.Header().Set(key, kv[i+1])
		//}
	}
}

func ValidateKVHeaders(kv ...string) error {
	if (len(kv) & 1) == 1 {
		return errors.New("invalid number of kv items: number is odd, missing a value")
	}
	return nil
}
