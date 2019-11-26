package schema

type Test struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	System      []SystemComponent `yaml:"system,omitempty" json:"system,omitempty"`
	Phases      []Phase           `yaml:"phases,omitempty" json:"phases,omitempty"`
}
