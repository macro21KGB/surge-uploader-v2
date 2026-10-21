package templates

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"surge-uploader/m/v2/internal/surge"
)

type Template struct {
	Name string
	Path string
}

func GetDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".surge-templates")
}

// GetAll returns the list of templates from ~/.surge-templates.
// If the directory does not exist, it creates it and adds a default template.
func GetAll() []Template {
	tplDir := GetDirectory()

	if _, err := os.Stat(tplDir); os.IsNotExist(err) {
		err := os.MkdirAll(tplDir, 0755)
		if err != nil {
			return nil
		}

		defaultTpl := filepath.Join(tplDir, "default.html")
		_ = os.WriteFile(defaultTpl, []byte(getDefaultTemplate()), 0644)

		return []Template{{Name: "default", Path: defaultTpl}}
	}

	files, err := os.ReadDir(tplDir)
	if err != nil {
		return nil
	}

	var templates []Template
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".html" {
			continue
		}

		templates = append(templates, Template{
			Name: strings.TrimSuffix(file.Name(), ".html"),
			Path: filepath.Join(tplDir, file.Name()),
		})
	}

	return templates
}

func Render(templatePath string, websites []surge.Website) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	htmlTemplate := string(content)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlTemplate))
	if err != nil {
		return "", err
	}

	contentDiv := doc.Find("div#content").First()
	if contentDiv.Length() == 0 {
		return "", errors.New("div with id 'content' not found in the template file")
	}

	cardTemplateHtml := doc.Find("template#card").First()
	cardHtml, err := cardTemplateHtml.Html()
	if err != nil {
		return "", errors.New("template with id 'card' not found in the template file")
	}

	if cardHtml == "" {
		return "", errors.New("template with id 'card' is empty in the template file")
	}

	cardTemplate, err := template.New("card").Parse(cardHtml)
	if err != nil {
		return "", err
	}

	for _, website := range websites {
		var renderedCard strings.Builder
		err := cardTemplate.Execute(&renderedCard, website)
		if err != nil {
			return "", err
		}

		contentDiv.AppendHtml(renderedCard.String())
	}

	cardTemplateHtml.Remove()

	renderedHtml, err := doc.Html()
	if err != nil {
		return "", err
	}

	return renderedHtml, nil
}
