/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/wangyi/MgHash/dao/mysql"
	"github.com/wangyi/MgHash/dao/redis"
	"github.com/wangyi/MgHash/model"
	"github.com/wangyi/MgHash/tools"
	"github.com/wangyi/MgHash/util"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

func NotificationTelegram(Message string) {
	TelegramToken := ":" //token
	TelegramChatId := "" //聊天id
	url := "https://api.telegram.org/bot" + TelegramToken + "/sendMessage?chat_id=" + TelegramChatId + "&text=" + Message
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

type AutoGenerated struct {
	Game string `json:"game"`
	Tx   struct {
		TxHash      string `json:"txHash"`
		BlockNumber int    `json:"blockNumber"`
		BlockHash   string `json:"blockHash"`
		Timestamp   int64  `json:"timestamp"`
		From        string `json:"from"`
		To          string `json:"to"`
		Amount      int    `json:"amount"`
		Token       string `json:"token"`
		Win         int    `json:"win"`
		Status      int    `json:"status"`
		Bonus       int    `json:"bonus"`
		SettleTx    string `json:"settleTx"`
		WagerType   string `json:"wagerType"`
		LotteryType string `json:"lotteryType"`
	} `json:"tx"`
}

/**
推送
*/
func ReceivePush(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	//fmt.Printf(string(data))
	zap.L().Debug(string(data))
	var mobile AutoGenerated
	err := json.Unmarshal(data, &mobile)
	if err != nil {
		util.JsonWrite(ctx, -101, nil, err.Error())
		return
	}

	//-1:投注无效, 0:输, 1:赢
	var win string
	status := 0
	if mobile.Tx.Win == -1 {
		win = "投注无效"
		status = 3
	} else if mobile.Tx.Win == 0 {
		win = "输"
		status = 1
	} else if mobile.Tx.Win == 1 {
		win = "赢"
		status = 2
	} else if mobile.Tx.Win == 2 {
		win = "平"
		status = 6
	}
	//投注金额			  /1000000
	acc := strconv.Itoa(mobile.Tx.Amount)
	account := ToDecimal(acc, 6).String()
	Bonus := strconv.Itoa(mobile.Tx.Bonus)
	bcc := ToDecimal(Bonus, 6).String()
	//PP, _ := util.ToDecimal(amount, 6).Float64()
	BlockNumber := strconv.Itoa(mobile.Tx.BlockNumber)

	var Tx4 string
	if len(mobile.Tx.From) == 34 {
		Tx := mobile.Tx.From
		Tx2 := Tx[:3]
		Tx3 := Tx[len(Tx)-3:]
		Tx4 = Tx2 + "********************" + Tx3
	} else {
		Tx4 = mobile.Tx.From
	}

	lickGame := mobile.Game
	if mobile.Game == "尾数" {
		lickGame = "幸运"
	}

	message := "🏆玩家地址：" + Tx4 + "%0A" +
		"🎲区块哈希（BLock hash）开奖结果:" + mobile.Tx.BlockHash + "%0A" +
		"区块号:" + BlockNumber + "%0A" +
		"💹资产类型：" + mobile.Tx.Token + "%0A" +
		"🎰游戏类型：" + lickGame + "%0A" +
		"🏅开奖状态：" + win + "%0A" +
		"💰投注金额：" + account + "%0A" +
		"投注地址:" + mobile.Tx.To + "%0A" +
		"奖交易哈希:" + mobile.Tx.SettleTx + "%0A" +
		"💵发放金额：" + bcc + "%0A" + "%0A" +
		"庄家开奖结果:" + mobile.Tx.LotteryType + "%0A" +
		"玩家开奖结果:" + mobile.Tx.WagerType + "%0A" +
		"欢迎大家加入官方交流群，每天晚上福利送不停，重要的话说三遍，返水返水返水只在交流群发放"
	go NotificationTelegram(message)
	//添加到数据库
	BlockNum := strconv.Itoa(mobile.Tx.BlockNumber)
	money, _ := ToDecimal(acc, 6).Float64()
	bet := money
	if mobile.Game == "牛牛" {
		bet = money / 10
	}
	BackMoney, _ := ToDecimal(Bonus, 6).Float64()
	//是否派发   1 已经派发 2无需派发 3没有派发  4派发失败
	Distribute := 1
	if BackMoney > 0 && mobile.Tx.BlockHash == "" {
		Distribute = 3
	}
	if BackMoney == 0 {
		Distribute = 2
	}
	if mobile.Tx.Status == 11 {
		Distribute = 4
	}
	tm := time.Unix(mobile.Tx.Timestamp/1000, 0)
	date := tm.Format("2006-01-02")
	add := model.TransactionRecordModel{
		ToAddress:      mobile.Tx.To,
		Form:           mobile.Tx.From,
		BlockHash:      mobile.Tx.BlockHash,
		TransactionId:  mobile.Tx.TxHash,
		Symbol:         mobile.Tx.Token,
		BetName:        mobile.Game,
		Status:         status, ////状态 1输  2 赢 3无效  6 和
		BlockNum:       BlockNum,
		BlockTimestamp: mobile.Tx.Timestamp,
		Money:          money,
		Bet:            bet,
		BankerResult:   mobile.Tx.LotteryType,
		PlayerResult:   mobile.Tx.WagerType,
		BackMoney:      BackMoney,
		Result:         BackMoney - money,
		Created:        time.Now().Unix(),
		Updated:        time.Now().Unix(),
		Distribute:     Distribute,
		Date:           date,
		Week:           tools.ReturnTheWeek(),
		Month:          tools.ReturnTheMonth(),
	}
	err = mysql.DB.Where("transaction_id=?", add.TransactionId).First(&model.TransactionRecordModel{}).Error
	if err != nil { //没有错误  说明找到了
		err = mysql.DB.Save(&add).Error
		if err == nil { //添加成功
			redis.Rdb.Set("Myself_"+add.ToAddress, add.BlockTimestamp, 0) //将最后一次的时间戳
			//入库成功  通知用户
			u := model.UserModel{Trc20Address: add.Form}
			Tid, err := u.ReturnTelegramId(mysql.DB)
			if err == nil {
				var strBool bool
				var filePath string
				var kinds string
				filePath = "real/" + strconv.FormatInt(Tid, 10) + "/" + time.Now().Format("20060112")
				kinds = mobile.Tx.BlockHash + ".png"
				if lickGame == "牛牛" {
					strBool = util.RealTimeSendResult("picture/Niu/back.jpg", util.NiuReturnImageUrl(mobile.Tx.WagerType, 1), util.NiuReturnImageUrl(mobile.Tx.LotteryType, 2), filePath, kinds)

				} else if lickGame == "幸运" {
					if status == 1 {
						//输
						strBool = util.RealTimeSendResult("picture/Luckly/back.jpg", "picture/Luckly/lose.png", "picture/Luckly/win.png", filePath, kinds)
					} else if status == 2 {
						//赢
						strBool = util.RealTimeSendResult("picture/Luckly/back.jpg", "picture/Luckly/win.png", "picture/Luckly/win.png", filePath, kinds)
					}
				} else if lickGame == "单双" {
					pally, zJ := util.DsResultBetData(mobile.Tx.BlockHash, money)
					strBool = util.RealTimeSendResult("picture/singleOrDouble/back.jpg", util.DsReturnImageUrl(pally), util.DsReturnImageUrl(zJ), filePath, kinds)

				} else if lickGame == "百家乐" {
					pally, zJ := util.BjlBetResultForTelegram(mobile.Tx.BlockHash)
					strBool = util.RealTimeSendResult("picture/Baccarat/back.jpg", "picture/Baccarat/player/"+pally+".png", "picture/Baccarat/banker/"+zJ+".png", filePath, kinds)

				}

				if strBool == true {
					if status != 3 {
						betString := strconv.FormatFloat(bet, 'f', 0, 64)
						backString := strconv.FormatFloat(BackMoney, 'f', 0, 64)
						tm := time.Unix(mobile.Tx.Timestamp/1000, 0)
						util.SendPhotoTo(viper.GetString("hash.HashToken"), Tid, filePath+"/"+kinds, status, lickGame, betString+mobile.Tx.Token, mobile.Tx.BlockHash, backString+mobile.Tx.Token, tm.Format("2006-01-02 15:04:05"))

					}

				}

			}

		}

	} else { //更新

	}
	util.JsonWrite(ctx, 200, nil, "执行成功")
	return
}