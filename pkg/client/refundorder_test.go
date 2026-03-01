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

func TestRefundOrder(t *testing.T) {
	runWithTimeout(t, "TestRefundOrder", func(t *testing.T) {
		// 生成测试密钥对
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		tests := []struct {
			name           string
			request        *model.RefundOrderRequest
			response       *model.RefundOrderResponse
			responseStatus int
			wantErr        bool
			errContains    string
		}{
			{
				name: "successful full refund",
				request: &model.RefundOrderRequest{
					MainOrdrNo:   "20240101120000123",
					RefundOrdrNo: "REFUND20240101120000123",
					RefundAmt:    "100.00",
					RefundRsn:    "用户申请退款",
				},
				response: &model.RefundOrderResponse{
					MainOrdrNo:   "20240101120000123",
					RefundOrdrNo: "REFUND20240101120000123",
					RefundAmt:    "100.00",
					RefundTrnNo:  "RFND20240101120000123",
					RefundTm:     "20240101120000",
					RefundStcd:   "00",
					SvcRspSt:     "00",
				},
				responseStatus: http.StatusOK,
				wantErr:        false,
			},
			{
				name: "successful partial refund",
				request: &model.RefundOrderRequest{
					MainOrdrNo:   "20240101120000124",
					RefundOrdrNo: "REFUND20240101120000124",
					RefundAmt:    "50.00",
					RefundRsn:    "部分退款",
				},
				response: &model.RefundOrderResponse{
					MainOrdrNo:   "20240101120000124",
					RefundOrdrNo: "REFUND20240101120000124",
					RefundAmt:    "50.00",
					RefundTrnNo:  "RFND20240101120000124",
					RefundTm:     "20240101120000",
					RefundStcd:   "00",
					SvcRspSt:     "00",
				},
				responseStatus: http.StatusOK,
				wantErr:        false,
			},
			{
				name: "refund failed",
				request: &model.RefundOrderRequest{
					MainOrdrNo:   "20240101120000125",
					RefundOrdrNo: "REFUND20240101120000125",
					RefundAmt:    "100.00",
					RefundRsn:    "用户申请退款",
				},
				response: &model.RefundOrderResponse{
					MainOrdrNo:   "20240101120000125",
					RefundOrdrNo: "REFUND20240101120000125",
					RefundAmt:    "100.00",
					SvcRspSt:     "01",
					SvcRspCd:     "ERROR_CODE",
					RspInf:       "退款失败",
				},
				responseStatus: http.StatusOK,
				wantErr:        true,
				errContains:    "service error",
			},
			{
				name: "refund amount exceeds order amount",
				request: &model.RefundOrderRequest{
					MainOrdrNo:   "20240101120000126",
					RefundOrdrNo: "REFUND20240101120000126",
					RefundAmt:    "200.00",
					RefundRsn:    "退款金额超限",
				},
				response: &model.RefundOrderResponse{
					MainOrdrNo:   "20240101120000126",
					RefundOrdrNo: "REFUND20240101120000126",
					RefundAmt:    "200.00",
					SvcRspSt:     "01",
					SvcRspCd:     "AMOUNT_EXCEED",
					RspInf:       "退款金额超过订单金额",
				},
				responseStatus: http.StatusOK,
				wantErr:        true,
				errContains:    "service error",
			},
			{
				name: "http request failed",
				request: &model.RefundOrderRequest{
					MainOrdrNo:   "20240101120000127",
					RefundOrdrNo: "REFUND20240101120000127",
					RefundAmt:    "100.00",
					RefundRsn:    "用户申请退款",
				},
				response:       nil,
				responseStatus: http.StatusInternalServerError,
				wantErr:        true,
				errContains:    "failed to send request",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
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
					var req model.RefundOrderRequest
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					// 设置响应状态码
					if tt.responseStatus != http.StatusOK {
						w.WriteHeader(tt.responseStatus)
						return
					}

					// 为响应添加签名（使用相同的密钥对）
					signer := signature.NewRSAService(privateKey, &privateKey.PublicKey)

					// 将响应转换为 map 并签名
					params, err := tt.response.ToMap()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					sign, err := signer.SignParams(params)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					tt.response.SignInf = sign

					// 返回响应
					w.Header().Set(HeaderContentType, ContentTypeJSON)
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(tt.response)
				}

				client, server := setupTestClientWithKey(t, handler, privateKey)
				defer server.Close()

				// 调用 RefundOrder
				resp, err := client.RefundOrder(context.Background(), tt.request)

				if tt.wantErr {
					assert.Error(t, err)
					if tt.errContains != "" {
						assert.Contains(t, err.Error(), tt.errContains)
					}
					assert.Nil(t, resp)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, resp)
					assert.True(t, resp.IsSuccess())
				}
			})
		}
	})
}

func TestRefundOrder_WithAutoGeneratedFields(t *testing.T) {
	runWithTimeout(t, "TestRefundOrder_WithAutoGeneratedFields", func(t *testing.T) {
		// 生成测试密钥对
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		request := &model.RefundOrderRequest{
			MainOrdrNo:   "20240101120000123",
			RefundOrdrNo: "REFUND20240101120000123",
			RefundAmt:    "100.00",
			RefundRsn:    "用户申请退款",
			// 不提供 IttpartyJrnlNo 和 IttpartyTms，让客户端自动生成
		}

		response := &model.RefundOrderResponse{
			MainOrdrNo:   "20240101120000123",
			RefundOrdrNo: "REFUND20240101120000123",
			RefundAmt:    "100.00",
			RefundTrnNo:  "RFND20240101120000123",
			RefundTm:     "20240101120000",
			RefundStcd:   "00",
			SvcRspSt:     "00",
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
			var req model.RefundOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// 验证自动生成的字段
			if req.IttpartyJrnlNo == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.IttpartyTms == "" {
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

		// 调用 RefundOrder
		resp, err := client.RefundOrder(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.IsSuccess())
	})
}

func TestRefundOrderResponse_IsSuccess(t *testing.T) {
	runWithTimeout(t, "TestRefundOrderResponse_IsSuccess", func(t *testing.T) {
		tests := []struct {
			name string
			resp *model.RefundOrderResponse
			want bool
		}{
			{
				name: "success response",
				resp: &model.RefundOrderResponse{
					SvcRspSt: "00",
				},
				want: true,
			},
			{
				name: "failure response",
				resp: &model.RefundOrderResponse{
					SvcRspSt: "01",
				},
				want: false,
			},
			{
				name: "empty status",
				resp: &model.RefundOrderResponse{
					SvcRspSt: "",
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
