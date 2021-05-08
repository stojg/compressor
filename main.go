package main

import (
	"fmt"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func myLogger(messageDomain string, verbosity vips.LogLevel, message string) {
	var messageLevelDescription string
	switch verbosity {
	case vips.LogLevelError:
		messageLevelDescription = "error"
	case vips.LogLevelCritical:
		messageLevelDescription = "critical"
	case vips.LogLevelWarning:
		messageLevelDescription = "warning"
	case vips.LogLevelMessage:
		messageLevelDescription = "message"
	case vips.LogLevelInfo:
		messageLevelDescription = "info"
	case vips.LogLevelDebug:
		messageLevelDescription = "debug"
	}
	log.Printf("[%v.%v] %v", messageDomain, messageLevelDescription, message)
}

func main() {
	vips.LoggingSettings(myLogger, vips.LogLevelMessage)
	vips.Startup(nil)
	defer vips.Shutdown()
	err := download("https://public.stojg.se/public/gM14KsgT.gif", "input.gif")
	checkError(err)

	err = sprite("input.gif")
	checkError(err)

	err = convert("input.gif")
	checkError(err)
}

func sprite(source string) error {
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := gif.DecodeAll(f)
	if err != nil {
		return err
	}

	frames := float32(len(img.Image))
	step := frames / 32.0

	fmt.Printf("%f, %f\n", frames, step)

	for i := float32(0.0); i < frames; i += step {
		fmt.Println(int(i))
	}

	return nil
}

func convert(source string) error {
	image1, err := vips.NewImageFromFile(source)
	if err != nil {
		return err
	}
	if err = image1.AutoRotate(); err != nil {
		return err
	}
	image1bytes, _, err := image1.ExportNative()
	if err != nil {
		return err
	}
	filetype := vips.ImageTypes[image1.Metadata().Format]
	return ioutil.WriteFile("output."+filetype, image1bytes, 0644)
}

func download(url, dest string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
