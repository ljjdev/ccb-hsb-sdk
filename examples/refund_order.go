// Package examples 提供了建行对公专业结算综合服务平台 SDK 的使用示例。
//
// 本包包含了各种常见场景的使用示例,帮助开发者快速上手。
package examples

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/ljjdev/ccb-hsb-sdk/internal/utils"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/client"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/config"
	"github.com/ljjdev/ccb-hsb-sdk/pkg/model"
)

// RefundOrderExample 演示如何使用订单退款接口
//
// 本示例展示了完整的订单退款流程,包括:
// 1. 配置初始化
// 2. 客户端创建
// 3. 构建退款请求
// 4. 调用 API
// 5. 处理响应和错误
func RefundOrderExample() {
	// 1. 准备 RSA 密钥对
	// 实际使用时,应该从文件加载密钥
	// privateKey, err := utils.LoadPrivateKeyFromFile("path/to/private_key.pem")
	// publicKey, err := utils.LoadPublicKeyFromFile("path/to/public_key.pem")
	// 这里使用示例密钥(实际使用时替换为真实密钥)
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	// 2. 创建配置
	cfg, err := config.NewConfig(
		config.WithMarketID("12345678901234"),                            // 市场编号,由银行提供
		config.WithMerchantID("12345678901234567890"),                    // 商家编号,由银行提供
		config.WithGatewayURL("https://marketpay.ccb.com/online/direct"), // 接口网关地址
		config.WithPrivateKey(privateKey),                                // 商户私钥,用于签名
		config.WithPublicKey(publicKey),                                  // 银行公钥,用于验签
		config.WithTimeout(30),                                           // 请求超时时间(秒)
		config.WithDebug(true),                                           // 开启调试模式,输出请求和响应日志
	)
	if err != nil {
		log.Fatalf("创建配置失败: %v", err)
	}

	// 3. 创建客户端
	cli, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 4. 构建退款订单请求
	// 退款订单编号: 不允许重复,建议使用 UUID 或时间戳+随机数
	refundOrderNo := utils.GenerateSerialNumber("REF")
	// 主订单编号: 需要退款的原订单编号
	mainOrderNo := "ORD20240101120000123001"
	// 退款金额: 100.01 元
	refundAmount := "100.01"

	req := &model.RefundOrderRequest{
		// 基础信息
		MktId:        "12345678901234", // 市场编号
		MainOrdrNo:   mainOrderNo,      // 主订单编号(需要退款的原订单)
		RefundOrdrNo: refundOrderNo,    // 退款订单编号,不允许重复
		RefundAmt:    refundAmount,     // 退款金额
		RefundRsn:    "用户申请退款",         // 退款原因

		// 时间戳和流水号(如果不提供,SDK 会自动生成)
		IttpartyTms:    utils.CurrentTimestamp(),          // 发起方时间戳,格式: yyyyMMddHHmmssfff
		IttpartyJrnlNo: utils.GenerateSerialNumber("RFD"), // 发起方流水号,不允许重复
	}

	// 5. 调用订单退款接口
	// 该方法会自动进行签名、发送请求、验证签名等操作
	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		log.Fatalf("订单退款失败: %v", err)
	}

	// 6. 处理响应
	if resp.IsSuccess() {
		fmt.Println("订单退款成功!")
		fmt.Printf("主订单编号: %s\n", resp.MainOrdrNo)
		fmt.Printf("退款订单编号: %s\n", resp.RefundOrdrNo)
		fmt.Printf("退款金额: %s 元\n", resp.RefundAmt)
		fmt.Printf("退款流水号: %s\n", resp.RefundTrnNo)
		fmt.Printf("退款时间: %s\n", resp.RefundTm)
		fmt.Printf("退款状态: %s\n", resp.RefundStcd)
	} else {
		fmt.Printf("订单退款失败: %v\n", resp.GetError())
	}
}

// RefundOrderWithMinimalParams 演示使用最少参数进行退款
//
// 本示例展示了如何使用最少参数进行退款,SDK 会自动填充默认值
func RefundOrderWithMinimalParams(cli *client.Client, mainOrderNo string, refundAmount string) error {
	// 使用最少参数,其他参数会自动填充默认值
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,   // 必填: 主订单编号
		RefundOrdrNo: refundOrderNo, // 必填: 退款订单编号
		RefundAmt:    refundAmount,  // 必填: 退款金额
		RefundRsn:    "用户申请退款",      // 可选: 退款原因
	}

	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("订单退款失败: %w", err)
	}

	if resp.IsSuccess() {
		fmt.Printf("退款成功,流水号: %s, 金额: %s 元\n", resp.RefundTrnNo, resp.RefundAmt)
	}

	return nil
}

// RefundOrderWithErrorHandling 演示完整的错误处理
//
// 本示例展示了如何处理各种可能的错误情况
func RefundOrderWithErrorHandling(cli *client.Client, mainOrderNo string, refundAmount string) (*model.RefundOrderResponse, error) {
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    refundAmount,
		RefundRsn:    "用户申请退款",
	}

	// 调用 API
	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		// 处理错误
		// 可能的错误类型:
		// 1. 配置错误: 私钥/公钥无效
		// 2. 网络错误: 连接超时、网络不可达
		// 3. 签名错误: 签名生成失败
		// 4. 验签错误: 响应签名验证失败
		// 5. 业务错误: 退款失败
		return nil, fmt.Errorf("订单退款失败: %w", err)
	}

	// 检查响应是否成功
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("订单退款失败: %w", resp.GetError())
	}

	return resp, nil
}

