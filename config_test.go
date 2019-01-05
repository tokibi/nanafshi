package nanafshi

import (
	"fmt"
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshal(t *testing.T) {
	b, err := ioutil.ReadFile("./cmd/nanafshi/sample.yml")
	if err != nil {
		t.Error("Failed to read sample.yml")
	}
	c := &Config{}
	yaml.Unmarshal(b, c)
	fmt.Printf("%+v", c)
}
