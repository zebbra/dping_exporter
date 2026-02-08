package ping

import (
	"errors"
	"fmt"
)

const (
	RequestMaxCount      = 100
	RequestMaxInterval   = 10000
	RequestMaxTimeout    = 300000
	RequestMaxPacketSize = 4096
)

// Request is a ping request
type Request struct {
	Target         string `json:"target" query:"target"`
	Count          int    `json:"count" query:"count"`
	DontFragment   bool   `json:"dont_fragment" query:"dont_fragment"`
	Interval       int    `json:"interval" query:"interval"`
	PacketSize     int    `json:"packet_size" query:"packet_size"`
	Privileged     bool   `json:"privileged" query:"privileged"`
	ResolveTimeout int    `json:"resolve_timeout" query:"resolve_timeout"`
	Timeout        int    `json:"timeout" query:"timeout"`
}

// Validate validates the request
func (r *Request) Validate() error {
	if r.Target == "" {
		return errors.New("target is required")
	}

	if r.Count <= 0 || r.Count > RequestMaxCount {
		return fmt.Errorf("count has to be between 1 and %d", RequestMaxCount)
	}

	if r.Interval <= 0 || r.Interval > RequestMaxInterval {
		return fmt.Errorf("interval has to be between 1 and %d", RequestMaxInterval)
	}

	if (r.PacketSize <= 0 && r.PacketSize != -1) || r.PacketSize > RequestMaxPacketSize {
		return fmt.Errorf("packet size has to be between 1 and %d", RequestMaxPacketSize)
	}

	if r.ResolveTimeout < 0 || r.ResolveTimeout > RequestMaxTimeout {
		return fmt.Errorf("resolve timeout has to be between 0 and %d", RequestMaxTimeout)
	}

	if r.Timeout <= 0 || r.Timeout > RequestMaxTimeout {
		return fmt.Errorf("timeout has to be between 1 and %d", RequestMaxTimeout)
	}

	if (r.Count-1)*r.Interval >= r.Timeout {
		return fmt.Errorf(
			"timeout %dms too small to send %d pings with %dms interval, it will always fail",
			r.Timeout,
			r.Count,
			r.Interval,
		)
	}

	return nil
}

// NewRequest returns a new Request struct with default values
func NewRequest() *Request {
	return &Request{
		Count:          3,
		DontFragment:   false,
		Interval:       1000,
		PacketSize:     -1,
		Privileged:     false,
		ResolveTimeout: 0,
		Timeout:        5000,
	}
}
