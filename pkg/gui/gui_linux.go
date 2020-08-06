package gui

func LaunchIDPrompt() {
	teamID, err := shellCommandOutput(`
		#!/bin/bash
		teamid=$(
			zenity --entry= \
			--text="Enter in your TeamID here"
		)
		echo $teamid
	`)
	if err == nil {
		writeFile(mc.DirPath+"TeamID.txt", teamID)
	} else {
		FailPrint("Error saving TeamID!")
		sendNotification("Error saving TeamID!")
	}
}

func LaunchConfigGui() {
	WarnPrint("The script doesn't currently have the ability to add multiple check or fail conditions-- you must still do these manually.")
	_, err := shellCommandOutput("bash ./misc/gui_linux.sh")
	if err == nil {
		InfoPrint("Configuration successfully written to scoring.conf!")
	}
}
