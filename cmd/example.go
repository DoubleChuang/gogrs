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
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/DoubleChuang/gogrs/tradingdays"
	"github.com/DoubleChuang/gogrs/twse"
	"github.com/DoubleChuang/gogrs/utils"
	"github.com/spf13/cobra"
)

const shortForm = "20060102"

var fiNetBuyDay int
var itNetBuyDay int
var fiIncrementalBuy bool
var itIncrementalBuy bool
var minDataNum *int
var useMtss *bool
var useT38 *bool
var useT44 *bool
var useMa *bool
var useCp *bool
var useDate *string
var useNo *string

//global var
var (
	T38U *twse.TWT38U
	T44U *twse.TWT44U
	MTSS *twse.TWMTSS

	TpexFI *twse.TPEXT38U
	TpexIT *twse.TPEXT44U
)

//OTCCLASS OTCCLASS
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

//TWSECLASS TWSECLASS
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

		if err := utils.RecoveryStockBackup(*useDate); err != nil {
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

var testTpexT38Cmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  `test tpex t38`,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := time.Parse(shortForm, *useDate)
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, utils.TaipeiTimeZone)

		TpexT38U := twse.NewTPEXT38U(date)
		if v, err := TpexT38U.Get(); err == nil {
			utils.Dbgln(v["3105"])
		} else {
			utils.Dbgln(err.Error())
		}

		TpexT44U := twse.NewTPEXT44U(date)
		if v, err := TpexT44U.Get(); err == nil {
			utils.Dbgln(v["3105"])
		} else {
			utils.Dbgln(err.Error())
		}
	},
}

var getStockCmd = &cobra.Command{
	Use:   "gs",
	Short: "Get one stock",
	Long:  `Get One Stock of TWSE`,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := time.Parse(shortForm, *useDate)
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, utils.TaipeiTimeZone)

		if T38U == nil {
			T38U = twse.NewTWT38U(date)
		}
		if T44U == nil {
			T44U = twse.NewTWT44U(date)
		}
		if MTSS == nil {
			MTSS = twse.NewTWMTSS(date, "ALL")
		}

		if s, err := getOneTWSE(date, *useNo, &T38U, &T44U, &MTSS); err != nil {
			utils.Dbgln(err)
		} else {
			utils.Dbgln(s)
		}
	},
}

func getOneTWSE(date time.Time, stockNo string, t38 **twse.TWT38U, t44 **twse.TWT44U, mtss **twse.TWMTSS) (string, error) {
	var ret string

	t := twse.NewLists(date)
	tList := t.GetCategoryList("ALLBUT0999")
	found := false
	var pStock *twse.Data
	for _, v := range tList {
		if v.No == stockNo {
			pStock = twse.NewTWSE(stockNo, date)
			found = true
			break
		}
	}
	if !found {
		return fmt.Sprintf("%s沒有%s此股票", date.Format(shortForm), stockNo), errors.Errorf("%s沒有%s此股票", date.Format(shortForm), stockNo)
	}
	//}
	utils.Dbgln(pStock.Date)

	mtssMapData, err := (*mtss).SetDate(date).GetData()
	if err != nil {
		return fmt.Sprintf("融資融券資料錯誤"), errors.Errorf("融資融券資料錯誤")
	}
	if err := prepareStock(pStock, 20); err == nil {
		var i int
		var d time.Time
		for i, d = range pStock.GetDateList() {
			if d == date {
				break
			}
		}
		utils.Dbgln(i)
		isT38OverBought, _, x := (*t38).IsOverBoughtDates(stockNo, 3)
		isT44OverBought, _, y := (*t44).IsOverBoughtDates(stockNo, 3)
		if s, err := showStock(pStock, 20); err == nil {
			ret = fmt.Sprintf("漲跌: %.2f\n成交價: %.2f\n漲跌幅: %.2f%%\n20MA:%.2f\n突破MA:%t\n外資增：%t %d\n投信增:%t %d\n融資增：%t %d\n融券增：%t %d\n=========\n",
				s.todayRange,
				s.todayPrice,
				s.todayGain,
				s.NDayAvg,
				s.overMA,
				isT38OverBought, x[0],
				isT44OverBought, y[0],
				mtssMapData[stockNo].MT.Total > 0, mtssMapData[stockNo].MT.Total,
				mtssMapData[stockNo].SS.Total > 0, mtssMapData[stockNo].SS.Total,
			)
		}
	} else {
		ret = fmt.Sprintf("資料錯誤")
	}

	return ret, nil

}

