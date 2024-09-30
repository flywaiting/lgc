package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
)

var mux sync.RWMutex

func InitLog(c *Workspace) {
	// 初始化日志文件夹
	log := filepath.Join(c.Root, c.Log)
	if _, err := os.Stat(log); err == nil || os.IsExist(err) {
		os.RemoveAll(log)
	}
	if err := os.MkdirAll(log, 0775); err != nil {
		panic("日志文件夹创建出错")
	}

}

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

	// branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := bytes.Split(bytes.TrimSpace(output), []byte("\n"))
	tag := []byte("HEAD")
	split := []byte("/")
	var res []string
	for _, branch := range branches {
		if bytes.Contains(branch, tag) {
			continue
		}
		arr := bytes.Split(branch, split)
		if len(arr) > 1 {
			res = append(res, string(arr[1]))
		}
	}
	return res, nil
}

// branch: alias
func GetBranchAliasMap(c *Project) (res map[string]string) {
	mux.RLock()
	defer mux.RUnlock()

	file := filepath.Join(c.Root, c.Config)
	input, err := os.Open(file)
	if err != nil {
		return
	}
	defer input.Close()

	reg := regexp.MustCompile(`".*?"`)
	target := []byte("verIndex")
	res = make(map[string]string)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !bytes.Contains(line, target) {
			continue
		}

		match := reg.FindAll(line, 2)
		if match == nil || len(match) < 2 || bytes.Equal(match[0], match[1]) {
			continue
		}
		res[string(match[1])] = string(match[0])
	}
	return
}

// 检测配置文件是否存在 key 字符串
func CheckKeyExist(c *Project, sub []byte) bool {
	mux.RLock()
	defer mux.RUnlock()

	file := filepath.Join(c.Root, c.Config)
	input, err := os.Open(file)
	if err != nil {
		return false
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// line :=scanner.Bytes()
		if bytes.Contains(scanner.Bytes(), sub) {
			return true
		}
	}
	return false
}

func UpFile(config *Project, s string) error {
	mux.Lock()
	defer mux.Unlock()

	file := filepath.Join(config.Root, config.Config)
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer input.Close()

	tmp := filepath.Join(config.Root, "tmp")
	output, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	target := []byte(config.Flag)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.Contains(line, target) {
			if _, err := writer.WriteString(s + "\n"); err != nil {
				return err
			}
		}
		if _, err := writer.Write(line); err != nil {
			return err
		}
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	writer.Flush()
	input.Close()
	output.Close()

	err = os.Rename(tmp, file)
	return err
}

func GetLog(config *Workspace, id int) ([]byte, error) {
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
