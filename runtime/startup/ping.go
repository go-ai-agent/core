package startup

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

const (
	maxWait = time.Second * 2
)

var pingLocation = PkgUri + "/Ping"

// Ping - templated function to "ping" a startup
func Ping[E runtime.ErrorHandler](ctx context.Context, uri string) (status *runtime.Status) {
	var e E

	if uri == "" {
		return e.Handle(runtime.ContextRequestId(ctx), pingLocation, errors.New("invalid argument: startup uri is empty"))
	}
	cache := NewMessageCache()
	msg := Message{To: uri, From: HostName, Event: PingEvent, Status: nil, ReplyTo: NewMessageCacheHandler(cache)}
	err := directory.Send(msg)
	if err != nil {
		return e.Handle(runtime.ContextRequestId(ctx), pingLocation, err)
	}
	duration := maxWait
	for wait := time.Duration(float64(duration) * 0.20); duration >= 0; duration -= wait {
		time.Sleep(wait)
		result, err1 := cache.Get(uri)
		if err1 != nil {
			continue //return e.HandleWithContext(ctx, pingLocation, err1)
		}
		if result.Status == nil {
			return e.Handle(runtime.ContextRequestId(ctx), pingLocation, errors.New(fmt.Sprintf("ping response status not available: [%v]", uri))).SetCode(runtime.StatusNotProvided)
		}
		return result.Status
	}
	return e.Handle(runtime.ContextRequestId(ctx), pingLocation, errors.New(fmt.Sprintf("ping response time out: [%v]", uri))).SetCode(runtime.StatusDeadlineExceeded)
}