package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var mux sync.RWMutex

// 获取仓库分支列表
func BranchList(path string) ([]string, error) {
	// 记录工作环境 完成分支抓取之后进行还原
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	absPwd, err := filepath.Abs(pwd)
	if err != nil {
		return nil, err
	}
	defer os.Chdir(absPwd)

	err = os.Chdir(path)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "branch", "-r")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var res []string
	for _, branch := range branches {
		if strings.Contains(branch, "HEAD") {
			continue
		}
		arr := strings.Split(strings.TrimSpace(branch), "/")
		if len(arr) > 1 {
			res = append(res, arr[1])
		}
	}
	return res, nil
}

// 检测配置文件是否存在 key 字符串
func CheckKeyExist(c *PorductConfig, key string) bool {
	mux.RLock()
	defer mux.RUnlock()

	file := filepath.Join(c.Workspace, c.ConfigFile)
	input, err := os.Open(file)
	if err != nil {
		return false
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// line :=scanner.Bytes()
		if bytes.ContainsAny(scanner.Bytes(), key) {
			return true
		}
	}
	return false
}

func UpFile(config *PorductConfig, s string) error {
	mux.Lock()
	defer mux.Unlock()

	file := filepath.Join(config.Workspace, config.ConfigFile)
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer input.Close()

	tmp := filepath.Join(config.Workspace, "tmp")
	output, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.ContainsAny(line, config.InsertFlag) {
			if _, err := writer.WriteString(s + "\n"); err != nil {
				return err
			}
		}
		if _, err := writer.Write(line); err != nil {
			return err
		}
	}
	writer.Flush()

	err = os.Rename(tmp, file)
	return err
}

func GetLog(config *ServerConfig, id int) ([]byte, error) {
	f := filepath.Join(config.Root, config.Log, fmt.Sprintf("%d.log", id))
	return ioutil.ReadFile(f)
}

// func ErrCheck(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func Atoi(reader io.Reader) (id int) {
// 	data, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	id, err = strconv.Atoi(string(data))
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	return id
// }

// func CatchBranchs(data []byte) (info [][]string) {
// 	catchReg := regexp.MustCompile(`(?s)verCfg\s+=\s+({.*?});`)
// 	catchInfo := catchReg.FindSubmatch(data)
// 	if len(catchInfo) != 2 {
// 		return
// 	}

// 	info = [][]string{}
// 	reg := regexp.MustCompile(`".*?"`) // 抓取引号内文本
// 	for _, s := range strings.Split(string(catchInfo[1]), "\n") {
// 		res := reg.FindAllString(s, 2)
// 		if len(res) == 2 {
// 			for i := range res {
// 				res[i] = strings.Trim(res[i], `"`)
// 			}
// 			info = append(info, res)
// 		}
// 	}
// 	return
// }

// func CatchTeams(data []byte) (info []string) {
// 	catchReg := regexp.MustCompile(`(?s)langsMap\s+=\s+{(.*?)};`)
// 	catchInfo := catchReg.FindSubmatch(data)
// 	if len(catchInfo) != 2 {
// 		return
// 	}

// 	reg := regexp.MustCompile(`team\d+`)
// 	info = reg.FindAllString(string(catchInfo[1]), -1)
// 	return
// }
