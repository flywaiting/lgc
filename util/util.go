package util

import (
	"regexp"
	"strings"
)

func ErrCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func CatchBranchs(data []byte) (info [][]string) {
	catchReg := regexp.MustCompile(`(?s)verCfg\s+=\s+({.*?});`)
	catchInfo := catchReg.FindSubmatch(data)
	if len(catchInfo) != 2 {
		return
	}

	info = [][]string{}
	reg := regexp.MustCompile(`".*?"`) // 抓取引号内文本
	for _, s := range strings.Split(string(catchInfo[1]), "\n") {
		res := reg.FindAllString(s, 2)
		if len(res) == 2 {
			for i := range res {
				strings.Trim(res[i], `"`)
			}
			info = append(info, res)
		}
	}
	return
}

func CatchTeams(data []byte) (info []string) {
	catchReg := regexp.MustCompile(`(?s)langsMap\s+=\s+{(.*?)};`)
	catchInfo := catchReg.FindSubmatch(data)
	if len(catchInfo) != 2 {
		return
	}

	reg := regexp.MustCompile(`team\d+`)
	info = reg.FindAllString(string(catchInfo[1]), -1)
	return
}
