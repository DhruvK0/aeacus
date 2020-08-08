package gui

import "aeacus/pkg/utilities"

func LaunchIDPrompt() {
	teamID, err := utilities.ShellCommandOutput(`
		#!/bin/bash
		teamid=$(
			zenity --entry= \
			--text="Enter in your TeamID here"
		)
		echo $teamid
	`)
	if err == nil {
		WriteFile(mc.DirPath+"TeamID.txt", teamID)
	} else {
		FailPrint("Error saving TeamID!")
		SendNotification("Error saving TeamID!")
	}
}

func LaunchConfigGui() {
	WarnPrint("The script doesn't currently have the ability to add multiple check or fail conditions-- you must still do these manually.")
	_, err := utilities.ShellCommandOutput("bash ./misc/gui_linux.sh")
	if err == nil {
		InfoPrint("Configuration successfully written to scoring.conf!")
	}
}
