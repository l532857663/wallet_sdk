package elastic

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	elasticHost = "https://localhost:9200"
	defaultType = "_doc"
	username    = ""
	password    = ""
	isHttps     bool
	caCrtPath   string
	httpClient  *http.Client
)

func InitElasticInfo(conf ElasticConfig) {
	elasticHost = conf.Host
	username = conf.Username
	password = conf.Password
	caCrtPath = conf.CaCrt
	// 某些https节点配置需要做一些特殊处理
	if strings.HasPrefix(elasticHost, "https://") {
		isHttps = true
		elasticHost = strings.TrimPrefix(elasticHost, "https://")
	}
	InitHttps()
}

func InitHttps() {
	timeout := 100 * time.Second
	transport := &http.Transport{
		// 设置为短连接请求模式
		DisableKeepAlives: false,
	}
	if isHttps {
		// 读取根证书文件
		caCert, err := os.ReadFile(caCrtPath)
		if err != nil {
			log.Fatalln("无法读取根证书文件:", err)
		}
		// 创建根证书池
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		transport.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
		log.Printf("初始化[https]客户端")
	} else {
		log.Printf("初始化[http]客户端")
	}
	// 创建自定义的 http.Client
	httpClient = &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

type UrlFilter struct {
	Index  string
	Type   string
	Id     string
	Action string
}

func GetElasticUrl(filter UrlFilter) string {
	if filter.Type == "" {
		filter.Type = defaultType
	}
	url := elasticHost + "/" + filepath.Join(filter.Index, filter.Type, filter.Id, filter.Action) // + "?pretty"
	// fmt.Printf("wch------ url: %+v\n", url)
	return url
}

func AskHttp(filter UrlFilter, respBody interface{}) error {
	url := GetElasticUrl(filter)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("http.NewRequest error", err)
		return err
	}
	return HttpRequest(req, respBody)
}

func AskHttpJson(method string, filter UrlFilter, reqBody, respBody interface{}) error {
	url := GetElasticUrl(filter)
	var content []byte
	if reqBody != nil {
		content, _ = json.Marshal(reqBody)
	}
	// fmt.Printf("wch------ askContent: %+v\n", string(content))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	return HttpRequest(req, respBody)
}

func HttpRequest(req *http.Request, respBody interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("httpClient.Do error", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("io.ReadAll error", err)
		return err
	}
	// fmt.Printf("wch------ body: %+v\n", string(body))
	err = json.Unmarshal(body, respBody)
	if err != nil {
		log.Println("json.Unmarshal error", err)
		return err
	}
	return nil
}
