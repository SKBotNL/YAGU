package main

import (
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	"github.com/ucwong/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type latestjson struct {
	Name string `json:"name"`
}

func main() {
	idout, err := exec.Command("id", "-u").Output()
	if err != nil {
		log.Fatal(err)
	}
	idoutput := string(idout[:len(idout)-1])

	if idoutput != "0" {
		fmt.Println("Please run the program with sudo")
		return
	}

	commandout, err := exec.Command("whereis", "tar").Output()
	if err != nil {
		log.Fatal(err)
	}
	commandoutput := strings.TrimSuffix(string(commandout), "\n")

	if commandoutput == "tar:" {
		fmt.Println("Tar not found")
		return
	}

	var input string
ask:
	color.Blue("-----YARU-----")
	fmt.Println("What do you want to update or install?")
	color.Green("(1) Waterfox")
	color.Green("(2) VSCode Insiders")
	color.Green("(3) MultiMC Development")
	color.Red("(u) Uninstall")
	fmt.Println("(q) Quit")
	fmt.Print("Please input what you want to do: ")
	fmt.Scanln(&input)

	if input == "q" {
		return
	}
	if input == "u" {
		color.Blue("-----YARU-----")
		fmt.Println("What do you want to uninstall?")
		color.Red("(1) Waterfox")
		color.Red("(2) VSCode Insiders")
		color.Red("(3) MultiMC Development")
		fmt.Println("(q) Quit")
		fmt.Print("Please input what you want to do: ")
		fmt.Scanln(&input)

		if input == "q" {
			return
		}
		if input == "1" {
			uninstallwaterfox()
		}
		if input == "2" {
			uninstallvscode()
		}
		if input == "3" {
			uninstallmultimc()
		}
		goto ask
	}

	if input != "1" && input != "2" && input != "3" {
		fmt.Println("That is not a valid option")
		goto ask
	}
	if input == "1" {
		updatewaterfox()
	}
	if input == "2" {
		updatevscode()
	}
	if input == "3" {
		updatemultimc()
	}
}

func updatewaterfox() {

	fmt.Println("Downloading Waterfox...")
	githuburl := "https://api.github.com/repos/WaterfoxCo/Waterfox/releases/latest"
	httpclient := http.Client{}

	req, err := http.NewRequest(http.MethodGet, githuburl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "YAU")
	res, error := httpclient.Do(req)
	if error != nil {
		log.Fatal(error)
	}

	body, readerr := ioutil.ReadAll(res.Body)
	if readerr != nil {
		log.Fatal(readerr)
	}

	latestjson := latestjson{}
	json.Unmarshal(body, &latestjson)

	releasename := latestjson.Name
	url := fmt.Sprintf("https://github.com/WaterfoxCo/Waterfox/releases/download/%s/waterfox-%s.en-US.linux-x86_64.tar.bz2", releasename, releasename)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("waterfox.tar.bz2")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	color.Green("Downloaded Waterfox")
	fmt.Println("Unzipping...")
	exec.Command("bsdtar", "-xf", "waterfox.tar.bz2").Run()
	color.Green("Unzipped Waterfox")
	fmt.Println("Installing Waterfox...")
	os.Mkdir("/opt/waterfox", 0777)
	err = copy.Copy("waterfox/", "/opt/waterfox/")
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Create("/usr/share/applications/Waterfox.desktop")
	if err != nil {
		log.Fatal(err)
	}
	desktopfile := "[Desktop Entry]\nName=Waterfox\nComment=Striking the perfect balance between privacy and usability\nExec=/opt/waterfox/waterfox\nIcon=/opt/waterfox/browser/chrome/icons/default/default128.png\nCategories=Network;\nTerminal=false\nType=Application"
	err = ioutil.WriteFile("/usr/share/applications/Waterfox.desktop", []byte(desktopfile), 0777)
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Installed Waterfox")
	fmt.Println("Cleaning up...")
	err = os.Remove("waterfox.tar.bz2")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("waterfox/")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Done")
	main()
}

func uninstallwaterfox() {
	fmt.Println("Uninstalling Waterfox...")
	err := os.RemoveAll("/opt/waterfox/")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("/usr/share/applications/Waterfox.desktop")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Uninstalled Waterfox")
}

