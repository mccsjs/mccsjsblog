package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ipAPIURL = "http://ip-api.com/json/%s?lang=zh-CN"
	ipMu     sync.RWMutex
)

// SetIPApiURL 设置 IP 归属地查询 API URL
func SetIPApiURL(url string) {
	if url != "" {
		ipMu.Lock()
		ipAPIURL = url
		ipMu.Unlock()
	}
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

type ipAPIResponse struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
	Message    string `json:"message"`
}

// 省份和直辖市中文映射
var provinceCNMap = map[string]string{
	"Beijing":      "北京",
	"Shanghai":     "上海",
	"Tianjin":      "天津",
	"Chongqing":    "重庆",
	"Guangdong":    "广东",
	"Jiangsu":      "江苏",
	"Zhejiang":     "浙江",
	"Shandong":     "山东",
	"Henan":        "河南",
	"Sichuan":      "四川",
	"Hubei":        "湖北",
	"Hunan":        "湖南",
	"Hebei":        "河北",
	"Fujian":       "福建",
	"Anhui":        "安徽",
	"Shaanxi":      "陕西",
	"Liaoning":     "辽宁",
	"Jiangxi":      "江西",
	"Yunnan":       "云南",
	"Guangxi":      "广西",
	"Guizhou":      "贵州",
	"Shanxi":       "山西",
	"Jilin":        "吉林",
	"Heilongjiang": "黑龙江",
	"Inner Mongolia": "内蒙古",
	"Xinjiang":     "新疆",
	"Gansu":        "甘肃",
	"Hainan":       "海南",
	"Ningxia":      "宁夏",
	"Qinghai":      "青海",
	"Tibet":        "西藏",
	"Hong Kong":    "香港",
	"Macau":        "澳门",
	"Taiwan":       "台湾省",
}

// translateProvince 翻译省份名称为中文
func translateProvince(province string) string {
	if cn, ok := provinceCNMap[province]; ok {
		return cn
	}
	return province
}

// InitIPSearcher 初始化 IP 地址搜索器（在线模式无需初始化）
func InitIPSearcher(_ string) error {
	return nil
}

// CloseIPSearcher 关闭 IP 搜索器（在线模式无需关闭）
func CloseIPSearcher() {}

// GetIPLocation 获取 IP 地址的地理位置
func GetIPLocation(ip string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return "未知"
	}

	if ipAddr.IsLoopback() || ipAddr.IsPrivate() {
		return "本地"
	}

	resp, err := httpClient.Get(fmt.Sprintf(ipAPIURL, ip))
	if err != nil {
		return "未知"
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "未知"
	}

	if result.Status != "success" {
		return "未知"
	}

	// 统一转为中文映射
	country := translateProvince(result.Country)
	region := translateProvince(result.RegionName)

	// 特殊地区处理
	switch country {
	case "香港", "澳门", "台湾省":
		return country
	}

	// 中国大陆：只显示省份/直辖市
	if country == "中国" {
		if region != "" {
			return region
		}
		return "中国"
	}

	// 国外：显示国家 + 城市
	var parts []string
	if country != "" {
		parts = append(parts, country)
	}
	if result.City != "" && result.City != result.RegionName && result.City != result.Country {
		parts = append(parts, result.City)
	}

	if len(parts) == 0 {
		return "未知"
	}
	return strings.Join(parts, "")
}