package schema

type TaskRunner struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	OutputFiles []string          `yaml:"output-files,omitempty" json:"output-files,omitempty"`
	Resources   Resources         `yaml:"resources,omitempty" json:"resources,omitempty"`
	Args        []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Script      Script            `yaml:"script,omitempty" json:"script,omitempty"`
	InputFiles  []InputFile       `yaml:"input-files,omitempty" json:"input-files,omitempty"`
}
