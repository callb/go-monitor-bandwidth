package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"github.com/andlabs/ui"
	"time"
)

type bandwidthInfo struct {
	interfaceName string
	receivedData map[string] int
	transmittedData map[string] int
}

func main() {
	initBandwidthMonitoring()
}

func initBandwidthMonitoring() {
	networkInfoList := parseNetworkDataFromFile(readNetworkFile())

	err := ui.Main(func() {
		box := ui.NewVerticalBox()
		window := ui.NewWindow("Network Stats", 200, 200, false)

		for _, info := range networkInfoList {
			interfaceLabel := ui.NewLabel("Interface: " + info.interfaceName)
			receivedLabel := ui.NewLabel("")
			transmittedLabel := ui.NewLabel("")

			box.Append(interfaceLabel, false)
			box.Append(receivedLabel, false)
			box.Append(transmittedLabel, false)

			go startUpdateLabelsProcessForInterface(info.interfaceName, receivedLabel, transmittedLabel)
		}

		window.SetChild(box)
		window.SetMargined(true)
		window.Show()

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
	})
	checkForError(err)
}

// Starts process of continually updating stats for a particular interface
func startUpdateLabelsProcessForInterface(interfaceName string, receivedLabel *ui.Label, transmittedLabel *ui.Label) {
	initialInfo := getBandwidthInfoForInterface(interfaceName)
	previousReceivedBytes := initialInfo.receivedData["bytes"]
	previousTransmittedBytes := initialInfo.transmittedData["bytes"]
	for {
		info := getBandwidthInfoForInterface(interfaceName)

		receivedDeltaValue := info.receivedData["bytes"] - previousReceivedBytes
		transmittedDeltaValue := info.transmittedData["bytes"] - previousTransmittedBytes

		receivedText := "Receiving: " + strconv.Itoa(receivedDeltaValue) + " bytes/sec"
		transmittedText := "Transmitting: " + strconv.Itoa(transmittedDeltaValue) + " bytes/sec\n"
		ui.QueueMain(func() {
			receivedLabel.SetText(receivedText)
			transmittedLabel.SetText(transmittedText)
		})

		previousReceivedBytes = info.receivedData["bytes"]
		previousTransmittedBytes = info.transmittedData["bytes"]
		time.Sleep(time.Second)
	}
}

// Read in the network information file
func readNetworkFile() string {
	bytes, err := ioutil.ReadFile("/proc/net/dev")
	checkForError(err)
	return string(bytes)
}

// Parse the data from the network file into a bandwidthInfo struct
func parseNetworkDataFromFile(data string) []bandwidthInfo {
	var receiveCols, transmitCols []string
	lines := strings.Split(strings.TrimSpace(data), "\n")
	var allInfo []bandwidthInfo

	for i, line := range lines {
		// Skip the first line
		if i == 0 {
			continue
		}
		if strings.Contains(line, "|") {
			receiveCols, transmitCols = getColumnNamesFromRow(line)
		} else {
			bandwidthInfo := getBandwidthInfoFromRow(line, receiveCols, transmitCols)
			allInfo = append(allInfo, bandwidthInfo)
		}
	}
	return allInfo
}

// Get the column names for receiving and transmitting data
func getColumnNamesFromRow(line string)  ([]string, []string) {
	var receiveCols, transmitCols []string

	for i, colsSet := range strings.Split(line, "|") {
		// Skip first column which is part of "inter-face" string
		if i == 0 {
			continue
		}
		// First set of columns are for Receive data
		if i == 1 {
			receiveCols = strings.Fields(colsSet)
		}
		// Second set of columns are for Transmit data
		if i == 2 {
			transmitCols = strings.Fields(colsSet)
		}
	}
	return receiveCols, transmitCols
}

// Map the bandwidth data of a row to it's respective column,
// and then map to the respective interface
func getBandwidthInfoFromRow(line string, receiveCols []string, transmitCols []string) bandwidthInfo {
	bandwidthInfo := bandwidthInfo {
		receivedData: make(map[string] int),
		transmittedData: make(map[string] int),
	}

	for i, value := range strings.Fields(line) {
		// first data point in the row is the connection name
		if i == 0 {
			bandwidthInfo.interfaceName = strings.Trim(value, ":")
			continue
		}
		dataPoint, err := strconv.Atoi(value)
		checkForError(err)

		if i <= len(receiveCols) {
			col := receiveCols[i - 1]
			bandwidthInfo.receivedData[col] = dataPoint

		} else {
			col := transmitCols[i - len(receiveCols) - 1]
			bandwidthInfo.transmittedData[col] = dataPoint
		}
	}

	return bandwidthInfo
}

// Gets a particular interface from the list of bandwidth info
func getBandwidthInfoForInterface(interfaceName string) bandwidthInfo {
	networkInfoList := parseNetworkDataFromFile(readNetworkFile())
	for _, info := range networkInfoList {
		if info.interfaceName == interfaceName {
			return info
		}
	}
	return bandwidthInfo{}
}

// Check for an error, panic if one has occurred
func checkForError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
