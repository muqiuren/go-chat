/**
 * @Author Hatch
 * @Date 2021/01/01 22:18
**/
package config

import "time"

type App struct {
	MaxRoomNum				int					`yaml:"maxRoomNum"`
	MaxPersonOfRoom			int					`yaml:"maxPersonOfRoom"`
	MaxDisconnectTime		int					`yaml:"maxDisconnectTime"`
	HeartBeatTime			time.Duration		`yaml:"heartBeatTime"`
	HeartBeatMaxRetry		int					`yaml:"heartBeatMaxRetry"`
	PingWaitTime			time.Duration		`yaml:"pingWaitTime"`
	Port					string				`yaml:"port"`
	Tls						bool				`yaml:"tls"`
	RetryConnectWaitTime	int					`yaml:"retryConnectWaitTime"`
	StaticFolder			string				`yaml:"staticFolder"`
	AllowExts				map[string]bool		`yaml:"allowExts"`
	InitRoomCount			int					`yaml:"initRoomCount"`
	MaxRoomCapacity			int					`yaml:"maxRoomCapacity"`
}
