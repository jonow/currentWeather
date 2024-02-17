package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/jonow/currentWeather/server"
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "currentWeather",
	Short: "A simple HTTP server that returns the general weather conditions " +
		"at the provided coordinates.",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.StartServer(port, apiKey)
		if err != nil {
			jww.FATAL.Panicf("Failed to start the server: %+v", err)
		}
	},
}

var (
	logLevel              int
	logPath, port, apiKey string
)

// init initializes the command line flags.
func init() {
	// Initialize all startup flags
	cobra.OnInitialize(initLog)

	rootCmd.PersistentFlags().IntVarP(&logLevel, "logLevel", "v", 0,
		"Verbosity level for log printing (2+ = Trace, 1 = Debug, 0 = Info).")

	rootCmd.PersistentFlags().StringVarP(&logPath, "logPath", "l", "-",
		"File path to save log file to.")

	rootCmd.Flags().StringVarP(&port, "port", "p", "9090",
		"Port the server listens on.")

	rootCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "",
		"OpenWeather API key.")

	err := rootCmd.MarkFlagRequired("apiKey")
	if err != nil {
		jww.ERROR.Printf("Failed to mark flag %q required: %+v", "apiKey", err)
	}
}

// initLog initializes logging thresholds and the log path. If the path is not
// provided, the log output is not set. Possible values for logLevel:
//
//	<0 = warn
//	0  = info
//	1  = debug
//	2+ = trace
func initLog() {
	// Select the level of logs to display
	var threshold jww.Threshold
	if logLevel < 0 {
		threshold = jww.LevelWarn
	} else if logLevel == 0 {
		threshold = jww.LevelInfo
	} else if logLevel == 1 {
		threshold = jww.LevelDebug
	} else {
		// Turn on trace logs
		threshold = jww.LevelTrace
	}
	// Set logging thresholds
	jww.SetLogThreshold(threshold)
	jww.SetStdoutThreshold(threshold)

	// Set log file output
	if logPath == "-" {
		jww.INFO.Print("Setting log output to stdout")
	} else if logPath != "" {
		logFile, err := os.OpenFile(
			logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			jww.FATAL.Panicf("Could not open log file %q: %+v", logPath, err)
		} else {
			jww.SetStdoutOutput(io.Discard)
			jww.SetLogOutput(logFile)
			jww.INFO.Printf("Setting log output to %q", logPath)
		}
	} else {
		jww.INFO.Printf("No log output set: no log path provided")
	}

	jww.INFO.Printf("Log level set to: %s", threshold)
}
