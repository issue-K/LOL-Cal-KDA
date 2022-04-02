package admin

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"github.com/yusufpapurcu/wmi"
)

var (
	lolCommandlineReg = regexp.MustCompile(`--remoting-auth-token=(.+?)" "--app-port=(\d+)"`)
	ErrLolProcessNotFound = errors.New("未找到lol进程")
)

func GetLolClientApiInfoV2() (int, string, error) {
	type Process struct {
		Commandline string
	}
	var port int
	var token string
	var dst []Process
	err := wmi.Query("select * from Win32_Process  WHERE name='LeagueClientUx.exe'", &dst)

	if err != nil || len(dst) == 0 {
		return port, token, ErrLolProcessNotFound
	}

	btsChunk := lolCommandlineReg.FindSubmatch([]byte(dst[0].Commandline))
	if len(btsChunk) < 3 {
		return port, token, ErrLolProcessNotFound
	}
	fmt.Println( string( btsChunk[0] ) )

	token = string(btsChunk[1])
	port, err = strconv.Atoi(string(btsChunk[2]))
	return port, token, err
}
