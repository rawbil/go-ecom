package email

import (
	"bytes"
	"html/template"
)

func RenderHtml(html string, data any) (string, error) {
	tmpl, err := template.New("welcome").ParseFiles(html)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
