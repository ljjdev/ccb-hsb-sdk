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

func TestPlaceOrder(t *testing.T) {
	runWithTimeout(t, "TestPlaceOrder", func(t *testing.T) {
		tests := []struct {
			name           string
			request        *model.CreateOrderRequest
			response       *model.CreateOrderResponse
			responseStatus int
			wantErr        bool
			errContains    string
		}{
			{
				name: "successful order creation",
				request: &model.CreateOrderRequest{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123001",
					MainOrdrNo:     "20240101120000123",
					PymdCd:         model.PaymentMethodMobileH5,
					PyOrdrTpcd:     model.OrderTypeNormal,
					OrdrTamt:       "100.01",
					TxnTamt:        "100.01",
					PayDsc:         "商品",
					OrderTimeOut:   "1800",
				},
				response: &model.CreateOrderResponse{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123001",
					MainOrdrNo:     "20240101120000123",
					PyTrnNo:        "PAY20240101120000123001",
					PrimOrdrNo:     "PRIM20240101120000123",
					OrdrGenTm:      "20240101120000",
					OrdrOvtmTm:     "20240101123000",
					CshdkUrl:       "https://pay.ccb.com/cashier?param=test",
					OrdrStcd:       model.OrderStatusPending,
					SvcRspSt:       "00",
				},
				responseStatus: http.StatusOK,
				wantErr:        false,
			},
			{
				name: "successful order with URL encoded",
				request: &model.CreateOrderRequest{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123002",
					MainOrdrNo:     "20240101120000124",
					PymdCd:         model.PaymentMethodMobileH5,
					PyOrdrTpcd:     model.OrderTypeNormal,
					OrdrTamt:       "100.01",
					TxnTamt:        "100.01",
					PayDsc:         "商品",
					OrderTimeOut:   "1800",
				},
				response: &model.CreateOrderResponse{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123002",
					MainOrdrNo:     "20240101120000124",
					PyTrnNo:        "PAY20240101120000123002",
					PrimOrdrNo:     "PRIM20240101120000124",
					OrdrGenTm:      "20240101120000",
					OrdrOvtmTm:     "20240101123000",
					CshdkUrl:       "https%3A%2F%2Fpay.ccb.com%2Fcashier%3Fparam%3Dtest",
					OrdrStcd:       model.OrderStatusPending,
					SvcRspSt:       "00",
				},
				responseStatus: http.StatusOK,
				wantErr:        false,
			},
			{
				name: "order creation failed",
				request: &model.CreateOrderRequest{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123003",
					MainOrdrNo:     "20240101120000125",
					PymdCd:         model.PaymentMethodMobileH5,
					PyOrdrTpcd:     model.OrderTypeNormal,
					OrdrTamt:       "100.01",
					TxnTamt:        "100.01",
					PayDsc:         "商品",
					OrderTimeOut:   "1800",
				},
				response: &model.CreateOrderResponse{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123003",
					MainOrdrNo:     "20240101120000125",
					OrdrStcd:       model.OrderStatusFailed,
					SvcRspSt:       "01",
					SvcRspCd:       "ERROR_CODE",
					RspInf:         "订单创建失败",
				},
				responseStatus: http.StatusOK,
				wantErr:        true,
				errContains:    "order creation failed",
			},
			{
				name: "invalid payment URL",
				request: &model.CreateOrderRequest{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123004",
					MainOrdrNo:     "20240101120000126",
					PymdCd:         model.PaymentMethodMobileH5,
					PyOrdrTpcd:     model.OrderTypeNormal,
					OrdrTamt:       "100.01",
					TxnTamt:        "100.01",
					PayDsc:         "商品",
					OrderTimeOut:   "1800",
				},
				response: &model.CreateOrderResponse{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123004",
					MainOrdrNo:     "20240101120000126",
					PyTrnNo:        "PAY20240101120000123004",
					PrimOrdrNo:     "PRIM20240101120000126",
					OrdrGenTm:      "20240101120000",
					OrdrOvtmTm:     "20240101123000",
					CshdkUrl:       "http://invalid.url",
					OrdrStcd:       model.OrderStatusPending,
					SvcRspSt:       "00",
				},
				responseStatus: http.StatusOK,
				wantErr:        true,
				errContains:    "invalid payment URL",
			},
			{
				name: "empty payment URL",
				request: &model.CreateOrderRequest{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123005",
					MainOrdrNo:     "20240101120000127",
					PymdCd:         model.PaymentMethodMobileH5,
					PyOrdrTpcd:     model.OrderTypeNormal,
					OrdrTamt:       "100.01",
					TxnTamt:        "100.01",
					PayDsc:         "商品",
					OrderTimeOut:   "1800",
				},
				response: &model.CreateOrderResponse{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123005",
					MainOrdrNo:     "20240101120000127",
					PyTrnNo:        "PAY20240101120000123005",
					PrimOrdrNo:     "PRIM20240101120000127",
					OrdrGenTm:      "20240101120000",
					OrdrOvtmTm:     "20240101123000",
					CshdkUrl:       "",
					OrdrStcd:       model.OrderStatusPending,
					SvcRspSt:       "00",
				},
				responseStatus: http.StatusOK,
				wantErr:        true,
				errContains:    "invalid payment URL",
			},
			{
				name: "http request failed",
				request: &model.CreateOrderRequest{
					IttpartyTms:    "20240101120000123",
					IttpartyJrnlNo: "20240101120000123006",
					MainOrdrNo:     "20240101120000128",
					PymdCd:         model.PaymentMethodMobileH5,
					PyOrdrTpcd:     model.OrderTypeNormal,
					OrdrTamt:       "100.01",
					TxnTamt:        "100.01",
					PayDsc:         "商品",
					OrderTimeOut:   "1800",
				},
				response:       nil,
				responseStatus: http.StatusInternalServerError,
				wantErr:        true,
				errContains:    "failed to send request",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// 生成测试密钥对
				privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
				require.NoError(t, err)

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
					var req model.CreateOrderRequest
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

				// 调用 PlaceOrder
				payURL, err := client.PlaceOrder(context.Background(), tt.request)

				if tt.wantErr {
					assert.Error(t, err)
					if tt.errContains != "" {
						assert.Contains(t, err.Error(), tt.errContains)
					}
					assert.Empty(t, payURL)
				} else {
					assert.NoError(t, err)
					assert.NotEmpty(t, payURL)
					assert.Contains(t, payURL, "https://")
				}
			})
		}
	})
}

