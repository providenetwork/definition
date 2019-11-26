package schema

type SystemComponent struct {
	Type      string  `yaml:"type,omitempty" json:"type,omitempty"`
	Name      string  `yaml:"name,omitempty" json:"name,omitempty"`
	Count     float64 `yaml:"count,omitempty" json:"count,omitempty"`
	Resources struct {
		Cpus     float64   `yaml:"cpus,omitempty" json:"cpus,omitempty"`
		Memory   string    `yaml:"memory,omitempty" json:"memory,omitempty"`
		Storage  string    `yaml:"storage,omitempty" json:"storage,omitempty"`
		Networks []Network `yaml:"networks,omitempty" json:"networks,omitempty"`
	}
	Sidecars struct {
		Type        []string          `yaml:"type,omitempty" json:"type,omitempty"`
		Name        []string          `yaml:"name,omitempty" json:"name,omitempty"`
		Resources   []Resources       `yaml:"resources,omitempty" json:"resources,omitempty"`
		Args        [][]string        `yaml:"args,omitempty" json:"args,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	}
	Args        []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
}
