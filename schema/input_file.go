package schema

type InputFile struct {
	SourcePath      string `yaml:"source-path,omitempty" json:"source-path,omitempty"`
	DestinationPath string `yaml:"destination-path,omitempty" json:"destination-path,omitempty"`
	Template        bool   `yaml:"template,omitempty" json:"template,omitempty"`
}
