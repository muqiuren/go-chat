/**
 * @Author Hatch
 * @Date 2021/01/01 10:54
**/
package service

import (
	"fmt"
	"math/rand"
	"os"
	"path"
)

var (
	nameMap = []string{
		"黄","河","之","水","天","上","来",
		"奔","流","到","海","不","复","回",
		"高","堂","明","镜","悲","白","发",
		"朝","如","青","丝","暮","成","雪",
		"人","生","得","意","须","尽","欢",
		"莫","使","金","樽","空","对","月",
		"天","生","我","材","必","有","用",
		"千","金","散","尽","还","复","来",
	}
	nameMapLen = len(nameMap)
)

// 随机生成昵称
func RandomName(length int) string {
	if length <= 0 {
		length = 4
	}
	str := ""
	for i := 0; i < length; i++ {
		str += fmt.Sprintf(nameMap[rand.Intn(nameMapLen - 1)])
	}

	return str
}

// 文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 是否是目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// 是否是文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// 是否是允许后缀
func IsAllowExtFile(filename string, allowExts map[string]bool) bool {
	if !IsFile(filename) {
		return false
	}

	suffix := path.Ext(filename)
	return len(suffix) > 0 && allowExts[suffix[1:]]
}
