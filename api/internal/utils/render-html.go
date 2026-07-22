package utils

import (
	"bytes"
	"errors"
	"html/template"
	"path"
	"path/filepath"
	"runtime"
)

type HtmlParams struct {
	T    string
	Data any
}

func RenderHtml(h HtmlParams) (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to resolve template path")
	}
	tmpl, err := template.ParseFiles(path.Join(filepath.Dir(file), "../templates", h.T))

	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, h.Data); err != nil {
		return "", err
	}

	return body.String(), nil
}
