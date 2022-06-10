package telegram

import (
	"errors"
	"regexp"
	"strings"
)

var titleMatcher = regexp.MustCompile(`^([^=]+)(=.*)$`)
var undesiredSections = regexp.MustCompile(`(?s)====?(?:Conjugation|Declension|Derived terms|Pronunciation)====?[^=]*`)
var mainDefinitionSearcher = regexp.MustCompile(`(?s)===([^=]+)===[^#]*# ([^\n]*)`)
var removeTransitiveness = regexp.MustCompile(`{{indtr\|[^}|]*\|([^}])}}\s*`)
var removeCurlyLink = regexp.MustCompile(`{{[^}]*[|=]([^|}=]+)}}`)
var removeSquareLink = regexp.MustCompile(`\[\[(?:[^|]*\|)?([^|\]]*)\]\]`)

type Word struct {
	title            string
	grammaticalClass string
	mainDefinition   string
	err              error
}

func Parse(s string) Word {
	titleSize := strings.Index(s, "=")
	if titleSize == 0 {
		return Word{
			err: errors.New("Invalid title"),
		}
	}

	title := s[:titleSize]

	// replace escaped line brakes with newline
	s = strings.Replace(s[titleSize:], "\\n", "\n", -1)

	// remove undesired sections
	s = undesiredSections.ReplaceAllString(s, "")

	section := mainDefinitionSearcher.FindStringSubmatch(s)
	if len(section) < 3 {
		return Word{
			err: errors.New("No mainDefinition found"),
		}
	}

	mainDefinition := removeTransitiveness.ReplaceAllString(section[2], "")
	mainDefinition = removeCurlyLink.ReplaceAllString(mainDefinition, "$1")
	mainDefinition = removeSquareLink.ReplaceAllString(mainDefinition, "$1")
	return Word{
		title:            title,
		grammaticalClass: section[1],
		mainDefinition:   mainDefinition,
	}
}
