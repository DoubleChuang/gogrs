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
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/DoubleChuang/gogrs/tradingdays"
	"github.com/DoubleChuang/gogrs/twse"
	"github.com/DoubleChuang/gogrs/utils"
	"github.com/spf13/cobra"
)
const shortForm = "20060102"
var minDataNum *int
var useMtss *bool
var useT38 *bool
var useT44 *bool
var useMa *bool
var useCp *bool
var useDate *string

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

		if err := utils.RecoveryStockBackup(*useDate);err != nil {
			utils.Dbgln(err)
		}
		//getT38(tradingdays.FindRecentlyOpened(time.Now()))
		//getT44(tradingdays.FindRecentlyOpened(time.Now()))

		//getTWSE("ALLBUT0999", *minDataNum)
		//getTWSE("26", *minDataNum)
		//getOTC("EW", *minDataNum)
		/*if v, err := getMTSS(tradingdays.FindRecentlyOpened(time.Now())); err == nil {
			utils.Dbgln(v)
		}*/
		//isMTSSOverBought, _ := getMTSSByDate("1215", 1)
		//fmt.Println(isMTSSOverBought)
	},
}

var getTWSECmd = &cobra.Command{
	Use:   "gtw",
	Short: "Get ALL TWSE",
	Long:  `Get All Stock of TWSE`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := getTWSE("ALLBUT0999", *minDataNum); err != nil {
			utils.Dbgln(err)
		}
	},
}

var getTPEXCmd = &cobra.Command{
	Use:   "gtp",
	Short: "Get ALL TPEX",
	Long:  `Get All Stock of TPEX`,

	Run: func(cmd *cobra.Command, args []string) {
		getOTC("EW", *minDataNum)
	},
}
var getT38Cmd = &cobra.Command{
	Use:   "gf",
	Short: "Get All Foreign Investor",
	Long:  `Get All Stock of Foreign Investor`,

	Run: func(cmd *cobra.Command, args []string) {
		getT38(tradingdays.FindRecentlyOpened(time.Now()))
		getT38(tradingdays.FindRecentlyOpened(time.Now()))
	},
}
var getT44Cmd = &cobra.Command{
	Use:   "gt",
	Short: "Get All Investment Trust",
	Long:  `Get All Stock of Investment Trust`,

	Run: func(cmd *cobra.Command, args []string) {
		if v, err := getT44(tradingdays.FindRecentlyOpened(time.Now())); err == nil {
			utils.Dbg("%v\n", v)
		}

	},
}

func checkFirstDayOfMonth(stock *twse.Data) error {
	year, month, day := stock.Date.Date()
	//d := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	s := twse.NewTWSE(stock.No, time.Date(year, month, day, 0, 0, 0, 0, time.UTC))

	hash := md5.New()
	io.WriteString(hash, s.URL())
	io.WriteString(hash, "")
	filehash := fmt.Sprintf("%x", hash.Sum(nil))
	//utils.Dbg("filehash:%s\n", filehash)
	str := utils.GetOSRamdiskPath("") + utils.TempFolderName + "/" + filehash
	if err := os.Remove(str); err != nil {
		utils.Dbg("Remove: %s %s\n", str, err)
		return err
	} else {
		fmt.Println("Remove: ", str)
	}
	return nil

	//fmt.Println("checkFirstDayOfMonth:", stock.Date, d, stock.Date==d)
}

var (
	T38DataMap map[time.Time]map[string]TXXData = make(map[time.Time]map[string]TXXData)
	T44DataMap map[time.Time]map[string]TXXData = make(map[time.Time]map[string]TXXData)
)

func getMTSSByDate(stockNo string, day int) (bool, []int64) {
	var (
		overbought int
		getDay     int
	)

	data := make([]int64, day)
	RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	//從最近的天數開始抓取 day 天的 資料 到 前(10+day)天 如果沒有抓到 day 天資料則錯誤
	for i := RecentlyOpendtoday; RecentlyOpendtoday.AddDate(0, 0, -10-day).Before(i) && getDay < day; i = i.AddDate(0, 0, -1) {
		if v, err := getMTSS(i); err == nil {
			getDay++
			if v[stockNo].MT.Total > 0 && v[stockNo].SS.Total > 0 {
				data[overbought] = v[stockNo].MT.Total
				overbought++
			}
		}
	}
	if getDay == day {
		return overbought == day, data
	} else {
		return false, nil
	}
}

