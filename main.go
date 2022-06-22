package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	LOW_THRESHOLD      = 15
	STATUS_CHARGING    = "Charging"
	STATUS_DISCHARGING = "Discharging"
	STATUS_NOTCHARGING = "Not charging"
)

var statusRe = regexp.MustCompile(`^Battery\s[0-9]+:\s(Charging|Discharging|Not charging),\s([0-9]+)%`)

func main() {

	var out bytes.Buffer
	cmd := exec.Command("acpi", "-b")

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal("failed to poll acpi")
	}
	matches := statusRe.FindStringSubmatch(out.String())

	status := matches[1]
	percent := matches[2]

	percentI, err := strconv.Atoi(percent)
	if err != nil {
		log.Fatal("could not parse battery percent")
	}

	if percentI < LOW_THRESHOLD && (status == STATUS_DISCHARGING || status == STATUS_NOTCHARGING) {
		fmt.Println("sending notification")
		notifyCmd := exec.Command("notify-send", "-u", "critical", "Battery", "LOW BATTERY")
		notifyCmd.Run()
	}
}
