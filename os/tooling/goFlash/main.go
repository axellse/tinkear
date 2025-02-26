package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"go.bug.st/serial"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: goFlash <file> <serial-port>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	selectedPort := os.Args[2]

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", filePath)
		os.Exit(1)
	}

	tmpOutput := filepath.Join(os.TempDir(), "bundle.js")

	cmd := exec.Command("rollup", filePath, "--file", tmpOutput, "--format", "iife")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Failed to run rollup: %v\n", err)
		os.Exit(1)
	}

	bundledCode, err := os.ReadFile(tmpOutput)
	if err != nil {
		fmt.Printf("Error: Failed to read bundled output: %v\n", err)
		os.Exit(1)
	}
	//os.Remove(tmpOutput)

	encoded := base64.StdEncoding.EncodeToString(bundledCode)
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(selectedPort, mode)
	if err != nil {
		fmt.Printf("Error: Failed to open serial port: %v\n", err)
		os.Exit(1)
	}
	defer port.Close()

	_, err = port.Write([]byte(`serialInterface.installOs('tinkearOS via goFlash', atob('` + encoded + `'));\n`))
	if err != nil {
		fmt.Printf("Error: Failed to write to serial port: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully sent encoded data to serial port")
	fmt.Println(`serialInterface.installOs('tinkearOS via goFlash', atob('` + encoded + `'));\n`)
}
