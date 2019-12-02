package schema

type RootSchema struct {
	Services    []Service    `yaml:"services,omitempty" json:"services,omitempty"`
	Sidecars    []Sidecar    `yaml:"sidecars,omitempty" json:"sidecars,omitempty"`
	TaskRunners []TaskRunner `yaml:"task-runners,omitempty" json:"task-runners,omitempty"`
	Tests       []Test       `yaml:"tests,omitempty" json:"tests,omitempty"`
}
