package zincsearch

type ZincSearchOption func(*zincSearch)

type options struct {
	getSuffix func() string
}

func WithSuffix(f func() string) ZincSearchOption {
	return func(z *zincSearch) {
		z.options.getSuffix = f
	}
}
