package main

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Lego-Fan9/MarqueeReminder/assetextractor"
	"github.com/Lego-Fan9/MarqueeReminder/comlinkevent"
	"github.com/Lego-Fan9/MarqueeReminder/env"
	"github.com/Lego-Fan9/MarqueeReminder/httpclient"
	log "github.com/sirupsen/logrus"
)

var (
	ErrBadDiscordStatus = errors.New("bad status code from discord")
)

func PostMarqueeDiscord(input comlinkevent.ComlinkEvent, localization comlinkevent.ComlinkLocalization, units []comlinkevent.ComlinkUnit) error {
	var unit comlinkevent.ComlinkUnit

	var foundUnit bool

	for _, v := range units {
		if v.BaseID == input.MarqueeUnitBaseID {
			unit = v
			foundUnit = true

			break
		}
	}

	if !foundUnit {
		log.Warnf("Failed to find baseid for %s", unit.BaseID)
	}

	nameKeyCorrected, ok := localization[unit.NameKey]
	if !ok {
		log.Warnf("Failed to find baseId for %s", unit.BaseID)
		nameKeyCorrected = unit.NameKey
	}

	log.Infof("Found unit id: %s, nameKey: %s, localized: %s", unit.BaseID, unit.NameKey, nameKeyCorrected)

	imgData, found := assetextractor.GetEventTex(unit.EraThumbnailName)
	if found {
		log.Infof("Found asset %s using image post", unit.EraThumbnailName)

		return PostMarqueeDiscordImg(nameKeyCorrected, imgData)
	}

	data, err := env.GetMarqueeDiscordPostTemplate(env.MarqueeTemplateData{
		Role:    env.PING_ROLE,
		NameKey: nameKeyCorrected,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, env.DISCORD_WEBHOOK, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.Discord(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Errorf("Bad status from discord: %s: %s", resp.Status, string(body))

		return ErrBadDiscordStatus
	}

	_, _ = io.ReadAll(resp.Body)

	return nil
}

func PostMarqueeDiscordImg(nameKeyCorrected string, image []byte) error {
	var imageName = "marquee.png"

	data, err := env.GetMarqueeDiscordPostTemplateImg(env.MarqueeTemplateData{
		Role:     env.PING_ROLE,
		NameKey:  nameKeyCorrected,
		Filename: imageName,
	})
	if err != nil {
		return err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("payload_json", data); err != nil {
		return err
	}

	part, err := writer.CreateFormFile("files", imageName)
	if err != nil {
		return err
	}

	if _, err := part.Write(image); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	log.Infof("payload_json: %s", data)

	req, err := http.NewRequest(http.MethodPost, env.DISCORD_WEBHOOK, &body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := httpclient.Discord(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Errorf("Bad status from discord: %s: %s", resp.Status, string(body))
		return ErrBadDiscordStatus
	}

	_, _ = io.ReadAll(resp.Body)

	return nil
}
