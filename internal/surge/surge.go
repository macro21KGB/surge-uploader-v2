package surge

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Website struct {
	Name string
	URL  string
}

func GetWebsites() ([]Website, error) {
	var websites []Website

	out, err := exec.Command("surge", "list").Output()
	if err != nil {
		return nil, err
	}

	elements := strings.Split(string(out), "\n")

	if len(elements) < 2 {
		return nil, errors.New("unexpected output from surge list command, Check if surge is installed and logged in")
	}

	for _, line := range elements {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			log.Printf("Skipping line: %s", line)
			continue
		}
		cleanedPart := stripAnsiAndControl(parts[1])
		url := "https://" + cleanedPart
		name := convertURLToName(url)

		websites = append(websites, Website{Name: name, URL: url})
	}

	return websites, nil
}

func stripAnsiAndControl(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	s = ansiRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "[39m", "")
	s = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, s)
	return s
}

func convertURLToName(url string) string {
	url = stripAnsiAndControl(url)

	// Remove the protocol (http:// or https://)
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// Remove any trailing slashes
	url = strings.TrimSuffix(url, "/")

	// Remove the .surge.sh suffix if present
	url = strings.TrimSuffix(url, ".surge.sh")

	// Replace dots and slashes with hyphens
	name := strings.ReplaceAll(url, ".", "-")
	name = strings.ReplaceAll(name, "/", "-")

	return name
}
