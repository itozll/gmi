package conf

type FuncOption func(*Options)

type Options struct {
	Prefix string // has prefix
	Suffix string // has suffix

	IncludeDir      bool   // include dir, default false
	Depth           int    // depth
	FileName        string // file name
	ExcludeFileName string // exclude file name

	Recursion bool // recursion, default false,
}

// defaultOptions default Options
func defaultOptions() *Options { return &Options{} }

func WithPrefix(prefix string) FuncOption {
	return func(opt *Options) {
		opt.Prefix = prefix
	}
}

func WithFileName(FileName string) FuncOption {
	return func(opt *Options) {
		opt.FileName = FileName
	}
}

func WithSuffix(suffix string) FuncOption {
	return func(opt *Options) {
		opt.Suffix = suffix
	}
}

func WithRecursion(recursion bool) FuncOption {
	return func(opt *Options) {
		opt.Recursion = recursion
	}
}

func WithIncludeDir(includeDir bool) FuncOption {
	return func(opt *Options) {
		opt.IncludeDir = includeDir
	}
}

func WithDepth(depth int) FuncOption {
	return func(opt *Options) {
		opt.Depth = depth
	}
}

func WithExcludeFileName(ExcludeFileName string) FuncOption {
	return func(opt *Options) {
		opt.ExcludeFileName = ExcludeFileName
	}
}
