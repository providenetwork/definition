package schema

type Script struct {
	SourcePath string `yaml:"source-path,omitempty" json:"source-path,omitempty"`
	Inline     string `yaml:"inline,omitempty" json:"inline,omitempty"`
}