// RefundOrderFull 演示全额退款
//
// 本示例展示了如何进行全额退款
func RefundOrderFull(cli *client.Client, mainOrderNo string, totalAmount string) error {
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    totalAmount, // 退款金额等于订单总金额
		RefundRsn:    "全额退款",
	}

	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("全额退款失败: %w", err)
	}

	if resp.IsSuccess() {
		fmt.Printf("全额退款成功,流水号: %s, 金额: %s 元\n", resp.RefundTrnNo, resp.RefundAmt)
	}

	return nil
}

// RefundOrderPartial 演示部分退款
//
// 本示例展示了如何进行部分退款
// 注意: 部分退款可以多次进行,但累计退款金额不能超过订单总金额
func RefundOrderPartial(cli *client.Client, mainOrderNo string, refundAmount string) error {
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    refundAmount, // 部分退款金额
		RefundRsn:    "部分退款",
	}

	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("部分退款失败: %w", err)
	}

	if resp.IsSuccess() {
		fmt.Printf("部分退款成功,流水号: %s, 金额: %s 元\n", resp.RefundTrnNo, resp.RefundAmt)
	}

	return nil
}

// RefundOrderWithReason 演示带退款原因的退款
//
// 本示例展示了如何指定不同的退款原因
func RefundOrderWithReason(cli *client.Client, mainOrderNo string, refundAmount string, reason string) error {
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    refundAmount,
		RefundRsn:    reason, // 指定退款原因
	}

	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("退款失败: %w", err)
	}

	if resp.IsSuccess() {
		fmt.Printf("退款成功,原因: %s, 金额: %s 元\n", reason, resp.RefundAmt)
	}

	return nil
}

// RefundOrderAndQuery 演示退款后查询退款结果
//
// 本示例展示了如何进行退款并查询退款结果
func RefundOrderAndQuery(cli *client.Client, mainOrderNo string, refundAmount string) error {
	// 1. 发起退款
	refundOrderNo := utils.GenerateSerialNumber("REF")

	refundReq := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    refundAmount,
		RefundRsn:    "用户申请退款",
	}

	refundResp, err := cli.RefundOrder(context.Background(), refundReq)
	if err != nil {
		return fmt.Errorf("订单退款失败: %w", err)
	}

	if !refundResp.IsSuccess() {
		return fmt.Errorf("订单退款失败: %w", refundResp.GetError())
	}

	fmt.Printf("退款申请成功,退款流水号: %s\n", refundResp.RefundTrnNo)

	// 2. 查询退款结果
	queryReq := &model.QueryRefundRequest{
		CustRfndTrcno: refundOrderNo, // 使用客户方退款流水号查询
	}

	queryResp, err := cli.QueryRefund(context.Background(), queryReq)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	if queryResp.IsSuccess() {
		if queryResp.RfndAmt != nil {
			fmt.Printf("退款成功,金额: %s 元\n", *queryResp.RfndAmt)
		} else {
			fmt.Println("退款成功")
		}
	}

	return nil
}

// RefundOrderWithValidation 演示退款前的参数验证
//
// 本示例展示了如何在退款前进行参数验证
func RefundOrderWithValidation(cli *client.Client, mainOrderNo string, refundAmount string) error {
	// 1. 参数验证
	if mainOrderNo == "" {
		return fmt.Errorf("主订单编号不能为空")
	}

	if refundAmount == "" {
		return fmt.Errorf("退款金额不能为空")
	}

	// 2. 验证退款金额格式
	// 实际业务中,应该验证退款金额是否为有效的金额格式
	// 这里仅作示例

	// 3. 发起退款
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    refundAmount,
		RefundRsn:    "用户申请退款",
	}

	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("订单退款失败: %w", err)
	}

	if resp.IsSuccess() {
		fmt.Printf("退款成功,流水号: %s, 金额: %s 元\n", resp.RefundTrnNo, resp.RefundAmt)
	}

	return nil
}

// RefundOrderWithRetry 演示带重试机制的退款
//
// 本示例展示了如何实现退款的重试机制,适用于网络不稳定的情况
func RefundOrderWithRetry(cli *client.Client, mainOrderNo string, refundAmount string, maxRetries int) (*model.RefundOrderResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		refundOrderNo := utils.GenerateSerialNumber("REF")

		req := &model.RefundOrderRequest{
			MainOrdrNo:   mainOrderNo,
			RefundOrdrNo: refundOrderNo,
			RefundAmt:    refundAmount,
			RefundRsn:    "用户申请退款",
		}

		resp, err := cli.RefundOrder(context.Background(), req)
		if err != nil {
			lastErr = err
			fmt.Printf("第 %d 次退款失败: %v\n", i+1, err)
			continue
		}

		if resp.IsSuccess() {
			return resp, nil
		}

		lastErr = resp.GetError()
		fmt.Printf("第 %d 次退款失败: %v\n", i+1, lastErr)

		// 等待一段时间后重试
		// 实际业务中,应该使用更合理的重试间隔
		// 这里仅作示例,实际使用时应该使用 time.Sleep()
	}

	return nil, fmt.Errorf("退款失败,已重试 %d 次,最后错误: %w", maxRetries, lastErr)
}

// RefundOrderWithCallback 演示退款后的回调处理
//
// 本示例展示了如何在退款成功后执行回调处理
func RefundOrderWithCallback(cli *client.Client, mainOrderNo string, refundAmount string, callback func(*model.RefundOrderResponse)) error {
	refundOrderNo := utils.GenerateSerialNumber("REF")

	req := &model.RefundOrderRequest{
		MainOrdrNo:   mainOrderNo,
		RefundOrdrNo: refundOrderNo,
		RefundAmt:    refundAmount,
		RefundRsn:    "用户申请退款",
	}

	resp, err := cli.RefundOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("订单退款失败: %w", err)
	}

	if resp.IsSuccess() {
		// 执行回调处理
		if callback != nil {
			callback(resp)
		}
	}

	return nil
}
