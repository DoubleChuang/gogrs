package twse

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/DoubleChuang/gogrs/tradingdays"
	"github.com/DoubleChuang/gogrs/utils"
)

// 錯誤資訊
var (
	errorFileNoData = errors.New("File No Data")
)

// BaseSellBuy 買進賣出合計
type BaseSellBuy struct {
	No    string
	Name  string
	Buy   int64 // 買進
	Sell  int64 // 賣出
	Total int64 // 合計
}

// QFIISTOP20 取得「外資及陸資持股比率前二十名彙總表」
type QFIISTOP20 struct {
	Date time.Time
}

// URL 擷取網址
func (q QFIISTOP20) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, fmt.Sprintf(utils.QFIISTOP20, q.Date.Year(), q.Date.Month(), q.Date.Day()))
}

// Get 擷取資料
func (q QFIISTOP20) Get() ([][]string, error) {
	var (
		err  error
		data []byte
	)
	if data, err = hCache.PostForm(q.URL(), nil); err == nil {
		if len(data) > 0 {
			csvArrayContent := strings.Split(string(data), "\n")[2:]
			for i, v := range csvArrayContent {
				csvArrayContent[i] = strings.Replace(v, "=", "", -1)
			}
			return csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll()
		}
	}
	return nil, err
}

// BFI82U 取得「三大法人買賣金額統計表」
type BFI82U struct {
	Begin time.Time
	End   time.Time
}

// NewBFI82U 三大法人買賣金額統計表
func NewBFI82U(begin, end time.Time) *BFI82U {
	return &BFI82U{Begin: begin, End: end}
}

// URL 擷取網址
func (b BFI82U) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.BFI82U, b.Begin.Year(), b.Begin.Month(), b.Begin.Day()))
}

// Get 擷取資料
func (b BFI82U) Get() ([]BaseSellBuy, error) {
	var (
		csvdata [][]string
		data    []byte
		err     error
		result  []BaseSellBuy
	)
	if data, err = hCache.PostForm(b.URL(), nil); err != nil {
		return nil, err
	}
	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(strings.Split(string(data), "\n")[2:7], "\n"))).ReadAll(); err == nil {
		result = make([]BaseSellBuy, len(csvdata))
		for i, v := range csvdata {
			result[i].Name = v[0]
			result[i].Buy, _ = strconv.ParseInt(strings.Replace(v[1], ",", "", -1), 10, 64)
			result[i].Sell, _ = strconv.ParseInt(strings.Replace(v[2], ",", "", -1), 10, 64)
			result[i].Total, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
		}
	}
	return result, err
}

// T86 取得「三大法人買賣超日報(股)」
type T86 struct {
	Date time.Time
}

// URL 擷取網址
func (t T86) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, fmt.Sprintf(utils.T86, t.Date.Year(), t.Date.Month(), t.Date.Day()))
}

// T86Data 各欄位資料
type T86Data struct {
	No     string
	Name   string
	FII    BaseSellBuy // 外資
	SIT    BaseSellBuy // 投信
	DProp  BaseSellBuy // 自營商(自行買賣)
	DHedge BaseSellBuy // 自營商(避險)
	Diff   int64       // 三大法人買賣超股數
}

