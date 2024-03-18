package helpers

import (
	"errors"
	"regexp"
	"strings"
)

func ParseTaskDescription(desc string) (string, string) {
	re := regexp.MustCompile(`\[(.*?)\](.*)`)
	matches := re.FindStringSubmatch(desc)

	if len(matches) == 3 {
		return matches[1], strings.TrimSpace(matches[2])
	}

	return "", strings.TrimSpace(desc)
}

func ValidateJiraIDAndTitle(jiraID, title string) error {
	if strings.ReplaceAll(jiraID, " ", "") == "" {
		return errors.New("jira id must be set")
	}

	if containsBracketSequence(title) {
		return errors.New("jira id is not be in title")
	}

	return nil
}

func containsBracketSequence(s string) bool {
	re := regexp.MustCompile(`.*\[.*\].*`)
	return re.MatchString(s)
}
