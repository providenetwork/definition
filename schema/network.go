package schema

type Network struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	Bandwidth  string `yaml:"bandwidth,omitempty" json:"bandwidth,omitempty"`
	Latency    string `yaml:"latency,omitempty" json:"latency,omitempty"`
	PacketLoss string `yaml:"packet-loss,omitempty" json:"packet-loss,omitempty"`
}
