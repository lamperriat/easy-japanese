package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Metadata struct {
	LatestID int `json:"latestID"`
}

var (
	metadataLock = sync.Mutex{}
	metadataPath = filepath.Join("data", "metadata.json")
)

func GetNextID() (int, error) {
	metadataLock.Lock()
	defer metadataLock.Unlock()

	metadata := Metadata{LatestID: 0}
	file, err := os.ReadFile(metadataPath)
	if err == nil {
		if err := json.Unmarshal(file, &metadata); err != nil {
			return 0, err
		}
	}

	metadata.LatestID++

	newData, err := json.Marshal(metadata)
	if err != nil {
		return 0, err
	}

    // Create temp file and rename for atomic operation
    tempFile := metadataPath + ".tmp"
    if err := os.WriteFile(tempFile, newData, 0644); err != nil {
        return 0, err
    }
    if err := os.Rename(tempFile, metadataPath); err != nil {
        return 0, err
    }

    return metadata.LatestID, nil
}