// Get 擷取資料
func (t T86) Get(cate string) ([]T86Data, error) {
	data, err := hCache.PostForm(t.URL(), nil)
	if err != nil {
		return nil, err
	}
	csvArrayContent := strings.Split(string(data), "\n")[2:]
	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}
	var result []T86Data
	if csvdata, err := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[2:len(csvArrayContent)-8], "\n"))).ReadAll(); err == nil {
		result = make([]T86Data, len(csvdata))
		for i, v := range csvdata {
			if len(v) >= 11 {
				result[i].No = v[0]
				result[i].Name = v[1]

				result[i].FII.Buy, _ = strconv.ParseInt(strings.Replace(v[2], ",", "", -1), 10, 64)
				result[i].FII.Sell, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
				result[i].FII.Total = result[i].FII.Buy - result[i].FII.Sell

				result[i].SIT.Buy, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
				result[i].SIT.Sell, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)
				result[i].SIT.Total = result[i].SIT.Buy - result[i].SIT.Sell

				result[i].DProp.Buy, _ = strconv.ParseInt(strings.Replace(v[6], ",", "", -1), 10, 64)
				result[i].DProp.Sell, _ = strconv.ParseInt(strings.Replace(v[7], ",", "", -1), 10, 64)
				result[i].DProp.Total = result[i].DProp.Buy - result[i].DProp.Sell

				result[i].DHedge.Buy, _ = strconv.ParseInt(strings.Replace(v[8], ",", "", -1), 10, 64)
				result[i].DHedge.Sell, _ = strconv.ParseInt(strings.Replace(v[9], ",", "", -1), 10, 64)
				result[i].DHedge.Total = result[i].DHedge.Buy - result[i].DHedge.Sell

				result[i].Diff, _ = strconv.ParseInt(strings.Replace(v[10], ",", "", -1), 10, 64)
			}
		}
	} else {
		return nil, err
	}
	return result, err
}

type unixMapMTSSData map[int64]map[string]BaseMTSS
type unixMapT38UData map[int64]map[string]BaseT38U
type unixMapT43UData map[int64]map[string]BaseT43U
type unixMapT44UData map[int64]map[string]BaseT44U

type buySellTotalData map[int64]map[string]BaseBuySellTotalData

type TradingVolume struct {
	Buy   int64 // 買進
	Sell  int64 // 賣出
	Total int64 // 合計
}
type BaseMTSS struct {
	No   string
	Name string
	MT   TradingVolume
	SS   TradingVolume
}
type BaseT38U struct {
	No     string
	Name   string
	Volume TradingVolume
}

type BaseT43U struct {
	No     string
	Name   string
	Volume TradingVolume
}

type BaseT44U struct {
	No     string
	Name   string
	Volume TradingVolume
}

type BaseBuySellTotalData struct {
	No     string
	Name   string
	Volume TradingVolume
}

type TPEXT38U struct {
	Date     time.Time
	UnixData buySellTotalData
}
type TPEXT43U struct {
	Date     time.Time
	UnixData buySellTotalData
}
type TPEXT44U struct {
	Date     time.Time
	UnixData buySellTotalData
}

type TWT38U struct {
	Date            time.Time
	UnixMapT38UData unixMapT38UData
}
type TWT43U struct {
	Date            time.Time
	UnixMapT43UData unixMapT43UData
}
type TWT44U struct {
	Date            time.Time
	UnixMapT44UData unixMapT44UData
}

type TWMTSS struct {
	Date            time.Time
	Category        string
	UnixMapMTSSData unixMapMTSSData
}

// NewTWMTSS 融資融券匯總 Margin Trading and Short Selling
func NewTWMTSS(date time.Time, category string) *TWMTSS {
	return &TWMTSS{
		Date:            date,
		Category:        category,
		UnixMapMTSSData: make(unixMapMTSSData),
	}
}
func (t TWMTSS) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.TWMTSS,
			t.Date.Year(), t.Date.Month(), t.Date.Day(),
			t.Category))

}
func (t *TWMTSS) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date.AddDate(0, 0, -1))
}

func (t *TWMTSS) SetDate(date time.Time) *TWMTSS {
	t.Date = date
	return t
}
func (t *TWT38U) SetDate(date time.Time) *TWT38U {
	t.Date = date
	return t
}
func (t *TWT43U) SetDate(date time.Time) *TWT43U {
	t.Date = date
	return t
}
func (t *TWT44U) SetDate(date time.Time) *TWT44U {
	t.Date = date
	return t
}

