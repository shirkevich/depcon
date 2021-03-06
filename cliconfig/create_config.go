package cliconfig

import (
	"fmt"
	"strings"
)

const (
	StrDefaultY         string = "[Y/n]"
	StrDefaultN         string = "[y/N]"
	rootServiceQuestion        = `When only a single environment has been defined, should we root the service command.

For example: "depcon mar app list" would become "depcon app list", eliminating the marathon service command since
we are only dealing with a single service and know what it is.

Root single environment`
)

func CreateNewConfigFromUserInput() *ConfigFile {
	fmt.Println("\n-------------------------------[   Generating Initital Configuration   ]-------------------------------")

	configFile, _ := Load("")
	configFile.RootService = getBoolAnswer(rootServiceQuestion, true)
	configFile.Format = getDefaultFormatOption()
	serviceEnv := createEnvironment()
	configEnv := &ConfigEnvironment{
		Marathon: serviceEnv,
	}
	configFile.Environments[serviceEnv.Name] = configEnv
	configFile.Save()

	return configFile
}

func getDefaultFormatOption() string {

	var response string
	fmt.Println("Default output format (can be overridden via runtime flag)")
	fmt.Println("1 - column")
	fmt.Println("2 - json")
	fmt.Println("3 - yaml\n")
	fmt.Printf("Option: ")

	fmt.Scanf("%s", &response)
	fmt.Println("")

	if strings.HasPrefix(response, "2") {
		return "json"
	}
	if strings.HasPrefix(response, "3") {
		return "yaml"
	}
	return "column"
}

// Asks a yes or no question and returns the boolean equivalent
func getBoolAnswer(question string, defaultTrue bool) bool {
	var response string
	var yn string

	if defaultTrue {
		yn = StrDefaultY
	} else {
		yn = StrDefaultN
	}

	fmt.Printf("\n%s %s? ", question, yn)
	fmt.Scanf("%s", &response)

	if response == "" {
		if defaultTrue {
			return true
		}
		return false
	}

	response = strings.ToLower(response)
	if strings.HasPrefix(response, "y") {
		return true
	} else if strings.HasPrefix(response, "n") {
		return false
	}

	fmt.Printf("\nERROR: Must response with 'y' or 'no'\n", response)
	return getBoolAnswer(question, defaultTrue)
}
