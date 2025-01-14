package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

const (
	ayameVersion = "2021.2.1"
	// timeout は暫定的に 10 sec
	readHeaderTimeout = 10 * time.Second
)

var (
	config          *ayameConfig
	logger          *zerolog.Logger
	signalingLogger *zerolog.Logger
	webhookLogger   *zerolog.Logger
)

// 初期化処理
func init() {
	testing.Init()
	configFilePath := flag.String("c", "./ayame.yaml", "ayame の設定ファイルへのパス(yaml)")
	flag.Parse()
	// yaml ファイルを読み込み
	buf, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		// 読み込めない場合 Fatal で終了
		log.Fatal("Cannot open config file, err=", err)
	}
	if err := yaml.Unmarshal(buf, &config); err != nil {
		log.Fatal("Cannot parse config file, err=", err)
	}

	// グローバルの logger に代入する
	logger, err = initLogger()
	if err != nil {
		log.Fatal(err)
	}

	// バージョンをロギング
	logger.Info().Str("version", ayameVersion).Msg("AyameVersion")

	setDefaultsConfig()

	// グローバルの signalingLogger に代入する
	signalingLogger, err = initSignalingLogger()
	if err != nil {
		log.Fatal(err)
	}

	webhookLogger, err = initWebhookLogger()
	if err != nil {
		log.Fatal(err)
	}

	if config.AuthnWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.AuthnWebhookURL); err != nil {
			log.Fatal(err)
		}
	}

	if config.DisconnectWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.DisconnectWebhookURL); err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	args := flag.Args()
	// 引数の処理
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s", ayameVersion)
			return
		}
	}

	// URL の生成
	url := fmt.Sprintf("%s:%d", config.ListenIPv4Address, config.ListenPortNumber)

	go server()

	http.HandleFunc("/signaling", func(w http.ResponseWriter, r *http.Request) {
		signalingHandler(w, r)
	})
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: readHeaderTimeout}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err)
	}
}
