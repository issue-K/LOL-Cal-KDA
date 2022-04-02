package lol

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var(
	httpCli = &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,  //使用http2.0
			TLSClientConfig: &tls.Config{ InsecureSkipVerify: true, },  //不去验证服务端的数字证书
		},
	}
	cli = &client{}
)

type client struct {
	port    int
	token string
	baseUrl string
}

func InitClient(port int,token string){
	cli.token = token
	cli.port = port
	cli.baseUrl = fmt.Sprintf("https://riot:%s@127.0.0.1:%d", cli.token, cli.port) //后续的请求都要在此路径后添加
}

func (cli client) httpGet(url string) ([]byte, error) {
	return cli.req(http.MethodGet, url, nil)
}
func (cli client) httpPost(url string, body interface{}) ([]byte, error) {
	return cli.req(http.MethodPost, url, body)
}
func (cli client) httpDel(url string) ([]byte, error) {
	return cli.req(http.MethodDelete, url, nil)
}
func (cli client) req(method string, url string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {  //有需要携带的数据,就放在body中
		bts, err := json.Marshal(data)  //转为二进制
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bts)
	}
	req, _ := http.NewRequest(method, cli.baseUrl+url, body)
	fmt.Printf( "url = %v\n",url )
	if req.Body != nil {
		req.Header.Add("ContentType", "application/json")
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}