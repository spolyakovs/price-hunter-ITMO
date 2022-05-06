package apiStore

import (
	"regexp"
	"strings"
)

func cleanTagName(tagNameRaw string) string {
	var reOther = regexp.MustCompile(`[^a-zA-Z ]`)

	tagNameClean := reOther.ReplaceAllString(tagNameRaw, "")

	tagNameClean = strings.Replace(tagNameClean, "  ", " ", -1)
	tagNameClean = strings.TrimSpace(tagNameClean)

	tagNameClean = strings.ToLower(tagNameClean)
	tagNameClean = strings.Replace(tagNameClean, " ", "_", -1)

	return tagNameClean
}

func cleanGameName(gameNameRaw string) string {
	var reBrackets = regexp.MustCompile(`\[.*\]`)
	var reTM = regexp.MustCompile(`\(TM\)`)
	var reOther = regexp.MustCompile(`[^a-zA-Z0-9 :\-!$%&'"*+=?^_|\.\(\)]`)

	gameNameClean := reBrackets.ReplaceAllString(gameNameRaw, "")
	gameNameClean = reTM.ReplaceAllString(gameNameClean, "")
	gameNameClean = reOther.ReplaceAllString(gameNameClean, "")

	gameNameClean = strings.Replace(gameNameClean, "  ", " ", -1)
	gameNameClean = strings.TrimSpace(gameNameClean)

	return gameNameClean
}

func checkGameName(gameName string) bool {
	var regexes = []*regexp.Regexp{
		regexp.MustCompile(`.*\(.*\).*`),
		regexp.MustCompile(`.*\/.*`),
		regexp.MustCompile(`.* - .*`),
		regexp.MustCompile(`.*dlc.*`),
		regexp.MustCompile(`.*ost.*`),
		regexp.MustCompile(`.*demo.*`),
		regexp.MustCompile(`.*edition.*`),
		regexp.MustCompile(`.*pack.*`),
	}

	gameNameLower := strings.ToLower(gameName)

	for _, re := range regexes {
		if re.MatchString(gameNameLower) {
			return false
		}
	}

	return true
}
