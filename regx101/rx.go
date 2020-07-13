package regx101

import (
	"regexp"
)

/*
Usage:
var compRegEx = regexp.MustCompile(`101Parse\sstatus:\s(?P<ParseStatus>\w+)\s\|\sfolderId:(?P<FolderID>\w+)\s\|\sUserID:(?P<UserID>\w+)\s\|\sPath:(?P<PathCfg>\w+)`)
					paramsMap := regx101.ParseRxNammed(compRegEx,msg1.MyStdout)
val := paramsMap["ParseStatus"]
 */
func ParseRxNammed(rx *regexp.Regexp, str string) map[string]string  {
	var paramsMap map[string]string
	match := rx.FindStringSubmatch(str)
	paramsMap = make(map[string]string)
	for i, name := range rx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return  paramsMap
}
