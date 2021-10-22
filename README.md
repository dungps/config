# Go Config
Go config is a simple configuration solution for Go application using [Viper](https://github.com/spf13/viper) and [Consul KV](https://www.consul.io/docs/dynamic-app-config/kv)
- Reading config from file or Consul KV
- Support JSON, TOML, YAML, HCL, envfile and Java properties config files

# Install
```
go get github.com/dungps/config
```

# How to use
1. Reading config from file
```go
err := config.LoadConfigFrom(config.FileType, "./sample.yaml")
if err != nil {
    // handle err
}

viper.GetBool("hello")
```

1. Reading config from Consul KV
```go
err := config.LoadConfigFrom(config.ConsulType, "/sample.yaml,/local.yaml")
if err != nil {
    // handle err
}

viper.GetBool("hello")
```