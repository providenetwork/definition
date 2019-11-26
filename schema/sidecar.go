package schema

type Sidecar struct {
	Name           string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description    string            `yaml:"description,omitempty" json:"description,omitempty"`
	SidecarTo      []string          `yaml:"sidecar-to,omitempty" json:"sidecar-to,omitempty"`
	MountedVolumes []MountedVolume   `yaml:"mounted-volumes,omitempty" json:"mounted-volumes,omitempty"`
	Resources      Resources         `yaml:"resources,omitempty" json:"resources,omitempty"`
	Args           []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment    map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Image          string            `yaml:"image,omitempty" json:"image,omitempty"`
	Script         Script            `yaml:"script,omitempty" json:"script,omitempty"`
	InputFiles     []InputFile       `yaml:"input-files,omitempty" json:"input-files,omitempty"`
}
