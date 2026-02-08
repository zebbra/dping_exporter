package canary

import (
	"context"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// ProbeResult contains results of a single canary probe
type ProbeResult struct {
	Target     string
	Error      error
	Statistics *probing.Statistics
}

// Canary is a canary prober
type Canary struct {
	Count             int
	Interval          time.Duration
	Timeout           time.Duration
	Targets           []string
	LivenessThreshold int
}

func (c Canary) Alive(ctx context.Context) (bool, []*ProbeResult) {
	wg := sync.WaitGroup{}
	results := make([]*ProbeResult, len(c.Targets))

	for i, target := range c.Targets {
		wg.Add(1)

		go func() {
			defer wg.Done()
			res := &ProbeResult{
				Target: target,
			}

			p, err := probing.NewPinger(target)

			if err != nil {
				res.Error = err
				return
			}

			p.Count = c.Count
			p.Interval = c.Interval
			p.Timeout = c.Timeout

			err = p.RunWithContext(ctx)

			if err != nil {
				res.Error = err
				return
			}

			res.Statistics = p.Statistics()
			results[i] = res
		}()
	}

	wg.Wait()

	unreachable := 0

	for _, result := range results {
		if result.Error != nil {
			unreachable++
			continue
		}

		if result.Statistics != nil {
			stats := result.Statistics

			if stats.PacketsRecv == 0 {
				unreachable++
				continue
			}
		}
	}

	if len(c.Targets)-unreachable >= c.LivenessThreshold {
		return true, results
	}

	return false, results
}

// New creates a new Canary prober
func New(targets []string) *Canary {
	return &Canary{
		Count:             5,
		Interval:          500 * time.Millisecond,
		Timeout:           5 * time.Second,
		Targets:           targets,
		LivenessThreshold: 1,
	}
}
