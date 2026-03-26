package assetextractor

import (
	"github.com/Lego-Fan9/MarqueeReminder/env"
	"github.com/Lego-Fan9/MarqueeReminder/httpclient"
	"strings"
	log "github.com/sirupsen/logrus"
	"bytes"
	"net/http"
	"errors"
	"io"
	"encoding/json"
	"fmt"
)

var (
	BadHttpStatus = errors.New("bad http status code")
)

func GetEventTex(iconName string) ([]byte, bool) {
	if !strings.HasPrefix(iconName, "tex.charui") {
		log.Warnf("Could not download %s (bad prefix)", iconName)

		return nil, false
	}

	var assetName = strings.TrimPrefix(iconName, "tex.")

	assetVersion, err := GetAssetVersion()
	if err != nil {
		return nil, false
	}

	url := fmt.Sprintf("%s/Asset/single?version=%d&assetName=%s", env.SWGOH_AE_URL, assetVersion, assetName)

	resp, err := httpclient.Get(url)
	if err != nil {
		log.Errorf("Falied to get asset from %s: %v", url, err)

		return nil, false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("io read fail: %v", err)

		return nil, false
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnf("Bad status code getting asset: %s: %s", resp.Status, string(body))

		return nil, false
	}

	return body, true
}

func GetAssetVersion() (int, error) {
	resp, err := httpclient.Post(env.COMLINK_URL + "/metadata", "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		log.Errorf("Failed to make metadata call: %v", err)

		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("io read fail: %v", err)

		return 0, BadHttpStatus
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Bad status code getting metadata: %s: %s", resp.Status, string(body))

		return 0, BadHttpStatus
	}

	var metadata assetVersionMetadata

	err = json.Unmarshal(body, &metadata)
	if err != nil {
		log.Errorf("Failed to unmarshal metadata: %v", err)

		return 0, err
	}

	return metadata.AssetVersion, nil
}

type assetVersionMetadata struct {
	AssetVersion int `json:"assetVersion"`
}