//IsOverBought 獲取該日某股的融資融券是否為增
func (t *TWMTSS) IsOverBought(date time.Time, stockNo string) (bool, []int64, error) {
	var (
		data []int64
		err  error
		ans  bool
	)

	defer t.SetDate(t.Date)
	t.SetDate(date)

	if v, err := t.GetData(); err != nil {
		errors.Wrapf(err, "Fail to GetData in %s", date.Format("2006-01-02"))
	} else {
		if v[stockNo].MT.Total > 0 && v[stockNo].SS.Total > 0 {
			data := make([]int64, 2)
			data[0] = v[stockNo].MT.Total
			data[1] = v[stockNo].SS.Total
			ans = true
		}
	}
	return ans, data, err
}
func (t *TWMTSS) Get() (map[string]BaseMTSS, error) {
	dateUnix := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location()).Unix()
	if v, ok := t.UnixMapMTSSData[dateUnix]; ok {
		return v, nil
	}
	var (
		csvdata   [][]string
		data      []byte
		err       error
		resultMap map[string]BaseMTSS
	)
	//fmt.Println(t.URL())
	if data, err = hCache.PostForm(t.URL(), nil); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	if len(csvArrayContent) < 14 {
		if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
			return nil, err
		}
		return nil, errorFileNoData
	}
	//從第八列開始 然後刪掉最後面的八行(注意可能會有空白的行)
	csvArrayContent = csvArrayContent[7 : len(csvArrayContent)-8]

	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
		resultMap = make(map[string]BaseMTSS, len(csvdata))
		for i, v := range csvdata {
			if i == 0 {
				if false == checkCsvDataFormat("MTSS", v) {
					return nil, errors.New("Wrong MTSS Csv Data Format")
				}
				continue
			}
			var r BaseMTSS
			no := strings.Replace(v[0], " ", "", -1)

			r.Name = strings.Replace(v[1], " ", "", -1)

			r.MT.Buy, _ = strconv.ParseInt(strings.Replace(v[2], ",", "", -1), 10, 64)
			r.MT.Sell, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
			//TODO:確認是否是這樣計算總數
			r.MT.Total = r.MT.Buy - r.MT.Sell

			r.SS.Buy, _ = strconv.ParseInt(strings.Replace(v[8], ",", "", -1), 10, 64)
			r.SS.Sell, _ = strconv.ParseInt(strings.Replace(v[9], ",", "", -1), 10, 64)
			//TODO:確認是否是這樣計算總數
			r.SS.Total = r.SS.Sell - r.SS.Buy
			resultMap[no] = r

		}
		t.UnixMapMTSSData[dateUnix] = resultMap
	}
	return resultMap, err
}

func (t *TWMTSS) GetData() (map[string]BaseMTSS, error) {
	if v, err := t.Get(); err == nil {
		return v, err
	} else {
		t.Round()
		/*if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
		return nil, errors.Wrap(err, "TWMTSS Remove Cache File Fail")		}*/
		return t.GetData()
	}
}

func (t *TWT38U) GetData() (map[string]BaseT38U, error) {
	if v, err := t.Get(); err == nil {
		return v, err
	} else {
		t.Round()
		/*if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
		return nil, errors.Wrap(err, "TWMTSS Remove Cache File Fail")		}*/
		return t.GetData()
	}
}

