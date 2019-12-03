package schema

type Resources struct {
	Cpus    int64  `yaml:"cpus,omitempty" json:"cpus,omitempty"`
	Memory  string `yaml:"memory,omitempty" json:"memory,omitempty"`
	Storage string `yaml:"storage,omitempty" json:"storage,omitempty"`
}
