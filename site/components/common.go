package components

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type BaseHTMLProps struct {
	ID      string
	Classes templ.CSSClasses
	Styles  templ.SafeCSS
}

type BaseFormProps struct {
	BaseHTMLProps
	Name []string
}

func (p BaseHTMLProps) GetBaseHTMLProps() templ.Attributes {
	return templ.Attributes{
		"id":    p.ID,
		"class": p.Classes,
		"style": p.Styles,
	}
}

func GenerateNestedJsonFormPropName(name []string) (s string) {
	for i, n := range name {
		if i == 0 {
			s = n
			continue
		}

		if i == len(name)-1 && len(n) >= 1 && n[:1] == ":" {
			s += n
			break
		}

		s += fmt.Sprintf("[%s]", n)
	}

	return
}

func RenderMarkdown(markdown string) templ.Component {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	if err := md.Convert([]byte(markdown), &buf); err != nil {
		logrus.Errorf("failed to convert markdown to HTML: %v", err)
	}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, buf.String())
		return
	})
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func Append[T any](array []T, vals ...T) []T {
	return append(array, vals...)
}
