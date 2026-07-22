package utils

import (
	"bytes"
	"errors"
	"html/template"
	"path"
	"path/filepath"
	"runtime"
)

func RenderHtml(html string, data any) (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to resolve template path")
	}
	tmpl, err := template.ParseFiles(path.Join(filepath.Dir(file), "../templates", html))

	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
