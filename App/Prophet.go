package App

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"lolapi/pkg/lcu"
	admin "lolapi/pkg/window"
	"lolapi/service/lol"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type (
	Prophet struct {
		ctx       context.Context
		lcuPort   int
		lcuToken  string
		lcuActive bool
		cancel    func()
		mu        *sync.Mutex
	}
	WsMsg struct {
		Data      interface{} `json:"data"`
		EventType string      `json:"event_type"`
		Uri       string      `json:"uri"`
	}
)

const (
	onJsonApiEventPrefixLen = len(`[8,"OnJsonApiEvent",`)
	gameFlowChangedEvt      = "/lol-gameflow/v1/gameflow-phase"
)

func NewProphet() *Prophet {
	ctx, cancel := context.WithCancel(context.Background())  //上下文,用于取消一系列的goroutine
	return &Prophet{
		ctx:    ctx,
		cancel: cancel,
		mu:     &sync.Mutex{},
	}
}
func ( p *Prophet ) initLcuClient(port int,token string){
	lol.InitClient( port,token )
}

func ( p *Prophet ) Run() error{
	go p.UpdPortAndToken()  //检测lol进程
	time.Sleep( time.Second*100000000)
	return nil
}

func ( p *Prophet )UpdPortAndToken(){
	for{
		if !p.lcuActive{
			port, token, err := admin.GetLolClientApiInfoV2()
			if err != nil {
				if !errors.Is(admin.ErrLolProcessNotFound, err) {  //如果不是因为没有找到lol进程的错误,那就记录下来
					log.Println("获取lcu info 失败")
				}
				continue
			}
			p.initLcuClient(port, token)  //用port,token初始化cli
			err = p.WsConnectLol(port, token)
			if err != nil {
				log.Fatalln("游戏流程监视器 err:", err)
			}
		}
		time.Sleep( time.Second )
	}
}

func ( p *Prophet )WsConnectLol(port int,token string) error {
	//跳过证书检测
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{ InsecureSkipVerify: true }
	//构造header
	header := http.Header{}
	token = fmt.Sprintf("riot:%s", token )
	authToken := "Basic "+ base64.StdEncoding.EncodeToString( []byte( token ) )
	fmt.Printf("toekn = %v, authToken = %v, port = %v\n",token,authToken,port )
	header.Set("Authorization",authToken )
	//构造WsUrl
	WsUrl := fmt.Sprintf("wss://127.0.0.1:%d/",port)
	url,_ := url.Parse( WsUrl )
	WsUrl = url.String()
	//连接到lcu
	conn,_, err := dialer.Dial( WsUrl,header )
	if err != nil{
		log.Println("websocket连接到lcu失败, err = ",err )
		return err
	}
	p.lcuActive = true //成功连接了
	conn.WriteMessage( websocket.TextMessage, []byte("[5, \"OnJsonApiEvent\"]"))  //lcu的规则,先发这条消息,后续才会发给我们消息
	//下面就开始接收lcu发来的消息, 一切动作(开始匹配,回到大厅等消息都会发过来)
	for{
		messageType,msg,err := conn.ReadMessage()
		if err != nil{
			log.Println( "websocket接收消息失败,err = ",err )
			return err
		}
		//fmt.Printf( "找到websocket信息 %v\n",string(msg) )
		if messageType !=websocket.TextMessage || len(msg) < onJsonApiEventPrefixLen+1{
			continue
		}
		wsMsg := &WsMsg{}
		json.Unmarshal( msg[onJsonApiEventPrefixLen:len(msg)-1],wsMsg )

		if wsMsg.Uri == gameFlowChangedEvt{ //这是切换时间
			log.Printf("切换状态为,%v\n",wsMsg.Data )
			if wsMsg.Data.(string) == lcu.GameStatusChampionSelect{  //如果当前时间是匹配成功进入英雄选择界面
				log.Println("进入英雄选择阶段,正在计算队友分数")
				go lol.ChampionSelectStart()
			}
		}
	}
	return nil
}