var getTWSECmd = &cobra.Command{
	Use:   "gtw",
	Short: "Get ALL TWSE",
	Long:  `Get All Stock of TWSE`,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := time.Parse(shortForm, *useDate)
		if T38U == nil {
			T38U = twse.NewTWT38U(date)
		}
		if T44U == nil {
			T44U = twse.NewTWT44U(date)
		}
		if MTSS == nil {
			MTSS = twse.NewTWMTSS(date, "ALL")
		}

		if err := getTWSE(date, "ALLBUT0999", *minDataNum, T38U, T44U, MTSS); err != nil {
			utils.Dbgln(err)
		}
	},
}

var getTPEXCmd = &cobra.Command{
	Use:   "gtp",
	Short: "Get ALL TPEX",
	Long:  `Get All Stock of TPEX`,

	Run: func(cmd *cobra.Command, args []string) {

		date, _ := time.Parse(shortForm, *useDate)
		if TpexFI == nil {
			TpexFI = twse.NewTPEXT38U(date)
		}
		if TpexIT == nil {
			TpexIT = twse.NewTPEXT44U(date)
		}

		if err := getOTC(date, "EW", *minDataNum, TpexFI, TpexIT); err != nil {
			utils.Dbgln(err)
		}

	},
}
var getT38Cmd = &cobra.Command{
	Use:   "gf",
	Short: "Get All Foreign Investor",
	Long:  `Get All Stock of Foreign Investor`,

	Run: func(cmd *cobra.Command, args []string) {
		T38U.GetData()
	},
}
var getT44Cmd = &cobra.Command{
	Use:   "gt",
	Short: "Get All Investment Trust",
	Long:  `Get All Stock of Investment Trust`,

	Run: func(cmd *cobra.Command, args []string) {
		T44U.GetData()

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
	}
	fmt.Println("Remove: ", str)

	return nil

	//fmt.Println("checkFirstDayOfMonth:", stock.Date, d, stock.Date==d)
}

