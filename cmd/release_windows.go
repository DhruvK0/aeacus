package cmd

func WriteDesktopFiles() {
	firefoxBinary := `C:\Program Files (x86)\Mozilla Firefox\firefox.exe`
	if VerboseEnabled {
		InfoPrint("Writing ScoringReport.html shortcut to Desktop...")
	}
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ScoringReport.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\assets\ScoringReport.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if VerboseEnabled {
		InfoPrint("Writing ReadMe.html shortcut to Desktop...")
	}
	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ReadMe.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\assets\ReadMe.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if VerboseEnabled {
		InfoPrint("Creating or emptying TeamID.txt file...")
	}
	cmdString = "echo 'YOUR-TEAMID-HERE' > C:\\aeacus\\TeamID.txt"
	shellCommand(cmdString)
	if VerboseEnabled {
		InfoPrint("Writing TeamID shortcut to Desktop...")
	}
	powershellPermission := `
	$ACL = Get-ACL C:\aeacus\TeamID.txt
	$ACL.SetOwner([System.Security.Principal.NTAccount] $env:USERNAME)
	Set-Acl -Path C:\aeacus\TeamID.txt -AclObject $ACL
	`
	shellCommand(powershellPermission)
	if VerboseEnabled {
		InfoPrint("Changing Permissions of TeamID")
	}

	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\TeamID.lnk"); $Shortcut.TargetPath = "C:\aeacus\phocus.exe"; $Shortcut.Arguments = "-i yes"; $Shortcut.Save()`
	shellCommand(cmdString)

	// domain compatibility? doubt
}

func ConfigureAutologin() {
	if VerboseEnabled {
		InfoPrint("Setting Up autologin for " + mc.Config.User + "...")
	}
	powershellAutoLogin := `
	function Test-RegistryValue {

		param (

		 [parameter(Mandatory=$true)]
		 [ValidateNotNullOrEmpty()]$Path,

		[parameter(Mandatory=$true)]
		 [ValidateNotNullOrEmpty()]$Value
		)

		try {

		Get-ItemProperty -Path $Path | Select-Object -ExpandProperty $Value -ErrorAction Stop | Out-Null
		 return $true
		 }

		catch {

		return $false

		}

	}
	$RegPath1Exists = Test-RegistryValue -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -Value "DefaultUsername"
	if ($RegPath1Exists -eq $false) {
		New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "DefaultUsername" -Value $env:USERNAME -type String
	}
	elseif ($RegPath1Exists -eq $true) {
		Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "DefaultUsername" -Value $env:USERNAME -type String
	}

	$RegPath2Exists = Test-RegistryValue -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -Value "AutoAdminLogon"
	if ($RegPath2Exists -eq $false) {
		New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "AutoAdminLogon" -Value 1 -type String
	}
	elseif ($RegPath2Exists -eq $true) {
		Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "AutoAdminLogon" -Value 1 -type String
	}
	`
	shellCommand(powershellAutoLogin)
}

func InstallService() {
	if VerboseEnabled {
		InfoPrint("Installing service with sc.exe...")
	}
	cmdString := `sc.exe create CSSClient binPath= "C:\aeacus\phocus.exe" start= "auto" DisplayName= "CSSClient"`
	shellCommand(cmdString)
	if VerboseEnabled {
		InfoPrint("Setting service description...")
	}
	cmdString = `sc.exe description CSSClient "This is Aeacus's Competition Scoring System client. Don't stop or mess with this unless you want to not get points, and maybe have your registry deleted."`
	shellCommand(cmdString)
}

func CleanUp() {
	if VerboseEnabled {
		InfoPrint("Removing scoring.conf and ReadMe.conf...")
	}
	shellCommand("Remove-Item -Force C:\\aeacus\\scoring.conf")
	shellCommand("Remove-Item -Force C:\\aeacus\\ReadMe.conf")
	if VerboseEnabled {
		InfoPrint("Removing previous.txt...")
	}
	shellCommand("Remove-Item -Force C:\\aeacus\\previous.txt")
	if VerboseEnabled {
		InfoPrint("Emptying recycle bin...")
	}
	shellCommand("Clear-RecycleBin -Force")
	if VerboseEnabled {
		InfoPrint("Clearing recently used...")
	}
	shellCommand("Remove-Item -Force '${env:USERPROFILE}\\AppData\\Roaming\\Microsoft\\Windows\\Recent‌​*.lnk'")
	if VerboseEnabled {
		WarnPrint("Done with automatic CleanUp! You need to remove aeacus.exe manually. The only things you need in the C:\\aeacus directory is phocus, scoring.dat, TeamID.txt, and the assets directory.")
	}
}
