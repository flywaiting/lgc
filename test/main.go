package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"lgc/util"
	"os"
	"os/exec"
	"path"
)

func catchBranches() {

	cmd := exec.Command("bash", "-c", "git branch -r")
	data, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	// reg := regexp.MustCompile(`^\s*origin/(?!HEAD)(\S+)`)
	// res := reg.FindSubmatch(data)
	// for _, b := range res {
	//	fmt.Println(string(b))
	// }

	res := bytes.FieldsFunc(data, func(r rune) bool { return r == '\n' })
	head := []byte("HEAD")
	origin := []byte("origin/")
	for _, b := range res {
		if bytes.Contains(b, head) {
			continue
		}
		b = bytes.TrimPrefix(bytes.TrimSpace(b), origin)
		fmt.Println(string(b))
	}
}

type mapConf map[interface{}]interface{}
type branchConf struct {
	Dir   string
	Steps []string
}

type sec struct {
	Dir   string
	Steps []string
	Port  int
}

func conf() {
	f, err := os.Open("conf.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var m mapConf
	// var m mapConf
	err = yaml.NewDecoder(f).Decode(&m)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("parse: %v\n", m)

	var s sec
	err = util.Decode(m["branch"], &s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", s)

	for i, v := range s.Steps {
		fmt.Println(i, v)
	}
}

func none() {
	// fmt.Println(os.Getwd())
	p, _ := os.Getwd()
	f, err := os.Stat(path.Join(p, "conf.yml"))
	fmt.Println(err)
	fmt.Println(f.Name())
	fmt.Println(p + "conf.yml")
}
func main() {
	conf()
	// none()
}
