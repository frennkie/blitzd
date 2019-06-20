package metric

import (
	"github.com/frennkie/blitzd/internal/data"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/jsonpb"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var MetricsOld data.CacheOld
var MetricsOldMux sync.Mutex

var Lnd data.Lnd
var LndMux sync.Mutex

var Network data.Network
var NetworkMux sync.Mutex

var System data.System

var Metrics data.Cache
var MetricsMux sync.Mutex

var MetricsAPI []v1.Metric
var MetricsAPIMux sync.Mutex

var MetricsFoo4 v1.Metric

var Cache = cache.New(5*time.Minute, 10*time.Minute)

var PrettyPrintJson = jsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: true,
	Indent:       "  ",
	OrigName:     false,
	AnyResolver:  nil,
}
