package configure

import (
	ut "github.com/zdnscloud/cement/unittest"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

type dbConf struct {
	Name     string
	User     string `default:"root"`
	Password string `required:"true"`
	Port     uint   `default:"3306"`
}

type contact struct {
	Name  string
	Email string `required:"true"`
}

type appConfig struct {
	APPName  string `default:"myapp"`
	DB       dbConf
	Contacts []contact
}

func generateDefaultConfig() appConfig {
	config := appConfig{
		APPName: "myapp",
		DB: dbConf{
			Name:     "zdns",
			User:     "zdns",
			Password: "zdns",
			Port:     3306,
		},
		Contacts: []contact{
			contact{
				Name:  "benjamin",
				Email: "hf@zdns.cn",
			},

			contact{
				Name:  "zdns",
				Email: "zdns@zdns.cn",
			},
		},
	}
	return config
}

func TestLoadNormalConfig(t *testing.T) {
	config := generateDefaultConfig()
	bytes, _ := yaml.Marshal(config)
	ut.WithTempFile(t, "configure_unit_test", func(t *testing.T, file *os.File) {
		file.Write(bytes)
		var result appConfig
		Load(&result, file.Name())
		ut.Assert(t, reflect.DeepEqual(result, config), "")
	})
}

func TestDefaultValue(t *testing.T) {
	config := generateDefaultConfig()
	config.APPName = ""
	config.DB.Port = 0

	bytes, _ := yaml.Marshal(config)
	ut.WithTempFile(t, "configure_unit_test", func(t *testing.T, file *os.File) {
		file.Write(bytes)
		var result appConfig
		Load(&result, file.Name())
		ut.Assert(t, reflect.DeepEqual(result, generateDefaultConfig()), "")
	})
}

func TestMissingRequiredValue(t *testing.T) {
	config := generateDefaultConfig()
	config.DB.Password = ""

	bytes, _ := yaml.Marshal(config)
	ut.WithTempFile(t, "configure_unit_test", func(t *testing.T, file *os.File) {
		file.Write(bytes)

		var result appConfig
		err := Load(&result, file.Name())
		ut.Assert(t, err != nil, "")
	})
}
