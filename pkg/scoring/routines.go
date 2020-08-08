package scoring

import (
	"aeacus/pkg/misc"
	"errors"
	"os"
	"runtime"
	"time"
)

// ReadScoringData is a convenience function around readData and decodeString,
// which parses the encrypted scoring configuration file.
func ReadScoringData() error {
	decryptedData, err := readData(ScoringData)
	if err != nil {
		if VerboseEnabled || DebugEnabled {
			misc.FailPrint("Error reading in scoring data: " + err.Error())
		}
		return err
	} else if decryptedData == "" {
		misc.FailPrint("Scoring data is empty! Is the file corrupted?")
		return errors.New("Scoring data is empty!")
	}
	parseConfig(decryptedData)
	return nil
}

// CheckConfig parses and checks the validity of the current
// `scoring.conf` file.
func CheckConfig(fileName string) {
	fileContent, err := ReadFile(mc.DirPath + fileName)
	if err != nil {
		misc.FailPrint("Configuration file (" + fileName + ") not found!")
		os.Exit(1)
	}
	parseConfig(fileContent)
	if VerboseEnabled {
		printConfig()
	}
}

// FillConstants determines the correct constants, such as DirPath, for the
// given runtime and environment.
func FillConstants() {
	if runtime.GOOS == "linux" {
		mc.DirPath = linuxDir
	} else if runtime.GOOS == "windows" {
		mc.DirPath = windowsDir
	} else {
		FailPrint("This operating system (" + runtime.GOOS + ") is not supported!")
		os.Exit(1)
	}
}

// RunningPermsCheck is a convenience function wrapper around
// adminCheck, which prints an error indicating that admin
// permissions are needed.
func RunningPermsCheck() {
	if !adminCheck() {
		misc.FailPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}

// timeCheck calls destroyImage if the configured EndDate for the image has
// passed. Its purpose is to dissuade or prevent people using an image after
// the round ends.
func timeCheck() {
	if mc.Config.EndDate != "" {
		endDate, err := time.Parse("2006/01/02 15:04:05 MST", mc.Config.EndDate)
		if err != nil {
			misc.FailPrint("Your EndDate value in the configuration is invalid.")
		} else {
			if time.Now().After(endDate) {
				destroyImage()
			}
		}
	}
}
