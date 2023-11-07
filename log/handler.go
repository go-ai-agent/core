package log

import (
	"github.com/felixge/httpsnoop"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

// Configure as last handler in chain
//middleware2.ControllerHttpHostMetricsHandler(mux, ""), status

// HttpHostMetricsHandler - handler that applies access logging
func HttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		clone := r
		fn := Access()

		if fn != nil {
			ctx := NewAccessContext(r.Context())
			clone = r.Clone(ctx)
			requestId := runtime.GetOrCreateRequestId(r)
			if r.Header.Get(runtime.XRequestId) == "" {
				r.Header.Set(runtime.XRequestId, requestId)
			}
		}
		m := httpsnoop.CaptureMetrics(appHandler, w, clone)
		// log.Printf("%s %s (code=%d dt=%s written=%d)", r.Method, r.URL, m.Code, m.Duration, m.Written)
		if fn != nil {
			fn(IngressTraffic, start, time.Since(start), clone, &http.Response{StatusCode: m.Code, ContentLength: m.Written}, -1, "")
		}
	})
	return wrappedH
}
