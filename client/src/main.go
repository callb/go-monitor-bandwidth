package main

import (
	"strconv"
	"github.com/andlabs/ui"
	"time"
	"./utils"
	"./bandwidth"
)

func main() {
	initBandwidthMonitoring()
}

func initBandwidthMonitoring() {
	networkInfoList := bandwidth.GetBandwidthUsageData()

	err := ui.Main(func() {
		box := ui.NewVerticalBox()
		window := ui.NewWindow("Network Stats", 200, 200, false)

		for _, info := range networkInfoList {
			interfaceLabel := ui.NewLabel("Interface: " + info.InterfaceName)
			receivedLabel := ui.NewLabel("")
			transmittedLabel := ui.NewLabel("")

			box.Append(interfaceLabel, false)
			box.Append(receivedLabel, false)
			box.Append(transmittedLabel, false)

			go startUpdateLabelsProcessForInterface(info.InterfaceName, receivedLabel, transmittedLabel)
		}

		window.SetChild(box)
		window.SetMargined(true)
		window.Show()

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
	})
	utils.CheckForError(err)
}

// Starts process of continually updating stats for a particular interface
func startUpdateLabelsProcessForInterface(interfaceName string, receivedLabel *ui.Label, transmittedLabel *ui.Label) {
	initialInfo := bandwidth.GetBandwidthInfoForInterface(interfaceName)
	previousReceivedBytes := initialInfo.ReceivedData["bytes"]
	previousTransmittedBytes := initialInfo.TransmittedData["bytes"]

	var batchDataToUpload []bandwidth.BandwidthInfo
	batchDataToUpload = bandwidth.AddBandwidthDataToUpload(batchDataToUpload, initialInfo)

	for {
		info := bandwidth.GetBandwidthInfoForInterface(interfaceName)

		receivedDeltaValue := info.ReceivedData["bytes"] - previousReceivedBytes
		transmittedDeltaValue := info.TransmittedData["bytes"] - previousTransmittedBytes

		receivedText := "Receiving: " + strconv.Itoa(receivedDeltaValue) + " bytes/sec"
		transmittedText := "Transmitting: " + strconv.Itoa(transmittedDeltaValue) + " bytes/sec\n"
		ui.QueueMain(func() {
			receivedLabel.SetText(receivedText)
			transmittedLabel.SetText(transmittedText)
		})

		batchDataToUpload = bandwidth.AddBandwidthDataToUpload(batchDataToUpload, info)

		previousReceivedBytes = info.ReceivedData["bytes"]
		previousTransmittedBytes = info.TransmittedData["bytes"]
		time.Sleep(time.Second)
	}
}
