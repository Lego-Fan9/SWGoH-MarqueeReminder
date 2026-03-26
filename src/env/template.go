package env

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

const DefaultTemplate = `
{
  "content": "<@&{{ .Role }}>",
  "embeds": [
    {
      "title": "Marquee Reminder",
      "description": "Reminder to do the {{ .NameKey }} marquee"
    }
  ]
}
`

const DefaultTemplateImg = `
{
  "content": "<@&{{ .Role }}>",
  "embeds": [
    {
      "title": "Marquee Reminder",
      "description": "Reminder to do the {{ .NameKey }} marquee",
      "thumbnail": {
        "url": "attachment://{{ .Filename }}"
      }
    }
  ]
}
`

var (
	Template    *template.Template
	ImgTemplate *template.Template
)

type MarqueeTemplateData struct {
	Role     string
	NameKey  string
	Filename string
}

var ErrTemplateUse = errors.New("template.go: unknown error using template")

func LoadTemplate() {
	LoadTemplateStd()
	LoadTemplateImg()
}

func LoadTemplateStd() {
	funcMap := template.FuncMap{
		"upper":   strings.ToUpper,
		"lower":   strings.ToLower,
		"default": defaultTmplFunc,
	}

	var tmplStr = DefaultTemplate
	if CUSTOM_FORMAT != "" {
		tmplStr = CUSTOM_FORMAT
	}

	var err error

	Template, err = template.New("marquee_webhook_post_template").
		Funcs(funcMap).
		Option("missingkey=zero").
		Parse(tmplStr)
	if err != nil {
		log.Fatalf("Failed to load template string \"%s\" reason: %v", tmplStr, err)
	}
}

func LoadTemplateImg() {
	funcMap := template.FuncMap{
		"upper":   strings.ToUpper,
		"lower":   strings.ToLower,
		"default": defaultTmplFunc,
	}

	var tmplStr = DefaultTemplateImg
	if CUSTOM_FORMAT_IMG != "" {
		tmplStr = CUSTOM_FORMAT_IMG
	}

	var err error

	ImgTemplate, err = template.New("marquee_webhook_post_template").
		Funcs(funcMap).
		Option("missingkey=zero").
		Parse(tmplStr)
	if err != nil {
		log.Fatalf("Failed to load image template string \"%s\" reason: %v", tmplStr, err)
	}
}

func defaultTmplFunc(v any, d string) string {
	if v == nil {
		return d
	}

	s := fmt.Sprint(v)
	if s == "" {
		return d
	}

	return s
}

func GetMarqueeDiscordPostTemplate(input MarqueeTemplateData) (string, error) {
	var buf bytes.Buffer

	err := Template.Execute(&buf, input)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateUse, err)
	}

	return buf.String(), nil
}

func GetMarqueeDiscordPostTemplateImg(input MarqueeTemplateData) (string, error) {
	var buf bytes.Buffer

	err := ImgTemplate.Execute(&buf, input)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateUse, err)
	}

	return buf.String(), nil
}