func (t *TWT43U) GetData() (map[string]BaseT43U, error) {
	if v, err := t.Get(); err == nil {
		return v, err
	} else {
		t.Round()
		/*if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
		return nil, errors.Wrap(err, "TWMTSS Remove Cache File Fail")		}*/
		return t.GetData()
	}
}
func (t *TWT44U) GetData() (map[string]BaseT44U, error) {
	if v, err := t.Get(); err == nil {
		return v, err
	} else {
		t.Round()
		/*if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
		return nil, errors.Wrap(err, "TWMTSS Remove Cache File Fail")		}*/
		return t.GetData()
	}
}
func (t *TWT38U) IsOverBoughtDates(stockNo string, days int) (bool, bool, []int64) {
	var (
		overbought  int
		getDays     int
		increment   bool
		onlyOneZero bool = true
	)
	if days < 0 {
		return false, false, nil
	}
	data := make([]int64, days)
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	RecentlyOpendtoday := t.Date
	bkDate := t.Date
	//從最近的天數開始抓取 days 天的 資料 到 前(10+days)天 如果沒有抓到 days 天資料則錯誤
	for ; RecentlyOpendtoday.AddDate(0, 0, -10-days).Before(t.Date) && getDays < days; t.Round() {
		if v, err := t.GetData(); err == nil {
			getDays++
			if onlyOneZero || v[stockNo].Volume.Total > 0 {
				if v[stockNo].Volume.Total <= 0 {
					onlyOneZero = false
				}
				data[overbought] = v[stockNo].Volume.Total
				overbought++
			}
		}
	}
	t.SetDate(bkDate)
	if getDays == days {
		for i := 1; i < len(data); i++ {
			if data[i-1] >= data[i] {
				increment = true
			} else {
				increment = false
				break
			}
		}
		return overbought == days, increment, data
	} else {
		return false, false, nil
	}
}
func (t *TWT43U) IsOverBoughtDates(stockNo string, days int) (bool, []int64) {
	var (
		overbought int
		getDays    int
	)

	data := make([]int64, days)
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	RecentlyOpendtoday := t.Date
	bkDate := t.Date
	//從最近的天數開始抓取 days 天的 資料 到 前(10+days)天 如果沒有抓到 days 天資料則錯誤
	for ; RecentlyOpendtoday.AddDate(0, 0, -10-days).Before(t.Date) && getDays < days; t.Round() {
		if v, err := t.GetData(); err == nil {
			getDays++
			if v[stockNo].Volume.Total > 0 {
				data[overbought] = v[stockNo].Volume.Total
				overbought++
			}
		}
	}
	t.SetDate(bkDate)
	if getDays == days {
		return overbought == days, data
	} else {
		return false, nil
	}
}
func (t *TWT44U) IsOverBoughtDates(stockNo string, days int) (bool, bool, []int64) {
	var (
		overbought  int
		getDays     int
		increment   bool
		onlyOneZero = true
	)

	data := make([]int64, days)
	//RecentlyOpendtoday := tradingdays.FindRecentlyOpened(time.Now())
	RecentlyOpendtoday := t.Date
	bkDate := t.Date
	//從最近的天數開始抓取 days 天的 資料 到 前(10+days)天 如果沒有抓到 days 天資料則錯誤
	for ; RecentlyOpendtoday.AddDate(0, 0, -10-days).Before(t.Date) && getDays < days; t.Round() {
		if v, err := t.GetData(); err == nil {
			getDays++
			if onlyOneZero || v[stockNo].Volume.Total > 0 {
				if v[stockNo].Volume.Total <= 0 {
					onlyOneZero = false
				}
				data[overbought] = v[stockNo].Volume.Total
				overbought++

			}
		}
	}
	t.SetDate(bkDate)
	if getDays == days {
		for i := 1; i < len(data); i++ {
			if data[i-1] >= data[i] {
				increment = true
			} else {
				increment = false
				break
			}
		}
		return overbought == days, increment, data
	} else {
		return false, false, nil
	}
}

// TWTXXU 產生 自營商、投信、外資及陸資買賣超彙總表
type TWTXXU struct {
	Date time.Time
	fund string
}

func (t *TWTXXU) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

//台灣證券櫃買中心

//Round 外資
func (t *TPEXT38U) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

//Round 投信
func (t *TPEXT43U) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

//Round 自營商
func (t *TPEXT44U) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

//台灣證交所

//Round 外資
func (t *TWT38U) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

//Round 投信
func (t *TWT43U) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

//Round 自營商
func (t *TWT44U) Round() {
	t.Date = tradingdays.FindRecentlyOpened(t.Date)
}

func NewTPEXT38U(date time.Time) *TPEXT38U {
	return &TPEXT38U{
		Date:     date,
		UnixData: make(buySellTotalData),
	}
}

