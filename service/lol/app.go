package lol

import (
	"fmt"
	"log"
	"lolapi/service/models"
	"sync"
	"time"
)
var(
	defaultScore float64 = 100
)

type (
	UserScore struct {
		SummonerID   int64    `json:"summonerID"`
		SummonerName string   `json:"summonerName"`
		Score        float64  `json:"score"`
		AvgKDA      [][3]int `json:"currKDA"`
	}
)

func ChampionSelectStart() (error) {
	time.Sleep( time.Second )  //开始睡一会,因为服务端没那么快更新数据
	// 1、获取对战房间id
	roomId,err := GetRoomId()
	if err != nil{
		return err
	}
	//2、根据房间id获取5个召唤师id
	summonerIds := GetSummonerListByRoomId( roomId )
	// 3.根据5个召唤师id,分别查出各自的召唤师信息存起来
	//summonerInfos := make( []SummonerInfo,0 )
	//for _,summonerId := range summonerIds{
	//	summonerInfo,err := GetSummonerInfoById( summonerId )
	//	if err != nil{
	//		log.Printf("id = %v的用户信息查询失败, err = %v\n",summonerId,err )
	//		continue
	//	}
	//	summonerInfos = append( summonerInfos,summonerInfo )
	//	fmt.Printf("召唤师信息: %v\n",summonerInfo )
	//}
	// 5.根据5个召唤师信息,去计算各自的得分
	userScoreMap := map[int64]UserScore{}
	mu := sync.Mutex{}
	var wg = &sync.WaitGroup{}
	for _,summonerId := range summonerIds{
		wg.Add(1)
		go func(id int64){
    		score,err := CalUserScoreById( id )
    		if err != nil{
    			log.Printf("查询用户%v分数出错\n",id )
				return
    		}
    		mu.Lock()
    		userScoreMap[id] = score
    		mu.Unlock()
    		defer wg.Done()
		}(summonerId)
	}
	wg.Wait()
	// 6.把计算好的得分发送到聊天框内
	for _,msg := range userScoreMap {
		err := SendConversationMsg( msg , roomId)
		if err!=nil{
			log.Printf("发送消息时出现错误,err = %v\n",err )
		}
	}
	return nil
}
// 根据召唤师id计算得分
func CalUserScoreById( summonerId int64 ) (UserScore,error ) {
	userScoreInfo := &UserScore{
		SummonerID: summonerId,
		Score:      defaultScore,
	}
	// 获取用户信息
	summoner, err := GetSummonerInfoById(summonerId)
	if err != nil {
		return *userScoreInfo, err
	}
	userScoreInfo.SummonerName = summoner.DisplayName
	// 接下来只需要知道用户得分就行了.于是我们查询该召唤师的战绩列表
	gameList, err := listGameHistory(summonerId)
	if err != nil {
		log.Println("获取用户战绩失败,summonerId = %v, err = %v\n",summonerId,err )
		return *userScoreInfo, nil
	}
	for _,game := range gameList{
		var kda [3]int
		kda[0] = game.Participants[0].Stats.Kills
		kda[1] = game.Participants[0].Stats.Deaths
		kda[2] = game.Participants[0].Stats.Assists
		userScoreInfo.AvgKDA = append( userScoreInfo.AvgKDA,kda )
		if len( userScoreInfo.AvgKDA )==5 {
			break
		}
	}
	userScoreInfo.Score = CalSocre( gameList )
	return *userScoreInfo,nil
}
func CalSocre( gameList []GameInfo) float64 {
	var sum int64 = 0
	for _,game := range gameList{
		sum = sum+game.GameId
	}
	return float64( sum )
}
func listGameHistory( summonerId int64 ) ( []GameInfo,error ){
	fmtList := make([]GameInfo, 0, 20)
	resp, err := ListGamesBySummonerID(summonerId, 0, 20)
	fmt.Printf("用户id为%v\n", summonerId )
	fmt.Printf("用户战绩为%+v\n", resp )
	if err != nil {
		log.Printf("查询用户战绩失败,id = %v, err = %v\n",summonerId,err )
		return nil, err
	}
	for _, gameItem := range resp.Games.Games {  //遍历每一局游戏
		if gameItem.QueueId != models.NormalQueueID &&
			gameItem.QueueId != models.RankSoleQueueID &&
			gameItem.QueueId != models.ARAMQueueID &&
			gameItem.QueueId != models.RankFlexQueueID {
			continue
		}
		fmtList = append(fmtList, gameItem)
	}
	return fmtList, nil
}