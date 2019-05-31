package metric

import (
	"github.com/frennkie/blitzd/internal/data"
	"sync"
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
