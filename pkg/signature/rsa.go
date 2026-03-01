// Package signature 提供了 RSA 签名和验签功能。
//
// 该包实现了建行对公专业结算综合服务平台所需的签名算法,
// 包括请求签名生成和响应签名验证。
package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	// SignatureAlgorithm 签名算法
	SignatureAlgorithm = "SHA256withRSA"
)

// Signer 定义了签名器的接口
type Signer interface {
	// Sign 对数据进行签名
	Sign(data string) (string, error)
}

// Verifier 定义了验签器的接口
type Verifier interface {
	// Verify 验证签名
	Verify(data string, signature string) error
}

// RSAService 提供了 RSA 签名和验签服务
type RSAService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSAService 创建一个新的 RSA 服务实例
//
// privateKey 用于签名,publicKey 用于验签。
func NewRSAService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RSAService {
	return &RSAService{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Sign 对数据进行 RSA 签名
//
// 该函数使用 SHA256withRSA 算法对数据进行签名,返回 Base64 编码的签名结果。
//
// 签名流程:
// 1. 对原始数据进行 SHA256 哈希
// 2. 使用私钥对哈希值进行 RSA 签名
// 3. 对签名结果进行 Base64 编码
func (s *RSAService) Sign(data string) (string, error) {
	if s.privateKey == nil {
		return "", fmt.Errorf("private key is not set")
	}

	// 计算 SHA256 哈希
	hashed := sha256.Sum256([]byte(data))

	// 使用私钥进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// Base64 编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify 验证 RSA 签名
//
// 该函数验证签名是否有效,如果签名无效则返回错误。
//
// 验签流程:
// 1. 对 Base64 编码的签名进行解码
// 2. 对原始数据进行 SHA256 哈希
// 3. 使用公钥验证签名
func (s *RSAService) Verify(data string, signature string) error {
	if s.publicKey == nil {
		return fmt.Errorf("public key is not set")
	}

	// Base64 解码
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// 计算 SHA256 哈希
	hashed := sha256.Sum256([]byte(data))

	// 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, hashed[:], sigBytes)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// BuildSignatureString 构建待签名字符串
//
// 该函数按照建行规范将参数拼接成待签名字符串。
//
// 拼接规则:
// 1. 将所有参数按字典序排序
// 2. 使用 key=value&key=value 的格式拼接
// 3. 忽略值为空的参数
// 4. 忽略 Sign_Inf、Svc_Rsp_St、Svc_Rsp_Cd、Rsp_Inf 参数
func BuildSignatureString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	// 提取并排序参数键
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "Sign_Inf" && k != "Svc_Rsp_St" && k != "Svc_Rsp_Cd" && k != "Rsp_Inf" && params[k] != "" {
			keys = append(keys, k)
		}
	}

	// 按字典序排序
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// 拼接参数
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(k)
		builder.WriteByte('=')
		builder.WriteString(params[k])
	}

	return builder.String()
}

// SignParams 对参数进行签名
//
// 该函数是一个便捷方法,将参数构建为待签名字符串后进行签名。
func (s *RSAService) SignParams(params map[string]string) (string, error) {
	signatureString := BuildSignatureString(params)
	return s.Sign(signatureString)
}

// VerifyParams 验证参数签名
//
// 该函数是一个便捷方法,将参数构建为待签名字符串后进行验签。
func (s *RSAService) VerifyParams(params map[string]string, signature string) error {
	signatureString := BuildSignatureString(params)
	return s.Verify(signatureString, signature)
}

// LoadPrivateKeyFromPEM 从 PEM 格式加载私钥
//
// 该函数从 PEM 格式的字符串中解析 RSA 私钥。
func LoadPrivateKeyFromPEM(pemData []byte) (*rsa.PrivateKey, error) {
	// TODO: 实现 PEM 格式私钥的加载逻辑
	return nil, fmt.Errorf("not implemented")
}

// LoadPublicKeyFromPEM 从 PEM 格式加载公钥
//
// 该函数从 PEM 格式的字符串中解析 RSA 公钥。
func LoadPublicKeyFromPEM(pemData []byte) (*rsa.PublicKey, error) {
	// TODO: 实现 PEM 格式公钥的加载逻辑
	return nil, fmt.Errorf("not implemented")
}
