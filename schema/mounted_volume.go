package schema

type MountedVolume struct {
	DestinationPath string `yaml:"destination-path,omitempty" json:"destination-path,omitempty"`
	VolumeName      string `yaml:"volume-name,omitempty" json:"volume-name,omitempty"`
	Permissions     string `yaml:"permissions,omitempty" json:"permissions,omitempty"`
}
