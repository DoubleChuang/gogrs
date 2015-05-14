package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

type base interface {
	MA(days int) []float64
	Len() int
	PlusData()
}

type check01 struct{}

func (check01) String() string {
	return "MA 3 > 6 > 18"
}

func (check01) CheckFunc(b ...base) bool {
	defer wg.Done()
	var d = b[0]
	var start = d.Len()
	for {
		if d.Len() >= 18 {
			break
		}
		d.PlusData()
		if (d.Len() - start) == 0 {
			break
		}
	}
	var ma3 = d.MA(3)
	var ma6 = d.MA(6)
	var ma18 = d.MA(18)
	//log.Println(ma3[len(ma3)-1], ma6[len(ma6)-1], ma18[len(ma18)-1])
	if ma3[len(ma3)-1] > ma6[len(ma6)-1] && ma6[len(ma6)-1] > ma18[len(ma18)-1] {
		return true
	}
	return false
}

type check02 struct{}

func (check02) String() string {
	return "check02"
}
func (check02) CheckFunc(b ...base) bool {
	defer wg.Done()
	days, up := utils.CountCountineFloat64(utils.DeltaFloat64(b[0].MA(3)))
	if up && days > 1 {
		return true
	}
	return false
}

type CheckGroup interface {
	String() string
	CheckFunc(...base) bool
}

var wg sync.WaitGroup
var twseNo = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")

func main() {
	flag.Parse()
	var datalist []*twse.Data

	if *twseNo != "" {
		twselist := strings.Split(*twseNo, ",")
		datalist = make([]*twse.Data, len(twselist))

		for i, no := range twselist {
			datalist[i] = twse.NewTWSE(no, tradingdays.FindRecentlyOpened(time.Now()))
		}
	}

	if len(datalist) > 0 {
		for _, check := range []CheckGroup{CheckGroup(check01{}), CheckGroup(check02{})} {
			fmt.Printf("----- %v -----\n", check)
			for _, stock := range datalist {
				wg.Add(1)
				go func(check CheckGroup, stock *twse.Data) {
					runtime.Gosched()
					if check.CheckFunc(stock) {
						fmt.Printf("%s\n", stock.No)
					}
				}(check, stock)
			}
			wg.Wait()
		}
	} else {
		flag.PrintDefaults()
	}
}
