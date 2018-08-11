package bandwidth

import (
	"io/ioutil"
	"strings"
	"strconv"
	"../utils"
	"time"
)

type BandwidthInfo struct {
	InterfaceName string
	ReceivedData map[string] int
	TransmittedData map[string] int
	TimeRecorded time.Time
}

// Read in the network information file
func GetBandwidthUsageData() []BandwidthInfo {
	bytes, err := ioutil.ReadFile("/proc/net/dev")
	utils.CheckForError(err)
	return parseNetworkDataFromFile(string(bytes))
}

// Parse the data from the network file into a BandwidthInfo struct
func parseNetworkDataFromFile(data string) []BandwidthInfo {
	var receiveCols, transmitCols []string
	lines := strings.Split(strings.TrimSpace(data), "\n")
	var allInfo []BandwidthInfo

	for i, line := range lines {
		// Skip the first line
		if i == 0 {
			continue
		}
		if strings.Contains(line, "|") {
			receiveCols, transmitCols = getColumnNamesFromRow(line)
		} else {
			BandwidthInfo := getBandwidthInfoFromRow(line, receiveCols, transmitCols)
			allInfo = append(allInfo, BandwidthInfo)
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

// Map the Bandwidth data of a row to it's respective column,
// and then map to the respective interface
func getBandwidthInfoFromRow(line string, receiveCols []string, transmitCols []string) BandwidthInfo {
	BandwidthInfo := BandwidthInfo {
		ReceivedData: make(map[string] int),
		TransmittedData: make(map[string] int),
	}

	for i, value := range strings.Fields(line) {
		// first data point in the row is the connection name
		if i == 0 {
			BandwidthInfo.InterfaceName = strings.Trim(value, ":")
			continue
		}
		dataPoint, err := strconv.Atoi(value)
		utils.CheckForError(err)

		if i <= len(receiveCols) {
			col := receiveCols[i - 1]
			BandwidthInfo.ReceivedData[col] = dataPoint

		} else {
			col := transmitCols[i - len(receiveCols) - 1]
			BandwidthInfo.TransmittedData[col] = dataPoint
		}
	}

	return BandwidthInfo
}

// Gets a particular interface from the list of Bandwidth info
func GetBandwidthInfoForInterface(interfaceName string) BandwidthInfo {
	networkInfoList := GetBandwidthUsageData()
	for _, info := range networkInfoList {
		if info.InterfaceName == interfaceName {
			return info
		}
	}
	return BandwidthInfo{}
}


