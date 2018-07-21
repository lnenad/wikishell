package wikiapi

import (
	"regexp"

	strip "github.com/grokify/html-strip-tags-go"
)

type Content string

const (
	validChars     = `[+&\$\-!%<>\(\)\|/\-–_\"\'\=\:\;\,\.[\]\p{L}\s\|\d]`
	sectionsRegExp = `===(?P<section>` + validChars + `*)===`
	linksRegExp    = `\[\[(?P<link>` + validChars + `*)\]\]`
	refRegExp      = `<ref.*>.*</ref>`
	ctxRegExp      = `\{\{(?P<context>` + validChars + `*)\}\}`
	commentsRegExp = `\<\-\-\s(?P<comment>` + validChars + `*)\s\-\-\>`
	imageRegExp    = `\[\[([/\-_\'\=\:\;\,\.\n\t\r[\]\p{L}\s\|\d]+)\]\]`
)

func (c *Content) GetSections() []string {
	sectionsR := regexp.MustCompile(sectionsRegExp)

	return sectionsR.FindAllString(string(*c), -1)
}

func (c *Content) GetLinks() []string {
	linkR := regexp.MustCompile(linksRegExp)

	return linkR.FindAllString(string(*c), -1)
}

func (c *Content) Cleanup() string {
	//sectionsR := regexp.MustCompile(sectionsRegExp)
	linkR := regexp.MustCompile(linksRegExp)
	imageR := regexp.MustCompile(imageRegExp)
	commentsR, _ := regexp.Compile(commentsRegExp)
	ctxR, _ := regexp.Compile(ctxRegExp)
	refR, _ := regexp.Compile(refRegExp)

	var str string

	str = linkR.ReplaceAllString(string(*c), "$link")
	str = imageR.ReplaceAllString(str, "")
	str = commentsR.ReplaceAllString(str, "")
	str = ctxR.ReplaceAllString(str, "")
	str = refR.ReplaceAllString(str, "")

	return str
	//str := linkR.FindAllStringSubmatch(string(*c), -1)
}

func (c *Content) StripHTML() *Content {
	nc := Content(strip.StripTags(string(*c)))
	return &nc
}

func (contents Content) SplitContentByLines(width int, height int) []string {
	var lines []string
	runes := []rune(string(contents))

	var line string
	for i := 0; i < len(runes); i++ {

		line += string(runes[i])
		if runes[i] == 0x000A || runes[i] == 0x000B || runes[i] == 0x000C || runes[i] == 0x000D {
			lines = append(lines, line)
			line = ""
			continue
		}
		if len(line) == width {
			lines = append(lines, line)
			line = ""
		}
	}

	return lines
}
