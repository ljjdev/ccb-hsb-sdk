// Package model 的单元测试
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrderResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name string
		resp *CreateOrderResponse
		want bool
	}{
		{
			name: "success response",
			resp: &CreateOrderResponse{
				SvcRspSt: "00",
			},
			want: true,
		},
		{
			name: "failure response",
			resp: &CreateOrderResponse{
				SvcRspSt: "01",
			},
			want: false,
		},
		{
			name: "empty status",
			resp: &CreateOrderResponse{
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
}

func TestCreateOrderResponse_GetError(t *testing.T) {
	tests := []struct {
		name string
		resp *CreateOrderResponse
		want bool
	}{
		{
			name: "success response",
			resp: &CreateOrderResponse{
				SvcRspSt: "00",
			},
			want: false,
		},
		{
			name: "failure response with code and message",
			resp: &CreateOrderResponse{
				SvcRspSt: "01",
				SvcRspCd: "ERROR_CODE",
				RspInf:   "error message",
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
}

func TestCreateOrderRequest_ToMap(t *testing.T) {
	req := &CreateOrderRequest{
		IttpartyStmId:  "00000",
		PyChnlCd:       "0000000000000000000000000",
		IttpartyTms:    "20240101120000123",
		IttpartyJrnlNo: "20240101120000123001",
		MktId:          "12345678901234",
		MainOrdrNo:     "20240101120000123",
		PymdCd:         PaymentMethodPC,
		PyOrdrTpcd:     OrderTypeNormal,
		Ccy:            "156",
		OrdrTamt:       "100.00",
		TxnTamt:        "100.00",
	}

	result, err := req.ToMap()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "00000", result["Ittparty_Stm_Id"])
	assert.Equal(t, "0000000000000000000000000", result["Py_Chnl_Cd"])
	assert.Equal(t, "20240101120000123", result["Ittparty_Tms"])
	assert.Equal(t, "20240101120000123001", result["Ittparty_Jrnl_No"])
	assert.Equal(t, "12345678901234", result["Mkt_Id"])
	assert.Equal(t, "20240101120000123", result["Main_Ordr_No"])
	assert.Equal(t, string(PaymentMethodPC), result["Pymd_Cd"])
	assert.Equal(t, string(OrderTypeNormal), result["Py_Ordr_Tpcd"])
	assert.Equal(t, "156", result["Ccy"])
	assert.Equal(t, "100.00", result["Ordr_Tamt"])
	assert.Equal(t, "100.00", result["Txn_Tamt"])
}

func TestOrderStatus_String(t *testing.T) {
	tests := []struct {
		status OrderStatus
		want   string
	}{
		{OrderStatusPending, "1"},
		{OrderStatusSuccess, "2"},
		{OrderStatusFailed, "3"},
		{OrderStatusPolling, "9"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.status))
		})
	}
}

func TestPaymentMethod_String(t *testing.T) {
	tests := []struct {
		method PaymentMethod
		want   string
	}{
		{PaymentMethodPC, "01"},
		{PaymentMethodOffline, "02"},
		{PaymentMethodMobileH5, "03"},
		{PaymentMethodWechatMini, "05"},
		{PaymentMethodOnlineBank, "06"},
		{PaymentMethodQRCode, "07"},
		{PaymentMethodDragonPay, "08"},
		{PaymentMethodScan, "09"},
		{PaymentMethodDigitalWallet, "11"},
		{PaymentMethodContactless, "12"},
		{PaymentMethodSharedWallet, "13"},
		{PaymentMethodAlipayMini, "14"},
		{PaymentMethodSilentPay, "15"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.method))
		})
	}
}

func TestOrderType_String(t *testing.T) {
	tests := []struct {
		orderType OrderType
		want      string
	}{
		{OrderTypeCoupon, "02"},
		{OrderTypeTransit, "03"},
		{OrderTypeNormal, "04"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.orderType))
		})
	}
}

func TestQueryOrderRequest_ToMap(t *testing.T) {
	req := &QueryOrderRequest{
		IttpartyStmId:  "00000",
		PyChnlCd:       "0000000000000000000000000",
		IttpartyTms:    "20240101120000123",
		IttpartyJrnlNo: "20240101120000123001",
		MktId:          "12345678901234",
		MainOrdrNo:     "20240101120000123",
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
	assert.Equal(t, "20240101120000123", result["Main_Ordr_No"])
	assert.Equal(t, "4", result["Vno"])
}

func TestQueryOrderResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name string
		resp *QueryOrderResponse
		want bool
	}{
		{
			name: "success response",
			resp: &QueryOrderResponse{
				SvcRspSt: "00",
			},
			want: true,
		},
		{
			name: "failure response",
			resp: &QueryOrderResponse{
				SvcRspSt: "01",
			},
			want: false,
		},
		{
			name: "empty status",
			resp: &QueryOrderResponse{
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
}

func TestQueryOrderResponse_IsPaid(t *testing.T) {
	tests := []struct {
		name string
		resp *QueryOrderResponse
		want bool
	}{
		{
			name: "paid order",
			resp: &QueryOrderResponse{
				OrdrStcd: OrderStatusSuccess,
			},
			want: true,
		},
		{
			name: "pending order",
			resp: &QueryOrderResponse{
				OrdrStcd: OrderStatusPending,
			},
			want: false,
		},
		{
			name: "failed order",
			resp: &QueryOrderResponse{
				OrdrStcd: OrderStatusFailed,
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
}

func TestQueryOrderResponse_GetError(t *testing.T) {
	tests := []struct {
		name string
		resp *QueryOrderResponse
		want bool
	}{
		{
			name: "success response",
			resp: &QueryOrderResponse{
				SvcRspSt: "00",
			},
			want: false,
		},
		{
			name: "failure response with code and message",
			resp: &QueryOrderResponse{
				SvcRspSt: "01",
				SvcRspCd: "ERROR_CODE",
				RspInf:   "error message",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.resp.GetError()
			if tt.want {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "ERROR_CODE")
				assert.Contains(t, err.Error(), "error message")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
