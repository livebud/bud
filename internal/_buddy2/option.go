package buddy

// WithEmbed option
func WithEmbed(embed bool) embedOption {
	return embedOption(embed)
}

type embedOption bool

func (o embedOption) compile(cfg *compileConfig) {
	cfg.Embed = bool(o)
}

func (o embedOption) build(cfg *buildConfig) {
	cfg.Embed = bool(o)
}

func (o embedOption) run(cfg *runConfig) {
	cfg.Embed = bool(o)
}

// WithMinify option
func WithMinify(minify bool) minifyOption {
	return minifyOption(minify)
}

type minifyOption bool

func (o minifyOption) compile(cfg *compileConfig) {
	cfg.Minify = bool(o)
}

func (o minifyOption) build(cfg *buildConfig) {
	cfg.Minify = bool(o)
}

func (o minifyOption) run(cfg *runConfig) {
	cfg.Minify = bool(o)
}

// WithHot option
func WithHot(hot bool) hotOption {
	return hotOption(hot)
}

type hotOption bool

func (o hotOption) compile(cfg *compileConfig) {
	cfg.Hot = bool(o)
}

func (o hotOption) run(cfg *runConfig) {
	cfg.Hot = bool(o)
}

// WithPort option
func WithPort(port string) portOption {
	return portOption(port)
}

type portOption string

func (o portOption) run(cfg *runConfig) {
	cfg.Port = string(o)
}
