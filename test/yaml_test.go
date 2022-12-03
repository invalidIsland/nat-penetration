package test

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"nat-penetration/conf"
	"os"
	"testing"
)

func TestUnmarshalYaml(t *testing.T) {
	s := new(conf.Server)
	b, err := os.ReadFile("../conf/server.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
