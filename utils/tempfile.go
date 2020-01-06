package utils

import (
	"archive/tar"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	iconv "github.com/djimenez/iconv-go"
)

// TempFolderName 快取資料夾名稱
const TempFolderName = ".gogrscache"

//與TWSE訪問需間隔的時間差
const TWSEDURTION = 5

//與OTC訪問需間隔的時間差
const OTCDURTION = 0

//上一次造訪TWSE時間
var visitTwseTime time.Time = time.Now()

//上一次造訪OTC時間
var visitOtcTime time.Time = time.Now()

// HTTPCache net/http 快取功能
type HTTPCache struct {
	Dir            string
	fullpath       string
	iconvConverter func([]byte) []byte
}

func Dbgln(args ...interface{}) {
	programCounter, _, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	//prefix := fmt.Sprintf("[%s:%s %d] %s", file, fn.Name(), line, fmt_)
	prefix := fmt.Sprintf("[%s %d]", fn.Name(), line)

	fmt.Printf("%s", prefix)
	fmt.Println(args...)
}
func Dbg(fmt_ string, args ...interface{}) {
	programCounter, _, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	//prefix := fmt.Sprintf("[%s:%s %d] %s", file, fn.Name(), line, fmt_)
	prefix := fmt.Sprintf("[%s %d] %s", fn.Name(), line, fmt_)
	fmt.Printf(prefix, args...)
	fmt.Println()
}

// NewHTTPCache New 一個 HTTPCache.
//
// dir 為暫存位置，fromEncoding 來源檔案的編碼，一律轉換為 utf8
func NewHTTPCache(dir string, fromEncoding string) *HTTPCache {
	dir = makeCacheDir(dir)
	return &HTTPCache{
		Dir:            dir,
		fullpath:       filepath.Join(dir, TempFolderName),
		iconvConverter: renderIconvConverter(fromEncoding)}
}

// makeCacheDir 建立快取資料夾
func makeCacheDir(dir string) string {
	var fullpath = filepath.Join(dir, TempFolderName)

	if err := os.Mkdir(fullpath, 0700); os.IsNotExist(err) {
		dir = os.TempDir()
		fullpath = filepath.Join(os.TempDir(), TempFolderName)
		os.Mkdir(fullpath, 0700)
	}
	return dir
}

//checkAndSyncVisitTime 在 Post 或 Get 時檢查是否太頻繁
//訪問伺服器 如果時間小於過設定的時間差距則等待並同步這
//次存取的時間

func checkAndSyncVisitTime(urlWeb string) {
	rand.Seed(time.Now().UnixNano())
	ms := rand.Intn(1000)
	switch urlWeb {
	case "twse":
		t := time.Now().Sub(visitTwseTime)
		if t.Seconds() < TWSEDURTION*time.Second.Seconds() {
			//fmt.Println("Sleep:",TWSEDURTION * time.Second-t)
			time.Sleep(TWSEDURTION*time.Second - t + time.Duration(ms)*time.Millisecond)
		}
		visitTwseTime = time.Now()
	case "otc":
		t := time.Now().Sub(visitOtcTime)
		if t.Seconds() < OTCDURTION*time.Second.Seconds() {
			//fmt.Println("Sleep:",TWSEDURTION * time.Second-t)
			time.Sleep(TWSEDURTION*time.Second - t)
		}
		visitOtcTime = time.Now()
	default:
	}
}
func whereUrl(url string) string {
	if strings.Contains(url, TWSEHOST) {
		return "twse"
	} else if strings.Contains(url, OTCHOST) {
		return "otc"
	} else {
		return "Unknown"
	}
}

// Get 透過 http.Get 取得檔案或從暫存中取得檔案
//
// rand 為是否支援網址帶入亂數值，url 需有 '%d' 格式。
func (hc HTTPCache) Get(url string, rand bool) ([]byte, error) {
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	var (
		content []byte
		err     error
	)

	//fmt.Printf("file:%s%s/%s\n", GetOSRamdiskPath(""), TempFolderName, filehash)
	if content, err = hc.readFile(filehash); err != nil {
		checkAndSyncVisitTime(whereUrl(url))
		return hc.saveFile(url, filehash, rand, nil)
	} else {
		csvArrayContent := strings.Split(string(content), "\n")
		if len(csvArrayContent) < 1 {
			err = os.Remove(filepath.Join(hc.fullpath, filehash))
			if err != nil {
				return nil, err
			} else {
				return hc.saveFile(url, filehash, rand, nil)
			}
		}
	}
	return content, nil
}