func getT38ByDate(stockNo string, day int) (bool, []int64) {
	var (
		overbought int
		getDay     int
	)

	data := make([]int64, day)
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	RecentlyOpendtoday, _ := time.Parse(shortForm, *useDate)
	//從最近的天數開始抓取 day 天的 資料 到 前(10+day)天 如果沒有抓到 day 天資料則錯誤
	for i := RecentlyOpendtoday; RecentlyOpendtoday.AddDate(0, 0, -10-day).Before(i) && getDay < day; i = tradingdays.FindRecentlyOpened(i) {
		if v, err := getT38(i); err == nil {
			getDay++
			if v[stockNo].Total > 0 {
				data[overbought] = v[stockNo].Total
				overbought++
			}
		}
	}
	if getDay == day {
		return overbought == day, data
	} else {
		return false, nil
	}
}
func getT44ByDate(stockNo string, day int) (bool, []int64) {
	var (
		overbought int
		getDay     int
	)

	data := make([]int64, day)
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	RecentlyOpendtoday, _ := time.Parse(shortForm, *useDate)
	for i := RecentlyOpendtoday; RecentlyOpendtoday.AddDate(0, 0, -10-day).Before(i) && getDay < day; i = tradingdays.FindRecentlyOpened(i) {
		if v, err := getT44(i); err == nil {
			getDay++
			if v[stockNo].Total > 0 {
				data[overbought] = v[stockNo].Total
				overbought++
			}
		}
	}
	if getDay == day {
		return overbought == day, data
	} else {
		return false, nil
	}
}

func getTWSE(category string, minDataNum int) error {

	RecentlyOpendtoday, _ := time.Parse(shortForm, *useDate)
	utils.Dbgln(RecentlyOpendtoday)

	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())

	t := twse.NewLists(RecentlyOpendtoday)
	tList := t.GetCategoryList(category)
	year, month, day := RecentlyOpendtoday.Date()

	csvFile, err := os.OpenFile(fmt.Sprintf("%d%02d%02d.csv", year, month, day), os.O_CREATE|os.O_RDWR, 0666)
	defer csvFile.Close()
	if err != nil {
		utils.Dbg("error: %s\n", err)
		return err
	}
	csvWriter := csv.NewWriter(csvFile)
	//	t38 ,err := getT38(RecentlyOpendtoday)
	//	if err != nil{
	//		return err
	//	}

	mtssMapData, err := twse.NewTWMTSS(RecentlyOpendtoday, "ALL").GetData()
	if err != nil {
		return errors.Wrap(err, "MTSS GetData Fail.")
	}
	for _, v := range tList {
		//fmt.Printf("No:%s\n", v.No)
		stock := twse.NewTWSE(v.No, RecentlyOpendtoday)
		//checkFirstDayOfMonth(stock)
		if err := prepareStock(stock, minDataNum); err == nil {
			var output bool = true
			isT38OverBought, _ := getT38ByDate(v.No, 3)
			isT44OverBought, _ := getT44ByDate(v.No, 3)
			isMTSSOverBought := mtssMapData[v.No].MT.Total > 0 && mtssMapData[v.No].SS.Total > 0
			if res, err := showStock(stock, minDataNum); err == nil {
				if(*useCp){ 
					if(res.todayGain >= 3.5){
						output = true
					}else{
						output = false
					}
				}
				if(*useMa){
					if(!res.overMA){
						output = false
					}
				}
				if(*useT38){
					if(!isT38OverBought){
						output = false
					}
				}
				if(*useT44){
					if(!isT44OverBought){
						output = false
					}
				}
				if(*useMtss){
					if(!isMTSSOverBought){
						output = false
					}
				}
				if(output){
					err = csvWriter.Write([]string{v.No,
						v.Name,
						fmt.Sprintf("%.2f", res.todayRange),
						fmt.Sprintf("%.2f", res.todayPrice),
						fmt.Sprintf("%.2f", res.todayGain),
						fmt.Sprintf("%.2f", res.NDayAvg),
						fmt.Sprintf("%t", res.overMA),
						fmt.Sprintf("%t", isT38OverBought),
						fmt.Sprintf("%t", isT44OverBought),
						fmt.Sprintf("%t", isMTSSOverBought)})
					if err != nil {
						return err
					}
					csvWriter.Flush()
					err = csvWriter.Error()
					if err != nil {
						return err
					}
					fmt.Printf("No:%6s Range: %6.2f Price: %6.2f Gain: %6.2f%% NDayAvg:%6.2f overMA:%t T38OverBought:%t T44OverBought:%t MTSSOverBought:%t\n",
						v.No,
						res.todayRange,
						res.todayPrice,
						res.todayGain,
						res.NDayAvg,
						res.overMA,
						isT38OverBought,
						isT44OverBought,
						isMTSSOverBought)
				}

				
			}
		} else {
			fmt.Println(err)
		}
	}
	return nil

}

