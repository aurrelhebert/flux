package options

import (
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/builtin/functions"
)

func init() {
	flux.RegisterBuiltInOption("now", functions.SystemTime())
}
