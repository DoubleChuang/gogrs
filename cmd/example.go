// Copyright © 2017 Toomore Chiang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

var minDataNum *int
var OTCCLASS = map[string]string{
	"02": "食品工業",
	"03": "塑膠工業",
	"04": "紡織纖維",
	"05": "電機機械",
	"06": "電器電纜",
	"08": "玻璃陶瓷",
	"10": "鋼鐵工業",
	"11": "橡膠工業",
	"14": "建材營造",
	"15": "航運業",
	"16": "觀光事業",
	"17": "金融保險",
	"18": "貿易百貨",
	"20": "其他",
	"21": "化學工業",
	"22": "生技醫療業",
	"23": "油電燃氣業",
	"24": "半導體業",
	"25": "電腦及週邊設備業",
	"26": "光電業",
	"27": "通信網路業",
	"28": "電子零組件業",
	"29": "電子通路業",
	"30": "資訊服務業",
	"31": "其他電子業",
	"32": "文化創意業",
	"80": "管理股票",
	"AA": "受益證券",
	"AL": "所有證券",
	"BC": "牛證熊證",
	"EE": "上櫃指數股票型基金(ETF)",
	"EW": "所有證券(不含權證、牛熊證)",
	"GG": "認股權憑證",
	"TD": "台灣存託憑證(TDR)",
	"WW": "認購售權證",
}
var TWSECLASS = map[string]string{
	"MS":         "大盤統計資訊",
	"0049":       "封閉式基金",
	"0099P":      "ETF",
	"01":         "水泥工業",
	"019919T":    "受益證券",
	"02":         "食品工業",
	"03":         "塑膠工業",
	"04":         "紡織纖維",
	"05":         "電機機械",
	"06":         "電器電纜",
	"07":         "化學生技醫療",
	"08":         "玻璃陶瓷",
	"09":         "造紙工業",
	"0999":       "認購權證",
	"0999B":      "熊證",
	"0999C":      "牛證",
	"0999G9":     "認股權憑證",
	"0999GA":     "附認股權特別股",
	"0999GD":     "附認股權公司債",
	"0999P":      "認售權證",
	"0999X":      "可展延牛證",
	"0999Y":      "可展延熊證",
	"10":         "鋼鐵工業",
	"11":         "橡膠工業",
	"12":         "汽車工業",
	"13":         "電子工業",
	"14":         "建材營造",
	"15":         "航運業",
	"16":         "觀光事業",
	"17":         "金融保險",
	"18":         "貿易百貨",
	"19":         "綜合",
	"20":         "其他",
	"21":         "化學工業",
	"22":         "生技醫療業",
	"23":         "油電燃氣業",
	"24":         "半導體業",
	"25":         "電腦及週邊設備業",
	"26":         "光電業",
	"27":         "通信網路業",
	"28":         "電子零組件業",
	"29":         "電子通路業",
	"30":         "資訊服務業",
	"31":         "其他電子業",
	"9299":       "存託憑證",
	"ALL":        "全部",
	"ALLBUT0999": "全部(不含權證、牛熊證、可展延牛熊證)",
	"CB":         "可轉換公司債",
}

// exampleCmd represents the example command

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Show example",
	Long:  `Show gogrs example.`,
	Run: func(cmd *cobra.Command, args []string) {
		var stock = twse.NewOTC("6129", tradingdays.FindRecentlyOpened(time.Now()))
		stock.Get()
		showAll(stock)
		fmt.Println("-----------------------------")
		stock.PlusData()
		showAll(stock)
		fmt.Println("-----------------------------")
		fmt.Println(tradingdays.IsOpen(2015, 5, 1))
	},
}

var getAllStockCmd = &cobra.Command{
	Use:   "ga",
	Short: "Get all stock",
	Long:  `Get All.`,
	Run: func(cmd *cobra.Command, args []string) {

		getTWSE("ALLBUT0999", *minDataNum)
		getOTC("EW", *minDataNum)
	},
}

func getTWSE(category string, minDataNum int) {
	var t = twse.NewLists(tradingdays.FindRecentlyOpened(time.Now()))

	tList := t.GetCategoryList(category)
	for _, v := range tList {
		fmt.Printf("No:%s\n", v.No)
		stock := twse.NewTWSE(v.No, tradingdays.FindRecentlyOpened(time.Now()))
		if prepareStock(stock, minDataNum) == true {
			showStock(stock, minDataNum)
		} else {
			fmt.Println("Fail to get enough data")
		}
	}

}

