/**
 * @Author Hatch
 * @Date 2021/01/02 15:18
**/
package service

import (
	"encoding/base64"
	"fmt"
	"github.com/holys/initials-avatar"
)

const (
	FontFilePath = "resource/assert/font/Hiragino_Sans_GB_W3.ttf"
)

var canvas = &avatar.InitialsAvatar{}

func init() {
	canvas = avatar.NewWithConfig(avatar.Config{
		FontFile: FontFilePath,
		FontSize: 80,
		MaxItems: 1024,
	})
}

// 生成头像
func GenAvatar(name string, size int, encoding string) string {
	raw, err := canvas.DrawToBytes(name, size, encoding)
	if err != nil{
		panic("Generate avatar fail:" + err.Error())
	}
	img := base64.StdEncoding.EncodeToString(raw)
	return fmt.Sprintf("data:image/%s;base64,%s", encoding, img)
}