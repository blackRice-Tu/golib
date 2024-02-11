package xcommon

import (
	"bytes"
	"strings"
	"text/template"
)

// copy from Python string lib
const (
	Whitespace        = " \t\n\r\v\f"
	AsciiLowercase    = "abcdefghijklmnopqrstuvwxyz"
	AsciiUppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLetters      = AsciiLowercase + AsciiUppercase
	Digits            = "0123456789"
	Hexdigits         = Digits + "abcdef" + "ABCDEF"
	Octdigits         = "01234567"
	Punctuation       = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	PunctuationSimple = "!@#$%^&*()"
	Printable         = Digits + AsciiLetters + Punctuation + Whitespace
)

func TrimSpaceAndSplitString(s string, sep string) []string {
	resultList := make([]string, 0)
	s = strings.TrimSpace(s)
	if s != "" {
		for _, subString := range strings.Split(s, sep) {
			result := strings.TrimSpace(subString)
			if result != "" {
				resultList = append(resultList, result)
			}
		}
	}
	return resultList
}

func SplitUrl(s string, sep string) []string {
	urlList := make([]string, 0)
	for _, url := range TrimSpaceAndSplitString(s, sep) {
		if strings.HasPrefix(url, "http") {
			urlList = append(urlList, url)
		}
	}
	return urlList
}

func RenderTemplate(tmpl *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
