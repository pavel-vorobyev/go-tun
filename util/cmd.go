package util

import (
	"bytes"
	"github.com/mbndr/figlet4go"
	"log"
	"os/exec"
)

var (
	blue   = "\033[34m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func RunCommand(command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if stderr.String() != "" {
		return stderr.String(), err
	}
	return stdout.String(), err
}

func PrintHelloNidde() {
	greenColor, _ := figlet4go.NewTrueColorFromHexString("32DE84")
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{greenColor}
	renderer := figlet4go.NewAsciiRender()
	renderStr, _ := renderer.RenderOpts("Hello, Niddle!", options)
	log.Println("\n" + renderStr)
}

func PrintMessage(msg string) {
	log.Println(blue + msg + reset)
}

func PrintWarning(msg string) {
	log.Println(yellow + msg + reset)
}
