package main

import (
  "bytes"
  "fmt"
  "os/exec"
  "strings"
)

func getCurrentDefaultSink() (string, error) {
	cmd := exec.Command("pactl", "info")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Default Sink:") {
			fields := strings.Fields(line)
			return fields[len(fields)-1], nil
		}
	}
	return "", fmt.Errorf("default sink not found")
}

func cycleDefaultSink(sinks []string, currentSink string) (string, error) {
	for i, sink := range sinks {
		if sink == currentSink {
			// Get the next sink, or circle back to the first if at the end of the list
			nextSink := sinks[(i+1)%len(sinks)]
			return nextSink, setDefaultSink(nextSink)
		}
	}
	return "", fmt.Errorf("current default sink not found in the list")
}

func setDefaultSink(sinkName string) error {
	cmd := exec.Command("pactl", "set-default-sink", sinkName)
	return cmd.Run()
}

func moveAllInputsToSink(sinkName string) error {
	cmd := exec.Command("pactl", "list", "short", "sink-inputs")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			moveCmd := exec.Command("pactl", "move-sink-input", fields[0], sinkName)
			err := moveCmd.Run()
			if err != nil {
				return err
			}
		}
	}
	return nil
}


func main() {

  sinks := []string{"alsa_output.usb-SteelSeries_Arctis_7_-00.analog-stereo", "alsa_output.pci-0000_0a_00.1.hdmi-stereo"}
	currentSink, err := getCurrentDefaultSink()
	if err != nil {
		return
	}

  var nextSink string
  if currentSink == sinks[0]{
    nextSink = sinks[1]
  }else{
    nextSink = sinks[0]
  }

  setDefaultSink(nextSink)

	moveAllInputsToSink(nextSink)
}
