package json

import (
	"encoding/json"
	"github.com/go-ai-agent/core/runtime"
)

func MarshalType(t any) ([]byte, *runtime.Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusJsonEncodeError, runtime.PkgUri+"/MarshalType", err)
	}
	return buf, runtime.NewStatusOK()
}

func UnmarshalType[T any](buf []byte) (T, *runtime.Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, runtime.PkgUri+"/UnmarshalType", err)
	}
	return t, runtime.NewStatusOK()
}