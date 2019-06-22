//go:generate ...

// Package v1 implements a server for Greeter service.
package v1

import (
	"context"
	"errors"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
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
	if config.C.Service.Shutdown.Enabled {
		// ToDo(frennkie) implement this
		return &v1.ShutdownResponse{Message: "Received: ShutdownRequest"}, nil
	} else {
		return &v1.ShutdownResponse{}, errors.New("service disabled")
	}

}

type metricServer struct{}

// NewShutdownServer creates
func NewMetricServer() v1.MetricServiceServer {
	return &metricServer{}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *metricServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (s *metricServer) GetMetricByPath(ctx context.Context, req *v1.GetMetricByPathRequest) (*v1.GetMetricResponse, error) {
	// ToDo(frennkie): Check whether req.Path is set?!
	log.Printf("Received: GetMetricByPath (Path: %s)", req.Path)

	var m v1.Metric
	if x, found := data.Cache.Get(req.Path); found {
		m = x.(v1.Metric)

		expiredAfter, _ := ptypes.Timestamp(m.ExpiredAfter)
		if time.Now().After(expiredAfter) {
			m.Expired = true
		} else {
			m.Expired = false
		}

		return &v1.GetMetricResponse{Api: "v1", Metric: &m}, nil
	}

	return &v1.GetMetricResponse{}, errors.New("not found")

}

func (s *metricServer) GetMetricAll(context.Context, *v1.EmptyRequest) (*v1.GetMetricAllResponse, error) {
	log.Printf("Received: GetMetricAll")

	var mSlice []*v1.Metric
	var m = data.Cache.Items()

	// ToDo(frennkie) try-catch anything here..?! Also this would be nice as {"module.title": Metric, ...}
	for _, v := range m {
		metricObject := interface{}(v.Object).(v1.Metric)

		expiredAfter, _ := ptypes.Timestamp(metricObject.ExpiredAfter)
		if time.Now().After(expiredAfter) {
			metricObject.Expired = true
		} else {
			metricObject.Expired = false
		}

		mSlice = append(mSlice, &metricObject)
	}

	return &v1.GetMetricAllResponse{Api: "v1", Metrics: mSlice}, nil

}

func (s *metricServer) GetMetricFoo(_ context.Context, req *v1.GetMetricFooRequest) (*v1.GetMetricResponse, error) {
	log.Printf("Received: GetMetricFoo")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	var m v1.Metric
	if x, found := data.Cache.Get("lnd.foo5"); found {
		m = x.(v1.Metric)

		expiredAfter, _ := ptypes.Timestamp(m.ExpiredAfter)
		if time.Now().After(expiredAfter) {
			m.Expired = true
		} else {
			m.Expired = false
		}

		return &v1.GetMetricResponse{Api: "v1", Metric: &m}, nil
	}

	return &v1.GetMetricResponse{}, errors.New("not found")
}
