// Package config 的单元测试
package config

import (
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "valid config",
			opts: []Option{
				WithMarketID("12345678901234"),
				WithMerchantID("12345678901234567890"),
				WithGatewayURL("https://marketpay.ccb.com/online/direct"),
				WithPrivateKey(&rsa.PrivateKey{}),
				WithPublicKey(&rsa.PublicKey{}),
			},
			wantErr: false,
		},
		{
			name: "missing market id",
			opts: []Option{
				WithMerchantID("12345678901234567890"),
				WithGatewayURL("https://marketpay.ccb.com/online/direct"),
				WithPrivateKey(&rsa.PrivateKey{}),
				WithPublicKey(&rsa.PublicKey{}),
			},
			wantErr: true,
		},
		{
			name: "missing merchant id",
			opts: []Option{
				WithMarketID("12345678901234"),
				WithGatewayURL("https://marketpay.ccb.com/online/direct"),
				WithPrivateKey(&rsa.PrivateKey{}),
				WithPublicKey(&rsa.PublicKey{}),
			},
			wantErr: true,
		},
		{
			name: "missing private key",
			opts: []Option{
				WithMarketID("12345678901234"),
				WithMerchantID("12345678901234567890"),
				WithGatewayURL("https://marketpay.ccb.com/online/direct"),
				WithPublicKey(&rsa.PublicKey{}),
			},
			wantErr: true,
		},
		{
			name: "missing public key",
			opts: []Option{
				WithMarketID("12345678901234"),
				WithMerchantID("12345678901234567890"),
				WithGatewayURL("https://marketpay.ccb.com/online/direct"),
				WithPrivateKey(&rsa.PrivateKey{}),
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			opts: []Option{
				WithMarketID("12345678901234"),
				WithMerchantID("12345678901234567890"),
				WithGatewayURL("https://marketpay.ccb.com/online/direct"),
				WithPrivateKey(&rsa.PrivateKey{}),
				WithPublicKey(&rsa.PublicKey{}),
				WithTimeout(-1 * time.Second),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig(tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				MarketID:   "12345678901234",
				MerchantID: "12345678901234567890",
				GatewayURL: "https://marketpay.ccb.com/online/direct",
				PrivateKey: &rsa.PrivateKey{},
				PublicKey:  &rsa.PublicKey{},
				Timeout:    30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty market id",
			cfg: &Config{
				MerchantID: "12345678901234567890",
				GatewayURL: "https://marketpay.ccb.com/online/direct",
				PrivateKey: &rsa.PrivateKey{},
				PublicKey:  &rsa.PublicKey{},
				Timeout:    30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "empty merchant id",
			cfg: &Config{
				MarketID:   "12345678901234",
				GatewayURL: "https://marketpay.ccb.com/online/direct",
				PrivateKey: &rsa.PrivateKey{},
				PublicKey:  &rsa.PublicKey{},
				Timeout:    30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "empty gateway url",
			cfg: &Config{
				MarketID:   "12345678901234",
				MerchantID: "12345678901234567890",
				PrivateKey: &rsa.PrivateKey{},
				PublicKey:  &rsa.PublicKey{},
				Timeout:    30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "nil private key",
			cfg: &Config{
				MarketID:   "12345678901234",
				MerchantID: "12345678901234567890",
				GatewayURL: "https://marketpay.ccb.com/online/direct",
				PublicKey:  &rsa.PublicKey{},
				Timeout:    30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "nil public key",
			cfg: &Config{
				MarketID:   "12345678901234",
				MerchantID: "12345678901234567890",
				GatewayURL: "https://marketpay.ccb.com/online/direct",
				PrivateKey: &rsa.PrivateKey{},
				Timeout:    30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			cfg: &Config{
				MarketID:   "12345678901234",
				MerchantID: "12345678901234567890",
				GatewayURL: "https://marketpay.ccb.com/online/direct",
				PrivateKey: &rsa.PrivateKey{},
				PublicKey:  &rsa.PublicKey{},
				Timeout:    -1 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// TODO: 实现环境变量加载的测试
	t.Skip("not implemented")
}
