package definition

import (
	"gopkg.in/yaml.v2"
)

// ParseYAML takes a raw set of bytes and
// de-serializes them to a Definition structure
func ParseYAML(raw []byte) (Definition, error) {
	data := Definition{}
	err := yaml.Unmarshal(raw, &data)

	if err != nil {
		return data, err
	}

	return data, nil
}