func TestPlaceOrder_WithSubOrders(t *testing.T) {
	runWithTimeout(t, "TestPlaceOrder_WithSubOrders", func(t *testing.T) {
		// 生成测试密钥对
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		request := &model.CreateOrderRequest{
			IttpartyTms:    "20240101120000123",
			IttpartyJrnlNo: "20240101120000123001",
			MainOrdrNo:     "20240101120000123",
			PymdCd:         model.PaymentMethodMobileH5,
			PyOrdrTpcd:     model.OrderTypeNormal,
			OrdrTamt:       "100.01",
			TxnTamt:        "100.01",
			PayDsc:         "商品",
			OrderTimeOut:   "1800",
			Orderlist: []model.SubOrder{
				{
					CmdtyOrdrNo: "20240101120000123001",
					OrdrAmt:     "100.01",
					Txnamt:      "100.01",
					CmdtyDsc:    "商品",
					ClrgRuleId:  "123456",
					Parlist: []model.Participant{
						{
							SeqNo:     1,
							MktMrchId: "12345678901234567890",
						},
					},
				},
			},
		}

		response := &model.CreateOrderResponse{
			IttpartyTms:    "20240101120000123",
			IttpartyJrnlNo: "20240101120000123001",
			MainOrdrNo:     "20240101120000123",
			PyTrnNo:        "PAY20240101120000123001",
			PrimOrdrNo:     "PRIM20240101120000123",
			OrdrGenTm:      "20240101120000",
			OrdrOvtmTm:     "20240101123000",
			CshdkUrl:       "https://pay.ccb.com/cashier?param=test",
			OrdrStcd:       model.OrderStatusPending,
			SvcRspSt:       "00",
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
			var req model.CreateOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// 验证子订单
			if len(req.Orderlist) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.Orderlist[0].CmdtyOrdrNo != "20240101120000123001" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(req.Orderlist[0].Parlist) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.Orderlist[0].Parlist[0].SeqNo != 1 {
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

		// 调用 PlaceOrder
		payURL, err := client.PlaceOrder(context.Background(), request)
		assert.NoError(t, err)
		assert.NotEmpty(t, payURL)
		assert.Contains(t, payURL, "https://")
	})
}

func TestPlaceOrder_WithCoupons(t *testing.T) {
	runWithTimeout(t, "TestPlaceOrder_WithCoupons", func(t *testing.T) {
		// 生成测试密钥对
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		request := &model.CreateOrderRequest{
			IttpartyTms:    "20240101120000123",
			IttpartyJrnlNo: "20240101120000123001",
			MainOrdrNo:     "20240101120000123",
			PymdCd:         model.PaymentMethodMobileH5,
			PyOrdrTpcd:     model.OrderTypeNormal,
			OrdrTamt:       "100.01",
			TxnTamt:        "100.01",
			PayDsc:         "商品",
			OrderTimeOut:   "1800",
			Orderlist: []model.SubOrder{
				{
					CmdtyOrdrNo: "20240101120000123001",
					OrdrAmt:     "100.01",
					Txnamt:      "100.01",
					CmdtyDsc:    "商品",
					ClrgRuleId:  "123456",
					Cpnlist: []model.Coupon{
						{
							CnsmpNoteOrdrId: "COUPON20240101120000123",
						},
					},
					Parlist: []model.Participant{
						{
							SeqNo:     1,
							MktMrchId: "12345678901234567890",
						},
					},
				},
			},
		}

		response := &model.CreateOrderResponse{
			IttpartyTms:    "20240101120000123",
			IttpartyJrnlNo: "20240101120000123001",
			MainOrdrNo:     "20240101120000123",
			PyTrnNo:        "PAY20240101120000123001",
			PrimOrdrNo:     "PRIM20240101120000123",
			OrdrGenTm:      "20240101120000",
			OrdrOvtmTm:     "20240101123000",
			CshdkUrl:       "https://pay.ccb.com/cashier?param=test",
			OrdrStcd:       model.OrderStatusPending,
			SvcRspSt:       "00",
			Orderlist: []model.SubOrderResponse{
				{
					CmdtyOrdrNo: "20240101120000123001",
					SubOrdrId:   "SUB20240101120000123",
					Cpnlist: []model.UsedCoupon{
						{
							CnsmpNoteOrdrId: "COUPON20240101120000123",
							Amt:             "10.00",
							BalAmt:          "0.00",
						},
					},
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
			var req model.CreateOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// 验证消费券
			if len(req.Orderlist) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(req.Orderlist[0].Cpnlist) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.Orderlist[0].Cpnlist[0].CnsmpNoteOrdrId != "COUPON20240101120000123" {
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

		// 调用 PlaceOrder
		payURL, err := client.PlaceOrder(context.Background(), request)
		assert.NoError(t, err)
		assert.NotEmpty(t, payURL)
		assert.Contains(t, payURL, "https://")
	})
}
