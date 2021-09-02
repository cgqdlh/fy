package fy

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type baidu struct {
	from   string
	target string
}

type baiduResult struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []TransResult `json:"trans_result"`
	ErrorCode   string        `json:"error_code"`
	SrcTts      string        `json:"src_tts"`
	DstTts      string        `json:"dst_tts"`
	Dict        string        `json:"dict"`
}

func NewBaidu(from, target string) baidu {
	if from == "" {
		from = "auto"
	}
	if target == "" {
		target = "zh"
	}
	return baidu{from, target}
}

func zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func (b *baidu) Translation(query string) (Result, error) {
	url := "https://fanyi-api.baidu.com/api/trans/vip/translate"
	salt := strconv.FormatInt(time.Now().UnixMilli(), 16)
	appid, sign := baiduSign(query, salt)
	reqData := fmt.Sprintf("appid=%s&q=%s&salt=%s&sign=%s&from=%s&to=%s", appid, query, salt, sign, b.from, b.target)
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(reqData))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	if response.StatusCode != http.StatusOK {
		return Result{}, errors.New("response error. status: " + response.Status)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Result{}, err
	}
	if body, err = zhToUnicode(body); err != nil {
		return Result{}, err
	}
	result := new(baiduResult)
	if err = json.Unmarshal(body, result); err != nil {
		return Result{}, err
	}

	if result.ErrorCode != "" {
		return Result{}, errors.New(fmt.Sprintf("translation error. code: %s", result.ErrorCode))
	}
	return Result{
		From:        result.From,
		Target:      result.To,
		TransResult: result.TransResult,
	}, nil
}

func baiduSign(query string, salt string) (string, string) {
	var (
		hash = md5.New()
		buf  bytes.Buffer
	)
	cfg := GetConfig()
	appid := cfg.Baidu.Appid
	secret := cfg.Baidu.Secret
	buf.WriteString(appid)
	buf.WriteString(query)
	buf.WriteString(salt)
	buf.WriteString(secret)
	hash.Write(buf.Bytes())
	return appid, hex.EncodeToString(hash.Sum(nil))
}