func getTWSE(date time.Time, category string, minDataNum int, t38 *twse.TWT38U, t44 *twse.TWT44U, mtss *twse.TWMTSS) error {
	utils.Dbgln(date.Format(shortForm))

	t := twse.NewLists(date)
	tList := t.GetCategoryList(category)

	/*year, month, day := date.Date()
	csvFile, err := os.OpenFile(fmt.Sprintf("%d%02d%02d.csv", year, month, day), os.O_CREATE|os.O_RDWR, 0666)
	defer csvFile.Close()
	if err != nil {
		utils.Dbg("error: %s\n", err)
		return err
	}
	csvWriter := csv.NewWriter(csvFile)*/

	mtssMapData, err := mtss.GetData()
	if err != nil {
		return errors.Wrap(err, "MTSS GetData Fail.")
	}

	for _, v := range tList {
		//fmt.Printf("No:%s\n", v.No)
		stock := twse.NewTWSE(v.No, date)
		//checkFirstDayOfMonth(stock)
		if err := prepareStock(stock, minDataNum); err == nil {
			output := true
			isT38OverBought, t38Increase, t38ValList := t38.IsOverBoughtDates(v.No, fiNetBuyDay)
			isT44OverBought, t44Increase, t44ValList := t44.IsOverBoughtDates(v.No, itNetBuyDay)
			isMTSSOverBought := mtssMapData[v.No].MT.Total > 0 && mtssMapData[v.No].SS.Total > 0
			if res, err := showStock(stock, minDataNum); err == nil {
				if *useCp {
					if res.todayGain >= 3.5 {
						output = true
					} else {
						output = false
					}
				}
				if *useMa {
					if !res.overMA {
						output = false
					}
				}
				if *useT38 {
					if !isT38OverBought {
						output = false
					}
				}
				if *useT44 {
					if !isT44OverBought {
						output = false
					}
				}
				if *useMtss {
					if !isMTSSOverBought {
						output = false
					}
				}
				if fiIncrementalBuy {
					if !t38Increase {
						output = false
					}
				}
				if itIncrementalBuy {
					if !t44Increase {
						output = false
					}
				}
				if output {
					/*err = csvWriter.Write([]string{v.No,
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
					}*/
					fmt.Printf("No:%6s Range: %6.2f Price: %6.2f Gain: %6.2f%% NDayAvg:%6.2f overMA:%t T38OverBought:%t(%v) T44OverBought:%t(%v) MTSSOverBought:%t\n",
						v.No,
						res.todayRange,
						res.todayPrice,
						res.todayGain,
						res.NDayAvg,
						res.overMA,
						isT38OverBought,
						t38ValList,
						isT44OverBought,
						t44ValList,
						isMTSSOverBought)
				}

			}
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	return nil

}

func getOTC(date time.Time, category string, minDataNum int, t38 *twse.TPEXT38U, t44 *twse.TPEXT44U) error {
	utils.Dbgln(date.Format(shortForm))
	otc := twse.NewOTCLists(date)

	oList := otc.GetCategoryList(category)

	/*year, month, day := date.Date()
	csvFile, err := os.OpenFile(fmt.Sprintf("%d%02d%02d.csv", year, month, day), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		utils.Dbg("error: %s\n", err)
		return err
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)*/

	for _, v := range oList {
		stock := twse.NewOTC(v.No, date)
		if err := prepareStock(stock, minDataNum); err == nil {

			output := true
			isT38OverBought, t38Increase, t38ValList := t38.IsOverBoughtDates(v.No, fiNetBuyDay)
			isT44OverBought, t44Increase, t44ValList := t44.IsOverBoughtDates(v.No, itNetBuyDay)
			if res, err := showStock(stock, minDataNum); err == nil {
				if *useCp {
					if res.todayGain >= 3.5 {
						output = true
					} else {
						output = false
					}
				}
				if *useMa {
					if !res.overMA {
						output = false
					}
				}
				if *useT38 {
					if !isT38OverBought {
						output = false
					}
				}
				if *useT44 {
					if !isT44OverBought {
						output = false
					}
				}
				if fiIncrementalBuy {
					if !t38Increase {
						output = false
					}
				}
				if itIncrementalBuy {
					if !t44Increase {
						output = false
					}
				}
				if output {

					/*err = csvWriter.Write([]string{v.No,
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
					}*/
					fmt.Printf("No:%6s Range: %6.2f Price: %6.2f Gain: %6.2f%% NDayAvg:%6.2f overMA:%t T38OverBought:%t(%v) T44OverBought:%t(%v)\n",
						v.No,
						res.todayRange,
						res.todayPrice,
						res.todayGain,
						res.NDayAvg,
						res.overMA,
						isT38OverBought,
						t38ValList,
						isT44OverBought,
						t44ValList)
				}
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
	//utils.Dbgln("stock len:", stock.Len())
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
			return errors.New("Can't prepare enough data please check file has data or remove cache file")
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
	RootCmd.AddCommand(exampleCmd)
	RootCmd.AddCommand(testCmd)
	RootCmd.AddCommand(getAllStockCmd)
	//RootCmd.AddCommand(getTWSECmd)
	//RootCmd.AddCommand(getTPEXCmd)
	RootCmd.AddCommand(getT38Cmd)
	RootCmd.AddCommand(getT44Cmd)

	minDataNum = getAllStockCmd.PersistentFlags().IntP("num", "N", 3, "min date num")
	getAllStockCmd.PersistentFlags().IntVarP(&fiNetBuyDay, "fiBuyDay", "k", 3, "外資買超天數")
	getAllStockCmd.PersistentFlags().IntVarP(&itNetBuyDay, "itBuyDay", "j", 3, "投信買超天數")
	getAllStockCmd.PersistentFlags().BoolVarP(&fiIncrementalBuy, "fiIncrementalBuy", "g", false, "外資買超增加")
	getAllStockCmd.PersistentFlags().BoolVarP(&itIncrementalBuy, "itIncrementalBuy", "l", false, "投信買超增加")

	useMtss = getAllStockCmd.PersistentFlags().BoolP("mtss", "m", false, "使用融資融券篩選")
	useT38 = getAllStockCmd.PersistentFlags().BoolP("fi", "f", false, "使用外資篩選")
	useT44 = getAllStockCmd.PersistentFlags().BoolP("it", "i", false, "使用投信篩選")
	useMa = getAllStockCmd.PersistentFlags().BoolP("ma", "M", false, "使用移動平均篩選")
	useCp = getAllStockCmd.PersistentFlags().BoolP("cp" /*closing price*/, "c", false, "使用收盤價篩選")
	date := tradingdays.FindRecentlyOpenedTaipeiZone(time.Now()).Format(shortForm)
	useDate = getAllStockCmd.PersistentFlags().StringP("date", "d", date, "使用自訂日期")
	useNo = getStockCmd.Flags().StringP("No", "t", "2330", "stock No")

	getAllStockCmd.AddCommand(getTPEXCmd)
	getAllStockCmd.AddCommand(getTWSECmd)
	getAllStockCmd.AddCommand(getStockCmd)
	getAllStockCmd.AddCommand(testTpexT38Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
