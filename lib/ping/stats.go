package ping

import (
	"math"
	"slices"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"gonum.org/v1/gonum/stat"
)

// Statistics wraps probing.Statistics for improved serialization
type Statistics struct {
	PacketsRecv           int             `json:"packets_recv"`
	PacketsSent           int             `json:"packets_sent"`
	PacketsRecvDuplicates int             `json:"packets_recv_duplicates"`
	PacketLoss            float64         `json:"packet_loss"`
	Addr                  string          `json:"addr"`
	Rtts                  []time.Duration `json:"rtts"`
	TTLs                  []int           `json:"ttls"`
	MinRtt                time.Duration   `json:"min_rtt"`
	MaxRtt                time.Duration   `json:"max_rtt"`
	AvgRtt                time.Duration   `json:"avg_rtt"`
	StdDevRtt             time.Duration   `json:"std_dev_rtt"`
}

// ConvertStatistics converts probing.Statistics returned by pinger to Statistics
func ConvertStatistics(ps *probing.Statistics) *Statistics {
	s := &Statistics{
		PacketsRecv:           ps.PacketsRecv,
		PacketsSent:           ps.PacketsSent,
		PacketsRecvDuplicates: ps.PacketsRecvDuplicates,
		PacketLoss:            ps.PacketLoss,
		Addr:                  ps.Addr,
		Rtts:                  make([]time.Duration, 0),
		TTLs:                  make([]int, 0),
		MinRtt:                ps.MinRtt,
		MaxRtt:                ps.MaxRtt,
		AvgRtt:                ps.AvgRtt,
		StdDevRtt:             ps.StdDevRtt,
	}

	s.Rtts = append(s.Rtts, ps.Rtts...)

	for _, i := range ps.TTLs {
		s.TTLs = append(s.TTLs, int(i))
	}

	return s
}

// MergeStatistics merges statistics of multiple probes
func MergeStatistics(stats ...*Statistics) *Statistics {
	s := &Statistics{
		PacketsRecv:           0,
		PacketsSent:           0,
		PacketsRecvDuplicates: 0,
		PacketLoss:            0,
		Addr:                  "",
		Rtts:                  make([]time.Duration, 0),
		TTLs:                  make([]int, 0),
		MinRtt:                0,
		MaxRtt:                0,
		AvgRtt:                0,
		StdDevRtt:             0,
	}

	for _, stat := range stats {
		s.PacketsRecv += stat.PacketsRecv
		s.PacketsSent += stat.PacketsSent
		s.PacketsRecvDuplicates += stat.PacketsRecvDuplicates
		s.Addr = stat.Addr
		s.Rtts = append(s.Rtts, stat.Rtts...)
		s.TTLs = append(s.TTLs, stat.TTLs...)
	}

	if s.PacketsSent > 0 {
		s.PacketLoss = 100

		if s.PacketsRecv > 0 {
			s.PacketLoss = math.Round((1 - (float64(s.PacketsRecv) / float64(s.PacketsSent))) * 100)
		}
	}

	if len(s.Rtts) > 0 {
		s.MinRtt = slices.Min(s.Rtts)
		s.MaxRtt = slices.Max(s.Rtts)

		rtts := make([]float64, 0)

		for _, r := range s.Rtts {
			rtts = append(rtts, float64(r))
		}

		s.AvgRtt = time.Duration(stat.Mean(rtts, nil))
		s.StdDevRtt = time.Duration(stat.StdDev(rtts, nil))
	}

	return s
}