func updatevscode() {
	fmt.Println("Downloading VSCode Insiders...")
	resp, err := http.Get("https://code.visualstudio.com/sha/download?build=insider&os=linux-x64")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("code-insider.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	color.Green("Downloaded VSCode Insiders")
	fmt.Println("Unzipping...")
	exec.Command("bsdtar", "-xf", "code-insider.tar.gz").Run()
	color.Green("Unzipped VSCode Insiders")
	fmt.Println("Installing VSCode Insiders...")
	os.Mkdir("/opt/vscode-insiders", 0777)
	err = copy.Copy("VSCode-linux-x64/", "/opt/vscode-insiders/")
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Create("/usr/share/applications/VSCode-Insiders.desktop")
	if err != nil {
		log.Fatal(err)
	}
	desktopfile := "[Desktop Entry]\nName=Visual Studio Code - Insiders\nComment=Code Editing. Redefined.\nExec=/opt/vscode-insiders/bin/code-insiders\nIcon=/opt/vscode-insiders/resources/app/resources/linux/code.png\nCategories=Utility;TextEditor;Development;IDE;\nTerminal=false\nType=Application\n\n[Desktop Action new-empty-window]\nName=New Empty Window\nExec=/usr/share/code-insiders/code-insiders --new-window %F\nIcon=/opt/vscode-insiders/resources/app/resources/linux/code.png"
	err = ioutil.WriteFile("/usr/share/applications/VSCode-Insiders.desktop", []byte(desktopfile), 0777)
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Installed VSCode Insiders")
	fmt.Println("Cleaning up...")
	err = os.Remove("code-insider.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("VSCode-linux-x64/")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Done")
	main()
}

func uninstallvscode() {
	fmt.Println("Uninstalling VSCode Insiders...")
	err := os.RemoveAll("/opt/vscode-insiders/")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("/usr/share/applications/VSCode-Insiders.desktop")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Uninstalled VSCode Insiders")
}

func updatemultimc() {
	fmt.Println("Downloading MultiMC Development...")
	resp, err := http.Get("https://files.multimc.org/downloads/mmc-develop-lin64.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("multimc-dev.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	color.Green("Downloaded MultiMC Development")
	fmt.Println("Unzipping...")
	exec.Command("bsdtar", "-xf", "multimc-dev.tar.gz").Run()
	color.Green("Unzipped MultiMC Development")
	fmt.Println("Installing MultiMC Development...")
	os.Mkdir("/opt/multimc", 0777)
	err = copy.Copy("MultiMC/", "/opt/multimc-dev/")
	if err != nil {
		log.Fatal(err)
	}
	resp, err = http.Get("https://files.multimc.org/branding/multimc.svg")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err = os.Create("/opt/multimc-dev/MultiMC.svg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	_, err = os.Create("/usr/share/applications/MultiMC-Development.desktop")
	if err != nil {
		log.Fatal(err)
	}
	desktopfile := "[Desktop Entry]\nName=MultiMC Development\nGenericName=Minecraft Launcher\nComment=Free, open source launcher and instance manager for Minecraft.\nType=Application\nTerminal=false\nExec=/opt/multimc-dev/MultiMC\nIcon=/opt/multimc-dev/MultiMC.svg\nCategories=Game\nKeywords=game;minecraft;"
	err = ioutil.WriteFile("/usr/share/applications/MultiMC-Development.desktop", []byte(desktopfile), 0777)
	if err != nil {
		log.Fatal(err)
	}
	os.Chmod("/opt/multimc-dev", 0777)
	color.Green("Installed MultiMC Development")
	fmt.Println("Cleaning up...")
	err = os.Remove("multimc-dev.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("MultiMC/")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Done")
	color.Yellow("Make sure Qt5 is installed\nArch: qt5-base\nOpenSuse: libqt5-qtbase\nFedora/RHEL: qt5-qtbase\nUbuntu/Debian: qt5-default")
	main()
}

func uninstallmultimc() {
	fmt.Println("Uninstalling MultiMC Development...")
	err := os.RemoveAll("/opt/multimc-dev/")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("/usr/share/applications/MultiMC-Development.desktop")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("Uninstalled MultiMC Development")
}
