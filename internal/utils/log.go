package utils

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/DeeStarks/lime/configs"
)

// LogMessage write message to the log file 'lime.log'
func LogMessage(message string, args ...interface{}) {
	message = fmt.Sprintf("[%v]: %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(message, args...))
	// Read the log file
	logFile, err := ioutil.ReadFile(configs.LogFileName)
	if err != nil {
		// File does not exist
		ioutil.WriteFile(configs.LogFileName, []byte(message), 0644)
		return
	}
	// Append the message to the log file
	ioutil.WriteFile(configs.LogFileName, append(logFile, []byte(message)...), 0644)
}
