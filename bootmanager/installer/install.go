//This is the first time ive used golang, so expect the biggest piece of hot garbage you'l ever see

package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/artdarek/go-unzip"
	"go.bug.st/serial"
)

func installBootManager(port string, config string) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	log.Println("Connecting to device...")

	device, _serialErr := serial.Open(port, mode)
	if _serialErr != nil {
		log.Fatal(_serialErr)
	}

	log.Println("Downloading boot manager...")
	_bootManagerReq, err := http.Get("https://raw.githubusercontent.com/axellse/tinkear/refs/heads/main/bootmanager/main.js")
	if err != nil {
		log.Fatal("Could not download boot manager.")
	}

	_bootManager, err := io.ReadAll(_bootManagerReq.Body)
	if err != nil {
		log.Fatal("Could not download boot manager.")
		os.Exit(1)
	}

	bootManager := strings.Replace(string(_bootManager), "$HWConfig_ReplacedByInstaller$", config, 1)
	log.Println("Downloaded.")

	log.Println("Installing boot manager...")

	device.Write([]byte("let bootManagerInstallationBuffer = '';\n"))
	time.Sleep(100 * time.Millisecond)
	encodedBootManager := base64.StdEncoding.EncodeToString([]byte(bootManager))

	totalChunks := float64(len(encodedBootManager)) / float64(40)
	log.Println("Final boot manager size: " + strconv.Itoa(len(encodedBootManager)) + " (" + strconv.Itoa(int(totalChunks)) + " chunks)")
	for i := 0; i < len(encodedBootManager); i += 40 {
		log.Println("Sending chunk " + strconv.Itoa(i/40) + " (" + strconv.Itoa(int(math.Round((float64(i/40)/totalChunks)*100))) + "%)")
		end := i + 40
		if end > len(encodedBootManager) {
			end = len(encodedBootManager)
		}
		device.Write([]byte("bootManagerInstallationBuffer += atob('" + encodedBootManager[i:end] + "');\n"))
		time.Sleep(50 * time.Millisecond)
	}
	log.Println("Configuring boot manager to run at boot...")
	device.Write([]byte("E.setBootCode(bootManagerInstallationBuffer, true);\n"))
	time.Sleep(1000 * time.Millisecond)
	log.Println("Closing connection...")
	device.Close()
	log.Println("Update complete.")
}

