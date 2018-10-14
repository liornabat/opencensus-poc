package stats

import "time"

type statsOptions struct {
	exportInterval         time.Duration
	enableInternalExporter bool
	enablePrometheus       bool
	namespace              string
	errFunc                func(err error)
}

type StateOption interface {
	apply(*statsOptions)
}

type funcStatsOption struct {
	f func(*statsOptions)
}

func (fso *funcStatsOption) apply(do *statsOptions) {
	fso.f(do)
}

func newFuncDialOption(f func(*statsOptions)) *funcStatsOption {
	return &funcStatsOption{
		f: f,
	}
}

func WithInternalExporter() StateOption {
	return newFuncDialOption(func(o *statsOptions) {
		o.enableInternalExporter = true
	})
}

func WithExportInterval(t time.Duration) StateOption {
	return newFuncDialOption(func(o *statsOptions) {
		o.exportInterval = t
	})
}
func WithPrometheus(namespace string, errFunc func(err error)) StateOption {
	return newFuncDialOption(func(o *statsOptions) {
		o.namespace = namespace
		o.errFunc = errFunc
	})
}