func NewTPEXT44U(date time.Time) *TPEXT44U {
	return &TPEXT44U{
		Date:     date,
		UnixData: make(buySellTotalData),
	}
}

func NewTPEXT43U(date time.Time) *TPEXT43U {
	return &TPEXT43U{
		Date:     date,
		UnixData: make(buySellTotalData),
	}
}

// NewTWT38U 外資及陸資買賣超彙總表
func NewTWT38U(date time.Time) *TWT38U {
	return &TWT38U{
		Date:            date,
		UnixMapT38UData: make(unixMapT38UData),
	}
}

// NewTWT43U 自營商買賣超彙總表
func NewTWT43U(date time.Time) *TWT43U {
	return &TWT43U{
		Date:            date,
		UnixMapT43UData: make(unixMapT43UData),
	}
}

// NewTWT44U 投信買賣超彙總表
func NewTWT44U(date time.Time) *TWT44U {
	return &TWT44U{
		Date:            date,
		UnixMapT44UData: make(unixMapT44UData),
	}
}

// URL 擷取 外資 網址
func (t TPEXT38U) URL() string {

	return fmt.Sprintf("%s%s", utils.OTCHOST,
		fmt.Sprintf(utils.OTC3instiCSV, fmt.Sprintf("%d/%02d/%02d", t.Date.Year()-1911, t.Date.Month(), t.Date.Day())))
}

// URL 擷取 自營商 網址
func (t TPEXT43U) URL() string {

	return fmt.Sprintf("%s%s", utils.OTCHOST,
		fmt.Sprintf(utils.OTC3instiCSV, fmt.Sprintf("%d/%02d/%02d", t.Date.Year()-1911, t.Date.Month(), t.Date.Day())))
}

// URL 擷取 投信 網址
func (t TPEXT44U) URL() string {

	return fmt.Sprintf("%s%s", utils.OTCHOST,
		fmt.Sprintf(utils.OTC3instiCSV, fmt.Sprintf("%d/%02d/%02d", t.Date.Year()-1911, t.Date.Month(), t.Date.Day())))
}

// URL 擷取 外資 網址
func (t TWT38U) URL() string {

	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.TWTXXU, "TWT38U", t.Date.Year(), t.Date.Month(), t.Date.Day()))
}

// URL 擷取 自營商 網址
func (t TWT43U) URL() string {

	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.TWTXXU, "TWT43U", t.Date.Year(), t.Date.Month(), t.Date.Day()))
}

// URL 擷取 投信 網址
func (t TWT44U) URL() string {

	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.TWTXXU, "TWT44U", t.Date.Year(), t.Date.Month(), t.Date.Day()))
}

func checkCsvDataFormat(t string, data []string) bool {

	switch t {
	case "TeckCsvDataFormatMTSS":
		return "賣出" == strings.Replace(data[9], " ", "", -1) &&
			"買進" == strings.Replace(data[8], " ", "", -1) &&
			"賣出" == strings.Replace(data[3], " ", "", -1) &&
			"買進" == strings.Replace(data[2], " ", "", -1) &&
			"股票名稱" == strings.Replace(data[1], " ", "", -1) &&
			"股票代號" == strings.Replace(data[0], " ", "", -1)
	case "TpexT38":
		return "三大法人買賣超股數合計" == strings.Replace(data[23], " ", "", -1) &&
			"外資及陸資-買賣超股數" == strings.Replace(data[10], " ", "", -1) &&
			"外資及陸資-賣出股數" == strings.Replace(data[9], " ", "", -1) &&
			"外資及陸資-買進股數" == strings.Replace(data[8], " ", "", -1) &&
			"名稱" == strings.Replace(data[1], " ", "", -1) &&
			"代號" == strings.Replace(data[0], " ", "", -1)
	case "TpexT43":
		return "三大法人買賣超股數合計" == strings.Replace(data[23], " ", "", -1) &&
			"自營商-買賣超股數" == strings.Replace(data[22], " ", "", -1) &&
			"自營商-賣出股數" == strings.Replace(data[21], " ", "", -1) &&
			"自營商-買進股數" == strings.Replace(data[20], " ", "", -1) &&
			"名稱" == strings.Replace(data[1], " ", "", -1) &&
			"代號" == strings.Replace(data[0], " ", "", -1)
	case "TpexT44":
		return "三大法人買賣超股數合計" == strings.Replace(data[23], " ", "", -1) &&
			"投信-買賣超股數" == strings.Replace(data[13], " ", "", -1) &&
			"投信-賣出股數" == strings.Replace(data[12], " ", "", -1) &&
			"投信-買進股數" == strings.Replace(data[11], " ", "", -1) &&
			"名稱" == strings.Replace(data[1], " ", "", -1) &&
			"代號" == strings.Replace(data[0], " ", "", -1)
	default:
		return true
	}
}

