// Config
// 1. Read and parse yaml format config file.
// 2. Have an interface that returns the service based on the config.

package nanafshi

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Shell    string    `yaml:"shell"`
	Services []Service `yaml:"services"`
}

func LoadConfig(filePath string, config *Config) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, config); err != nil {
		return err
	}
	return nil
}

func (c *Config) NewRoot() *Root {
	now := time.Now()
	return &Root{
		Node: NewNode(),
		NodeInfo: NodeInfo{
			Mode:     os.ModeDir | 0777,
			Creation: now,
			LastMod:  now,
		},
		Services: c.Services,
	}
}

type Service struct {
	Name  string `yaml:"name"`
	Files []File `yaml:"files"`
}

type File struct {
	Name         string       `yaml:"name"`
	ReadCommand  ReadCommand  `yaml:"read"`
	WriteCommand WriteCommand `yaml:"write"`
}

type ReadCommand struct {
	Command   Command `yaml:"command"`
	CacheTime int     `yaml:"cache"` // TODO: WIP
}

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
