package ping

import (
	"net"
	"reflect"
	"testing"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func TestMergeStatistics(t *testing.T) {
	type args struct {
		stats []*Statistics
	}

	tests := []struct {
		name string
		args args
		want *Statistics
	}{
		{
			name: "test-unreachable",
			args: args{
				stats: []*Statistics{
					{
						PacketsRecv:           0,
						PacketsSent:           0,
						PacketsRecvDuplicates: 0,
						PacketLoss:            0,
						Addr:                  "127.0.0.1",
						Rtts:                  nil,
						TTLs:                  nil,
						MinRtt:                0,
						MaxRtt:                0,
						AvgRtt:                0,
						StdDevRtt:             0,
					},
					{
						PacketsRecv:           0,
						PacketsSent:           0,
						PacketsRecvDuplicates: 0,
						PacketLoss:            0,
						Addr:                  "127.0.0.1",
						Rtts:                  nil,
						TTLs:                  nil,
						MinRtt:                0,
						MaxRtt:                0,
						AvgRtt:                0,
						StdDevRtt:             0,
					},
				},
			},
			want: &Statistics{
				PacketsRecv:           0,
				PacketsSent:           0,
				PacketsRecvDuplicates: 0,
				PacketLoss:            0,
				Addr:                  "127.0.0.1",
				Rtts:                  []time.Duration{},
				TTLs:                  []int{},
				MinRtt:                0,
				MaxRtt:                0,
				AvgRtt:                0,
				StdDevRtt:             0,
			},
		}, {
			name: "test-single-arg",
			args: args{
				stats: []*Statistics{
					{
						PacketsRecv:           5,
						PacketsSent:           5,
						PacketsRecvDuplicates: 1,
						PacketLoss:            42,
						Addr:                  "127.0.0.1",
						Rtts: []time.Duration{
							5 * time.Millisecond,
							3 * time.Millisecond,
							3 * time.Millisecond,
							1 * time.Millisecond,
							2 * time.Millisecond,
						},
						TTLs:      []int{64, 62, 63, 64, 62},
						MinRtt:    0,
						MaxRtt:    0,
						AvgRtt:    0,
						StdDevRtt: 0,
					},
				},
			},
			want: &Statistics{
				PacketsRecv:           5,
				PacketsSent:           5,
				PacketsRecvDuplicates: 1,
				PacketLoss:            0,
				Addr:                  "127.0.0.1",
				Rtts: []time.Duration{
					5 * time.Millisecond,
					3 * time.Millisecond,
					3 * time.Millisecond,
					1 * time.Millisecond,
					2 * time.Millisecond,
				},
				TTLs:      []int{64, 62, 63, 64, 62},
				MinRtt:    1 * time.Millisecond,
				MaxRtt:    5 * time.Millisecond,
				AvgRtt:    time.Duration(2800000),
				StdDevRtt: time.Duration(1483239),
			},
		},
		{
			name: "test-multi-arg",
			args: args{
				stats: []*Statistics{
					{
						PacketsRecv:           5,
						PacketsSent:           5,
						PacketsRecvDuplicates: 1,
						PacketLoss:            42,
						Addr:                  "127.0.0.1",
						Rtts: []time.Duration{
							5 * time.Millisecond,
							3 * time.Millisecond,
							3 * time.Millisecond,
							1 * time.Millisecond,
							2 * time.Millisecond,
						},
						TTLs:      []int{64, 62, 63, 64, 62},
						MinRtt:    0,
						MaxRtt:    0,
						AvgRtt:    0,
						StdDevRtt: 0,
					}, {
						PacketsRecv:           3,
						PacketsSent:           5,
						PacketsRecvDuplicates: 4,
						PacketLoss:            42,
						Addr:                  "127.0.0.1",
						Rtts: []time.Duration{
							7 * time.Millisecond,
							9 * time.Millisecond,
							4 * time.Millisecond,
						},
						TTLs:      []int{127, 126, 128},
						MinRtt:    0,
						MaxRtt:    0,
						AvgRtt:    0,
						StdDevRtt: 0,
					},
				},
			},
			want: &Statistics{
				PacketsRecv:           8,
				PacketsSent:           10,
				PacketsRecvDuplicates: 5,
				PacketLoss:            20,
				Addr:                  "127.0.0.1",
				Rtts: []time.Duration{
					5 * time.Millisecond,
					3 * time.Millisecond,
					3 * time.Millisecond,
					1 * time.Millisecond,
					2 * time.Millisecond,
					7 * time.Millisecond,
					9 * time.Millisecond,
					4 * time.Millisecond,
				},
				TTLs:      []int{64, 62, 63, 64, 62, 127, 126, 128},
				MinRtt:    1 * time.Millisecond,
				MaxRtt:    9 * time.Millisecond,
				AvgRtt:    time.Duration(4250000),
				StdDevRtt: time.Duration(2659215),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeStatistics(tt.args.stats...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStatistics(t *testing.T) {
	type args struct {
		ps *probing.Statistics
	}
	tests := []struct {
		name string
		args args
		want *Statistics
	}{
		{
			name: "test-simple",
			args: args{
				ps: &probing.Statistics{
					PacketsRecv:           9,
					PacketsSent:           10,
					PacketsRecvDuplicates: 1,
					PacketLoss:            90,
					IPAddr:                &net.IPAddr{IP: net.IP{127, 0, 0, 1}, Zone: ""},
					Addr:                  "127.0.0.1",
					Rtts: []time.Duration{
						1 * time.Millisecond,
						2 * time.Millisecond,
						3 * time.Millisecond,
						1 * time.Millisecond,
						2 * time.Millisecond,
						3 * time.Millisecond,
						1 * time.Millisecond,
						2 * time.Millisecond,
						3 * time.Millisecond,
					},
					TTLs:      []uint8{64, 64, 64, 64, 64, 64, 64, 64, 64},
					MinRtt:    1 * time.Millisecond,
					MaxRtt:    3 * time.Millisecond,
					AvgRtt:    2 * time.Millisecond,
					StdDevRtt: 866025,
				},
			},
			want: &Statistics{
				PacketsRecv:           9,
				PacketsSent:           10,
				PacketsRecvDuplicates: 1,
				PacketLoss:            90,
				Addr:                  "127.0.0.1",
				Rtts: []time.Duration{
					1 * time.Millisecond,
					2 * time.Millisecond,
					3 * time.Millisecond,
					1 * time.Millisecond,
					2 * time.Millisecond,
					3 * time.Millisecond,
					1 * time.Millisecond,
					2 * time.Millisecond,
					3 * time.Millisecond,
				},
				TTLs:      []int{64, 64, 64, 64, 64, 64, 64, 64, 64},
				MinRtt:    1 * time.Millisecond,
				MaxRtt:    3 * time.Millisecond,
				AvgRtt:    2 * time.Millisecond,
				StdDevRtt: 866025,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertStatistics(tt.args.ps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}
