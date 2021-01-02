// Package conf provide functions for loading config
// from yaml file or env variables.
//
//
package conf

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ReadFlagsHelper forces user to specify either yaml conf path
// or env prefix through flags and returns them
func ReadFlagsHelper() (confFile, envPrefix string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Printf("Please provide conf\n" +
			"Specify yaml file or prefix for env variables.\n" +
			"Preferably not both..\n")
		flag.PrintDefaults()
	}

	cp := flag.String("c", "", "path to yaml conf file")
	ep := flag.String("e", "", "prefix to env conf variables")
	flag.Parse()
	confFile = *cp
	envPrefix = *ep

	if confFile == "" && envPrefix == "" {
		flag.Usage()
		os.Exit(1)
	}

	return confFile, envPrefix
}

// initViper tries to create viper instance
// and propagates config struct defaults into the instance
// so it knows what keys exist..
func initViper(config interface{}) (*viper.Viper, error) {
	var err error
	v := viper.New()

	defaultConfig, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal default config %v", err)
	}
	v.SetConfigType("yaml")
	if err := v.MergeConfig(bytes.NewReader(defaultConfig)); err != nil {
		return nil, fmt.Errorf("unable to merge default config %v", err)
	}
	return v, nil
}

// unmarshal tries to dump viper config back to config struct
func unmarshal(v *viper.Viper, config interface{}) error {
	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	return nil
}

// LoadYAML fills config variable with data from config yaml file
func LoadYAML(file string, config interface{}) error {
	v, err := initViper(config)
	if err != nil {
		return err
	}

	// name of config file (without extension)
	name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	// path to look for the config file in
	path := filepath.Dir(file)

	// find and read the config file
	v.SetConfigName(name)
	v.AddConfigPath(path)
	if err := v.MergeInConfig(); err != nil {
		return fmt.Errorf("error: %s. while "+
			"reading config file: %s.* in path: %s", err, name, path)
	}

	return unmarshal(v, config)
}

// LoadENV fills config variable from envs based on supplied prefix
// NOTE: Slices in ENVs need value separation by comma: val1,val2,..
func LoadENV(envprefix string, config interface{}) error {
	envprefix = strings.TrimSuffix(envprefix, "_")

	v, err := initViper(config)
	if err != nil {
		return err
	}

	// tell viper to use env variables
	v.AutomaticEnv()
	v.SetEnvPrefix(envprefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return unmarshal(v, config)
}
