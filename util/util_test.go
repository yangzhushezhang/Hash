/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package util

import (
	"fmt"
	"testing"
)

func TestGetNiuNiuResultForBanker(t *testing.T) {
	//dataAaary := GetNiuNiuResultForBanker("0000000002601449928fb0af09b1e7d141e715d1359d746389192f995c2463f0")
	//banker := dataAaary[0]
	//player := dataAaary[1]

	//fmt.Println(10%10)
	str := "0000000002601449928fb0af09b1e7d141e715d1359d746389192f995c2463f0"
	start := len(str) - 2
	newHash := str[start:]

	fmt.Println(newHash)

}

type TransactionsJsonToStruct struct {
	Data []struct {
		TransactionID string `json:"transaction_id"`
		TokenInfo     struct {
			Symbol   string `json:"symbol"`
			Address  string `json:"address"`
			Decimals int    `json:"decimals"`
			Name     string `json:"name"`
		} `json:"token_info"`
		BlockTimestamp int64  `json:"block_timestamp"`
		From           string `json:"from"`
		To             string `json:"to"`
		Type           string `json:"type"`
		Value          string `json:"value"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At          int64  `json:"at"`
		Fingerprint string `json:"fingerprint"`
		Links       struct {
			Next string `json:"next"`
		} `json:"links"`
		PageSize int `json:"page_size"`
	} `json:"meta"`
}

type AutoGenerated struct {
	Data []struct {
		TransactionID string `json:"transaction_id"`
		TokenInfo     struct {
			Symbol   string `json:"symbol"`
			Address  string `json:"address"`
			Decimals int    `json:"decimals"`
			Name     string `json:"name"`
		} `json:"token_info"`
		BlockTimestamp int64  `json:"block_timestamp"`
		From           string `json:"from"`
		To             string `json:"to"`
		Type           string `json:"type"`
		Value          string `json:"value"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At          int64  `json:"at"`
		Fingerprint string `json:"fingerprint"`
		Links       struct {
			Next string `json:"next"`
		} `json:"links"`
		PageSize int `json:"page_size"`
	} `json:"meta"`
}

func TestWsResult(t *testing.T) {
	//minTimestamp:="1650341413596"
	//url := "https://api.trongrid.io/v1/accounts/TW2HWaLWy9pwiRN4yLju6YKW3aQ6Fw8888/transactions/trc20?only_to=true&limit=200&min_timestamp="+ minTimestamp
	//req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Add("Accept", "application/json")
	//res, _ := http.DefaultClient.Do(req)
	//body, _ := ioutil.ReadAll(res.Body)
	//var data  AutoGenerated
	//jsonErr := json.Unmarshal(body, &data)

	//

	//fmt.Println(GetNiuNiuResultForBanker("000000000262c304633d088cd78dc3e399c9648da73ec5f8778962c393365ed6"))

	//6a11f   6  1
	//fmt.Println(BjlResult("0000000002611f8f2612d308992b94faefc9458128b0514eb01fe5a42b65a16f", 2))


	fmt.Println(DsResult("0000000002611f8f2612d308992b94faefc9458128b0514eb01fe5a42b65a19f",101))




}
