//go:generate ...

// Package v1 implements a server for Greeter service.
package v1

import (
	"context"
	"errors"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"log"
	"time"
)

// server is used to implement v1....
type shutdownServer struct{}

// NewShutdownServer creates
func NewShutdownServer() v1.ShutdownServer {
	return &shutdownServer{}
}

// DoShutdown implements v1.ShutdownServer
func (s *shutdownServer) DoShutdown(ctx context.Context, req *v1.ShutdownRequest) (*v1.ShutdownResponse, error) {
	log.Printf("Received: ShutdownRequest")
	return &v1.ShutdownResponse{Message: "Received: ShutdownRequest"}, nil
}

type metricServer struct{}

// NewShutdownServer creates
func NewMetricServer() v1.MetricServiceServer {
	return &metricServer{}
}

func (s *metricServer) GetMetricByPath(ctx context.Context, req *v1.GetMetricByPathRequest) (*v1.GetMetricResponse, error) {
	// ToDo(frennkie): Check whether req.Path is set?!
	log.Printf("Received: GetMetricByPath (Path: %s)", req.Path)

	var m v1.Metric
	if x, found := metric.Cache.Get(req.Path); found {
		m = x.(v1.Metric)

		expiredAfter, _ := ptypes.Timestamp(m.ExpiredAfter)
		if time.Now().After(expiredAfter) {
			m.Expired = true
		} else {
			m.Expired = false
		}

		return &v1.GetMetricResponse{Api: "1", Metric: &m}, nil
	}

	return &v1.GetMetricResponse{}, errors.New("not found")

}

func (s *metricServer) GetMetricFoo5(context.Context, *v1.GetMetricRequest) (*v1.GetMetricResponse, error) {
	log.Printf("Received: GetMetricFoo5")

	var m v1.Metric
	if x, found := metric.Cache.Get("lnd.foo5"); found {
		m = x.(v1.Metric)

		expiredAfter, _ := ptypes.Timestamp(m.ExpiredAfter)
		if time.Now().After(expiredAfter) {
			m.Expired = true
		} else {
			m.Expired = false
		}

		return &v1.GetMetricResponse{Api: "1", Metric: &m}, nil
	}

	return &v1.GetMetricResponse{}, errors.New("not found")
}
