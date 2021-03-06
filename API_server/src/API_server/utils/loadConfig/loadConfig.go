package loadConfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"API_server/utils/logs"
)

const (
	pkgName = "utils/load-config"

	RETHINKDB_ADDR_ENV        = "RETHINKDB_ADDR"
	RETHINKDB_PORT_ENV        = "RETHINKDB_PORT"
	FACEBOOK_APP_ID_ENV       = "FACEBOOK_APP_ID"
	FACEBOOK_APP_SECRET_ENV   = "FACEBOOK_APP_SECRET"
	FACEBOOK_CALLBACK_URL_ENV = "FACEBOOK_CALLBACK_URL"
)

var (
	l = logs.New(pkgName)
)

func FromFileAndEnv(cfg interface{}, configPath string) error {
	err := FromFile(cfg, configPath)
	if err != nil {
		return err
	}

	FromEnv(cfg, "")
	return nil
}

func FromFile(cfg interface{}, configPath string) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	l.Println("Load config from file:", absPath)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}

func FromEnv(v interface{}, tag string) {
	if tag == "" {
		tag = "json"
	}

	vConfig := reflect.ValueOf(v)
	vConfig = reflect.Indirect(vConfig)
	if vConfig.Kind() != reflect.Struct {
		l.Panicf("%v: config must be a struct", pkgName)
	}

	tConfig := vConfig.Type()
	for n, i := vConfig.NumField(), 0; i < n; i++ {
		vField := vConfig.Field(i)
		tField := tConfig.Field(i)

		tag := tField.Tag.Get("json")

		if tag == "" || strings.HasPrefix(tag, "-") {
			continue
		}

		if vField.Kind() == reflect.Struct {
			FromEnv(vField.Addr().Interface(), tag)
			continue
		}
		if vField.Kind() != reflect.String {
			l.Panicf("%v: field %v must be a string", pkgName, tField.Name)
		}

		env := os.Getenv(tag)
		if env != "" {
			l.Printf("%v=%v\n", tag, env)
			vField.SetString(env)
		}
	}
}
