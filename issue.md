# 1 
	tradingdays/list.csv 更新問題
	https://www.twse.com.tw/holidaySchedule/holidaySchedule?response=csv&queryYear=108
# 2
	[Done] Cache GET timeout問題
	A:因為TWSE不能頻繁要資訊 需要延遲三秒
# 3
	拿到的檔案會用網址轉成md5檔名存成檔案,
	但如果剛好是每月的一號則下個月拿資料時會出錯,
	必須刪除快取的檔案並更新最新資料
	A:每次發現不是當月一號的資料則刪掉當月一號的檔案在下載
# 4
	cache資料夾需要分層不然會有資料上限問題
# 5
	[Done] T38資料獲取有錯誤 可能需要檢查時間並修改Parser格式 
# 6
	[Done] T38資料獲取可能空白
# 7
    [Done] 如果在收盤後資料尚未上傳就抓取資料 可能會有空白資料倒置會有抓不足資料的>問題
    需要再GET完之後檢查是否有此問題刪掉檔案 再抓幾次 如果依舊不行就提示 是否資
    料尚未更新 或是改成抓前一天
# 8
    [Done] prepareStock 修改成回err而不是bool
# 9 
	[ErrorHandler](https://github.com/EthanCai/goErrorHandlingSample)
#10 
	自動抓完發到LINE message 或是傳Email
