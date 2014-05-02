# gev

[gev](http://godoc.org/github.com/e-dard/gev) is a small package that has one job — to make it simple to unmarshal
environment variables to a struct, using tags.

In a nutshell:

```go
package gev

import "github.com/e-dard/gev"

// Config will contain an application configuration
type Config struct {
    User         string `env:"SERVICE_USER"`
    Token        string `env:"PASS"`
    NotSensitive int    `env:"-"`
}

// load the configuration using the process' environment
c := Config{}
err := gev.Unmarshal(&c)
if err != nil {
    panic(err)
}

// c.User contains the contents of SERVICE_USER
// c.Token contains the contents of PASS
```

A common use-case for `gev` is where you wish to load the majority of a
configuration from a file, and a sensitive portion from environment
variables.

```go
// Config will contain an application configuration
type Config struct {
    User    string `yaml:"-"   env:"SERVICE_USER"`
    Token   string `yaml:"-"   env:"PASS"`

    // LogPath will be unmarshaled from a yaml file
    LogPath int    `yaml:"pth" env:"-"`
}
```

`gev` will ignore exported fields using the `-` on the `env` tag, as
well as un-exported fields.

