package config

import (
	"io/ioutil"
	"testing"

	"gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	config := Config{}
	if err := LoadConfig("./example/config.yml", &config); err != nil {
		t.Error()
	}

	b, err := ioutil.ReadFile("./example/config.yml")
	if err != nil {
		t.Error("Failed to read sample.yml")
	}
	c := &Config{}
	yaml.Unmarshal(b, c)
	v := validator.New()
	if err := v.Struct(c); err != nil {
		t.Error(err)
	}
}
