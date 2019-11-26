package schema

type Task struct {
	Type           string            `yaml:"type,omitempty" json:"type,omitempty"`
	Description    string            `yaml:"description,omitempty" json:"description,omitempty"`
	IgnoreExitCode bool              `yaml:"ignore-exit-code,omitempty" json:"ignore-exit-code,omitempty"`
	Timeout        string            `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Args           []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment    map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
}
