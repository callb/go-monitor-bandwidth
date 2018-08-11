package bandwidth

import (
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"fmt"
)

var (
	MaxUploadSize = 5
	uploadUrl = "http://localhost:8080/bandwidth/upload"
)

// Add the data to the batch and upload if max upload size is reached
func AddBandwidthDataToUpload(batchData []BandwidthInfo, dataToAdd BandwidthInfo) []BandwidthInfo {
	if MaxUploadSize == 0 {
		return batchData
	}

	batchData = append(batchData, dataToAdd)

	if len(batchData) == MaxUploadSize {
		upload(batchData)
		batchData = nil
	}

	return batchData
}

func upload(batchData []BandwidthInfo) {
	var uploadClient = &http.Client{
		Timeout: time.Second * 10,
	}
	dataAsJson, _ := json.Marshal(batchData)
	byteReader := bytes.NewReader(dataAsJson)
	response, err := uploadClient.Post(uploadUrl, "application/json", byteReader)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)

}
