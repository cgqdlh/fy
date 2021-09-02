package fy

type Result struct {
	From        string        `json:"from"`
	Target      string        `json:"target"`
	TransResult []TransResult `json:"trans_result"`
}

type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type Translator interface {
	Translation(query string) (Result, error)
}

type Config struct {
	Log   string
	Baidu BaiduConfig
}

type BaiduConfig struct {
	Appid  string
	Secret string
}