//台灣櫃買中心

//Get 從網頁抓取 外資 資料
func (t *TPEXT38U) Get() (map[string]BaseBuySellTotalData, error) {
	dateUnix := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location()).Unix()
	//檢查map快取裡面有沒有資料
	if v, ok := t.UnixData[dateUnix]; ok {
		return v, nil
	}
	var (
		csvdata   [][]string
		data      []byte
		err       error
		resultMap map[string]BaseBuySellTotalData
	)
	utils.Dbgln(t.URL())
	if data, err = hCache.Get(t.URL(), false); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	if len(csvArrayContent) < 3 {
		if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
			return nil, err
		}
		return nil, errorFileNoData
	}
	//從第八列開始 然後刪掉最後面的八行(注意可能會有空白的行)
	csvArrayContent = csvArrayContent[1:]

	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
		resultMap = make(map[string]BaseBuySellTotalData, len(csvdata))
		for i, v := range csvdata {
			if i == 0 {
				if false == checkCsvDataFormat("TpexT38", v) {
					return nil, errors.New("Wrong TpexT38 Csv Data Format")
				}
				continue
			}
			var r BaseBuySellTotalData
			no := strings.Replace(v[0], " ", "", -1)

			r.Name = strings.Replace(v[1], " ", "", -1)

			r.Volume.Buy, _ = strconv.ParseInt(strings.Replace(v[8], ",", "", -1), 10, 64)
			r.Volume.Sell, _ = strconv.ParseInt(strings.Replace(v[9], ",", "", -1), 10, 64)
			r.Volume.Total, _ = strconv.ParseInt(strings.Replace(v[10], ",", "", -1), 10, 64)

			resultMap[no] = r

		}
		t.UnixData[dateUnix] = resultMap
	}
	return resultMap, err
}

//Get 從網頁抓取 外資 資料
func (t *TPEXT44U) Get() (map[string]BaseBuySellTotalData, error) {
	dateUnix := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location()).Unix()
	//檢查map快取裡面有沒有資料
	if v, ok := t.UnixData[dateUnix]; ok {
		return v, nil
	}
	var (
		csvdata   [][]string
		data      []byte
		err       error
		resultMap map[string]BaseBuySellTotalData
	)
	utils.Dbgln(t.URL())
	if data, err = hCache.Get(t.URL(), false); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	if len(csvArrayContent) < 3 {
		if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
			return nil, err
		}
		return nil, errorFileNoData
	}
	//從第八列開始 然後刪掉最後面的八行(注意可能會有空白的行)
	csvArrayContent = csvArrayContent[1:]

	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
		resultMap = make(map[string]BaseBuySellTotalData, len(csvdata))
		for i, v := range csvdata {
			if i == 0 {
				if false == checkCsvDataFormat("TpexT44", v) {
					return nil, errors.New("Wrong TpexT44 Csv Data Format")
				}
				continue
			}
			var r BaseBuySellTotalData
			no := strings.Replace(v[0], " ", "", -1)

			r.Name = strings.Replace(v[1], " ", "", -1)

			r.Volume.Buy, _ = strconv.ParseInt(strings.Replace(v[11], ",", "", -1), 10, 64)
			r.Volume.Sell, _ = strconv.ParseInt(strings.Replace(v[12], ",", "", -1), 10, 64)
			r.Volume.Total, _ = strconv.ParseInt(strings.Replace(v[13], ",", "", -1), 10, 64)

			resultMap[no] = r

		}
		t.UnixData[dateUnix] = resultMap
	}
	return resultMap, err
}

