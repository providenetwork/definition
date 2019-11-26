package schema

type Process struct {
	Args        []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Script      Script            `yaml:"script,omitempty" json:"script,omitempty"`
	InputFiles  []InputFile       `yaml:"input-files,omitempty" json:"input-files,omitempty"`
}
