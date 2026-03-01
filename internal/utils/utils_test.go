// Package utils 的单元测试
package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentTimestamp(t *testing.T) {
	timestamp := CurrentTimestamp()
	assert.NotEmpty(t, timestamp)
	assert.Len(t, timestamp, 17) // yyyyMMddHHmmssfff
}

func TestCurrentTimestampShort(t *testing.T) {
	timestamp := CurrentTimestampShort()
	assert.NotEmpty(t, timestamp)
	assert.Len(t, timestamp, 8) // yyyyMMdd
}

func TestCurrentTimestampLong(t *testing.T) {
	timestamp := CurrentTimestampLong()
	assert.NotEmpty(t, timestamp)
	assert.Len(t, timestamp, 14) // yyyyMMddHHmmss
}

func TestParseTimestamp(t *testing.T) {
	timestamp := "20240101120000123"
	parsed, err := ParseTimestamp(timestamp)
	assert.NoError(t, err)
	assert.Equal(t, 2024, parsed.Year())
	assert.Equal(t, time.January, parsed.Month())
	assert.Equal(t, 1, parsed.Day())
	assert.Equal(t, 12, parsed.Hour())
	assert.Equal(t, 0, parsed.Minute())
	assert.Equal(t, 0, parsed.Second())
}

func TestParseTimestampShort(t *testing.T) {
	timestamp := "20240101"
	parsed, err := ParseTimestampShort(timestamp)
	assert.NoError(t, err)
	assert.Equal(t, 2024, parsed.Year())
	assert.Equal(t, time.January, parsed.Month())
	assert.Equal(t, 1, parsed.Day())
}

