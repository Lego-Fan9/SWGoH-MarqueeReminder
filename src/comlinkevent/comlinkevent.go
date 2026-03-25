package comlinkevent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Lego-Fan9/MarqueeReminder/env"
	"github.com/Lego-Fan9/MarqueeReminder/httpclient"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFailedToGetEvents = errors.New("failed to get events from comlink")
	ErrFailedToGetLoc    = errors.New("failed to get localization from repo")
	ErrFailedToGetUnits  = errors.New("failed to get units from repo")
	ErrBadStatus         = errors.New("failed to get with status code")
)

// TODO: Make this check status code
func GetActiveMarquees() ([]ComlinkEvent, error) {
	allEvents, err := GetEvents()
	if err != nil {
		return nil, err
	}

	var activeMarquees []ComlinkEvent

	for _, event := range allEvents.GameEvent {
		if event.MarqueeUnitBaseID == "" {
			continue
		}

		if len(event.Instance) == 0 {
			continue
		}

		err := event.FixTimes(0)
		if err != nil {
			log.Warnf("Failed to parse a timestamp: %v", err)

			continue
		}

		now := time.Now().Unix()
		if now >= event.StartTime && now <= event.EndTime {
			activeMarquees = append(activeMarquees, event)
		}
	}

	return activeMarquees, nil
}

func GetEvents() (ComlinkEventResponse, error) {
	resp, err := httpclient.Post(env.COMLINK_URL+"/getEvents", "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return ComlinkEventResponse{}, fmt.Errorf("%w: %w", ErrFailedToGetEvents, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ComlinkEventResponse{}, fmt.Errorf("%w: %w", ErrFailedToGetEvents, err)
	}

	var response ComlinkEventResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return ComlinkEventResponse{}, fmt.Errorf("%w: %w", ErrFailedToGetEvents, err)
	}

	return response, nil
}

type locInternal struct {
	Version string              `json:"version"`
	Data    ComlinkLocalization `json:"data"`
}

type ComlinkLocalization map[string]string

func GetLocalization() (ComlinkLocalization, error) {
	var url = "https://raw.githubusercontent.com/swgoh-utils/gamedata/refs/heads/main/Loc_ENG_US.txt.json"

	resp, err := httpclient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetLoc, err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s: %s", ErrBadStatus, resp.Status, string(body))
	}

	var response locInternal

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetLoc, err)
	}

	return response.Data, nil
}

type unitsInternal struct {
	Version string        `json:"version"`
	Data    []ComlinkUnit `json:"data"`
}

type ComlinkUnit struct {
	BaseID  string `json:"baseId"`
	NameKey string `json:"nameKey"`
}

func GetUnits() ([]ComlinkUnit, error) {
	var url = "https://raw.githubusercontent.com/swgoh-utils/gamedata/refs/heads/main/units.json"

	resp, err := httpclient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUnits, err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s: %s", ErrBadStatus, resp.Status, string(body))
	}

	var response unitsInternal

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUnits, err)
	}

	return response.Data, nil
}
