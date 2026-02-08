package ping

import (
	"context"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Probe runs a pinger on a request and returns its results
func Probe(ctx context.Context, req *Request) (*Response, error) {
	resp := &Response{
		Status: StatusError,
	}

	pinger, err := probing.NewPinger(req.Target)

	if err != nil {
		resp.Message = err.Error()

		return resp, err
	}

	pinger.Count = req.Count
	pinger.Interval = time.Millisecond * time.Duration(req.Interval)
	pinger.ResolveTimeout = time.Millisecond * time.Duration(req.ResolveTimeout)
	pinger.Timeout = time.Millisecond * time.Duration(req.Timeout)
	pinger.SetDoNotFragment(req.DontFragment)
	pinger.SetPrivileged(req.Privileged)

	// if not set, set packet size in request to default packet size of pinger so it can be properly reflected back to client
	if req.PacketSize <= 0 {
		req.PacketSize = pinger.Size
	}
	pinger.Size = req.PacketSize

	err = pinger.RunWithContext(ctx)

	if err != nil {
		resp.Message = err.Error()

		return resp, err
	}

	resp.Request = req
	resp.Result = ConvertStatistics(pinger.Statistics())
	resp.Status = StatusSuccess

	return resp, nil
}
