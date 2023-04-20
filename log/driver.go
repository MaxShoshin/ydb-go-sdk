package log

import (
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
)

// Driver makes trace.Driver with logging events from details
func Driver(l Logger, d trace.Detailer) (t trace.Driver) { //nolint:gocyclo
	t.OnResolve = func(
		info trace.DriverResolveStartInfo,
	) func(
		trace.DriverResolveDoneInfo,
	) {
		if d.Details()&trace.DriverResolverEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "resolver", "update")
		target := info.Target
		addresses := info.Resolved
		ll.Log(TRACE, "start",
			String("target", target),
			Strings("resolved", addresses),
		)
		return func(info trace.DriverResolveDoneInfo) {
			if info.Error == nil {
				ll.Log(INFO, "done",
					String("target", target),
					Strings("resolved", addresses),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					String("target", target),
					Strings("resolved", addresses),
					version(),
				)
			}
		}
	}
	t.OnInit = func(info trace.DriverInitStartInfo) func(trace.DriverInitDoneInfo) {
		if d.Details()&trace.DriverEvents == 0 {
			return nil
		}
		endpoint := info.Endpoint
		database := info.Database
		secure := info.Secure
		ll := l.WithNames("driver", "resolver", "init")
		ll.Log(TRACE, "start",
			String("endpoint", endpoint),
			String("database", database),
			Bool("secure", secure),
		)
		start := time.Now()
		return func(info trace.DriverInitDoneInfo) {
			if info.Error == nil {
				ll.Log(INFO, "done",
					String("endpoint", endpoint),
					String("database", database),
					Bool("secure", secure),
					latency(start),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					String("endpoint", endpoint),
					String("database", database),
					Bool("secure", secure),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnClose = func(info trace.DriverCloseStartInfo) func(trace.DriverCloseDoneInfo) {
		if d.Details()&trace.DriverEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "resolver", "close")
		ll.Log(TRACE, "start")
		start := time.Now()
		return func(info trace.DriverCloseDoneInfo) {
			if info.Error == nil {
				ll.Log(INFO, "done",
					latency(start),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnConnDial = func(info trace.DriverConnDialStartInfo) func(trace.DriverConnDialDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "dial")
		endpoint := info.Endpoint
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()
		return func(info trace.DriverConnDialDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					Stringer("endpoint", endpoint),
					latency(start),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnConnStateChange = func(info trace.DriverConnStateChangeStartInfo) func(trace.DriverConnStateChangeDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "state", "change")
		endpoint := info.Endpoint
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
			Stringer("state", info.State),
		)
		start := time.Now()
		return func(info trace.DriverConnStateChangeDoneInfo) {
			ll.Log(TRACE, "done",
				Stringer("endpoint", endpoint),
				latency(start),
				Stringer("state", info.State),
			)
		}
	}
	t.OnConnPark = func(info trace.DriverConnParkStartInfo) func(trace.DriverConnParkDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "park")
		endpoint := info.Endpoint
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()
		return func(info trace.DriverConnParkDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					Stringer("endpoint", endpoint),
					latency(start),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnConnClose = func(info trace.DriverConnCloseStartInfo) func(trace.DriverConnCloseDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "close")
		endpoint := info.Endpoint
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()
		return func(info trace.DriverConnCloseDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					Stringer("endpoint", endpoint),
					latency(start),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnConnInvoke = func(info trace.DriverConnInvokeStartInfo) func(trace.DriverConnInvokeDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "invoke")
		endpoint := info.Endpoint
		method := string(info.Method)
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
			String("method", method),
		)
		start := time.Now()
		return func(info trace.DriverConnInvokeDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					Stringer("endpoint", endpoint),
					String("method", method),
					latency(start),
					Stringer("metadata", metadata(info.Metadata)),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					String("method", method),
					latency(start),
					Stringer("metadata", metadata(info.Metadata)),
					version(),
				)
			}
		}
	}
	t.OnConnNewStream = func(
		info trace.DriverConnNewStreamStartInfo,
	) func(
		trace.DriverConnNewStreamRecvInfo,
	) func(
		trace.DriverConnNewStreamDoneInfo,
	) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "new", "stream")
		endpoint := info.Endpoint
		method := string(info.Method)
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
			String("method", method),
		)
		start := time.Now()
		return func(info trace.DriverConnNewStreamRecvInfo) func(trace.DriverConnNewStreamDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "intermediate receive",
					Stringer("endpoint", endpoint),
					String("method", method),
					latency(start),
				)
			} else {
				ll.Log(WARN, "intermediate fail",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					String("method", method),
					latency(start),
					version(),
				)
			}
			return func(info trace.DriverConnNewStreamDoneInfo) {
				if info.Error == nil {
					ll.Log(TRACE, "done",
						Stringer("endpoint", endpoint),
						String("method", method),
						latency(start),
						Stringer("metadata", metadata(info.Metadata)),
					)
				} else {
					ll.Log(WARN, "failed",
						Error(info.Error),
						Stringer("endpoint", endpoint),
						String("method", method),
						latency(start),
						Stringer("metadata", metadata(info.Metadata)),
						version(),
					)
				}
			}
		}
	}
	t.OnConnBan = func(info trace.DriverConnBanStartInfo) func(trace.DriverConnBanDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "ban")
		endpoint := info.Endpoint
		ll.Log(WARN, "start",
			Stringer("endpoint", endpoint),
			NamedError("cause", info.Cause),
			version(),
		)
		start := time.Now()
		return func(info trace.DriverConnBanDoneInfo) {
			ll.Log(WARN, "done",
				Stringer("endpoint", endpoint),
				latency(start),
				Stringer("state", info.State),
				version(),
			)
		}
	}
	t.OnConnAllow = func(info trace.DriverConnAllowStartInfo) func(trace.DriverConnAllowDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "conn", "allow")
		endpoint := info.Endpoint
		ll.Log(TRACE, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()
		return func(info trace.DriverConnAllowDoneInfo) {
			ll.Log(INFO, "done",
				Stringer("endpoint", endpoint),
				latency(start),
				Stringer("state", info.State),
			)
		}
	}
	t.OnRepeaterWakeUp = func(info trace.DriverRepeaterWakeUpStartInfo) func(trace.DriverRepeaterWakeUpDoneInfo) {
		if d.Details()&trace.DriverRepeaterEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "repeater", "wake", "up")
		name := info.Name
		event := info.Event
		ll.Log(TRACE, "start",
			String("name", name),
			String("event", event),
		)
		start := time.Now()
		return func(info trace.DriverRepeaterWakeUpDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					String("name", name),
					String("event", event),
					latency(start),
				)
			} else {
				ll.Log(ERROR, "failed",
					Error(info.Error),
					String("name", name),
					String("event", event),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnBalancerInit = func(info trace.DriverBalancerInitStartInfo) func(trace.DriverBalancerInitDoneInfo) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "balancer", "init")
		ll.Log(TRACE, "start")
		start := time.Now()
		return func(info trace.DriverBalancerInitDoneInfo) {
			ll.Log(DEBUG, "done",
				latency(start),
			)
		}
	}
	t.OnBalancerClose = func(info trace.DriverBalancerCloseStartInfo) func(trace.DriverBalancerCloseDoneInfo) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "balancer", "close")
		ll.Log(TRACE, "start")
		start := time.Now()
		return func(info trace.DriverBalancerCloseDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					latency(start),
				)
			} else {
				ll.Log(ERROR, "failed",
					Error(info.Error),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnBalancerChooseEndpoint = func(
		info trace.DriverBalancerChooseEndpointStartInfo,
	) func(
		trace.DriverBalancerChooseEndpointDoneInfo,
	) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "balancer", "choose", "endpoint")
		ll.Log(TRACE, "start")
		start := time.Now()
		return func(info trace.DriverBalancerChooseEndpointDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					latency(start),
					Stringer("endpoint", info.Endpoint),
				)
			} else {
				ll.Log(WARN, "failed",
					Error(info.Error),
					latency(start),
					version(),
				)
			}
		}
	}
	t.OnBalancerUpdate = func(
		info trace.DriverBalancerUpdateStartInfo,
	) func(
		trace.DriverBalancerUpdateDoneInfo,
	) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "balancer", "update")
		ll.Log(TRACE, "start",
			Bool("needLocalDC", info.NeedLocalDC),
		)
		start := time.Now()
		return func(info trace.DriverBalancerUpdateDoneInfo) {
			ll.Log(INFO, "done",
				latency(start),
				Stringer("endpoints", endpoints(info.Endpoints)),
				String("detectedLocalDC", info.LocalDC),
			)
		}
	}
	t.OnGetCredentials = func(info trace.DriverGetCredentialsStartInfo) func(trace.DriverGetCredentialsDoneInfo) {
		if d.Details()&trace.DriverCredentialsEvents == 0 {
			return nil
		}
		ll := l.WithNames("driver", "credentials", "get")
		ll.Log(TRACE, "start")
		start := time.Now()
		return func(info trace.DriverGetCredentialsDoneInfo) {
			if info.Error == nil {
				ll.Log(TRACE, "done",
					latency(start),
					String("token", Secret(info.Token)),
				)
			} else {
				ll.Log(ERROR, "done",
					Error(info.Error),
					latency(start),
					String("token", Secret(info.Token)),
					version(),
				)
			}
		}
	}
	return t
}
