/**
 * @Author Hatch
 * @Date 2021/01/02 11:48
**/
package service

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go-chat/config"
	"gopkg.in/yaml.v2"
)

type Config struct {
	App		*config.App
}


// 初始化配置
func InitConfig() *Config {
    pwd, err := os.Getwd()
    if err != nil {
        panic(err)
        os.Exit(1)
    }

    absPath := strings.Replace(pwd, "\\", "/", -1)
    file := filepath.Join(absPath, "config.yaml")
    // 读取文件
    yamlFile, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatal(err)
    }
    c := &Config{}
    // 解析到struct
    if err := yaml.Unmarshal(yamlFile, &c); err != nil {
        log.Fatal(err)
    }
    log.Println("config initialize complete")
    return c
}