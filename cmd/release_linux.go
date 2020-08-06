package cmd

// WriteDesktopFiles creates TeamID.txt and its shortcut, as well as links
// to the ScoringReport, ReadMe, and other needed files.
func WriteDesktopFiles() {
	if VerboseEnabled {
		InfoPrint("Creating or emptying TeamID.txt...")
	}
	shellCommand("echo 'YOUR-TEAMID-HERE' > " + mc.DirPath + "TeamID.txt")
	shellCommand("chmod 666 " + mc.DirPath + "TeamID.txt")
	shellCommand("chown " + mc.Config.User + ":" + mc.Config.User + " " + mc.DirPath + "TeamID.txt")
	if VerboseEnabled {
		InfoPrint("Writing shortcuts to Desktop...")
	}
	shellCommand("cp " + mc.DirPath + "misc/*.desktop /home/" + mc.Config.User + "/Desktop/")
	shellCommand("chmod +x /home/" + mc.Config.User + "/Desktop/*.desktop")
	shellCommand("chown " + mc.Config.User + ":" + mc.Config.User + " /home/" + mc.Config.User + "/Desktop/*")
}

// ConfigureAutologin configures the auto-login capability for LightDM and
// GDM3, so that the image automatically logs in to the main user's account
// on boot.
func ConfigureAutologin() {
	lightdm, _ := pathExists("/usr/share/lightdm")
	gdm, _ := pathExists("/etc/gdm3/")
	if lightdm {
		if VerboseEnabled {
			InfoPrint("LightDM detected for autologin.")
		}
		shellCommand(`echo "autologin-user=` + mc.Config.User + `" >> /usr/share/lightdm/lightdm.conf.d/50-ubuntu.conf`)
	} else if gdm {
		if VerboseEnabled {
			InfoPrint("GDM3 detected for autologin.")
		}
		shellCommand(`echo -e "AutomaticLogin=True\nAutomaticLogin=` + mc.Config.User + `" >> /etc/gdm3/custom.conf`)
	} else {
		FailPrint("Unable to configure autologin! Please do so manually.")
	}
}

// InstallService for Linux installs and starts the CSSClient init.d service.
func InstallService() {
	if VerboseEnabled {
		InfoPrint("Installing service...")
	}
	shellCommand("cp " + mc.DirPath + "misc/CSSClient /etc/init.d/")
	shellCommand("chmod +x /etc/init.d/CSSClient")
	shellCommand("systemctl enable CSSClient")
	shellCommand("systemctl start CSSClient")
}

// CleanUp for Linux is primarily focused on removing cached files, history,
// and other pieces of forensic evidence. It also removes the non-required
// files in the aeacus directory.
func CleanUp() {
	findPaths := "/bin /etc /home /opt /root /sbin /srv /usr /mnt /var"

	if VerboseEnabled {
		InfoPrint("Changing perms to 755 in " + mc.DirPath + "...")
	}
	shellCommand("chmod 755 -R " + mc.DirPath)

	if VerboseEnabled {
		InfoPrint("Removing .viminfo and .swp files...")
	}
	shellCommand("find " + findPaths + " -iname '*.viminfo*' -delete -iname '*.swp' -delete")

	if VerboseEnabled {
		InfoPrint("Symlinking .bash_history and .zsh_history to /dev/null...")
	}
	shellCommand("find " + findPaths + " -iname '*.bash_history' -exec ln -sf /dev/null {} \\;")
	shellCommand("find " + findPaths + " -name '.zsh_history' -exec ln -sf /dev/null {} \\;")

	if VerboseEnabled {
		InfoPrint("Removing .local files...")
	}
	shellCommand("rm -rf /root/.local /home/*/.local/")

	if VerboseEnabled {
		InfoPrint("Removing cache...")
	}
	shellCommand("rm -rf /root/.cache /home/*/.cache/")

	if VerboseEnabled {
		InfoPrint("Removing temp root and Desktop files...")
	}
	shellCommand("rm -rf /root/*~ /home/*/Desktop/*~")

	if VerboseEnabled {
		InfoPrint("Removing crash and VMWare data...")
	}
	shellCommand("rm -f /var/VMwareDnD/* /var/crash/*.crash")

	if VerboseEnabled {
		InfoPrint("Removing apt and dpkg logs...")
	}
	shellCommand("rm -rf /var/log/apt/* /var/log/dpkg.log")

	if VerboseEnabled {
		InfoPrint("Removing logs (auth and syslog)...")
	}
	shellCommand("rm -f /var/log/auth.log* /var/log/syslog*")

	if VerboseEnabled {
		InfoPrint("Removing initial package list...")
	}
	shellCommand("rm -f /var/log/installer/initial-status.gz")

	if VerboseEnabled {
		InfoPrint("Removing scoring.conf...")
	}
	shellCommand("rm " + mc.DirPath + "scoring.conf*")

	if VerboseEnabled {
		InfoPrint("Removing other setup files...")
	}
	shellCommand("rm -rf " + mc.DirPath + "misc/ " + mc.DirPath + "ReadMe.conf " + mc.DirPath + "README.md " + mc.DirPath + "TODO.md")

	if VerboseEnabled {
		InfoPrint("Removing aeacus binary...")
	}
	shellCommand("rm " + mc.DirPath + "aeacus")

	if VerboseEnabled {
		InfoPrint("Overwriting timestamps to obfuscate changes...")
	}
	shellCommand("find /etc /home /var -exec  touch --date='2012-12-12 12:12' {} \\; 2>/dev/null")
}
