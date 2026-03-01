// Package client 的单元测试
package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ljjdev/ccb-hsb-sdk/pkg/config"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// maxTimeoutAttempts 最大超时尝试次数
const maxTimeoutAttempts = 5

// testTimeout 测试超时时间
const testTimeout = 5 * time.Second

// skipTestOnTimeout 如果测试超时超过指定次数则跳过
func skipTestOnTimeout(t *testing.T, timeoutCount *int32) {
	if atomic.LoadInt32(timeoutCount) >= maxTimeoutAttempts {
		t.Skipf("测试超时超过 %d 次，跳过该测试", maxTimeoutAttempts)
	}
}

// runWithTimeout 运行测试并处理超时
func runWithTimeout(t *testing.T, testName string, testFunc func(t *testing.T)) {
	timeoutCount := new(int32)

	// 使用 defer 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("测试 %s 发生 panic: %v", testName, r)
		}
	}()

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("测试 %s 发生 panic: %v", testName, r)
			}
		}()
		testFunc(t)
		done <- true
	}()

	select {
	case <-done:
		// 测试正常完成
	case <-ctx.Done():
		// 测试超时
		count := atomic.AddInt32(timeoutCount, 1)
		t.Errorf("测试 %s 超时 (第 %d 次)", testName, count)
		skipTestOnTimeout(t, timeoutCount)
	}
}

// generateTestKeys 生成测试用的RSA密钥对
func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

// createTestClient 创建测试客户端
func createTestClient(t *testing.T) *Client {
	privateKey, publicKey := generateTestKeys(t)

	cfg, err := config.NewConfig(
		config.WithMarketID("12345678901234"),
		config.WithMerchantID("12345678901234567890"),
		config.WithGatewayURL("https://marketpay.ccb.com/online/direct"),
		config.WithPrivateKey(privateKey),
		config.WithPublicKey(publicKey),
	)
	require.NoError(t, err)

	client, err := NewClient(cfg)
	require.NoError(t, err)
	return client
}

// setupTestClientWithKey 使用指定的密钥创建测试客户端
func setupTestClientWithKey(t *testing.T, handler http.HandlerFunc, privateKey *rsa.PrivateKey) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)

	cfg := &config.Config{
		MarketID:   "12345678901234",
		MerchantID: "12345678901234567890",
		GatewayURL: server.URL,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Timeout:    30,
		Debug:      false,
	}

	client, err := NewClient(cfg)
	require.NoError(t, err)

	return client, server
}

func TestQueryRefund(t *testing.T) {
	runWithTimeout(t, "TestQueryRefund", func(t *testing.T) {
		client := createTestClient(t)

		tests := []struct {
			name    string
			req     *model.QueryRefundRequest
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid request with customer refund trace no",
				req: &model.QueryRefundRequest{
					IttpartyStmId:  "00000",
					PyChnlCd:       "0000000000000000000000000",
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123001",
					MktId:          "12345678901234",
					CustRfndTrcno:  "REFUND20240101120000123",
					Vno:            "4",
				},
				wantErr: false,
			},
			{
				name: "valid request with refund trace no",
				req: &model.QueryRefundRequest{
					IttpartyStmId:  "00000",
					PyChnlCd:       "0000000000000000000000000",
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123001",
					MktId:          "12345678901234",
					RfndTrcno:      "RFND20240101120000123",
					Vno:            "4",
				},
				wantErr: false,
			},
			{
				name: "missing both customer refund trace no and refund trace no",
				req: &model.QueryRefundRequest{
					IttpartyStmId:  "00000",
					PyChnlCd:       "0000000000000000000000000",
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123001",
					MktId:          "12345678901234",
					Vno:            "4",
				},
				wantErr: true,
				errMsg:  "Cust_Rfnd_Trcno and Rfnd_Trcno must provide at least one",
			},
			{
				name: "request with default values",
				req: &model.QueryRefundRequest{
					IttpartyJrnlNo: "20240101120000123001",
					CustRfndTrcno:  "REFUND20240101120000123",
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// 注意: 由于这是单元测试,我们只测试请求构建和签名生成
				// 实际的HTTP请求需要mock服务器,这里我们只测试参数验证和签名生成

				// 验证必输参数
				if tt.req.CustRfndTrcno == "" && tt.req.RfndTrcno == "" {
					// 这种情况应该返回错误
					_, err := client.QueryRefund(context.Background(), tt.req)
					if tt.wantErr {
						assert.Error(t, err)
						if tt.errMsg != "" {
							assert.Contains(t, err.Error(), tt.errMsg)
						}
					}
				} else {
					// 这种情况下,我们只测试请求能够正确构建
					// 实际发送会失败,因为服务器不存在
					_, err := client.QueryRefund(context.Background(), tt.req)
					// 我们期望请求构建成功,但HTTP请求会失败
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "failed to send request")
				}
			})
		}
	})
}

