// Package utils 提供了内部使用的工具函数。
//
// 该包包含了字符串处理、时间格式化、数据转换等辅助函数,
// 仅在项目内部使用,不对外暴露。
package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// TimeFormat 时间格式: yyyyMMddHHmmssfff
	TimeFormat = "20060102150405.000"

	// TimeFormatNoDot 时间格式(无点): yyyyMMddHHmmssfff
	TimeFormatNoDot = "20060102150405000"

	// TimeFormatShort 短时间格式: yyyyMMdd
	TimeFormatShort = "20060102"

	// TimeFormatLong 长时间格式: yyyyMMddHHmmss
	TimeFormatLong = "20060102150405"
)

// CurrentTimestamp 获取当前时间戳
//
// 返回格式为 yyyyMMddHHmmssfff 的时间戳字符串。
func CurrentTimestamp() string {
	return time.Now().Format(TimeFormatNoDot)
}

// CurrentTimestampShort 获取当前短时间戳
//
// 返回格式为 yyyyMMdd 的时间戳字符串。
func CurrentTimestampShort() string {
	return time.Now().Format(TimeFormatShort)
}

// CurrentTimestampLong 获取当前长时间戳
//
// 返回格式为 yyyyMMddHHmmss 的时间戳字符串。
func CurrentTimestampLong() string {
	return time.Now().Format(TimeFormatLong)
}

// ParseTimestamp 解析时间戳
//
// 该函数解析 yyyyMMddHHmmssfff 格式的时间戳字符串。
func ParseTimestamp(timestamp string) (time.Time, error) {
	if len(timestamp) != 17 {
		return time.Time{}, fmt.Errorf("invalid timestamp format, expected 17 characters, got %d", len(timestamp))
	}

	// 手动解析时间戳
	year := timestamp[0:4]
	month := timestamp[4:6]
	day := timestamp[6:8]
	hour := timestamp[8:10]
	minute := timestamp[10:12]
	second := timestamp[12:14]
	millisecond := timestamp[14:17]

	// 构建标准时间格式字符串
	stdFormat := fmt.Sprintf("%s-%s-%s %s:%s:%s.%s", year, month, day, hour, minute, second, millisecond)

	// 使用标准格式解析
	return time.Parse("2006-01-02 15:04:05.000", stdFormat)
}

// ParseTimestampShort 解析短时间戳
//
// 该函数解析 yyyyMMdd 格式的时间戳字符串。
func ParseTimestampShort(timestamp string) (time.Time, error) {
	return time.Parse(TimeFormatShort, timestamp)
}

// ParseTimestampLong 解析长时间戳
//
// 该函数解析 yyyyMMddHHmmss 格式的时间戳字符串。
func ParseTimestampLong(timestamp string) (time.Time, error) {
	return time.Parse(TimeFormatLong, timestamp)
}

// FormatAmount 格式化金额
//
// 该函数将金额格式化为字符串,保留两位小数。
func FormatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// ParseAmount 解析金额
//
// 该函数解析金额字符串为浮点数。
func ParseAmount(amount string) (float64, error) {
	return strconv.ParseFloat(amount, 64)
}

// GenerateSerialNumber 生成流水号
//
// 该函数生成唯一的流水号,格式为: 前缀 + 时间戳 + 随机数。
func GenerateSerialNumber(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%06d", prefix, timestamp, time.Now().Nanosecond()%1000000)
}

// TrimSpace 去除字符串首尾空格
//
// 该函数去除字符串首尾的空格、制表符、换行符等空白字符。
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// ToUpperCase 转换为大写
//
// 该函数将字符串转换为大写。
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase 转换为小写
//
// 该函数将字符串转换为小写。
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// IsEmpty 判断字符串是否为空
//
// 该函数判断字符串是否为空或仅包含空白字符。
func IsEmpty(s string) bool {
	return TrimSpace(s) == ""
}

// MaskString 掩码字符串
//
// 该函数对字符串进行掩码处理,保留前 showLen 和后 showLen 个字符,中间用 * 替换。
func MaskString(s string, showLen int) string {
	if len(s) <= showLen*2 {
		return s
	}
	return s[:showLen] + strings.Repeat("*", len(s)-showLen*2) + s[len(s)-showLen:]
}

// MaskMobile 掩码手机号
//
// 该函数对手机号进行掩码处理,格式为: 138****1234。
func MaskMobile(mobile string) string {
	if len(mobile) != 11 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}

// MaskIDCard 掩码身份证号
//
// 该函数对身份证号进行掩码处理,保留前6位和后4位。
func MaskIDCard(idCard string) string {
	if len(idCard) != 18 {
		return idCard
	}
	return idCard[:6] + "********" + idCard[14:]
}

// MaskBankCard 掩码银行卡号
//
// 该函数对银行卡号进行掩码处理,保留前6位和后4位。
func MaskBankCard(bankCard string) string {
	if len(bankCard) < 10 {
		return bankCard
	}
	return bankCard[:6] + "********" + bankCard[len(bankCard)-4:]
}

// ValidateMobile 验证手机号
//
// 该函数验证手机号格式是否正确。
func ValidateMobile(mobile string) bool {
	if len(mobile) != 11 {
		return false
	}
	return mobile[0] == '1'
}

// ValidateIDCard 验证身份证号
//
// 该函数验证身份证号格式是否正确。
func ValidateIDCard(idCard string) bool {
	// TODO: 实现完整的身份证号验证逻辑
	return len(idCard) == 18
}

// ValidateBankCard 验证银行卡号
//
// 该函数验证银行卡号格式是否正确。
func ValidateBankCard(bankCard string) bool {
	// TODO: 实现完整的银行卡号验证逻辑
	return len(bankCard) >= 10 && len(bankCard) <= 19
}

// Contains 判断字符串是否包含子串
//
// 该函数判断字符串是否包含指定的子串。
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Split 分割字符串
//
// 该函数按指定的分隔符分割字符串。
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Join 连接字符串
//
// 该函数将字符串数组按指定的分隔符连接。
func Join(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// HasPrefix 判断字符串是否以指定前缀开头
//
// 该函数判断字符串是否以指定的前缀开头。
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix 判断字符串是否以指定后缀结尾
//
// 该函数判断字符串是否以指定的后缀结尾。
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// LoadPrivateKey 将没有头尾的 Base64 字符串解析为 RSA 私钥
func LoadPrivateKey(rawStr string) (*rsa.PrivateKey, error) {
	der, err := base64.StdEncoding.DecodeString(rawStr)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %v", err)
	}

	// 优先尝试 PKCS#8 格式
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		if priv, ok := key.(*rsa.PrivateKey); ok {
			return priv, nil
		}
	}

	// 备选尝试 PKCS#1 格式
	return x509.ParsePKCS1PrivateKey(der)
}

// LoadPublicKey 将没有头尾的 Base64 字符串解析为 RSA 公钥
func LoadPublicKey(rawStr string) (*rsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(rawStr)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %v", err)
	}

	// 优先尝试 PKIX 格式 (最常用)
	if pub, err := x509.ParsePKIXPublicKey(der); err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
	}

	// 备选尝试 PKCS#1 格式
	return x509.ParsePKCS1PublicKey(der)
}
