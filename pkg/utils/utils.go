package utils

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

var (
	// from kong to apisix
	WordMap = map[string]string{
		"round-robin":        "roundrobin",
		"consistent-hashing": "chash",
	}
	configFilePath        = "./repos/apisix-docker/example/apisix_conf/config.yaml"
	ErrMixedCommentStyles = errors.New("unexpected non-empty object")
)

func AddValueToYaml(value interface{}, path ...interface{}) error {
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	v := make(map[interface{}]interface{})

	err = yaml.Unmarshal(yamlFile, &v)
	if err != nil {
		return err
	}

	for i := 1; i <= len(path); i++ {
		fmt.Println(v)
		if i == len(path) {
			if err := dyno.Set(v, value, path...); err != nil {
				return err
			}
		} else {
			if value, err := dyno.Get(v, path[:i]...); err == nil && value != nil {
				// layer already have content
				value.(map[interface{}]interface{})[path[i]] = nil
				if err := dyno.Set(v, value, path[:i]...); err != nil {
					return err
				}
			} else {
				if err := dyno.Set(v, map[interface{}]interface{}{path[i]: nil}, path[:i]...); err != nil {
					return err
				}
			}
		}
	}

	out, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	ioutil.WriteFile(configFilePath, out, 0644)

	return nil
}