package httpx

import (
	"encoding/json"
	"errors"
	io2 "github.com/go-ai-agent/core/io"
	"github.com/go-ai-agent/core/runtime"
	"io"
)

var deserializeLoc = PkgUri + "/Deserialize"

// Deserialize - provide deserialization of a request/response body
func Deserialize[T any](body io.ReadCloser) (T, *runtime.Status) {
	var t T

	if body == nil {
		return t, runtime.NewStatusError(runtime.StatusInvalidContent, deserializeLoc, errors.New("body is nil"))
	}
	switch ptr := any(&t).(type) {
	case *[]byte:
		buf, status := io2.ReadAll(body)
		if !status.OK() {
			return t, status
		}
		*ptr = buf
	default:
		err := json.NewDecoder(body).Decode(&t)
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, deserializeLoc, err)
		}
	}
	return t, runtime.NewStatusOK()
}
