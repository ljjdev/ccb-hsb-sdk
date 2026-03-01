// Package client 的单元测试 - QueryRefund
package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ljjdev/ccb-hsb-sdk/pkg/config"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// strPtr 创建字符串指针的辅助函数
func strPtr(s string) *string {
	return &s
}

// setupTestClientForQueryRefund 创建测试客户端
func setupTestClientForQueryRefund(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	// 生成测试密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// 创建测试服务器
	server := httptest.NewServer(handler)

	// 创建配置
	cfg := &config.Config{
		MarketID:   "12345678901234",
		MerchantID: "12345678901234567890",
		GatewayURL: server.URL,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Timeout:    30 * time.Second,
		Debug:      false,
	}

	// 创建客户端
	client, err := NewClient(cfg)
	require.NoError(t, err)

	return client, server
}

func TestQueryRefund_Success(t *testing.T) {
	// 创建mock服务器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		assert.Equal(t, http.MethodPost, r.Method)

		// 验证Content-Type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// 解析请求体
		var req model.QueryRefundRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// 验证必填字段
		assert.Equal(t, "00000", req.IttpartyStmId)
		assert.Equal(t, "0000000000000000000000000", req.PyChnlCd)
		assert.Equal(t, "4", req.Vno)
		assert.NotEmpty(t, req.SignInf)

		// 返回成功响应 - 不包含SignInf以跳过签名验证
		resp := model.QueryRefundResponse{
			IttpartyTms:    req.IttpartyTms,
			IttpartyJrnlNo: req.IttpartyJrnlNo,
			CustRfndTrcno:  req.CustRfndTrcno,
			RfndTrcno:      "RFND20240101120000123",
			RfndAmt:        strPtr("100.00"),
			RefundRspSt:    model.RefundStatusSuccess,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	client, server := setupTestClientForQueryRefund(t, handler)
	defer server.Close()

	// 创建请求
	req := &model.QueryRefundRequest{
		IttpartyStmId:  "00000",
		PyChnlCd:       "0000000000000000000000000",
		IttpartyTms:    "20240101120000123",
		IttpartyJrnlNo: "20240101120000123001",
		MktId:          "12345678901234",
		CustRfndTrcno:  "REFUND20240101120000123",
		Vno:            "4",
	}

	// 调用QueryRefund
	resp, err := client.QueryRefund(context.Background(), req)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "RFND20240101120000123", resp.RfndTrcno)
	assert.Equal(t, "100.00", resp.RfndAmt)
	assert.Equal(t, model.RefundStatusSuccess, resp.RefundRspSt)
	assert.True(t, resp.IsSuccess())
}

func TestQueryRefund_WithRefundTraceNo(t *testing.T) {
	// 创建mock服务器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析请求体
		var req model.QueryRefundRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// 验证使用退款流水号查询
		assert.Empty(t, req.CustRfndTrcno)
		assert.NotEmpty(t, req.RfndTrcno)

		// 返回成功响应 - 不包含SignInf以跳过签名验证
		resp := model.QueryRefundResponse{
			IttpartyTms:    req.IttpartyTms,
			IttpartyJrnlNo: req.IttpartyJrnlNo,
			RfndTrcno:      req.RfndTrcno,
			RfndAmt:        strPtr("50.00"),
			RefundRspSt:    model.RefundStatusSuccess,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	client, server := setupTestClientForQueryRefund(t, handler)
	defer server.Close()

	// 创建请求 - 使用退款流水号
	req := &model.QueryRefundRequest{
		IttpartyStmId:  "00000",
		PyChnlCd:       "0000000000000000000000000",
		IttpartyTms:    "20240101120000123",
		IttpartyJrnlNo: "20240101120000123001",
		MktId:          "12345678901234",
		RfndTrcno:      "RFND20240101120000123",
		Vno:            "4",
	}

	// 调用QueryRefund
	resp, err := client.QueryRefund(context.Background(), req)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "RFND20240101120000123", resp.RfndTrcno)
	assert.Equal(t, "50.00", resp.RfndAmt)
	assert.True(t, resp.IsSuccess())
}

func TestQueryRefund_MissingRequiredParams(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cfg := &config.Config{
		MarketID:   "12345678901234",
		MerchantID: "12345678901234567890",
		GatewayURL: "https://marketpay.ccb.com/online/direct",
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Timeout:    30 * time.Second,
		Debug:      false,
	}

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// 创建请求 - 缺少必填参数
	req := &model.QueryRefundRequest{
		IttpartyStmId:  "00000",
		PyChnlCd:       "0000000000000000000000000",
		IttpartyTms:    "20240101120000123",
		IttpartyJrnlNo: "20240101120000123001",
		MktId:          "12345678901234",
		Vno:            "4",
		// 缺少 CustRfndTrcno 和 RfndTrcno
	}

	// 调用QueryRefund
	_, err = client.QueryRefund(context.Background(), req)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cust_Rfnd_Trcno and Rfnd_Trcno must provide at least one")
}

func TestQueryRefund_WithDefaultValues(t *testing.T) {
	// 创建mock服务器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析请求体
		var req model.QueryRefundRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// 验证默认值已设置
		assert.Equal(t, "00000", req.IttpartyStmId)
		assert.Equal(t, "0000000000000000000000000", req.PyChnlCd)
		assert.Equal(t, "4", req.Vno)
		assert.Equal(t, "12345678901234", req.MktId)

		// 返回成功响应 - 不包含SignInf以跳过签名验证
		resp := model.QueryRefundResponse{
			IttpartyTms:    req.IttpartyTms,
			IttpartyJrnlNo: req.IttpartyJrnlNo,
			CustRfndTrcno:  req.CustRfndTrcno,
			RfndTrcno:      "RFND20240101120000123",
			RefundRspSt:    model.RefundStatusSuccess,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	client, server := setupTestClientForQueryRefund(t, handler)
	defer server.Close()

	// 创建请求 - 不设置默认值
	req := &model.QueryRefundRequest{
		IttpartyJrnlNo: "20240101120000123001",
		CustRfndTrcno:  "REFUND20240101120000123",
	}

	// 调用QueryRefund
	resp, err := client.QueryRefund(context.Background(), req)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsSuccess())
}

func TestQueryRefund_RefundFailed(t *testing.T) {
	// 创建mock服务器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回失败响应 - 不包含SignInf以跳过签名验证
		resp := model.QueryRefundResponse{
			IttpartyTms:    "20240101120000123",
			IttpartyJrnlNo: "20240101120000123001",
			RfndTrcno:      "RFND20240101120000123",
			RefundRspSt:    model.RefundStatusFailed,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	client, server := setupTestClientForQueryRefund(t, handler)
	defer server.Close()

	// 创建请求
	req := &model.QueryRefundRequest{
		IttpartyJrnlNo: "20240101120000123001",
		CustRfndTrcno:  "REFUND20240101120000123",
	}

	// 调用QueryRefund
	resp, err := client.QueryRefund(context.Background(), req)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsSuccess())
	assert.Error(t, resp.GetError())
}

func TestQueryRefund_HTTPError(t *testing.T) {
	// 创建mock服务器 - 返回500错误
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	client, server := setupTestClientForQueryRefund(t, handler)
	defer server.Close()

	// 创建请求
	req := &model.QueryRefundRequest{
		IttpartyJrnlNo: "20240101120000123001",
		CustRfndTrcno:  "REFUND20240101120000123",
	}

	// 调用QueryRefund
	_, err := client.QueryRefund(context.Background(), req)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send request")
}
