// Config
// 1. Read and parse yaml format config file.
// 2. Have an interface that returns the service based on the config.

package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os/exec"

	yaml "gopkg.in/yaml.v2"
)

// Config is the top-level configuration for nanafshi.
type Config struct {
	Shell    string    `yaml:"shell" validate:"required"`
	Services []Service `yaml:"services"`
}

// Load parses given YAML input.
func Load(b []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.UnmarshalStrict(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadFile parses given YAML file.
func LoadFile(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(b)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Service represents setting of file division unit directly under the nanafshi directory.
type Service struct {
	Name  string `yaml:"name" validate:"required"`
	Files []File `yaml:"files"`
}

// File represents setting of command executed at access.
type File struct {
	Name         string       `yaml:"name" validate:"required"`
	ReadCommand  ReadCommand  `yaml:"read"`
	WriteCommand WriteCommand `yaml:"write"`
}

// ReadCommand represents the command to be executed on read and its options.
type ReadCommand struct {
	Command   Command `yaml:"command"`
	CacheTime int     `yaml:"cache"` // TODO: WIP
}

// ReadCommand represents the command to be executed on write and its options.
type WriteCommand struct {
	Command Command `yaml:"command"`
	Async   bool    `yaml:"async"`
}

type Command string

func (c Command) Build(shell string, env []string) (*exec.Cmd, error) {
	if c == "" {
		return nil, errors.New("emtpy command")
	}
	cmd := exec.Command(shell)
	cmd.Env = append(cmd.Env, env...)
	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, string(c))
	stdin.Close()
	return cmd, nil
}
