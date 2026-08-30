package scaffold

type Options struct {
	NoClobber bool `yaml:"no_clobber"`

	// OutputDir is the directory the caller renders into, as provided on the
	// command line. It is surfaced to templates through .Ctx and is not read
	// from scaffold.yaml.
	OutputDir string `yaml:"-"`
}
