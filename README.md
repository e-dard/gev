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
    Home         string `env:"-"`
}

// load the configuration using the process' environment
c := Config{}
if err := gev.Unmarshal(&c); err != nil {
    panic(err)
}

// c.User contains the contents of SERVICE_USER
// c.Token contains the contents of PASS
// c.Home will be ignored by gev because of the `env:"-"` tag.
```

`gev` will ignore exported fields using the `-` on the `env` tag, as
well as un-exported fields.

As well as extracting values from the environment and setting fields on
structs, `gev` will try to parse variables into non-string types.

Types supported:

 - `string`, `*string` and `[]byte`;
 - `int64`, `float64`, `*int64` and `*float64`;
 - `bool` and `*bool`.

In the case of pointer and slice types, `gev` will set them to `nil`
if the relevant environment variable cannot be found. For numeric and
bool types, `Unmarshal` returns an error if the value for the relevant
variable cannot be parsed into the correct type.

A common use-case for `gev` is where you wish to load the majority of a
configuration from a file, and a sensitive portion from environment
variables.

```go

import (
    "github.com/e-dard/gev"
    "gonuts.org/v1/yaml"
)

// Config will contain an application configuration
type Config struct {
    User    string `yaml:"-"   env:"SERVICE_USER"`
    Token   string `yaml:"-"   env:"PASS"`
    Acct    int    `yaml:"-"   env:"ACCT_NUM"`
    LogPath int    `yaml:"pth" env:"-"`
    AppName string `yaml:"app" env:"-"`
    Foo     string `yaml:"foo" env:"-"`
}

// load available Config contents from yaml file
c := Config{}
if err := yaml.Unmarshal(someYaml, &c); err != nil {
    panic(err)
}

// load sensitive part of config from environment using gev
if err := gev.Unmarshal(&c); err != nil {
    panic(err)
}

// c now fully initialised
```