func TestParseTimestampLong(t *testing.T) {
	timestamp := "20240101120000"
	parsed, err := ParseTimestampLong(timestamp)
	assert.NoError(t, err)
	assert.Equal(t, 2024, parsed.Year())
	assert.Equal(t, time.January, parsed.Month())
	assert.Equal(t, 1, parsed.Day())
	assert.Equal(t, 12, parsed.Hour())
	assert.Equal(t, 0, parsed.Minute())
	assert.Equal(t, 0, parsed.Second())
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   string
	}{
		{
			name:   "integer",
			amount: 100,
			want:   "100.00",
		},
		{
			name:   "decimal",
			amount: 100.5,
			want:   "100.50",
		},
		{
			name:   "small decimal",
			amount: 0.01,
			want:   "0.01",
		},
		{
			name:   "large number",
			amount: 999999.99,
			want:   "999999.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAmount(tt.amount)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  string
		want    float64
		wantErr bool
	}{
		{
			name:    "integer",
			amount:  "100.00",
			want:    100.00,
			wantErr: false,
		},
		{
			name:    "decimal",
			amount:  "100.50",
			want:    100.50,
			wantErr: false,
		},
		{
			name:    "invalid",
			amount:  "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAmount(tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGenerateSerialNumber(t *testing.T) {
	prefix := "TEST"
	serial := GenerateSerialNumber(prefix)
	assert.NotEmpty(t, serial)
	assert.True(t, len(serial) > len(prefix))
	assert.True(t, len(serial) < len(prefix)+30)
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "with spaces",
			s:    "  test  ",
			want: "test",
		},
		{
			name: "without spaces",
			s:    "test",
			want: "test",
		},
		{
			name: "only spaces",
			s:    "   ",
			want: "",
		},
		{
			name: "empty string",
			s:    "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimSpace(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToUpperCase(t *testing.T) {
	assert.Equal(t, "TEST", ToUpperCase("test"))
	assert.Equal(t, "TEST", ToUpperCase("TeSt"))
	assert.Equal(t, "", ToUpperCase(""))
}

func TestToLowerCase(t *testing.T) {
	assert.Equal(t, "test", ToLowerCase("TEST"))
	assert.Equal(t, "test", ToLowerCase("TeSt"))
	assert.Equal(t, "", ToLowerCase(""))
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "empty string",
			s:    "",
			want: true,
		},
		{
			name: "only spaces",
			s:    "   ",
			want: true,
		},
		{
			name: "with tabs",
			s:    "\t\t",
			want: true,
		},
		{
			name: "with newlines",
			s:    "\n\n",
			want: true,
		},
		{
			name: "non-empty",
			s:    "test",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmpty(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaskString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		showLen int
		want    string
	}{
		{
			name:    "normal string",
			s:       "1234567890",
			showLen: 2,
			want:    "12******90",
		},
		{
			name:    "short string",
			s:       "1234",
			showLen: 2,
			want:    "1234",
		},
		{
			name:    "empty string",
			s:       "",
			showLen: 2,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskString(tt.s, tt.showLen)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaskMobile(t *testing.T) {
	tests := []struct {
		name   string
		mobile string
		want   string
	}{
		{
			name:   "valid mobile",
			mobile: "13812345678",
			want:   "138****5678",
		},
		{
			name:   "invalid length",
			mobile: "123456",
			want:   "123456",
		},
		{
			name:   "empty string",
			mobile: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskMobile(tt.mobile)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		name   string
		idCard string
		want   string
	}{
		{
			name:   "valid id card",
			idCard: "123456199001011234",
			want:   "123456********1234",
		},
		{
			name:   "invalid length",
			idCard: "123456",
			want:   "123456",
		},
		{
			name:   "empty string",
			idCard: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskIDCard(tt.idCard)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaskBankCard(t *testing.T) {
	tests := []struct {
		name     string
		bankCard string
		want     string
	}{
		{
			name:     "valid bank card",
			bankCard: "1234567890123456789",
			want:     "123456********6789",
		},
		{
			name:     "short card",
			bankCard: "12345678901",
			want:     "123456********8901",
		},
		{
			name:     "empty string",
			bankCard: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskBankCard(tt.bankCard)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateMobile(t *testing.T) {
	tests := []struct {
		name   string
		mobile string
		want   bool
	}{
		{
			name:   "valid mobile",
			mobile: "13812345678",
			want:   true,
		},
		{
			name:   "invalid length",
			mobile: "123456",
			want:   false,
		},
		{
			name:   "invalid prefix",
			mobile: "23812345678",
			want:   false,
		},
		{
			name:   "empty string",
			mobile: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateMobile(tt.mobile)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateIDCard(t *testing.T) {
	tests := []struct {
		name   string
		idCard string
		want   bool
	}{
		{
			name:   "valid id card",
			idCard: "123456199001011234",
			want:   true,
		},
		{
			name:   "invalid length",
			idCard: "123456",
			want:   false,
		},
		{
			name:   "empty string",
			idCard: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateIDCard(tt.idCard)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateBankCard(t *testing.T) {
	tests := []struct {
		name     string
		bankCard string
		want     bool
	}{
		{
			name:     "valid bank card",
			bankCard: "123456789012345678",
			want:     true,
		},
		{
			name:     "short card",
			bankCard: "1234567890",
			want:     true,
		},
		{
			name:     "too short",
			bankCard: "123456789",
			want:     false,
		},
		{
			name:     "too long",
			bankCard: "123456789012345678901",
			want:     false,
		},
		{
			name:     "empty string",
			bankCard: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateBankCard(tt.bankCard)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContains(t *testing.T) {
	assert.True(t, Contains("hello world", "world"))
	assert.False(t, Contains("hello world", "test"))
}

func TestSplit(t *testing.T) {
	result := Split("a,b,c", ",")
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestJoin(t *testing.T) {
	result := Join([]string{"a", "b", "c"}, ",")
	assert.Equal(t, "a,b,c", result)
}

func TestHasPrefix(t *testing.T) {
	assert.True(t, HasPrefix("hello world", "hello"))
	assert.False(t, HasPrefix("hello world", "world"))
}

func TestHasSuffix(t *testing.T) {
	assert.True(t, HasSuffix("hello world", "world"))
	assert.False(t, HasSuffix("hello world", "hello"))
}
