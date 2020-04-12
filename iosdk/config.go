package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
)

// IoSDKConfig is the global configuration type
type IoSDKConfig struct {
	// WhiskApiHost is the openwhisk api host
	WhiskAPIHostLocal string `json:"whisk-apihost-local"`
	// WhiskApiHostDocker is the api host within docker
	WhiskAPIHostDocker string `json:"whisk-apihost-docker"`
	// WhiskAPIKey is the openwhisk api key
	WhiskAPIKey string `json:"whisk-apikey"`
	// WhiskNamespace is the openwhisk namespace
	WhiskNamespace string `json:"whisk-namespace"`
	// IoAPIKey is the io api key
	IoAPIKey string `json:"io-apikey"`
	// IoMessages is the io api key
	IoMessages string `json:"io-messages"`
	// AppDir is the application directory
	AppDir string `json:"app-dir"`
}

// ConfigMap returns a map of the configuration
func ConfigMap() map[string]string {
	return map[string]string{
		"whisk-apihost-local":  Config.WhiskAPIHostLocal,
		"whisk-apihost-docker": Config.WhiskAPIHostDocker,
		"whisk-apikey":         Config.WhiskAPIKey,
		"whisk-namespace":      Config.WhiskNamespace,
		"io-apikey":            Config.IoAPIKey,
		"io-messages":          Config.IoMessages,
		"app-dir":              Config.AppDir,
	}
}

// Config is the global configuration
var Config *IoSDKConfig
var ConfigFile string

// ConfigLoad loads the configuration
func ConfigLoad() error {
	var err error
	ConfigFile, err = homedir.Expand("~/.iosdk")
	if err != nil {
		return err
	}
	if _, err := os.Stat(configFile); err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	json.Unmarshal(buf, &Config)
	return nil
}

// ConfigSave save configuration file
func ConfigSave() error {
	if Config == nil {
		return errors.New("empty configuration")
	}
	configFile, err := homedir.Expand("~/.iosdk")
	if err != nil {
		return err
	}
	// saving
	json, err := json.MarshalIndent(Config, "", " ")
	err = ioutil.WriteFile(configFile, json, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Wrote", configFile)

	if *initWskPropsFlag {
		wskFile, err := homedir.Expand("~/.wskprops")
		if err != nil {
			return err
		}
		data := fmt.Sprintf("APIHOST=%s\nAUTH=%s\n", Config.WhiskAPIHostLocal, Config.WhiskAPIKey)
		err = ioutil.WriteFile(wskFile, []byte(data), 0644)
		if err != nil {
			return err
		}
		fmt.Println("Wrote", wskFile)
	}

	return nil
}

// configureDefaults sets defaults in configuration
func configureDefaults() {
	if Config.IoMessages == "" {
		Config.IoMessages = "https://api.cd.italia.it/api/v1/messages"
	}
	if Config.WhiskAPIHostLocal == "" {
		Config.WhiskAPIHostLocal = "http://localhost:3280"
	}
	if Config.WhiskAPIHostDocker == "" {
		Config.WhiskAPIHostDocker = "http://openwhisk:3280"
	}

	if Config.WhiskNamespace == "" {
		Config.WhiskNamespace = "guest"
	}

	// generate random key if not there
	if Config.WhiskAPIKey == "" {
		key := *initWhiskKeyFlag
		if key == "" {
			Config.WhiskAPIKey = fmt.Sprintf("%s:%s", uuid.New(), RandomString(64))
		} else {
			Config.WhiskAPIKey = key
		}
	}
}

func configureAsk() error {
	// ask or override api key
	key := *initIOKeyFlag
	if key == "" {
		key = Input("IO Api Key", Config.IoAPIKey)
	}
	if key == "" {
		return errors.New("You need to provide an api key")
	}
	Config.IoAPIKey = key
	return nil
}

// Configure asking values or setting defaults
func Configure(dir string) error {
	err := ConfigLoad()
	if err != nil {
		// initialized
		Config = &IoSDKConfig{}
	}
	// ignore errors
	Config.AppDir = dir
	configureDefaults()
	if err := configureAsk(); err != nil {
		return err
	}
	return ConfigSave()
}

// PropagateConfig propagate configurations to started services
func PropagateConfig() {

	fmt.Println("Configuring Whisk")
	WhiskUpdatePackageParameters("iosdk", ConfigMap())

	fmt.Println("Configuring IDE")
	cmd := fmt.Sprintf("docker exec ide-js wsk property set --apihost %s --auth %s", Config.WhiskAPIHostDocker, Config.WhiskAPIKey)
	err := Run(cmd)
	if err != nil {
		fmt.Println(err)
	}
}
