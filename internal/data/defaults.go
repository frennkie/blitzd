package data

import (
	"github.com/golang/protobuf/ptypes"
	"time"
)

const (
	DefaultExpireTime = 300 * time.Second // 5 minutes

	APIv1 = "/api/v1/"
)

var (
	// maxTime (Metric does not expire): "3000-01-01T00:00:00Z"
	MaxTime      = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	MaxTimeNg, _ = ptypes.TimestampProto(MaxTime)
)
