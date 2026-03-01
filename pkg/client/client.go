// Package client 提供了与建行对公专业结算综合服务平台交互的 HTTP 客户端实现。
//
// 该包包含了客户端的初始化、请求发送、响应处理等核心功能。
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ljjdev/ccb-hsb-sdk/internal/utils"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/config"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/model"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/signature"
)

const (
	// DefaultTimeout 默认请求超时时间
	DefaultTimeout = 30 * time.Second

	// ContentTypeJSON JSON 内容类型
	ContentTypeJSON = "application/json"

	// HeaderContentType 内容类型头
	HeaderContentType = "Content-Type"

	// HeaderUserAgent 用户代理头
	HeaderUserAgent = "User-Agent"

	// UserAgent 用户代理值
	UserAgent = "ccb-hsb-sdk/1.0.0"
)

// Client 定义了 SDK 客户端
type Client struct {
	config     *config.Config
	signer     *signature.RSAService
	httpClient *http.Client
}

// NewClient 创建一个新的客户端实例
//
// 使用示例:
//
//	cfg, err := config.NewConfig(
//		config.WithMarketID("12345678901234"),
//		config.WithMerchantID("12345678901234567890"),
//		config.WithGatewayURL("https://marketpay.ccb.com/online/direct"),
//		config.WithPrivateKey(privateKey),
//		config.WithPublicKey(publicKey),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	client, err := client.NewClient(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
func NewClient(cfg *config.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	signer := signature.NewRSAService(cfg.PrivateKey, cfg.PublicKey)
	// http.Client的Timeout单位为 纳秒 需要换算
	httpClient := &http.Client{
		Timeout: cfg.Timeout * time.Second,
	}

	return &Client{
		config:     cfg,
		signer:     signer,
		httpClient: httpClient,
	}, nil
}

