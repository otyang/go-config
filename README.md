# go-config
This package provides functions to read and parse configuration data from various sources into Go structs.

**Features:**
* Read configuration from files in various formats: YAML, JSON, TOML, ENV, EDN.
* Parse configuration data from a string in JSON, TOML, or YAML format.
* Populate a struct with configuration values using struct tags for mapping.
* Convenient function to load configuration from command-line flags or environment variables.
 

**Usage:**
```go
package main

import config "github.com/otyang/x-config"

func main() {

	type fileConfig struct {
		AppName string `env:"APP_NAME" env-default:"Auth"`
		AppPort int    `env:"APP_PORT" env-default:"Auth"`
	}

	 // Reads configuration data from one or more files and merges them into the provided struct. 
		var myConfig fileConfig
		err := config.LoadFromFile[fileConfig](&myConfig, "config1.json", "config2.yaml", ....)
		if err != nil {
			// handle error
		}
		println(myConfig.AppName) 

	 // Parses a configuration string data into the provided struct. 
		data := `{"key": "value"}`
		var mySConfig fileConfig
		err := config.LoadFromString(&mySConfig, data, "json")
		if err != nil {
			// handle error
		}
		println(mySConfig.AppName) 

	  // Parses configuration data from command-line flags or environment variables. 
		var myCConfig fileConfig
		err := config.LoadFromCLIFlagsOrENV(&myCConfig)
		if err != nil {
			// handle error
		}
		println(myCConfig.AppName) 
}
```

**Benefits:** 
* Simplifies configuration management by providing various loading options.
* Improves code maintainability by utilizing struct tags for mapping.
* Offers a convenient way to populate configuration from various sources.

**Supported String Formats:** Json, Toml, Yaml

**Supported File Formats:** YAML, JSON, TOML, ENV, EDN

**Further Information:** 
* Refer to the unit test for detailed usage information and examples.
* Feel free to report issues or suggest improvements on the project repository.