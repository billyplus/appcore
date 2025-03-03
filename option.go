package appcore

type Option func(*Mgr)

func WithLogger(logger Logger) Option {
	return func(mgr *Mgr) {
		mgr.logger = logger
	}
}
