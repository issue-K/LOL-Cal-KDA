package lol

import (
	"encoding/json"
	"fmt"
)
// GetRoomId 获取当前对战聊天房间的ID
func GetRoomId() (string,error){
	msg,err := cli.httpGet("/lol-chat/v1/conversations")
	if err != nil{
		return "",err
	}
	conversations := make( []Conversation,10 )
	json.Unmarshal( msg,&conversations )
	for _,conversation := range conversations{
		fmt.Printf("conversation: %+v\n",conversation )
		if conversation.Type == "championSelect" {
			return conversation.Id, nil
		}
	}
	return "",nil
}
// GetSummonerListByRoomId 通过聊天房间ID,查出聊天记录,从中找到5个召唤师的id值
func GetSummonerListByRoomId( roomId string ) []int64 {
	msg,_ := cli.httpGet(fmt.Sprintf("/lol-chat/v1/conversations/%v/messages",roomId)) //得到这个房间内的所有消息
	lolMsgs := make( []LolMessage,10 )
	json.Unmarshal( msg,&lolMsgs )
	summonerIds := make( []int64,0 )
	for _,lolmsg := range lolMsgs{
		fmt.Printf("房间消息为%v\n",lolmsg )
		if lolmsg.Type == "system"{ //系统发出的消息,形如xxx进入房间
			summonerIds = append( summonerIds,lolmsg.FromSummonerID )
		}
	}
	return summonerIds
}
// GetSummonerInfoById 根据召唤师id查找召唤师的完整信息
func GetSummonerInfoById(id int64) ( SummonerInfo, error) {
	msg,err := cli.httpGet(fmt.Sprintf("/lol-summoner/v2/summoners?ids=[%v]",id ))
	summoners := make( []SummonerInfo,1 )
	if err != nil{
		return SummonerInfo{},err
	}

	json.Unmarshal( msg,&summoners )
	return summoners[0],nil
}
// ListGamesBySummonerID 根据召唤师id,查询最近[begin,begin+limit-1]的游戏战绩
func ListGamesBySummonerID( summonerId int64,begin,limit int) (*GameListResp,error){
	bts, err := cli.httpGet(fmt.Sprintf("/lol-match-history/v3/matchlist/account/%d?begIndex=%d&endIndex=%d",
		summonerId, begin, begin+limit))
	if err != nil{
		return nil,err
	}
	data := &GameListResp{}
	json.Unmarshal( bts,data )
	return data,nil
}
// SendConversationMsg 根据房间id发送消息
func SendConversationMsg(msg interface{},roomId string) error {
	TempByte,_ := json.Marshal( msg )
	data := struct {  //发送消息时,服务端指定格式的数据
		Body string `json:"body"`
		Type string `json:"type"`
	}{
		Body: string( TempByte ),
		Type: "chat",
	}
	mess, err := cli.httpPost(fmt.Sprintf("/lol-chat/v1/conversations/%s/messages", roomId), data )
	fmt.Println("响应为",string( mess ) )
	return err
}