func installProcess(port string, config string) {
	log.Println("Performing some pre-install checks...")
	cmd := exec.Command("python", "--version")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Python is not installed. Please install Python and try again.")
	}
	log.Println("Everything OK, proceeding.")

	tempDir, err := os.MkdirTemp("", "esptool")
	if err != nil {
		log.Fatal("Could not create temporary directory.")
	}
	defer os.RemoveAll(tempDir)

	log.Println("Downloading esptool...")
	esptoolURL := "https://github.com/espressif/esptool/archive/refs/heads/master.zip"
	resp, err := http.Get(esptoolURL)
	if err != nil {
		log.Fatal("Could not download esptool.")
	}
	defer resp.Body.Close()

	zipPath := tempDir + "/esptool.zip"
	out, err := os.Create(zipPath)
	if err != nil {
		log.Fatal("Could not create zip file.")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal("Could not write to zip file.")
	}

	log.Println("Unzipping esptool...")
	uz := unzip.New(zipPath, tempDir)

	unzipError := uz.Extract()
	if unzipError != nil {
		log.Fatal("Could not unzip esptool.")
	}

	log.Println("Installing esptool dependencies")
	cmd = exec.Command("pip", "install", "esptool")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal("Could not install esptool via pip.")
	}

	log.Println("Downloading Espruino firmware...")
	firmwareURL := "https://www.espruino.com/binaries/espruino_2v25_esp8266_4mb_combined_4096.bin"
	firmwareResp, err := http.Get(firmwareURL)
	if err != nil {
		log.Fatal("Could not download Espruino firmware.")
	}
	defer firmwareResp.Body.Close()

	firmwarePath := tempDir + "/espruino_firmware.bin"
	firmwareOut, err := os.Create(firmwarePath)
	if err != nil {
		log.Fatal("Could not create firmware file.")
	}
	defer firmwareOut.Close()

	_, err = io.Copy(firmwareOut, firmwareResp.Body)
	if err != nil {
		log.Fatal("Could not write to firmware file.")
	}
	log.Println("Clearing device...")

	esptoolPath := tempDir + "/esptool-master/esptool.py"
	cmd = exec.Command("python", esptoolPath, "--port", port, "--baud", "460800", "write_flash", "--flash_mode", "dio", "--flash_size", "4MB", "0x00000", firmwarePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("──────────────────────────────────────────────────────────")
	err = cmd.Run()
	fmt.Println("──────────────────────────────────────────────────────────")

	if err != nil {
		log.Fatal("Could not clear device.")
	}

	log.Println("Flashing Espruino firmware to the device...")
	cmd = exec.Command("python", esptoolPath, "--port", port, "--baud", "460800", "write_flash", "--flash_mode", "dio", "--flash_size", "4MB", "0x00000", firmwarePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("──────────────────────────────────────────────────────────")
	err = cmd.Run()
	fmt.Println("──────────────────────────────────────────────────────────")
	if err != nil {
		log.Fatal("Could not flash Espruino firmware to the device.")
	}
	log.Println("Firmware flashed successfully.")
	log.Println("Waiting 20 sec for firmware to initalize.")
	time.Sleep(20000)
	installBootManager(port, config)
	fmt.Println("──────────────────────────────────────────────────────────")
	fmt.Println("\x1b[32mThe tinkearOS boot manager was successfully installed!\nVisit https://axell.me/getTiOs to get the latest version of tinkearOS running on your device!\x1b[0m")
}

func update() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		fmt.Println("❌ No serial ports found! Please make sure your device is connected and that the driver works properly.")
		os.Exit(1)
	}

	fmt.Println("🔌 Available ports:")
	for i, port := range ports {
		fmt.Println("" + strconv.Itoa(i) + ": " + port)
	}
	fmt.Println("\nSelect a port by typing it's corresponding number above, and pressing enter:")

	reader := bufio.NewReader(os.Stdin)
	_portIndex, _ := reader.ReadString('\n')

	index, err := strconv.Atoi(strings.TrimSpace(_portIndex))
	if err != nil {
		fmt.Println("Please input a number.")
		os.Exit(1)
	} else if int(index) > len(ports)-1 {
		fmt.Println("That's not a valid list item")
		os.Exit(1)
	}
	var selectedPort = ports[int(index)]
	fmt.Print("\033[1A")
	fmt.Println(strings.TrimSpace(_portIndex) + " ✅ (" + selectedPort + ")" + "\n")

	fmt.Println("📼 Choose your tinkear model:")

	_configsRequest, err := http.Get("https://raw.githubusercontent.com/axellse/tinkear/refs/heads/main/bootmanager/configs.json")
	if err != nil {
		fmt.Println("❌ Could not fetch hardware configurations. Please make sure you have a working internet connection.")
		os.Exit(1)
	}

	_configs, err := io.ReadAll(_configsRequest.Body)
	if err != nil {
		fmt.Println("❌ Could not fetch hardware configurations. Please make sure you have a working internet connection.")
		os.Exit(1)
	}

	var configs []string
	jsonErr := json.Unmarshal(_configs, &configs)

	if jsonErr != nil {
		fmt.Println("❌ Could not decode hardware configurations.")
	}

	for i, config := range configs {
		fmt.Println("" + strconv.Itoa(i) + ": " + config)
	}
	fmt.Println("\nChoose a model by typing it's corresponding number above, and pressing enter:")
	_configIndex, _ := reader.ReadString('\n')

	configIndex, err := strconv.Atoi(strings.TrimSpace(_configIndex))
	if err != nil {
		fmt.Println("Please input a number.")
		os.Exit(1)
	} else if int(configIndex) > len(configs)-1 {
		fmt.Println("That's not a valid list item")
		os.Exit(1)
	}
	var selectedConfig = configs[int(configIndex)]
	fmt.Print("\033[1A")
	fmt.Println(strings.TrimSpace(_configIndex) + " ✅ (" + selectedConfig + ")" + "\n")

	fmt.Println("Would you like to continue updating the boot manager? (y/N)")

	_installConfirmation, _ := reader.ReadByte()
	if _installConfirmation == []byte("y")[0] || _installConfirmation == []byte("Y")[0] {
		fmt.Println("\n\x1b[32mAlright, continuing with port " + selectedPort + "\x1b[0m")
		fmt.Println("\n\x1b[33mPlease do not unplug the device during the process.\x1b[0m")
		fmt.Println("──────────────────────────────────────────────────────────")
	} else {
		fmt.Println("Cancelled.")
	}

	//Actual install
	installBootManager(selectedPort, selectedConfig)
	fmt.Println("──────────────────────────────────────────────────────────")
	fmt.Println("\x1b[32mThe tinkearOS boot manager was successfully updated/reinstalled!\x1b[0m")
}

