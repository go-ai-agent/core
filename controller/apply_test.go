package controller

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	applyTestUri = "urn:postgresql.us-test-1:query.access-log"
)

type testStatus struct {
	code uint32
}

func newStatusOK() *testStatus {
	return &testStatus{code: 0}
}

func newStatusCode(code uint32) *testStatus {
	return &testStatus{code: code}
}

func (t *testStatus) Code() uint32 {
	return t.code
}

func init() {
	defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, statusCode int, routeName string, timeout int, limit rate.Limit, burst int, rateThreshold, statusFlags string) {
		_, host, path := runtime.ParseUri(req.URL.String())
		s := fmt.Sprintf("traffic:%v ,"+
			"route:%v ,"+
			"request-id:%v, "+
			"status-code:%v, "+
			"method:%v, "+
			"url:%v, "+
			"host:%v, "+
			"path:%v, "+
			"timeout:%v, "+
			"rate-limit:%v, "+
			"rate-burst:%v, "+
			"rate-threshold:%v, "+
			"status-flags:%v",
			traffic, routeName, req.Header.Get(RequestIdHeaderName), statusCode, req.Method, req.URL.String(), host, path,
			timeout,
			limit, burst, rateThreshold,
			statusFlags)
		fmt.Printf("{%v}\n", s)
	}
}

func ExampleEgressApply() {
	function(context.Background())

	//Output:
	//{traffic:egress ,route:* ,request-id:123-456-7890, status-code:0, method:GET, url:urn:postgresql.us-test-1:query.access-log, host:postgresql.us-test-1, path:query.access-log, timeout:-1, rate-limit:-1, rate-burst:-1, rate-threshold:, status-flags:}

}

func ExampleEgressApply_RateLimit() {
	name := "rate-limit-route"
	egressTable = NewEgressTable()

	route := NewRoute(name, EgressTraffic, "", false, NewRateLimiterConfig(true, 503, 1, 0, ""))
	errs := CtrlTable().AddController(route)
	fmt.Printf("test: CtrlTable().AddController(route) [errs:%v]\n", errs)
	CtrlTable().SetUriMatcher(func(uri string, method string) (string, bool) {
		return name, true
	})

	functionRateLimited(context.Background())

	//Output:
	//test: CtrlTable().AddController(route) [errs:[]]
	//{traffic:egress ,route:rate-limit-route ,request-id:123-456-7890, status-code:94, method:GET, url:urn:postgresql.us-test-1:query.access-log, host:postgresql.us-test-1, path:query.access-log, timeout:-1, rate-limit:1, rate-burst:0, rate-threshold:, status-flags:RL}

}

func ExampleEgressApply_Timeout() {
	name := "timeout-route"
	egressTable = NewEgressTable()

	route := NewRoute(name, EgressTraffic, "", false, NewTimeoutConfig(true, 504, time.Second))
	CtrlTable().AddController(route)
	CtrlTable().SetUriMatcher(func(uri string, method string) (string, bool) {
		return name, true
	})

	functionTimeout(context.Background())

	//Output:
	//{traffic:egress ,route:timeout-route ,request-id:123-456-7890, status-code:4, method:GET, url:urn:postgresql.us-test-1:query.access-log, host:postgresql.us-test-1, path:query.access-log, timeout:1000, rate-limit:-1, rate-burst:-1, rate-threshold:, status-flags:UT}

}

func function(ctx context.Context) (status *testStatus) {
	var fn func()

	fn, ctx, _ = Apply(ctx, func() int { return int((*(&status)).Code()) }, applyTestUri, "123-456-7890", "GET")
	defer fn()
	return newStatusOK()
}

func functionRateLimited(ctx context.Context) (status *testStatus) {
	var fn func()
	var limited = false

	fn, ctx, limited = Apply(ctx, func() int { return int((*(&status)).Code()) }, applyTestUri, "123-456-7890", "GET")
	defer fn()
	if limited {
		return newStatusCode(StatusRateLimited)
	}
	return newStatusOK()
}

func functionTimeout(ctx context.Context) (status *testStatus) {
	var fn func()
	var limited = false

	fn, ctx, limited = Apply(ctx, func() int { return int((*(&status)).Code()) }, applyTestUri, "123-456-7890", "GET")
	defer fn()
	if limited {
		return newStatusCode(StatusRateLimited)
	}
	done := make(chan struct{})
	panicChan := make(chan any, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		workFunction()
		close(done)
	}()
	// Waiting for events
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		break
	case <-ctx.Done():
		switch err := ctx.Err(); err {
		case context.DeadlineExceeded:
			return newStatusCode(StatusDeadlineExceeded)
		default:
		}
	}
	return newStatusOK()
}

func workFunction() {
	time.Sleep(time.Second * 2)
	return
}