func TestQueryRefundRequest_ToMap(t *testing.T) {
	runWithTimeout(t, "TestQueryRefundRequest_ToMap", func(t *testing.T) {
		req := &model.QueryRefundRequest{
			IttpartyStmId:  "00000",
			PyChnlCd:       "0000000000000000000000000",
			IttpartyTms:    "20240101120000123",
			IttpartyJrnlNo: "20240101120000123001",
			MktId:          "12345678901234",
			CustRfndTrcno:  "REFUND20240101120000123",
			Vno:            "4",
		}

		result, err := req.ToMap()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "00000", result["Ittparty_Stm_Id"])
		assert.Equal(t, "0000000000000000000000000", result["Py_Chnl_Cd"])
		assert.Equal(t, "20240101120000123", result["Ittparty_Tms"])
		assert.Equal(t, "20240101120000123001", result["Ittparty_Jrnl_No"])
		assert.Equal(t, "12345678901234", result["Mkt_Id"])
		assert.Equal(t, "REFUND20240101120000123", result["Cust_Rfnd_Trcno"])
		assert.Equal(t, "4", result["Vno"])
	})
}

func TestQueryRefundResponse_IsSuccess(t *testing.T) {
	runWithTimeout(t, "TestQueryRefundResponse_IsSuccess", func(t *testing.T) {
		tests := []struct {
			name string
			resp *model.QueryRefundResponse
			want bool
		}{
			{
				name: "success response",
				resp: &model.QueryRefundResponse{
					RefundRspSt: model.RefundStatusSuccess,
				},
				want: true,
			},
			{
				name: "failed response",
				resp: &model.QueryRefundResponse{
					RefundRspSt: model.RefundStatusFailed,
				},
				want: false,
			},
			{
				name: "delayed response",
				resp: &model.QueryRefundResponse{
					RefundRspSt: model.RefundStatusDelayed,
				},
				want: false,
			},
			{
				name: "uncertain response",
				resp: &model.QueryRefundResponse{
					RefundRspSt: model.RefundStatusUncertain,
				},
				want: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.resp.IsSuccess()
				assert.Equal(t, tt.want, got)
			})
		}
	})
}

func TestQueryRefundResponse_GetError(t *testing.T) {
	runWithTimeout(t, "TestQueryRefundResponse_GetError", func(t *testing.T) {
		tests := []struct {
			name string
			resp *model.QueryRefundResponse
			want bool
		}{
			{
				name: "success response",
				resp: &model.QueryRefundResponse{
					RefundRspSt: model.RefundStatusSuccess,
				},
				want: false,
			},
			{
				name: "failed response",
				resp: &model.QueryRefundResponse{
					RefundRspSt: model.RefundStatusFailed,
				},
				want: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.resp.GetError()
				if tt.want {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestRefundStatus_String(t *testing.T) {
	tests := []struct {
		name string
		st   model.RefundStatus
		want string
	}{
		{
			name: "success status",
			st:   model.RefundStatusSuccess,
			want: "00",
		},
		{
			name: "failed status",
			st:   model.RefundStatusFailed,
			want: "01",
		},
		{
			name: "delayed status",
			st:   model.RefundStatusDelayed,
			want: "02",
		},
		{
			name: "uncertain status",
			st:   model.RefundStatusUncertain,
			want: "03",
		},
		{
			name: "unknown status",
			st:   model.RefundStatus("99"),
			want: "99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.st)
			assert.Equal(t, tt.want, got)
		})
	}
}
