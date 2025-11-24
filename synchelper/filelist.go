package synchelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/schema"
	"encoding/json"
	"fmt"
	"net/http"
)

// Get a list of files in MinIO from master
func getFileListFromMaster() (result []schema.FileMINIO, err error) {
	// request
	url := config.MasterUrl + "/server/internal/filelist"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200")
	}

	// decode
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
