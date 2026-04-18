package comlinkevent

import (
	"strconv"
	"time"
)

//nolint:revive
type ComlinkEventResponse struct {
	GameEvent []ComlinkEvent `json:"gameEvent"`
}

type ComlinkEvent struct {
	Instance []struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"instance"`
	ID                string `json:"id"`
	NameKey           string `json:"nameKey"`
	SummaryKey        string `json:"summaryKey"`
	DescKey           string `json:"descKey"`
	Image             string `json:"image"`
	MarqueeUnitBaseID string `json:"marqueeUnitBaseId"`
	StartTime         int64
	EndTime           int64
}

func (c *ComlinkEvent) FixTimes(instanceCount int) error {
	var err error

	c.StartTime, err = strconv.ParseInt(c.Instance[instanceCount].StartTime, 10, 64)
	if err != nil {
		return err
	}

	c.EndTime, err = strconv.ParseInt(c.Instance[instanceCount].EndTime, 10, 64)
	if err != nil {
		return err
	}

	c.StartTime /= 1000
	c.EndTime /= 1000

	return nil
}

func GetTestMarquee() ComlinkEvent {
	now := time.Now()

	return ComlinkEvent{
		Instance: []struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}{
			{
				StartTime: strconv.FormatInt(now.Add(-24*time.Hour).UnixMilli(), 10),
				EndTime:   strconv.FormatInt(now.Add(24*time.Hour).UnixMilli(), 10),
			},
		},
		ID:                "EVENT_MARQUEE_KXSECURITYDROID",
		NameKey:           "EVENT_MARQUEE_KXSECURITYDROID_NAME",
		SummaryKey:        "EVENT_MARQUEE_KXSECURITYDROID_SUMMARY",
		DescKey:           "EVENT_MARQUEE_GENERIC_DAILY_DESC",
		Image:             "tex.events_kxsecuritydroid",
		MarqueeUnitBaseID: "KXSECURITYDROID",
	}
}
