package schema

type SharedVolume struct {
	SourcePath string `yaml:"source-path,omitempty" json:"source-path,omitempty"`
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
}
