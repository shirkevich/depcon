package commands

import (
	"fmt"
	"github.com/gondor/depcon/pkg/cli"
	"github.com/gondor/depcon/pkg/encoding"
	"github.com/gondor/depcon/pkg/logger"
	"os"
)

const (
	FLAG_FORMAT string = "output"
	TypeJSON    string = "json"
	TypeYAML    string = "yaml"
	TypeColumn  string = "column"
)

func init() {
	cli.Register(&cli.CLIWriter{FormatWriter: PrintFormat, ErrorWriter: PrintError})
	rootCmd.PersistentFlags().StringP(FLAG_FORMAT, "o", "column", "Specifies the output format [column | json | yaml]")
}

func getFormatType() string {
	if rootCmd.PersistentFlags().Changed(FLAG_FORMAT) {
		format, err := rootCmd.PersistentFlags().GetString(FLAG_FORMAT)
		if err == nil {
			return format
		}
	}
	if configFile != nil && configFile.Format != "" {
		return configFile.Format
	}
	return "column"
}

func PrintError(err error) {
	logger.Logger().Error("%v", err.Error())
	os.Exit(1)
}

func PrintFormat(formatter cli.Formatter) {
	switch getFormatType() {
	case TypeJSON:
		printEncodedType(formatter, encoding.JSON)
	case TypeYAML:
		printEncodedType(formatter, encoding.YAML)
	default:
		printColumn(formatter)
	}
}

func printEncodedType(formatter cli.Formatter, encoder encoding.EncoderType) {
	e, _ := encoding.NewEncoder(encoder)
	str, _ := e.MarshalIndent(formatter)
	fmt.Println(str)
}

func printColumn(formatter cli.Formatter) {
	err := formatter.ToColumns(os.Stdout)
	if err != nil {
		logger.Logger().Error("Error: %s", err.Error())
	}
}
