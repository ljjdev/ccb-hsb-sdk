// Package config 提供了 SDK 的配置管理功能。
//
// 该包包含了客户端配置、证书配置、接口地址等配置项的定义和管理。
package config

import (
	"crypto/rsa"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/ljjdev/ccb-hsb-sdk/internal/utils"
)

// Config 定义了 SDK 的配置项
type Config struct {
	// MarketID 市场编号,由银行提供
	MarketID string

	// MerchantID 商家编号,由银行提供
	MerchantID string

	// GatewayURL 接口网关地址
	GatewayURL string

	// PrivateKey 商户私钥,用于签名
	PrivateKey *rsa.PrivateKey

	// PublicKey 银行公钥,用于验签
	PublicKey *rsa.PublicKey

	// Timeout HTTP 请求超时时间
	Timeout time.Duration

	// Debug 是否开启调试模式
	Debug bool
}

// NewConfig 创建一个新的配置实例
//
// 使用示例:
//
//	cfg, err := config.NewConfig(
//		config.WithMarketID("12345678901234"),
//		config.WithMerchantID("12345678901234567890"),
//		config.WithGatewayURL("https://marketpay.ccb.com/online/direct"),
//		config.WithPrivateKey(privateKey),
//		config.WithPublicKey(publicKey),
//	)
func NewConfig(opts ...Option) (*Config, error) {
	cfg := &Config{
		GatewayURL: "https://marketpay.ccb.com/online/direct",
		Timeout:    30 * time.Second,
		Debug:      false,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Option 定义了配置的选项函数类型
type Option func(*Config)

// WithMarketID 设置市场编号
func WithMarketID(marketID string) Option {
	return func(c *Config) {
		c.MarketID = marketID
	}
}

// WithMerchantID 设置商家编号
func WithMerchantID(merchantID string) Option {
	return func(c *Config) {
		c.MerchantID = merchantID
	}
}

// WithGatewayURL 设置网关地址
func WithGatewayURL(gatewayURL string) Option {
	return func(c *Config) {
		c.GatewayURL = gatewayURL
	}
}

// WithPrivateKey 设置私钥
func WithPrivateKey(privateKey *rsa.PrivateKey) Option {
	return func(c *Config) {
		c.PrivateKey = privateKey
	}
}

// WithPublicKey 设置公钥
func WithPublicKey(publicKey *rsa.PublicKey) Option {
	return func(c *Config) {
		c.PublicKey = publicKey
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithDebug 设置调试模式
func WithDebug(debug bool) Option {
	return func(c *Config) {
		c.Debug = debug
	}
}

// LoadConfigFromFile 从文件加载配置
//
// 该函数从指定的文件路径加载配置文件,文件格式为 JSON。
func LoadConfigFromFile(filepath string) (*Config, error) {
	// TODO: 实现从文件加载配置的逻辑
	return nil, errors.New("not implemented")
}

// LoadConfigFromEnv 从环境变量加载配置
//
// 该函数从环境变量中读取配置项,支持以下环境变量:
//   - CCB_MARKET_ID: 市场编号
//   - CCB_MERCHANT_ID: 商家编号
//   - CCB_GATEWAY_URL: 网关地址
//   - CCB_PRIVATE_KEY: 私钥BASE64字符串
//   - CCB_PUBLIC_KEY: 公钥BASE64字符串
//   - CCB_TIMEOUT: 超时时间(秒)
//   - CCB_DEBUG: 调试模式(true/false)
func LoadConfigFromEnv() (*Config, error) {
	// 1. 准备 RSA 密钥对
	// 实际使用时,加载私钥的base64字符串
	privateKey, err := utils.LoadPrivateKey(os.Getenv("CCB_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}
	// 实际使用时,加载公钥的base64字符串
	publicKey, err := utils.LoadPublicKey(os.Getenv("CCB_PUBLIC_KEY"))
	if err != nil {
		return nil, err
	}
	timeout := 30
	if os.Getenv("CCB_TIMEOUT") != "" {
		timeout, err = strconv.Atoi(os.Getenv("CCB_TIMEOUT"))
	}
	cfg := &Config{
		MarketID:   os.Getenv("CCB_MARKET_ID"),
		MerchantID: os.Getenv("CCB_MERCHANT_ID"),
		GatewayURL: os.Getenv("CCB_GATEWAY_URL"),
		Timeout:    time.Duration(timeout),
		Debug:      os.Getenv("CCB_DEBUG") == "true",
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	if cfg.GatewayURL == "" {
		cfg.GatewayURL = "https://marketpay.ccb.com/online/direct"
	}
	return cfg, cfg.Validate()
}

// Validate 验证配置的有效性
//
// 该函数检查配置项是否完整且有效,如果配置无效则返回错误。
func (c *Config) Validate() error {
	if c.MarketID == "" {
		return errors.New("market id is required")
	}

	if c.MerchantID == "" {
		return errors.New("merchant id is required")
	}

	if c.GatewayURL == "" {
		return errors.New("gateway url is required")
	}

	if c.PrivateKey == nil {
		return errors.New("private key is required")
	}

	if c.PublicKey == nil {
		return errors.New("public key is required")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}

	return nil
}
