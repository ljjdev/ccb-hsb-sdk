// Package signature 的单元测试
package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSAService_Sign(t *testing.T) {
	runWithTimeout(t, "TestRSAService_Sign", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		service := NewRSAService(privateKey, nil)

		tests := []struct {
			name    string
			data    string
			wantErr bool
		}{
			{
				name:    "valid data",
				data:    "test data",
				wantErr: false,
			},
			{
				name:    "empty data",
				data:    "",
				wantErr: false,
			},
			{
				name:    "long data",
				data:    string(make([]byte, 10000)),
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				signature, err := service.Sign(tt.data)
				if tt.wantErr {
					assert.Error(t, err)
					assert.Empty(t, signature)
				} else {
					assert.NoError(t, err)
					assert.NotEmpty(t, signature)
				}
			})
		}
	})
}

func TestRSAService_Verify(t *testing.T) {
	runWithTimeout(t, "TestRSAService_Verify", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		publicKey := &privateKey.PublicKey
		service := NewRSAService(privateKey, publicKey)

		data := "test data"
		signature, err := service.Sign(data)
		require.NoError(t, err)

		tests := []struct {
			name      string
			data      string
			signature string
			wantErr   bool
		}{
			{
				name:      "valid signature",
				data:      data,
				signature: signature,
				wantErr:   false,
			},
			{
				name:      "invalid data",
				data:      "invalid data",
				signature: signature,
				wantErr:   true,
			},
			{
				name:      "invalid signature",
				data:      data,
				signature: "invalid signature",
				wantErr:   true,
			},
			{
				name:      "empty signature",
				data:      data,
				signature: "",
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := service.Verify(tt.data, tt.signature)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestRSAService_SignWithoutPrivateKey(t *testing.T) {
	runWithTimeout(t, "TestRSAService_SignWithoutPrivateKey", func(t *testing.T) {
		service := NewRSAService(nil, &rsa.PublicKey{})

		_, err := service.Sign("test data")
		assert.Error(t, err)
	})
}

func TestRSAService_VerifyWithoutPublicKey(t *testing.T) {
	runWithTimeout(t, "TestRSAService_VerifyWithoutPublicKey", func(t *testing.T) {
		service := NewRSAService(&rsa.PrivateKey{}, nil)

		err := service.Verify("test data", "signature")
		assert.Error(t, err)
	})
}

func TestBuildSignatureString(t *testing.T) {
	runWithTimeout(t, "TestBuildSignatureString", func(t *testing.T) {
		tests := []struct {
			name     string
			params   map[string]string
			expected string
		}{
			{
				name:     "empty params",
				params:   map[string]string{},
				expected: "",
			},
			{
				name:     "single param",
				params:   map[string]string{"key1": "value1"},
				expected: "key1=value1",
			},
			{
				name:     "multiple params",
				params:   map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
				expected: "key1=value1&key2=value2&key3=value3",
			},
			{
				name:     "params with empty value",
				params:   map[string]string{"key1": "value1", "key2": "", "key3": "value3"},
				expected: "key1=value1&key3=value3",
			},
			{
				name:     "params with Sign_Inf field",
				params:   map[string]string{"key1": "value1", "Sign_Inf": "signature", "key2": "value2"},
				expected: "key1=value1&key2=value2",
			},
			{
				name:     "params with special characters",
				params:   map[string]string{"key1": "value1", "key2": "value with spaces", "key3": "value3"},
				expected: "key1=value1&key2=value with spaces&key3=value3",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := BuildSignatureString(tt.params)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestRSAService_SignParams(t *testing.T) {
	runWithTimeout(t, "TestRSAService_SignParams", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		service := NewRSAService(privateKey, nil)

		params := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		signature, err := service.SignParams(params)
		assert.NoError(t, err)
		assert.NotEmpty(t, signature)
	})
}

func TestRSAService_VerifyParams(t *testing.T) {
	runWithTimeout(t, "TestRSAService_VerifyParams", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		publicKey := &privateKey.PublicKey
		service := NewRSAService(privateKey, publicKey)

		params := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		signature, err := service.SignParams(params)
		require.NoError(t, err)

		err = service.VerifyParams(params, signature)
		assert.NoError(t, err)

		params["key1"] = "invalid"
		err = service.VerifyParams(params, signature)
		assert.Error(t, err)
	})
}
