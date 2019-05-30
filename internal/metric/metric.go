package metric

import (
	"github.com/frennkie/blitzd/internal/data"
	"sync"
)

var Lnd data.Lnd
var Metrics data.Cache
var Network data.Network
var System data.System

var LndMux sync.Mutex
var MetricsMux sync.Mutex
var NetworkMux sync.Mutex
var SystemMux sync.Mutex
