package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go"
)

func main() {
	UploadTSDBBlocks("data")
	createMetricProfile()
	SendBulkAPIRequest()
}

func SendBulkAPIRequest() {
	type Include struct {
		Namespace  []string          `json:"namespace"`
		Workload   []string          `json:"workload"`
		Containers []string          `json:"containers"`
		Labels     map[string]string `json:"labels"`
	}
	type Exclude struct {
		Namespace  []string          `json:"namespace"`
		Workload   []string          `json:"workload"`
		Containers []string          `json:"containers"`
		Labels     map[string]string `json:"labels"`
	}

	type Filter struct {
		Include Include `json:"include"`
		Exclude Exclude `json:"exclude"`
	}

	type TimeRange struct {
		// Start string `json:"start"`
		// End   string `json:"end"`
	}

	type RequestPayload struct {
		Filter     Filter    `json:"filter"`
		TimeRange  TimeRange `json:"time_range"`
		Datasource string    `json:"datasource"`
	}

	payload := RequestPayload{
		Filter: Filter{
			Include: Include{
				Namespace:  []string{},
				Workload:   []string{},
				Containers: []string{},
				Labels:     map[string]string{},
			},
			Exclude: Exclude{
				Namespace:  []string{},
				Workload:   []string{},
				Containers: []string{},
				Labels:     map[string]string{},
			},
		},
		TimeRange:  TimeRange{},
		Datasource: "thanos",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println("Bulk API Request: ", string(jsonData))

	url := "http://127.0.0.1:8080/bulk"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	type ResponseData struct {
		JobID string `json:"job_id"`
	}

	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		return
	}
	fmt.Println("Received ID:", responseData.JobID)
}

func UploadTSDBBlocks(dataDir string) {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		log.Fatalln(err)
	}

	// List all directories in dataDir
	blocks, err := os.ReadDir(dataDir)
	if err != nil {
		log.Fatalln("Error reading data directory:", err)
	}

	for _, block := range blocks {
		if !block.IsDir() {
			continue
		}
		blockID := block.Name()
		blockDir := filepath.Join(dataDir, blockID)
		err = filepath.Walk(blockDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relPath, err := filepath.Rel(blockDir, filePath)
				if err != nil {
					return err
				}
				objectName := fmt.Sprintf("%s/%s", blockID, relPath)
				_, err = client.FPutObject("rosocp-tsdb", objectName, filePath, minio.PutObjectOptions{})
				if err != nil {
					log.Println("Error uploading file:", err)
					return err
				}
				fmt.Println("Uploaded:", filePath)
			}
			return nil
		})
		if err != nil {
			log.Printf("Error processing block %s: %v\n", blockID, err)
		} else {
			fmt.Printf("Block %s uploaded successfully!\n", blockID)
		}
	}
}

func checkURL(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode > 0
}

func createMetricProfile() {
	url := "http://127.0.0.1:8080/createMetricProfile"
	jsonFile := "resource_optimization_openshift_metric_profile.json"

	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	for {
		up := checkURL(url)
		if !up {
			fmt.Println("Metric endpoint is not reachable, retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Metric Profile Creation Status:", resp.Status)
}
