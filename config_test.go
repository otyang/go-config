package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	stringTOML = `
	userId = 1 
	title = "TOML delectus aut autem"
	completed = false
	
	[company]
	id = 12
	name = "Transform Inc"`

	stringJSON = `
	{
		"userId": 1, 
		"title": "JSON delectus aut autem",
		"completed": false,
		"company": {}
	  }`

	stringYAML = ` 
id: 1 
completed: false
company:
  name: YAML Transform Inc
`
)

func TestLoadFromString(t *testing.T) {
	type (
		Company struct {
			ID   int    `json:"id" yaml:"id"`
			Name string `json:"name" yaml:"name"`
		}

		stringData struct {
			UserID    int     `json:"userId" yaml:"userId"`
			ID        int     `json:"id" yaml:"id"`
			Title     string  `json:"title" yaml:"title"`
			Completed bool    `json:"completed" yaml:"completed"`
			Company   Company `json:"company" yaml:"company"`
		}
	)

	var (
		wantJson = stringData{
			UserID:    1,
			Title:     "JSON delectus aut autem",
			Completed: false,
			Company:   Company{},
		}

		wantToml = stringData{
			UserID:    1,
			Title:     "TOML delectus aut autem",
			Completed: false,
			Company: Company{
				ID:   12,
				Name: "Transform Inc",
			},
		}

		wantYaml = stringData{
			ID:        1,
			Completed: false,
			Company: Company{
				Name: "YAML Transform Inc",
			},
		}
	)

	tests := []struct {
		name    string
		config  *stringData
		data    string
		format  string
		want    stringData
		wantErr bool
	}{
		{name: "empty string", data: "", format: "json", want: stringData{}, wantErr: true},
		{name: "json string", data: stringJSON, format: "json", want: wantJson, wantErr: false},
		{name: "toml string", data: stringTOML, format: "toml", want: wantToml, wantErr: false},
		{name: "yaml string", data: stringYAML, format: "yaml", want: wantYaml, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf stringData

			err := LoadFromString(&cf, tt.data, tt.format)
			assert.Equal(t, (err != nil), tt.wantErr)
			assert.Equal(t, tt.want, cf)
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	var (
		fileNotExists             = "config_test_file_does-not-exist.env"
		fileExistAndConfigValid   = "config_test_valid.env"
		fileExistButConfigInvalid = "config_test_invalid.env.env"
	)

	type fileConfig struct {
		AppName string `env:"APP_NAME" env-default:"Auth"`
		AppPort int    `env:"APP_PORT" env-default:"Auth"`
	}

	wantFileExistConfigValid := fileConfig{
		AppName: "testing-service",
		AppPort: 9000,
	}

	tests := []struct {
		name      string
		config    fileConfig
		filePaths []string
		want      fileConfig
		wantErr   bool
	}{
		{
			name:      "file not exists",
			filePaths: []string{fileNotExists}, want: fileConfig{}, wantErr: true,
		},
		{
			name:      "file exists but config invalid",
			filePaths: []string{fileExistButConfigInvalid}, want: fileConfig{}, wantErr: true,
		},
		{
			name:      "file exists & config valid",
			filePaths: []string{fileExistAndConfigValid}, want: wantFileExistConfigValid, wantErr: false,
		},
		{
			name:      "multiple files, some exist, some dont",
			filePaths: []string{fileExistAndConfigValid, fileNotExists}, want: wantFileExistConfigValid, wantErr: true,
		},
		{
			name:      "multiple files: all exist, some config valid & some not",
			filePaths: []string{fileExistAndConfigValid, fileExistButConfigInvalid}, want: wantFileExistConfigValid, wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf fileConfig

			err := LoadFromFile(&cf, tt.filePaths...)
			assert.Equal(t, (err != nil), tt.wantErr)
			assert.Equal(t, tt.want, cf)
		})
	}
}

type customType string

func (c *customType) UnmarshalText(b []byte) error {
	*c = customType("custom_type_" + string(b))
	return nil
}

func TestLoadFromCLIFlagsOrENV(t *testing.T) {
	s := `./example --database="dsn" --timeout=1s --userids john=123 mary=456 -v=true --workers=8 --custom=only-custom`

	os.Args = strings.Split(s, " ")

	type CLIData struct {
		Custom   customType
		Database string
		Timeout  time.Duration
		UserIDs  map[string]int
		Verbose  bool `arg:"-v,--verbose,env:ENV_VERBOSE" default:"something" help:"verbosity level"`
		Workers  int
	}

	var cfg CLIData

	err := LoadFromCLIFlagsOrENV[CLIData](&cfg)
	assert.NoError(t, err)
	t.Logf("%+v", cfg)
	// t.Errorf("%+v", cfgPtr) //debug purpose
}