//台灣證交所

//Get 從網頁抓取 外資 資料
func (t *TWT38U) Get() (map[string]BaseT38U, error) {
	dateUnix := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location()).Unix()
	if v, ok := t.UnixMapT38UData[dateUnix]; ok {
		return v, nil
	}
	var (
		csvdata   [][]string
		data      []byte
		err       error
		resultMap map[string]BaseT38U
	)
	//fmt.Println(t.URL())
	if data, err = hCache.PostForm(t.URL(), nil); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	if len(csvArrayContent) < 9 {
		if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
			return nil, err
		}
		return nil, errorFileNoData
	}
	//從第八列開始 然後刪掉最後面的八行(注意可能會有空白的行)
	csvArrayContent = csvArrayContent[3 : len(csvArrayContent)-9]

	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
		resultMap = make(map[string]BaseT38U, len(csvdata))
		for i, v := range csvdata {
			if i == 0 {
				if false == checkCsvDataFormat("T38U", v) {
					return nil, errors.New("Wrong T38U Csv Data Format")
				}
				continue
			}
			var r BaseT38U
			no := strings.Replace(v[1], " ", "", -1)

			r.Name = strings.Replace(v[2], " ", "", -1)

			r.Volume.Buy, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
			r.Volume.Sell, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
			r.Volume.Total, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)

			resultMap[no] = r

		}
		t.UnixMapT38UData[dateUnix] = resultMap
	}
	return resultMap, err
}
func (t *TWT43U) Get() (map[string]BaseT43U, error) {
	dateUnix := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location()).Unix()
	if v, ok := t.UnixMapT43UData[dateUnix]; ok {
		return v, nil
	}
	var (
		csvdata   [][]string
		data      []byte
		err       error
		resultMap map[string]BaseT43U
	)
	//fmt.Println(t.URL())
	if data, err = hCache.PostForm(t.URL(), nil); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	if len(csvArrayContent) < 5 {
		if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
			return nil, err
		}
		return nil, errorFileNoData
	}
	//從第八列開始 然後刪掉最後面的八行(注意可能會有空白的行)
	csvArrayContent = csvArrayContent[3 : len(csvArrayContent)-5]

	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
		resultMap = make(map[string]BaseT43U, len(csvdata))
		for i, v := range csvdata {
			if i == 0 {
				if false == checkCsvDataFormat("T43U", v) {
					return nil, errors.New("Wrong T43U Csv Data Format")
				}
				continue
			}
			var r BaseT43U
			no := strings.Replace(v[0], " ", "", -1)

			r.Name = strings.Replace(v[1], " ", "", -1)

			r.Volume.Buy, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
			r.Volume.Sell, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
			r.Volume.Total, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)

			resultMap[no] = r

		}
		t.UnixMapT43UData[dateUnix] = resultMap
	}
	return resultMap, err
}
func (t *TWT44U) Get() (map[string]BaseT44U, error) {
	dateUnix := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location()).Unix()
	if v, ok := t.UnixMapT44UData[dateUnix]; ok {
		return v, nil
	}
	var (
		csvdata   [][]string
		data      []byte
		err       error
		resultMap map[string]BaseT44U
	)
	//fmt.Println(t.URL())
	if data, err = hCache.PostForm(t.URL(), nil); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	if len(csvArrayContent) < 9 {
		if err := os.Remove(utils.GetMD5FilePath(t)); err != nil {
			return nil, err
		}
		return nil, errorFileNoData
	}
	//從第八列開始 然後刪掉最後面的八行(注意可能會有空白的行)
	csvArrayContent = csvArrayContent[2 : len(csvArrayContent)-9]

	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
		resultMap = make(map[string]BaseT44U, len(csvdata))
		for i, v := range csvdata {
			if i == 0 {
				if false == checkCsvDataFormat("T44U", v) {
					return nil, errors.New("Wrong T44U Csv Data Format")
				}
				continue
			}
			var r BaseT44U
			no := strings.Replace(v[1], " ", "", -1)

			r.Name = strings.Replace(v[2], " ", "", -1)

			r.Volume.Buy, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
			r.Volume.Sell, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
			r.Volume.Total, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)

			resultMap[no] = r

		}
		t.UnixMapT44UData[dateUnix] = resultMap
	}
	return resultMap, err
}
func (t TWTXXU) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.TWTXXU, t.fund, t.Date.Year(), t.Date.Month(), t.Date.Day()))

}

