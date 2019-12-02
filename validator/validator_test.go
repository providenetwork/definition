package validator

/*import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/whiteblock/definition/schema"

	"gopkg.in/yaml.v2"
)

var validatorTests = []struct {
	doc         string
	expectation bool
}{
	{`{
		"firstName": "Biggie",
		"lastName": "Smalls",
		"age": 24
	  }`,
		true,
	},
	{`{
		"whatever": "crimson",
		"lastName": "ghost",
		"age": 1000
	  }`,
		true,
	},
	{`{
		"firstName": "ghoulish",
		"lastName": "being",
		"age": "gtfo"
	  }`,
		false,
	},
}

func TestDummyValidatorCanValidate(t *testing.T) {
	version := "dummy"

	for _, tt := range validatorTests {
		t.Run(tt.doc, func(t *testing.T) {
			v, err := NewValidatorByVersion(version)
			if err != nil {
				t.Error(err)
				return
			}

			err = v.Validate([]byte(tt.doc))

			if err != nil {
				if err != ErrValidationFailed {
					t.Error(err)
				}
			}

			if v.Valid != tt.expectation {
				t.Error("failed to validate")
			}
		})
	}
}

var definedVersionTests = []struct {
	version string
}{
	{
		"v1.0.0",
	},
}

func TestValidatorLoadsByVersion(t *testing.T) {
	for _, tt := range definedVersionTests {
		t.Run(tt.version, func(t *testing.T) {
			_, err := NewValidatorByVersion(tt.version)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestExamplesValidate(t *testing.T) {
	path := "../../../examples"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}

	version := "v1.0.0"

	for _, f := range files {
		t.Run(f.Name(), func(t *testing.T) {
			fp := fmt.Sprintf("%s/%s", path, f.Name())

			d := schema.RootSchema{}

			data, err := ioutil.ReadFile(fp)
			if err != nil {
				t.Error(err)
			}

			err = yaml.Unmarshal(data, &d)
			if err != nil {
				t.Error(err)
			}

			v, err := NewValidatorByVersion(version)
			if err != nil {
				t.Error(err)
				return
			}

			dataJSON, err := json.Marshal(d)
			if err != nil {
				t.Error(err)
			}

			err = v.Validate(dataJSON)

			if err != nil {
				if err != ErrValidationFailed {
					t.Error(err)
				}
			}

			if v.Valid == false {
				t.Error(v.Errors())
			}
		})
	}
}
*/
