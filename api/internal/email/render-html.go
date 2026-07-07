package email

import (
	"bytes"
	"html/template"
	"io/fs"
)

func RenderHtml(html fs.FS, data any) (string, error) {
	// tmpl, err := template.New("welcome").ParseFiles(html)
	tmpl, err := template.New("welcome").ParseFS(html, "welcome.html")
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
