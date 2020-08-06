package misc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

// parseConfig takes the config content as a string and attempts to parse it
// into the mc.Config struct based on the TOML spec.
func parseConfig(configContent string) {
	if configContent == "" {
		FailPrint("Configuration is empty!")
	}

	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		FailPrint("Error decoding TOML: " + err.Error())
		os.Exit(1)
	}

	// If there's no remote, local must be enabled.
	if mc.Config.Remote == "" {
		mc.Config.Local = true
	}
}

// WriteConfig reads the plaintext configuration from sourceFile, and writes
// the encrypted configuration into the destFile name passed.
func WriteConfig(sourceFile, destFile string) {
	if VerboseEnabled {
		InfoPrint("Reading configuration from " + mc.DirPath + sourceFile + "...")
	}

	configFile, err := readFile(mc.DirPath + sourceFile)
	if err != nil {
		FailPrint("Can't open scoring configuration file (" + sourceFile + "): " + err.Error())
		os.Exit(1)
	}
	encryptedConfig, err := encryptConfig(configFile)
	if err != nil {
		FailPrint("Encrypting config failed: " + err.Error())
		os.Exit(1)
	} else if VerboseEnabled {
		InfoPrint("Writing data to " + mc.DirPath + "...")
	}
	writeFile(mc.DirPath+destFile, encryptedConfig)
}

// readData is a wrapper around decryptData, taking the scoring data fileName,
// and reading its content. It returns the decrypt config.
func readData(fileName string) (string, error) {
	if VerboseEnabled {
		InfoPrint("Decrypting data from " + mc.DirPath + fileName + "...")
	}
	// Read in the encrypted configuration file.
	dataFile, err := readFile(mc.DirPath + ScoringData)
	if err != nil {
		return "", err
	} else if dataFile == "" {
		return "", errors.New("Scoring data is empty!")
	}
	decryptedConfig, err := decryptConfig(dataFile)
	if err != nil {
		return "", err
	}
	return decryptedConfig, nil
}

// printConfig offers a printed representation of the config, as parsed
// by readData and parseConfig.
func printConfig() {
	PassPrint("Configuration " + mc.DirPath + ScoringConf + " check passed!")
	fmt.Println("Title:", mc.Config.Title)
	fmt.Println("Name:", mc.Config.Name)
	fmt.Println("OS:", mc.Config.OS)
	fmt.Println("User:", mc.Config.User)
	fmt.Println("Remote:", mc.Config.Remote)
	fmt.Println("Local:", mc.Config.Local)
	fmt.Println("EndDate:", mc.Config.EndDate)
	fmt.Println("NoDestroy:", mc.Config.NoDestroy)
	fmt.Println("Checks:")
	for i, check := range mc.Config.Check {
		fmt.Printf("\tCheck %d (%d points):\n", i+1, check.Points)
		fmt.Println("\t\tMessage:", check.Message)
		if check.Pass != nil {
			fmt.Println("\t\tPassConditions:")
			for _, condition := range check.Pass {
				fmt.Printf("\t\t\t%s: %s %s %s %s\n", condition.Type, condition.Arg1, condition.Arg2, condition.Arg3, condition.Arg4)
			}
		}
		if check.PassOverride != nil {
			fmt.Println("\t\tPassOverrideConditions:")
			for _, condition := range check.PassOverride {
				fmt.Printf("\t\t\t%s: %s %s %s %s\n", condition.Type, condition.Arg1, condition.Arg2, condition.Arg3, condition.Arg4)
			}
		}
		if check.Fail != nil {
			fmt.Println("\t\tFailConditions:")
			for _, condition := range check.Fail {
				fmt.Printf("\t\t\t%s: %s %s %s %s\n", condition.Type, condition.Arg1, condition.Arg2, condition.Arg3, condition.Arg4)
			}
		}
	}
}

// ConfirmPrint will prompt the user with the given toPrint string, and
// exit the program if N or n is input.
func ConfirmPrint(toPrint string) {
	printer(color.FgYellow, "CONF", "")
	fmt.Print(toPrint + " [Y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(strings.TrimSpace(resp)) == "n" {
		os.Exit(1)
	}
}

func PassPrint(toPrint string) {
	printer(color.FgGreen, "PASS", toPrint)
}

func FailPrint(toPrint string) {
	printer(color.FgRed, "FAIL", toPrint)
}

func WarnPrint(toPrint string) {
	printer(color.FgYellow, "WARN", toPrint)
}

func InfoPrint(toPrint string) {
	printer(color.FgBlue, "INFO", toPrint)
}

func printer(colorChosen color.Attribute, messageType string, toPrint string) {
	printer := color.New(colorChosen, color.Bold)
	fmt.Printf("[")
	printer.Printf(messageType)
	fmt.Printf("] %s", toPrint)
	if toPrint != "" {
		fmt.Printf("\n")
	}
}

func xor(key string, plaintext string) string {
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		ciphertext[i] = key[i%len(key)] ^ plaintext[i]
	}
	return string(ciphertext)
}

func hexEncode(inputString string) string {
	return hex.EncodeToString([]byte(inputString))
}
