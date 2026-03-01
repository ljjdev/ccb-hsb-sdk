// Package signature 的单元测试
package signature

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	maxTimeoutAttempts = 5
	testTimeout        = 5 * time.Second
)

func runWithTimeout(t *testing.T, testName string, testFunc func(t *testing.T)) {
	timeoutCount := new(int32)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("测试 %s 发生 panic: %v", testName, r)
		}
	}()

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
		return
	case <-ctx.Done():
		atomic.AddInt32(timeoutCount, 1)
		if atomic.LoadInt32(timeoutCount) >= maxTimeoutAttempts {
			t.Skipf("测试超时超过 %d 次，跳过该测试", maxTimeoutAttempts)
		}
		t.Errorf("测试 %s 超时", testName)
	}
}

func TestBuildSignatureStringFromJSON(t *testing.T) {
	runWithTimeout(t, "TestBuildSignatureStringFromJSON", func(t *testing.T) {
		tests := []struct {
			name     string
			json     string
			expected string
			wantErr  bool
		}{
			{
				name:     "simple object",
				json:     `{"key1":"value1","key2":"value2"}`,
				expected: "key1=value1&key2=value2",
				wantErr:  false,
			},
			{
				name:     "object with empty value",
				json:     `{"key1":"value1","key2":"","key3":"value3"}`,
				expected: "key1=value1&key3=value3",
				wantErr:  false,
			},
			{
				name:     "object with SIGN_INF field",
				json:     `{"key1":"value1","SIGN_INF":"signature","key2":"value2"}`,
				expected: "key1=value1&key2=value2",
				wantErr:  false,
			},
			{
				name:     "object with Sign_Inf field",
				json:     `{"key1":"value1","Sign_Inf":"signature","key2":"value2"}`,
				expected: "key1=value1&key2=value2",
				wantErr:  false,
			},
			{
				name:     "object with Svc_Rsp_St field",
				json:     `{"key1":"value1","Svc_Rsp_St":"00","key2":"value2"}`,
				expected: "key1=value1&key2=value2",
				wantErr:  false,
			},
			{
				name:     "object with Svc_Rsp_Cd field",
				json:     `{"key1":"value1","Svc_Rsp_Cd":"error","key2":"value2"}`,
				expected: "key1=value1&key2=value2",
				wantErr:  false,
			},
			{
				name:     "object with Rsp_Inf field",
				json:     `{"key1":"value1","Rsp_Inf":"error message","key2":"value2"}`,
				expected: "key1=value1&key2=value2",
				wantErr:  false,
			},
			{
				name:     "nested object",
				json:     `{"key1":"value1","nested":{"key2":"value2","key3":"value3"}}`,
				expected: "key1=value1&key2=value2&key3=value3",
				wantErr:  false,
			},
			{
				name:     "array of objects",
				json:     `{"key1":"value1","array":[{"key2":"value2"},{"key3":"value3"}]}`,
				expected: "key2=value2&key3=value3&key1=value1",
				wantErr:  false,
			},
			{
				name:     "array of strings",
				json:     `{"array":["value1","value2","value3"]}`,
				expected: "",
				wantErr:  false,
			},
			{
				name:     "complex nested structure",
				json:     `{"key1":"value1","nested":{"array":[{"key2":"value2"},{"key3":"value3"}],"key4":"value4"}}`,
				expected: "key1=value1&key2=value2&key3=value3&key4=value4",
				wantErr:  false,
			},
			{
				name:     "empty object",
				json:     `{}`,
				expected: "",
				wantErr:  false,
			},
			{
				name:     "invalid json",
				json:     `{invalid}`,
				expected: "",
				wantErr:  true,
			},
			{
				name:     "order request example",
				json:     `{"Ittparty_Stm_Id":"00000","Py_Chnl_Cd":"0000000000000000000000000","Ittparty_Tms":"20240101120000123","Ittparty_Jrnl_No":"20240101120000123001","Mkt_Id":"12345678901234","Main_Ordr_No":"20240101120000123","Pymd_Cd":"03","Py_Ordr_Tpcd":"04","Ccy":"156","Ordr_Tamt":"100.01","Txn_Tamt":"100.01","Pay_Dsc":"商品","Order_Time_Out":"1800"}`,
				expected: "Ccy=156&Ittparty_Jrnl_No=20240101120000123001&Ittparty_Stm_Id=00000&Ittparty_Tms=20240101120000123&Main_Ordr_No=20240101120000123&Mkt_Id=12345678901234&Order_Time_Out=1800&Ordr_Tamt=100.01&Pay_Dsc=商品&Py_Chnl_Cd=0000000000000000000000000&Py_Ordr_Tpcd=04&Pymd_Cd=03&Txn_Tamt=100.01",
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := BuildSignatureStringFromJSON(tt.json)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, result)
				}
			})
		}
	})
}

