package exchange

import (
	"encoding/json"
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"io"
)

var deserializeLoc = PkgUri + "/deserialize"

// Deserialize - templated function, providing deserialization of a request/response body
func Deserialize[E runtime.ErrorHandler, T any](ctx any, body io.ReadCloser) (T, *runtime.Status) {
	var e E
	var t T

	if body == nil {
		return t, e.Handle(ctx, deserializeLoc, errors.New("body is nil")).SetCode(runtime.StatusInvalidContent)
	}
	switch ptr := any(&t).(type) {
	case *[]byte:
		buf, status := ReadAll[E](body)
		if !status.OK() {
			return t, status
		}
		*ptr = buf
	default:
		err := json.NewDecoder(body).Decode(&t)
		if err != nil {
			return t, e.Handle(ctx, deserializeLoc, err).SetCode(runtime.StatusJsonDecodeError)
		}
	}
	return t, runtime.NewStatusOK()
}
