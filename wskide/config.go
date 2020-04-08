package wskide

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
)

// IoSDKConfig is the global configuration type
type IoSDKConfig struct {
	// WhiskApiHost is the openwhisk api host
	WhiskAPIHost string `json:"whisk-apihost"`
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

// Config is the global configuration
var Config *IoSDKConfig

// ConfigLoad loads the configuration
func ConfigLoad() error {
	configFile, err := homedir.Expand("~/.iosdk")
	if err != nil {
		return Config, err
	}
	if _, err := os.Stat(configFile); err != nil {
		return Config, err
	}
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config, err
	}
	json.Unmarshal(buf, &Config)
	return Config, err
}

func config() {

	var ioSDKConfig IoSDKConfig
	var scanner *bufio.Scanner

	// fixed values.. do not ask
	var response = map[string]string{
		"WhiskAPIHost":   "http://localhost:3280",
		"WhiskAPIKey":    "",
		"WhiskNamespace": "guest",
	}

	// read .iosdk, no error.. file should not be present
	jsonFile, _ := os.Open(configFile)
	buf, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(buf, &ioSDKConfig)

	// if WhiskAPIKey is "" => generate a random one
	if ioSDKConfig.WhiskAPIKey == "" {
		pass := randomString(64)
		randomWhiskAPIKey := fmt.Sprintf("%s:%s", uuid.New(), pass)
		response["WhiskAPIKey"] = randomWhiskAPIKey
	} else {
		response["WhiskAPIKey"] = ioSDKConfig.WhiskAPIKey
	}

	// struct to interface
	v := reflect.ValueOf(ioSDKConfig)
	typeOfS := v.Type()

	// parse interface, ask/read user input and assign value to response[]
	for i := 0; i < v.NumField(); i++ {
		// ask for value if fields is not in response map
		if _, ok := response[typeOfS.Field(i).Name]; ok == false {
			fmt.Printf("Enter %s: (%s) ", typeOfS.Field(i).Name, v.Field(i).Interface())
			scanner = bufio.NewScanner(os.Stdin)
			scanner.Scan()
			response[typeOfS.Field(i).Name] = scanner.Text()
			if response[typeOfS.Field(i).Name] == "" {
				response[typeOfS.Field(i).Name] = v.Field(i).Interface().(string)
			}
		}
	}

	// the json
	res := &IoSDKConfig{
		WhiskAPIHost:   response["WhiskAPIHost"],
		WhiskAPIKey:    response["WhiskAPIKey"],
		WhiskNamespace: response["WhiskNamespace"],
		IoAPIKey:       response["IoAPIKey"]}
	json, err := json.MarshalIndent(res, "", " ")

	err = ioutil.WriteFile(configFile, json, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Wrote", configFile, "ok")
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
	return nil
}

// configureDefaults sets defaults in configuration
func configureDefaults() {
	if Config.IoMessages == "" {
		Config.IoMessages = "https://api.cd.italia.it/api/v1/messages"
	}
	if Config.WhiskAPIHost == "" {
		Config.WhiskAPIHost = "http://localhost:3280"
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
