package test

import (
	"testing"

	"github.com/itozll/ddm/pkg/conf"
	"github.com/itozll/ddm/pkg/env"
)

type vt struct {
	Name    string `json:"name,omitempty" yaml:"name"`
	Version int    `json:"version,omitempty" yaml:"version"`
}

func TestConf(t *testing.T) {
	for _, mode := range []string{"dev", "prd"} {
		var out0 vt
		env.WithEnv(mode)
		_, err := conf.AutoLoadYaml(&out0, "app")
		if err != nil {
			t.Fatal(err.Error())
		}

		var out1 vt
		err = conf.LoadYamlPath(&out1, nil, conf.PathRoot())
		if err != nil {
			t.Fatal(err.Error())
		}
		err = conf.LoadYamlPath(&out1, nil, conf.PathRoot()+mode+"/")
		if err != nil {
			t.Fatal(err.Error())
		}

		var out2 vt
		err = conf.LoadMultiYamlPath(&out2, nil, conf.PathRoot(), conf.PathRoot()+mode+"/")
		if err != nil {
			t.Fatal(err.Error())
		}

		if out1.Name != out2.Name || out0.Name != out1.Name {
			t.Fatal("name not eq")
		}

		if out1.Version != out2.Version || out0.Version != out1.Version {
			t.Fatal("version not eq")
		}

		t.Log(mode, out0)
	}
}
