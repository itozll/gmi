package conf

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const yamlSuffix = ".yaml"

var YamlOption = &Options{
	Suffix: yamlSuffix,
}

func AutoLoadYaml(out interface{}, file string) (f string, err error) {
	return AutoLoad(out, file, LoadYaml)
}

// LoadYaml load yaml file, and decode to out
func LoadYaml(out interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return LoadYamlBuffer(out, data)
}

// LoadYamlBuffer load yaml buffer, and decode to out
func LoadYamlBuffer(out interface{}, buffer []byte) error {
	if len(buffer) == 0 {
		return errors.New("empty buffer")
	}

	return yaml.Unmarshal(buffer, out)
}

// LoadYamlPath load yaml file, and decode to out
func LoadYamlPath(out interface{}, opt *Options, path string) error {
	if opt == nil {
		opt = YamlOption
	} else {
		WithSuffix(yamlSuffix)(opt)
	}

	return loadPath(path, out, LoadYaml, opt)
}

// LoadMultiYamlPath load multi paths
func LoadMultiYamlPath(out interface{}, opt *Options, paths ...string) (err error) {
	if opt == nil {
		opt = YamlOption
	} else {
		WithSuffix(yamlSuffix)(opt)
	}

	for _, path := range paths {
		err = loadPath(path, out, LoadYaml, opt)
		if err != nil {
			return
		}
	}

	return nil
}