func getOTC(category string, minDataNum int) error {
	RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	otc := twse.NewOTCLists(RecentlyOpendtoday)

	oList := otc.GetCategoryList(category)

	year, month, day := RecentlyOpendtoday.Date()

	csvFile, err := os.OpenFile(fmt.Sprintf("%d%02d%02d.csv", year, month, day), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	defer csvFile.Close()
	if err != nil {
		utils.Dbg("error: %s\n", err)
		return err
	}
	csvWriter := csv.NewWriter(csvFile)

	for _, v := range oList {
		stock := twse.NewOTC(v.No, RecentlyOpendtoday)
		if err := prepareStock(stock, minDataNum); err == nil {
			if res, err := showStock(stock, minDataNum); err == nil {
				err = csvWriter.Write([]string{v.No,
					fmt.Sprintf("%.2f", res.todayRange),
					fmt.Sprintf("%.2f", res.todayPrice),
					fmt.Sprintf("%.2f", res.todayGain),
					fmt.Sprintf("%.2f", res.NDayAvg),
					fmt.Sprintf("%t", res.overMA)})
				if err != nil {
					return err
				}
				csvWriter.Flush()
				err = csvWriter.Error()
				if err != nil {
					return err
				}
				fmt.Printf("No: %6s Range: %.2f Price: %.2f Gain: %.2f%% NDayAvg:%.2f overMA:%t\n",
					v.No,
					res.todayRange,
					res.todayPrice,
					res.todayGain,
					res.NDayAvg,
					res.overMA,
				)
			}
		} else {
			fmt.Println(err)
		}
	}
	return nil
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

func prepareStock(stock *twse.Data, mindata int) error {

	if _, err := stock.Get(); err != nil {
		return err
	}

	if stock.Len() < mindata {
		start := stock.Len()
		for {
			stock.PlusData()
			if stock.Len() > mindata {
				break
			}
			if stock.Len() == start {
				break
			}
			start = stock.Len()
		}
		if stock.Len() < mindata {
			return errors.New("Can't prepare enough data, please check file has data or remove cache file")
		}
	}
	return nil
}

type resData struct {
	todayRange float64
	todayPrice float64
	todayGain  float64
	NDayAvg    float64
	overMA     bool
}

func showStock(stock *twse.Data, minDataNum int) (*resData, error) {
	var todayRange float64
	var todayPrice float64
	res := new(resData)
	minData := minDataNum
	if len(stock.RawData) < minData {
		fmt.Println(stock.Name, "No Data")
		return nil, errors.New("No Data")
	}
	rangeList := stock.GetRangeList()
	priceList := stock.GetPriceList()
	if len(rangeList) >= minData && len(priceList) >= minData {
		todayRange = rangeList[len(rangeList)-1]
		todayPrice = priceList[len(priceList)-1]
		res.todayRange = todayRange
		res.todayPrice = todayPrice
		res.todayGain = todayRange / todayPrice * 100

		//fmt.Printf("%.2f%%\n", todayRange/todayPrice*100)
	} else {
		return nil, errors.New("No enough price data")
	}
	daysAvg := stock.MA(minData)
	if len(daysAvg) > 0 {
		NDayAvg := daysAvg[len(daysAvg)-1]
		//fmt.Println(NDayAvg, todayPrice, todayPrice > NDayAvg)
		res.NDayAvg = NDayAvg
		res.overMA = todayPrice > NDayAvg
	} else {
		return nil, errors.New("No enough avg data")
	}
	return res, nil

}

//獲取外資與陸資
type TXXData struct {
	Buy   int64
	Sell  int64
	Total int64
}

func getMTSS(date time.Time) (map[string]twse.BaseMTSS, error) {
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	mtss := twse.NewTWMTSS(date, "ALL")
	data, err := mtss.Get()
	if err != nil {
		if strings.Contains(err.Error(), "File No Data") {
			if err := os.Remove(utils.GetMD5FilePath(mtss)); err != nil {
				return nil, errors.Wrap(err, "Get MTSS Fail")
			} else {
				if data, err = mtss.Get(); err != nil {
					return nil, errors.Wrap(err, "Get MTSS Fail")
				}
			}
		}
	}
	return data, nil
}
func getT38(date time.Time) (map[string]TXXData, error) {
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	if v, ok := T38DataMap[date]; ok {
		//utils.Dbg("Reuse T38Data:%v\n", date)
		return v, nil
	}

	t38 := twse.NewTWT38U(date)
	//fmt.Println(t38.URL())
	t38Map := make(map[string]TXXData)
	if data, err := t38.Get(); err == nil {
		for _, v := range data {
			//	fmt.Printf("No: %s Buy %d Sell %d Total %d\n",
			//		v[0].No,
			//		v[0].Buy,
			//		v[0].Sell,
			//		v[0].Total)
			t38Map[v[0].No] = TXXData{v[0].Buy, v[0].Sell, v[0].Total}
		}

	} else {
		utils.Dbg("Error: %s\n", err.Error())
		if strings.Contains(err.Error(), "File No Data") {
			if err := os.Remove(utils.GetMD5FilePath(t38)); err != nil {
				return nil, err
			} else {
				if data, err = t38.Get(); err != nil {
					//if t38Map, err = getT38(date.AddDate(0,0,-1));err!=nil{
					return nil, err
					//}

				} else {
					for _, v := range data {
						//	fmt.Printf("No: %s Buy %d Sell %d Total %d\n",
						//		v[0].No,
						//		v[0].Buy,
						//		v[0].Sell,
						//		v[0].Total)
						t38Map[v[0].No] = TXXData{v[0].Buy, v[0].Sell, v[0].Total}
					}
				}
			}
		}
	}
	//fmt.Println(t38Map)
	T38DataMap[date] = t38Map
	return t38Map, nil
}
func getT44(date time.Time) (map[string]TXXData, error) {
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	if v, ok := T44DataMap[date]; ok {
		//utils.Dbg("Reuse T44Data:%v\n", date)
		return v, nil
	}

	t44 := twse.NewTWT44U(date)
	//fmt.Println(t44.URL())
	t44Map := make(map[string]TXXData)
	if data, err := t44.Get(); err == nil {
		for _, v := range data {
			//	fmt.Printf("No: %s Buy %d Sell %d Total %d\n",
			//		v[0].No,
			//		v[0].Buy,
			//		v[0].Sell,
			//		v[0].Total)
			t44Map[v[0].No] = TXXData{v[0].Buy, v[0].Sell, v[0].Total}
		}

	} else {
		utils.Dbg("Error: %s\n", err.Error())
		if strings.Contains(err.Error(), "File No Data") {
			if err := os.Remove(utils.GetMD5FilePath(t44)); err != nil {
				return nil, err
			} else {
				if data, err = t44.Get(); err != nil {
					//if t44Map, err = getT44(date.AddDate(0,0,-1));err!=nil{
					return nil, err
					//}

				} else {
					for _, v := range data {
						//	fmt.Printf("No: %s Buy %d Sell %d Total %d\n",
						//		v[0].No,
						//		v[0].Buy,
						//		v[0].Sell,
						//		v[0].Total)
						t44Map[v[0].No] = TXXData{v[0].Buy, v[0].Sell, v[0].Total}
					}
				}
			}
		}
	}
	//fmt.Println(t44Map)
	T44DataMap[date] = t44Map
	return t44Map, nil
}

func showMyAll(stock *twse.Data) {
	rangeList := stock.GetRangeList()
	priceList := stock.GetPriceList()
	todayRange := rangeList[len(rangeList)-1]
	todayPrice := priceList[len(priceList)-1]

	daysAvg := stock.MA(20)
	today20MA := daysAvg[len(daysAvg)-1]

	//t86 := &twse.T86{Date: time.Now()}
	//result := Weight(time.Date(2017, 3, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))
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
	/*minDataNum = getAllStockCmd.Flags().IntP("num", "n", 3, "min date num")

	minDataNumTWSE = getTWSECmd.Flags().IntP("num", "n", 3, "min date num")
	minDataNumTPEX = getTPEXCmd.Flags().IntP("num", "n", 3, "min date num")*/
	//minDataNum = RootCmd.Flags().IntP("num", "n", 3, "min date num")
	RootCmd.AddCommand(exampleCmd)
	RootCmd.AddCommand(testCmd)
	RootCmd.AddCommand(getAllStockCmd)
	//RootCmd.AddCommand(getTWSECmd)
	//RootCmd.AddCommand(getTPEXCmd)
	RootCmd.AddCommand(getT38Cmd)
	RootCmd.AddCommand(getT44Cmd)
	//RootCmd.AddCommand(getStockCmd)

	minDataNum = getAllStockCmd.PersistentFlags().IntP("num", "N", 3, "min date num")
	useMtss = getAllStockCmd.PersistentFlags().BoolP("mtss", "m", false, "使用融資融券篩選")
	useT38 = getAllStockCmd.PersistentFlags().BoolP("fi", "f", false, "使用外資篩選")
	useT44 = getAllStockCmd.PersistentFlags().BoolP("it", "i", false, "使用投信篩選")
	useMa = getAllStockCmd.PersistentFlags().BoolP("ma", "M", false, "使用移動平均篩選")
	useCp = getAllStockCmd.PersistentFlags().BoolP("cp" /*closing price*/, "c", false, "使用收盤價篩選")
	date:=tradingdays.FindRecentlyOpened(time.Now()).Format(shortForm)
	useDate = getAllStockCmd.PersistentFlags().StringP("date" , "d", date ,"使用自訂日期")

	getAllStockCmd.AddCommand(getTPEXCmd)
	getAllStockCmd.AddCommand(getTWSECmd)
	//useMa = RootCmd.PersistentFlags().BoolP("c", "cp"/*closing price*/, true, "使用收盤價篩選")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
