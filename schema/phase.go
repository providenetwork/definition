package schema

type Phase struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	System      []SystemComponent `yaml:"system,omitempty" json:"system,omitempty"`
	Remove      []string          `yaml:"remove,omitempty" json:"remove,omitempty"`
	Tasks       []Task            `yaml:"tasks,omitempty" json:"tasks,omitempty"`
}