// PostForm 透過 http.PostForm 取得檔案或從暫存中取得檔案
func (hc HTTPCache) PostForm(url string, data url.Values) ([]byte, error) {
	hash := md5.New()
	io.WriteString(hash, url)
	io.WriteString(hash, data.Encode())

	var (
		content []byte
		err     error
	)

	filehash := fmt.Sprintf("%x", hash.Sum(nil))

	//Dbg("file:%s%s/%s\nurl:%s\n", GetOSRamdiskPath(""), TempFolderName, filehash, url)
	if content, err = hc.readFile(filehash); err != nil {
		checkAndSyncVisitTime(whereUrl(url))

		return hc.saveFile(url, filehash, false, data)
	}
	return content, nil
}

// FlushAll 清除快取
func (hc *HTTPCache) FlushAll() {
	os.RemoveAll(hc.fullpath)
	hc.Dir = makeCacheDir(hc.Dir)
}

// readFile 從快取資料裡面取得
func (hc HTTPCache) readFile(filehash string) ([]byte, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(filepath.Join(hc.fullpath, filehash)); err == nil {
		defer f.Close()
		return ioutil.ReadAll(f)
	}
	return nil, err
}

// HTTPClient is default func and fixed http too many open files.
var HTTPClient = &http.Client{Transport: &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   0,
		KeepAlive: 0,
	}).Dial,
	TLSHandshakeTimeout: 1 * time.Second,
},
}

// saveFile 從網路取得資料後放入快取資料夾
func (hc HTTPCache) saveFile(url, filehash string, rand bool, data url.Values) ([]byte, error) {
	if rand {
		url = fmt.Sprintf(url, RandInt())
	}

	var (
		content []byte
		err     error
		f       *os.File
		out     []byte
		req     *http.Request
		resp    *http.Response
	)

	if len(data) == 0 {
		// http.Get
		req, err = http.NewRequest("GET", url, nil)
	} else {
		// http.PostForm
		req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return out, err
	}

	req.Header.Set("Connection", "close")
	if resp, err = HTTPClient.Do(req); err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		if content, err = ioutil.ReadAll(resp.Body); err != nil {
			return out, err
		}

		if f, err = os.Create(filepath.Join(hc.fullpath, filehash)); err != nil {
			return out, err
		}
		defer f.Close()

		out = hc.iconvConverter(content)
		f.Write(out)
	} else {
		return out, fmt.Errorf("Fail Get File: [%d] %s\n", resp.StatusCode, url)
	}

	return out, err
}

// renderIconvConverter wrapper function for iconv converter.
func renderIconvConverter(fromEncoding string) func([]byte) []byte {
	if fromEncoding == "utf8" || fromEncoding == "utf-8" {
		return func(str []byte) []byte {
			return str
		}
	}
	return func(content []byte) []byte {
		converter, _ := iconv.NewConverter(fromEncoding, "utf-8")
		var out []byte
		out = make([]byte, len(content)*2)
		_, outLen, _ := converter.Convert(content, out)
		return out[:outLen]
	}
}

// GetOSRamdiskPath try to get RamDisk path.
func GetOSRamdiskPath(goos string) string {
	switch goos {
	case "darwin":
		return "/Volumes/RamDisk/"
	case "linux":
		return "/run/shm/"
	case "":
		return GetOSRamdiskPath(runtime.GOOS)
	default:
		return os.TempDir()
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status < 200 || status >= 300 {
		return fmt.Errorf("Fail Get File: [%d] %s\n", status, url)
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Untar(dst string, r io.Reader) error {
	/*gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)*/

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)
		Dbgln(target)
		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
func RecoveryStockBackup(date string) error {
	dstFilePath := fmt.Sprintf("./dl_%s.tar", date)
	if _, err := os.Stat(dstFilePath); os.IsNotExist(err) {
		// file does not exist
		srcUrl := fmt.Sprintf("https://github.com/DoubleChuang/gogrs/raw/master/STOCK_DATA_BACKUP/%s.tar", date)

		if err := DownloadFile(dstFilePath, srcUrl); err != nil {
			return err
		}

		if file, err := os.Open(dstFilePath); err == nil {
			defer file.Close()
			var fileReader io.ReadCloser = file
			Untar(GetOSRamdiskPath("")+"../../", fileReader)
		} else {
			return err
		}
	}
	return nil
}
