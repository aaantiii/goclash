package clash

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	regexpTag = regexp.MustCompile("[^A-Z0-9]+")
)

// CorrectTag returns a valid Clash of Clans tag. It will be uppercase, have no special characters, and have a # at the beginning.
//
// Credit to: https://github.com/mathsman5133/coc.py/blob/master/coc/utils.py
func CorrectTag(tag string) string {
	return "#" + strings.ReplaceAll(regexpTag.ReplaceAllString(strings.ToUpper(tag), ""), "O", "0")
}

// TagURLSafe encodes a tag to be used in a URL.
func TagURLSafe(tag string) string {
	return url.PathEscape(tag)
}