func TestSplicingObject(t *testing.T) {
	runWithTimeout(t, "TestSplicingObject", func(t *testing.T) {
		tests := []struct {
			name     string
			obj      map[string]interface{}
			expected string
		}{
			{
				name:     "simple object",
				obj:      map[string]interface{}{"key1": "value1", "key2": "value2"},
				expected: "key1=value1&key2=value2&",
			},
			{
				name:     "object with empty value",
				obj:      map[string]interface{}{"key1": "value1", "key2": "", "key3": "value3"},
				expected: "key1=value1&key3=value3&",
			},
			{
				name:     "object with SIGN_INF field",
				obj:      map[string]interface{}{"key1": "value1", "SIGN_INF": "signature", "key2": "value2"},
				expected: "key1=value1&key2=value2&",
			},
			{
				name:     "object with nested object",
				obj:      map[string]interface{}{"key1": "value1", "nested": map[string]interface{}{"key2": "value2", "key3": "value3"}},
				expected: "key1=value1&key2=value2&key3=value3&",
			},
			{
				name:     "object with array",
				obj:      map[string]interface{}{"key1": "value1", "array": []interface{}{"value1", "value2"}},
				expected: "key1=value1&",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := splicingObject(tt.obj)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestSplicingArray(t *testing.T) {
	runWithTimeout(t, "TestSplicingArray", func(t *testing.T) {
		tests := []struct {
			name     string
			arr      []interface{}
			expected string
		}{
			{
				name:     "array of strings",
				arr:      []interface{}{"value1", "value2", "value3"},
				expected: "",
			},
			{
				name:     "array of objects",
				arr:      []interface{}{map[string]interface{}{"key1": "value1"}, map[string]interface{}{"key2": "value2"}},
				expected: "key1=value1&key2=value2&",
			},
			{
				name:     "empty array",
				arr:      []interface{}{},
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := splicingArray(tt.arr)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestSplicingSign(t *testing.T) {
	runWithTimeout(t, "TestSplicingSign", func(t *testing.T) {
		tests := []struct {
			name     string
			data     interface{}
			expected string
		}{
			{
				name:     "object",
				data:     map[string]interface{}{"key1": "value1", "key2": "value2"},
				expected: "key1=value1&key2=value2&",
			},
			{
				name:     "array of objects",
				data:     []interface{}{map[string]interface{}{"key1": "value1"}, map[string]interface{}{"key2": "value2"}},
				expected: "key1=value1&key2=value2&",
			},
			{
				name:     "string",
				data:     "value",
				expected: "",
			},
			{
				name:     "number",
				data:     123,
				expected: "",
			},
			{
				name:     "nil",
				data:     nil,
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := splicingSign(tt.data)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

// BenchmarkBuildSignatureStringFromJSON 性能测试
func BenchmarkBuildSignatureStringFromJSON(b *testing.B) {
	jsonStr := `{"Ittparty_Stm_Id":"00000","Py_Chnl_Cd":"0000000000000000000000000","Ittparty_Tms":"20240101120000123","Ittparty_Jrnl_No":"20240101120000123001","Mkt_Id":"12345678901234","Main_Ordr_No":"20240101120000123","Pymd_Cd":"03","Py_Ordr_Tpcd":"04","Ccy":"156","Ordr_Tamt":"100.01","Txn_Tamt":"100.01","Pay_Dsc":"商品","Order_Time_Out":"1800","Orderlist":[{"Mkt_Mrch_Id":"12345678901234567890","Cmdty_Ordr_No":"20240101120000123001","Ordr_Amt":"100.01","Txnamt":"100.01","Cmdty_Dsc":"商品","Clrg_Rule_Id":"123456","Parlist":[{"Seq_No":1,"Mkt_Mrch_Id":"12345678901234567890"}]}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := BuildSignatureStringFromJSON(jsonStr)
		require.NoError(b, err)
	}
}
