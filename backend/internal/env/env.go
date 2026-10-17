package env

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	LOCAL = "local"
	DEV   = "dev"
	PROD  = "prod"
)

var AlcoholIndex = 99

type env struct {
	EnvUnique
	EnvOperate
}

type EnvUnique struct {
	StoreID          string  `mapstructure:"STORE_ID"`
	Environment      string  `mapstructure:"ENVIRONMENT"`
	StoreName        string  `mapstructure:"STORE_NAME"`
	StoreAddress     string  `mapstructure:"STORE_ADDRESS"`
	URL              string  `mapstructure:"URL"`
	CertPath         string  `mapstructure:"CERT_PATH"`
	Commission       float64 `mapstructure:"COMMISSION"`
	ChargeBasePrice  int     `mapstructure:"CHARGE_BASE_PRICE"`
	ChargeBaseMinute int     `mapstructure:"CHARGE_BASE_MINUTE"`
	IsPermitAlcohol  bool    `mapstructure:"IS_PERMIT_ALCOHOL"`
	IsEnableCharge   bool    `mapstructure:"IS_ENABLE_CHARGE"`
	Version          string  `mapstructure:"VERSION"`
	OpenTime         string  `mapstructure:"OPEN_TIME"`
	SignboardProduct []int   `mapstructure:"SIGNBOARD_PRODUCT"`
}

type EnvOperate struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	// アプリケーション自身の認証設定
	BasicAuthUserName string `mapstructure:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword string `mapstructure:"BASIC_AUTH_PASSWORD"`
	FromGateKey       string `mapstructure:"FROM_GATE_KEY"`
	FromDeviceKey     string `mapstructure:"FROM_DEVICE_KEY"`
	// スマレジ連携設定
	ContractId   string `mapstructure:"CONTRACT_ID"`
	ClientId     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	AuthURL      string `mapstructure:"AUTH_URL"`
	ApiBaseURL   string `mapstructure:"API_BASE_URL"`
	// 取引情報ゲートウェイ連携設定
	GateTransactionInfoBaseURL string `mapstructure:"GATE_TRANSACTION_INFO_BASE_URL"`
	GateTransactionInfoPath    string `mapstructure:"GATE_TRANSACTION_INFO_PATH"`
	// Azure Document Intelligence連携設定
	AzureDocumentIntelligenceURL string `mapstructure:"AZURE_DOCUMENT_INTELLIGENCE_URL"`
	AzureSubscriptionKey         string `mapstructure:"AZURE_SUBSCRIPTION_KEY"`
}

var Env env

// init は、パッケージの初期状態を構成します。
func init() {
	Env = env{
		EnvUnique:  EnvUnique{},
		EnvOperate: EnvOperate{},
	}

	initEnvUnique()
	initEnvOperate()
}

// initEnvUnique は、店舗IDに対応する店舗固有設定をYAMLと環境変数から読み込みます。
func initEnvUnique() {
	v := viper.New()

	// 必須の環境変数
	v.AutomaticEnv()
	storeID := v.GetString("STORE_ID")
	if storeID == "" {
		if isTestProcess() {
			storeID = "0"
		} else {
			log.Fatal("STORE_ID is not set")
		}
	}

	// 店舗固有設定
	v.SetConfigName(storeID + "_environment")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")
	v.AddConfigPath(filepath.Join(sourceDir(), "..", "..", "config"))

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	if err := v.Unmarshal(&Env.EnvUnique); err != nil {
		log.Fatalf("error un marshalling config: %v", err)
	}

	Env.EnvUnique.StoreID = storeID
	Env.EnvUnique.Version = time.Now().Format("2006-01-02-15-04-05")
}

// initEnvOperate は、認証情報と外部サービス設定を環境変数または機密設定ファイルから読み込みます。
func initEnvOperate() {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	secretKeys := []string{
		"BASIC_AUTH_USERNAME",
		"BASIC_AUTH_PASSWORD",
		"FROM_GATE_KEY",
		"FROM_DEVICE_KEY",
		"CONTRACT_ID",
		"CLIENT_ID",
		"CLIENT_SECRET",
		"AZURE_SUBSCRIPTION_KEY",
	}
	for _, key := range secretKeys {
		if err := v.BindEnv(key); err != nil {
			log.Fatalf("error binding environment variable %s: %v", key, err)
		}
	}
	if err := v.BindEnv("AZURE_DOCUMENT_INTELLIGENCE_URL"); err != nil {
		log.Fatalf("error binding environment variable AZURE_DOCUMENT_INTELLIGENCE_URL: %v", err)
	}
	for _, key := range []string{"GATE_TRANSACTION_INFO_BASE_URL", "GATE_TRANSACTION_INFO_PATH"} {
		if err := v.BindEnv(key); err != nil {
			log.Fatalf("error binding environment variable %s: %v", key, err)
		}
	}

	configPath := getConfigPath()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	if err := v.Unmarshal(&Env.EnvOperate); err != nil {
		log.Fatalf("error un marshalling config: %v", err)
	}

	// 環境変数の値をファイル設定より優先します。
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		Env.ServerPort = serverPort
	}

	if !isTestProcess() {
		for _, key := range secretKeys {
			if v.GetString(key) == "" {
				log.Fatalf("%s is not set", key)
			}
		}
	}

	if Env.EnvUnique.Environment == LOCAL && !isTestProcess() {
		if ipAddress, err := GetOutboundIP(); err == nil {
			Env.EnvUnique.URL = fmt.Sprintf("%s:6789", ipAddress)
		}
	}
}

// getConfigPath は、実行環境に対応する機密設定ファイルのパスを返します。
func getConfigPath() string {
	fileName := "local.yaml"
	switch Env.EnvUnique.Environment {
	case DEV:
		fileName = "dev.yaml"
	case PROD:
		fileName = "prod.yaml"
	}

	return filepath.Join(sourceDir(), fileName)
}

// sourceDir は、このパッケージのソースディレクトリーを返します。
func sourceDir() string {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		return "internal/env"
	}

	dir := filepath.Dir(fileName)
	if filepath.IsAbs(dir) {
		return dir
	}

	// -trimpathで相対パスになる場合は、go.modを基準に絶対パスへ戻す。
	workingDir, err := os.Getwd()
	if err == nil {
		for currentDir := workingDir; ; currentDir = filepath.Dir(currentDir) {
			if _, statErr := os.Stat(filepath.Join(currentDir, "go.mod")); statErr == nil {
				return filepath.Join(currentDir, "internal", "env")
			}

			parentDir := filepath.Dir(currentDir)
			if parentDir == currentDir {
				break
			}
		}
	}

	return dir
}

// isTestProcess は、現在のプロセスがGoテストとして実行中かを判定します。
func isTestProcess() bool {
	return flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test")
}

// GetOutboundIP は、外部接続に使用されるローカルIPアドレスを取得します。
func GetOutboundIP() (string, error) {
	conn, err := (&net.Dialer{}).DialContext(context.Background(), "udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", fmt.Errorf("unexpected local address type %T", conn.LocalAddr())
	}
	return localAddr.IP.String(), nil
}
