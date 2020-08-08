package utilities

import (
	"io/ioutil"
	"regexp"
)

const (
	AeacusVersion = "1.5.0"
	ScoringConf   = "scoring.conf"
	ScoringData   = "scoring.dat"
	linuxDir      = "/opt/aeacus/"
	windowsDir    = "C:\\aeacus\\"
)

var (
	VerboseEnabled = false
	DebugEnabled   = false
	YesEnabled     = false
	mc             = &metaConfig{}
)

// writeFile wraps ioutil's WriteFile function, and prints
// the error the screen if one occurs.
func WriteFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0o644)
	if err != nil {
		FailPrint("Error writing file: " + err.Error())
	}
}

// grepString acts like grep, taking in a pattern to search for, and the
// fileText to search in. It returns the line which contains the string
// (if any).
func GrepString(patternText, fileText string) string {
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + patternText + ".*$")
	return string(re.Find([]byte(fileText)))
}
