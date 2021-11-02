package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type ConfigType string

const (
	FileType   ConfigType = "file"
	ConsulType ConfigType = "consul"
	BytesType  ConfigType = "bytes"
)

func LoadConfigFrom(configType ConfigType, args ...interface{}) error {
	switch configType {
	case FileType:
		if len(args) == 0 || cast.ToString(args[0]) == "" {
			err := ReadConfig("conf", "./conf", ".")
			if err != nil {
				return err
			}
		} else {
			err := ReadConfigFromFile(cast.ToString(args[0]))
			if err != nil {
				return err
			}
		}
	case ConsulType:
		remoteAddress := cast.ToString(args[0])
		var stringConfigRemoteKeys string
		if len(remoteAddress) > 0 && !strings.HasPrefix(remoteAddress, "http") {
			stringConfigRemoteKeys = remoteAddress
			remoteAddress = ""
		} else {
			stringConfigRemoteKeys = cast.ToString(args[1])
		}

		if len(stringConfigRemoteKeys) > 0 {
			configRemoteKeys := stringSlice(stringConfigRemoteKeys, ",")
			for index := 0; index < len(configRemoteKeys); index++ {
				isMerge := true
				if index == 0 {
					isMerge = false
				}

				remoteKey := configRemoteKeys[index]
				if len(remoteKey) == 0 {
					continue
				}
				if remoteKey[0:1] != "/" {
					remoteKey = "/" + remoteKey
				}

				valueBytes, err := ReadConfigFromConsulKV(remoteAddress, remoteKey)
				if err != nil {
					fmt.Printf("error when get config %v from remote: %s", remoteKey, err.Error())
					continue
				}

				err = LoadConfigFromByte(remoteKey, valueBytes, isMerge)
				if err != nil {
					fmt.Printf("error when load config %v: %s", remoteKey, err.Error())
					continue
				}

				fmt.Printf("config %v is loaded", remoteKey)
			}
		}
	case BytesType:
		configType := cast.ToString(args[0])
		valueBytes := []byte(cast.ToString(args[1]))

		err := LoadConfigFromByte(configType, valueBytes, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadConfig(fileName string, configPaths ...string) error {
	viper.SetConfigName(fileName)
	if len(configPaths) < 1 {
		viper.AddConfigPath(".")
	} else {
		for _, configPath := range configPaths {
			viper.AddConfigPath(configPath)
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
	})

	return nil
}

func ReadConfigFromFile(file string) error {
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
	})

	return nil
}

func ReadConfigFromConsulKV(endpoint, key string) ([]byte, error) {
	config := api.DefaultConfig()
	if endpoint != "" {
		config.Address = endpoint
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	kv := client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair != nil {
		return nil, fmt.Errorf("remote config key is not existed: %v", key)
	}

	return pair.Value, nil
}

func LoadConfigFromByte(configType string, value []byte, isMerge bool) (err error) {
	viper.SetConfigType(configType)
	if isMerge {
		err = viper.MergeConfig(bytes.NewBuffer(value))
	} else {
		err = viper.ReadConfig(bytes.NewBuffer(value))
	}

	return
}

func stringSlice(s, sep string) []string {
	var sl []string

	for _, p := range strings.Split(s, sep) {
		if str := strings.TrimSpace(p); len(str) > 0 {
			sl = append(sl, str)
		}
	}

	return sl
}