// Get 擷取資料
func (t TWTXXU) Get() ([][]BaseSellBuy, error) {
	var (
		csvdata  [][]string
		data     []byte
		datalist int
		err      error
		result   [][]BaseSellBuy
	)
	fmt.Println(t.URL())
	if data, err = hCache.PostForm(t.URL(), nil); err != nil {
		return nil, err
	}
	var csvArrayContent = strings.Split(string(data), "\n")
	switch t.fund {
	case "TWT38U":
		if len(csvArrayContent) < 9 {
			return nil, errorFileNoData
		}
		csvArrayContent = csvArrayContent[3 : len(csvArrayContent)-9]
		datalist = 1
	case "TWT43U":
		if len(csvArrayContent) < 5 {
			return nil, errorFileNoData
		}
		csvArrayContent = csvArrayContent[3 : len(csvArrayContent)-5]
		datalist = 3
	case "TWT44U":
		if len(csvArrayContent) < 9 {
			return nil, errorFileNoData
		}
		csvArrayContent = csvArrayContent[2 : len(csvArrayContent)-9]
		//fmt.Println(csvArrayContent)
		datalist = 1
	}

	for i, v := range csvArrayContent {

		csvArrayContent[i] = strings.Replace(v, "=", "", -1)

		//fmt.Println(csvArrayContent[i])
	}

	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {

		result = make([][]BaseSellBuy, len(csvdata))
		for i, v := range csvdata {
			var name, no string
			result[i] = make([]BaseSellBuy, datalist)
			switch {
			case datalist == 1:
				name = strings.Replace(v[2], " ", "", -1)
				no = strings.Replace(v[1], " ", "", -1)
			case datalist > 1:
				name = strings.Replace(v[1], " ", "", -1)
				no = strings.Replace(v[0], " ", "", -1)
			}
			result[i][0].Name = name
			result[i][0].No = no
			result[i][0].Buy, _ = strconv.ParseInt(
				strings.Replace(v[3], ",", "", -1), 10, 64)
			result[i][0].Sell, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
			result[i][0].Total, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)
			if datalist > 1 {
				result[i][1].Name = name
				result[i][1].No = no
				result[i][1].Buy, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)
				result[i][1].Sell, _ = strconv.ParseInt(strings.Replace(v[6], ",", "", -1), 10, 64)
				result[i][1].Total, _ = strconv.ParseInt(strings.Replace(v[7], ",", "", -1), 10, 64)
				result[i][2].Name = name
				result[i][2].No = no
				result[i][2].Buy, _ = strconv.ParseInt(strings.Replace(v[8], ",", "", -1), 10, 64)
				result[i][2].Sell, _ = strconv.ParseInt(strings.Replace(v[9], ",", "", -1), 10, 64)
				result[i][2].Total, _ = strconv.ParseInt(strings.Replace(v[10], ",", "", -1), 10, 64)
			}
		}
	}
	return result, err
}
