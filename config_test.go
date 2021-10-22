package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfigFromFile(t *testing.T) {
	err := LoadConfigFrom(FileType, "./sample.yaml")
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	t.Log(viper.GetBool("hello"))
}

func TestLoadConfigFromKV(t *testing.T) {
	err := LoadConfigFrom(ConsulType, "http://192.168.100.9:8500", "/sample.yaml")
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	t.Log(viper.GetBool("hello"))
}

func TestLoadConfigFromBytes(t *testing.T) {
	config := "hello: true"

	err := LoadConfigFrom(BytesType, "yaml", []byte(config))
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	t.Log(viper.GetBool("hello"))
}