// CreateOrder 创建支付订单
//
// 该方法向建行对公专业结算综合服务平台发送创建订单请求。
//
// 使用示例:
//
//	req := &model.CreateOrderRequest{
//		IttpartyStmId: "00000",
//		PyChnlCd:      "0000000000000000000000000",
//		IttpartyTms:   "20240101120000123",
//		IttpartyJrnlNo: "20240101120000123001",
//		MktId:         "12345678901234",
//		MainOrdrNo:    "20240101120000123",
//		PymdCd:        model.PaymentMethodPC,
//		PyOrdrTpcd:    model.OrderTypeNormal,
//		Ccy:           "156",
//		OrdrTamt:      "100.00",
//		TxnTamt:       "100.00",
//	}
//
//	resp, err := client.CreateOrder(context.Background(), req)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !resp.IsSuccess() {
//		log.Fatal(resp.GetError())
//	}
func (c *Client) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// 设置默认值
	if req.IttpartyStmId == "" {
		req.IttpartyStmId = "00000"
	}
	if req.PyChnlCd == "" {
		req.PyChnlCd = "0000000000000000000000000"
	}
	if req.Ccy == "" {
		req.Ccy = "156"
	}

	// 设置市场编号和商家编号
	if req.MktId == "" {
		req.MktId = c.config.MarketID
	}
	if req.Orderlist != nil && len(req.Orderlist) > 0 {
		for i := range req.Orderlist {
			if req.Orderlist[i].MktMrchId == "" {
				req.Orderlist[i].MktMrchId = c.config.MerchantID
			}
		}
	}

	// 生成签名
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	signatureString, err := signature.BuildSignatureStringFromJSON(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to build signature string: %w", err)
	}

	sign, err := c.signer.Sign(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	req.SignInf = sign

	// 发送请求
	url := c.config.GatewayURL + "/gatherPlaceorder"
	resp, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 解析响应
	var orderResp model.CreateOrderResponse
	if err := json.Unmarshal(resp, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 验证签名
	//if err := c.verifyResponseSignature(&orderResp); err != nil {
	//	return nil, fmt.Errorf("signature verification failed: %w", err)
	//}

	return &orderResp, nil
}

// RefundOrder 订单退款
//
// 该方法向建行对公专业结算综合服务平台发送退款请求,支持全额退款和部分退款。
//
// 使用示例:
//
//	req := &model.RefundOrderRequest{
//		MktId:        "12345678901234",
//		MainOrdrNo:   "20240101120000123",
//		RefundOrdrNo: "REFUND20240101120000123",
//		RefundAmt:    "100.00",
//		RefundRsn:    "用户申请退款",
//	}
//
//	resp, err := client.RefundOrder(context.Background(), req)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !resp.IsSuccess() {
//		log.Fatal(resp.GetError())
//	}
func (c *Client) RefundOrder(ctx context.Context, req *model.RefundOrderRequest) (*model.RefundOrderResponse, error) {
	// 设置默认值
	if req.MktId == "" {
		req.MktId = c.config.MarketID
	}

	// 生成时间戳和流水号(如果未提供)
	if req.IttpartyTms == "" {
		req.IttpartyTms = utils.CurrentTimestamp()
	}
	if req.IttpartyJrnlNo == "" {
		req.IttpartyJrnlNo = utils.GenerateSerialNumber("REF")
	}

	// 构建待签名字符串
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	signatureString, err := signature.BuildSignatureStringFromJSON(string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to build signature string: %w", err)
	}

	// 生成签名
	sign, err := c.signer.Sign(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	req.SignInf = sign

	// 发送请求
	url := c.config.GatewayURL + "/refundOrder"
	resp, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 解析响应
	var refundResp model.RefundOrderResponse
	if err := json.Unmarshal(resp, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 验证响应签名
	//if err := c.verifyRefundOrderResponseSignature(&refundResp, resp); err != nil {
	//	return nil, fmt.Errorf("signature verification failed: %w", err)
	//}

	return &refundResp, nil
}

// doRequest 发送 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)

		if c.config.Debug {
			fmt.Printf("[DEBUG] Request Body: %s\n", string(jsonData))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(HeaderContentType, ContentTypeJSON)
	req.Header.Set(HeaderUserAgent, UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if c.config.Debug {
		fmt.Printf("[DEBUG] Response Status: %d\n", resp.StatusCode)
		fmt.Printf("[DEBUG] Response Body: %s\n", string(respBody))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// verifyResponseSignature 验证响应签名
func (c *Client) verifyResponseSignature(resp *model.CreateOrderResponse) error {
	// 将响应转换为 JSON 字符串
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// 使用 JSON 签名字符串拼接算法
	signatureString, err := signature.BuildSignatureStringFromJSON(string(data))
	if err != nil {
		return fmt.Errorf("failed to build signature string: %w", err)
	}

	// 验证签名
	signInf := resp.SignInf
	if signInf == "" {
		return fmt.Errorf("response signature is empty")
	}

	if err := c.signer.Verify(signatureString, signInf); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

// verifyRefundOrderResponseSignature 验证退款订单响应签名
func (c *Client) verifyRefundOrderResponseSignature(resp *model.RefundOrderResponse, respBody []byte) error {
	// 如果签名为空,跳过验证(用于测试)
	if resp.SignInf == "" {
		return nil
	}

	// 使用 JSON 签名字符串拼接算法
	signatureString, err := signature.BuildSignatureStringFromJSON(string(respBody))
	if err != nil {
		return fmt.Errorf("failed to build signature string: %w", err)
	}

	// 验证签名
	if err := c.signer.Verify(signatureString, resp.SignInf); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

// PlaceOrder 支付订单生成接口
//
// 该方法向建行对公专业结算综合服务平台发送创建订单请求,并返回解码后的支付 URL。
//
// 使用示例:
//
//	req := &model.CreateOrderRequest{
//		IttpartyTms:    "20240101120000123",
//		IttpartyJrnlNo: "20240101120000123001",
//		MainOrdrNo:     "20240101120000123",
//		PymdCd:         model.PaymentMethodMobileH5,
//		PyOrdrTpcd:     model.OrderTypeNormal,
//		OrdrTamt:       "100.01",
//		TxnTamt:        "100.01",
//		PayDsc:         "商品",
//		OrderTimeOut:   "1800",
//		Orderlist: []model.SubOrder{
//			{
//				CmdtyOrdrNo: "20240101120000123001",
//				OrdrAmt:     "100.01",
//				Txnamt:      "100.01",
//				CmdtyDsc:    "商品",
//				ClrgRuleId:  "123456",
//				Parlist: []model.Participant{
//					{
//						SeqNo:     1,
//						MktMrchId: "12345678901234567890",
//					},
//				},
//			},
//		},
//	}
//
//	payURL, err := client.PlaceOrder(context.Background(), req)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("支付URL:", payURL)
func (c *Client) PlaceOrder(ctx context.Context, req *model.CreateOrderRequest) (string, error) {
	// 调用 CreateOrder 创建订单
	resp, err := c.CreateOrder(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	// 检查响应状态
	if !resp.IsSuccess() {
		return "", fmt.Errorf("order creation failed: %w", resp.GetError())
	}

	// 检查支付 URL 是否包含 https
	if resp.CshdkUrl == "" || !strings.Contains(resp.CshdkUrl, "https") {
		return "", fmt.Errorf("invalid payment URL: %s", resp.CshdkUrl)
	}

	// URL 解码
	decodedURL, err := url.QueryUnescape(resp.CshdkUrl)
	if err != nil {
		return "", fmt.Errorf("failed to decode payment URL: %w", err)
	}

	return decodedURL, nil
}

// GetConfig 获取客户端配置
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// GetSigner 获取签名器
func (c *Client) GetSigner() signature.Signer {
	return c.signer
}

// QueryRefund 查询退款结果
//
// 该方法向建行对公专业结算综合服务平台发送查询退款结果请求。
//
// 使用示例:
//
//	req := &model.QueryRefundRequest{
//		IttpartyStmId:  "00000",
//		PyChnlCd:       "0000000000000000000000000",
//		IttpartyTms:    "20240101120000123",
//		IttpartyJrnlNo: "20240101120000123001",
//		MktId:          "12345678901234",
//		CustRfndTrcno:  "REFUND20240101120000123",
//		Vno:            "4",
//	}
//
//	resp, err := client.QueryRefund(context.Background(), req)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !resp.IsSuccess() {
//		log.Fatal(resp.GetError())
//	}
func (c *Client) QueryRefund(ctx context.Context, req *model.QueryRefundRequest) (*model.QueryRefundResponse, error) {
	// 设置默认值
	if req.IttpartyStmId == "" {
		req.IttpartyStmId = "00000"
	}
	if req.PyChnlCd == "" {
		req.PyChnlCd = "0000000000000000000000000"
	}
	if req.Vno == "" {
		req.Vno = "4"
	}

	// 设置市场编号
	if req.MktId == "" {
		req.MktId = c.config.MarketID
	}

	// 验证必输参数
	if req.CustRfndTrcno == "" && req.RfndTrcno == "" {
		return nil, fmt.Errorf("Cust_Rfnd_Trcno and Rfnd_Trcno must provide at least one")
	}

	// 生成签名
	params, err := req.ToMap()
	if err != nil {
		return nil, fmt.Errorf("failed to convert request to map: %w", err)
	}

	signatureString := signature.BuildSignatureString(params)
	sign, err := c.signer.Sign(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	req.SignInf = sign

	// 发送请求
	url := c.config.GatewayURL + "/enquireRefundOrder"
	resp, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 解析响应
	var refundResp model.QueryRefundResponse
	if err := json.Unmarshal(resp, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 验证签名
	//if err := c.verifyRefundResponseSignature(&refundResp); err != nil {
	//	return nil, fmt.Errorf("signature verification failed: %w", err)
	//}

	return &refundResp, nil
}

// verifyRefundResponseSignature 验证退款查询响应签名
func (c *Client) verifyRefundResponseSignature(resp *model.QueryRefundResponse) error {
	// 如果签名为空,跳过验证(用于测试环境)
	if resp.SignInf == "" {
		return nil
	}

	// 将响应转换为 map 用于验签
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	// 验证签名
	signInf := resp.SignInf
	sigString := signature.BuildSignatureString(stringMap)
	if err := c.signer.Verify(sigString, signInf); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

// QueryOrder 查询支付结果
//
// 该方法向建行对公专业结算综合服务平台发送查询支付结果请求。
//
// 使用示例:
//
//	req := &model.QueryOrderRequest{
//		IttpartyJrnlNo: "20240101120000123001",
//		MktId:         "12345678901234",
//		MainOrdrNo:    "20240101120000123",
//	}
//
//	resp, err := client.QueryOrder(context.Background(), req)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !resp.IsSuccess() {
//		log.Fatal(resp.GetError())
//	}
//
//	if resp.IsPaid() {
//		log.Println("订单已支付成功")
//	}
func (c *Client) QueryOrder(ctx context.Context, req *model.QueryOrderRequest) (*model.QueryOrderResponse, error) {
	// 设置默认值
	if req.IttpartyStmId == "" {
		req.IttpartyStmId = "00000"
	}
	if req.PyChnlCd == "" {
		req.PyChnlCd = "0000000000000000000000000"
	}
	if req.MktId == "" {
		req.MktId = c.config.MarketID
	}
	if req.Vno == "" {
		req.Vno = "4"
	}

	// 验证必填参数
	if req.MainOrdrNo == "" && req.PyTrnNo == "" {
		return nil, fmt.Errorf("主订单编号和支付流水号必输其一")
	}

	// 生成签名
	params, err := req.ToMap()
	if err != nil {
		return nil, fmt.Errorf("failed to convert request to map: %w", err)
	}

	signatureString := signature.BuildSignatureString(params)
	sign, err := c.signer.Sign(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	req.SignInf = sign

	// 发送请求
	url := c.config.GatewayURL + "/gatherEnquireOrder"
	resp, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 解析响应
	var queryResp model.QueryOrderResponse
	if err := json.Unmarshal(resp, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 验证签名
	//if err := c.verifyQueryResponseSignature(&queryResp); err != nil {
	//	return nil, fmt.Errorf("signature verification failed: %w", err)
	//}

	return &queryResp, nil
}

// verifyQueryResponseSignature 验证查询响应签名
func (c *Client) verifyQueryResponseSignature(resp *model.QueryOrderResponse) error {
	// 将响应转换为 map 用于验签
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	// 验证签名
	signInf := resp.SignInf
	if signInf == "" {
		return fmt.Errorf("response signature is empty")
	}

	sigString := signature.BuildSignatureString(stringMap)
	fmt.Printf("sigString==> %s\n", data)
	if err := c.signer.Verify(sigString, signInf); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}