func getOTC(category string, minDataNum int) {
	var otc = twse.NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))

	oList := otc.GetCategoryList(category)

	for _, v := range oList {
		fmt.Printf("No:%s\n", v.No)
		stock := twse.NewOTC(v.No, tradingdays.FindRecentlyOpened(time.Now()))
		if prepareStock(stock, minDataNum) == true {
			showStock(stock, minDataNum)
		} else {
			fmt.Println("Fail to get enough data")
		}
	}
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Show Test",
	Long:  `Show gogrs Test.`,
	Run: func(cmd *cobra.Command, args []string) {
		//var stock = twse.NewOTC("3297", time.Date(2019, 2, 20, 0, 0, 0, 0, time.Local))
		var stock = twse.NewOTC("6188", tradingdays.FindRecentlyOpened(time.Now()))
		stock.Get()
		stock.PlusData()
		showMyAll(stock)
		fmt.Println("-----------------------------")
	},
}

func prepareStock(stock *twse.Data, mindata int) bool {

	var result bool = false
	if _, err := stock.Get(); err != nil {
		return result
	}

	if stock.Len() < mindata {
		start := stock.Len()
		for {
			stock.PlusData()
			if stock.Len() > mindata {
				result = true
				break
			}
			if stock.Len() == start {
				break
			}
			start = stock.Len()
		}
		if stock.Len() < mindata {
			result = false
		}
	} else {
		result = true
	}
	return result
}
func showStock(stock *twse.Data, minDataNum int) {
	var todayRange float64
	var todayPrice float64
	minData := minDataNum
	if len(stock.RawData) < minData {
		fmt.Println(stock.Name, "No Data")
		return
	}
	rangeList := stock.GetRangeList()
	priceList := stock.GetPriceList()
	if len(rangeList) >= minData && len(priceList) >= minData {
		todayRange = rangeList[len(rangeList)-1]
		todayPrice = priceList[len(priceList)-1]

		fmt.Printf("%.2f%%\n", todayRange/todayPrice*100)
	}
	daysAvg := stock.MA(minData)
	if len(daysAvg) > 0 {
		today20MA := daysAvg[len(daysAvg)-1]
		fmt.Println(today20MA, todayPrice, todayPrice > today20MA)
	}

}
func showMyAll(stock *twse.Data) {
	rangeList := stock.GetRangeList()
	priceList := stock.GetPriceList()
	todayRange := rangeList[len(rangeList)-1]
	todayPrice := priceList[len(priceList)-1]

	daysAvg := stock.MA(20)
	today20MA := daysAvg[len(daysAvg)-1]

	//t86 := &twse.T86{Date: time.Now()}
	//result := Weight(time.Date(2017, 3, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))a
	//d, _ := time.ParseDuration("-24h")
	//t38 := twse.NewTWT38U(time.Date(2015, 5, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))
	/*t38 := twse.NewTWT38U(tradingdays.FindRecentlyOpened(time.Now()))
	fmt.Println(t38.URL())
	if data, err := t38.Get(); err == nil {
		fmt.Println(len(data))
		fmt.Printf("data:%+v\n", data[len(data)-1])
	} else {
		fmt.Println("Error: ", err.Error())
	}*/

	fmt.Printf("%.2f%%\n", todayRange/todayPrice*100)
	fmt.Println(today20MA, todayPrice, todayPrice > today20MA)

}

func showAll(stock *twse.Data) {
	fmt.Println(stock.RawData)
	fmt.Println(stock.MA(6))
	fmt.Println(stock.MAV(6))
	fmt.Println(stock.GetPriceList())
	fmt.Println(utils.ThanPastFloat64(stock.GetPriceList(), 3, true))
	fmt.Println(utils.ThanPastFloat64(stock.GetPriceList(), 3, false))
	fmt.Println(stock.GetVolumeList())
	fmt.Println(utils.ThanPastUint64(stock.GetVolumeList(), 3, true))
	fmt.Println(utils.ThanPastUint64(stock.GetVolumeList(), 3, false))
	fmt.Println(stock.GetRangeList())
	fmt.Println(stock.IsRed())
}

func init() {
	minDataNum = getAllStockCmd.Flags().IntP("num", "n", 3, "min date num")
	RootCmd.AddCommand(exampleCmd)
	RootCmd.AddCommand(testCmd)
	RootCmd.AddCommand(getAllStockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