func installCommand() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		fmt.Println("❌ No serial ports found! Please make sure your device is connected and that the driver works properly.")
		os.Exit(1)
	}

	fmt.Println("🔌 Available ports:")
	for i, port := range ports {
		fmt.Println("" + strconv.Itoa(i) + ": " + port)
	}
	fmt.Println("\nSelect a port by typing it's corresponding number above, and pressing enter:")

	reader := bufio.NewReader(os.Stdin)
	_portIndex, _ := reader.ReadString('\n')

	index, err := strconv.Atoi(strings.TrimSpace(_portIndex))
	if err != nil {
		fmt.Println("Please input a number.")
		os.Exit(1)
	} else if int(index) > len(ports)-1 {
		fmt.Println("That's not a valid list item")
		os.Exit(1)
	}
	var selectedPort = ports[int(index)]
	fmt.Print("\033[1A")
	fmt.Println(strings.TrimSpace(_portIndex) + " ✅ (" + selectedPort + ")" + "\n")

	fmt.Println("📼 Choose your tinkear model:")

	_configsRequest, err := http.Get("https://raw.githubusercontent.com/axellse/tinkear/refs/heads/main/bootmanager/configs.json")
	if err != nil {
		fmt.Println("❌ Could not fetch hardware configurations. Please make sure you have a working internet connection.")
		os.Exit(1)
	}

	_configs, err := io.ReadAll(_configsRequest.Body)
	if err != nil {
		fmt.Println("❌ Could not fetch hardware configurations. Please make sure you have a working internet connection.")
		os.Exit(1)
	}

	var configs []string
	jsonErr := json.Unmarshal(_configs, &configs)

	if jsonErr != nil {
		fmt.Println("❌ Could not decode hardware configurations.")
	}

	for i, config := range configs {
		fmt.Println("" + strconv.Itoa(i) + ": " + config)
	}
	fmt.Println("\nChoose a model by typing it's corresponding number above, and pressing enter:")
	_configIndex, _ := reader.ReadString('\n')

	configIndex, err := strconv.Atoi(strings.TrimSpace(_configIndex))
	if err != nil {
		fmt.Println("Please input a number.")
		os.Exit(1)
	} else if int(configIndex) > len(configs)-1 {
		fmt.Println("That's not a valid list item")
		os.Exit(1)
	}
	var selectedConfig = configs[int(configIndex)]
	fmt.Print("\033[1A")
	fmt.Println(strings.TrimSpace(_configIndex) + " ✅ (" + selectedConfig + ")" + "\n")

	fmt.Println("\x1b[31mThis will erase all data on the device. Do you want to continue? (y/N)\x1b[0m")

	_installConfirmation, _ := reader.ReadByte()
	if _installConfirmation == []byte("y")[0] || _installConfirmation == []byte("Y")[0] {
		fmt.Println("\n\x1b[32mAlright, continuing with port " + selectedPort + "\x1b[0m")
		fmt.Println("\n\x1b[33mPlease do not unplug the device during the process.\x1b[0m")
		fmt.Println("──────────────────────────────────────────────────────────")
		installProcess(selectedPort, selectedConfig)
	} else {
		fmt.Println("Cancelled.")
	}
}

func main() {
	fmt.Println("👋 Welcome to the tinkearOS boot manager installation utility.\n")
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command. Type 'help' to view all commands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		installCommand()
	case "update":
		update()
	case "help":
		fmt.Printf(
			`
			Available commands:

			* install - installs the boot manager and the js runtime (espruino) onto an empty device.
			* update - updates the boot manager.
			`)
	default:
		fmt.Println("Invalid command. Type 'help' to view all commands.")
	}
}
