// Package client 的单元测试
package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ljjdev/ccb-hsb-sdk/pkg/model"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/signature"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryOrder_WithSubOrders(t *testing.T) {
	runWithTimeout(t, "TestQueryOrder_WithSubOrders", func(t *testing.T) {
		// 生成测试密钥对
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		request := &model.QueryOrderRequest{
			IttpartyJrnlNo: "20240101120000123001",
			MainOrdrNo:     "20240101120000123",
		}

		response := &model.QueryOrderResponse{
			MainOrdrNo: "20240101120000123",
			PyTrnNo:    "PAY20240101120000123001",
			Txnamt:     100.01,
			OrdrGenTm:  "20240101120000",
			OrdrStcd:   model.OrderStatusSuccess,
			SvcRspSt:   "00",
			Orderlist: []model.SubOrderResponse{
				{
					CmdtyOrdrNo: "20240101120000123001",
					SubOrdrId:   "SUB20240101120000123",
				},
			},
		}

		handler := func(w http.ResponseWriter, r *http.Request) {
			// 检查请求方法
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// 检查 Content-Type
			if r.Header.Get(HeaderContentType) != ContentTypeJSON {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// 检查请求体
			var req model.QueryOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// 为响应添加签名（使用相同的密钥对）
			signer := signature.NewRSAService(privateKey, &privateKey.PublicKey)

			// 将响应转换为 map 并签名
			params, err := response.ToMap()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sign, err := signer.SignParams(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.SignInf = sign

			// 返回响应
			w.Header().Set(HeaderContentType, ContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}

		client, server := setupTestClientWithKey(t, handler, privateKey)
		defer server.Close()

		// 调用 QueryOrder
		resp, err := client.QueryOrder(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.IsPaid())
		assert.Equal(t, 1, len(resp.Orderlist))
	})
}

func TestQueryOrderResponse_IsPaid(t *testing.T) {
	runWithTimeout(t, "TestQueryOrderResponse_IsPaid", func(t *testing.T) {
		tests := []struct {
			name string
			resp *model.QueryOrderResponse
			want bool
		}{
			{
				name: "paid order",
				resp: &model.QueryOrderResponse{
					OrdrStcd: model.OrderStatusSuccess,
				},
				want: true,
			},
			{
				name: "pending order",
				resp: &model.QueryOrderResponse{
					OrdrStcd: model.OrderStatusPending,
				},
				want: false,
			},
			{
				name: "failed order",
				resp: &model.QueryOrderResponse{
					OrdrStcd: model.OrderStatusFailed,
				},
				want: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.resp.IsPaid()
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
