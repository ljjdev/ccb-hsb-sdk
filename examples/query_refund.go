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

// QueryRefundExample 演示如何使用查询退款结果接口
//
// 本示例展示了完整的退款结果查询流程,包括:
// 1. 配置初始化
// 2. 客户端创建
// 3. 构建查询请求
// 4. 调用 API
// 5. 处理响应和错误
func QueryRefundExample() {
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

	// 4. 创建查询退款结果请求
	// 方式1: 使用客户方退款流水号查询
	req := &model.QueryRefundRequest{
		IttpartyStmId:  "00000",                           // 发起渠道编号,默认送5个0
		PyChnlCd:       "0000000000000000000000000",       // 支付渠道代码,默认送25个0
		IttpartyTms:    utils.CurrentTimestamp(),          // 发起方时间戳,格式: yyyyMMddHHmmssfff
		IttpartyJrnlNo: utils.GenerateSerialNumber("QRY"), // 发起方流水号,不允许重复
		MktId:          "12345678901234",                  // 市场编号
		CustRfndTrcno:  "REFUND20240101120000123",         // 客户方退款流水号
		Vno:            "4",                               // 版本号,填写版本为4
	}

	// 方式2: 使用退款流水号查询
	// req := &model.QueryRefundRequest{
	// 	IttpartyStmId:  "00000",
	// 	PyChnlCd:       "0000000000000000000000000",
	// 	IttpartyTms:    utils.CurrentTimestamp(),
	// 	IttpartyJrnlNo: utils.GenerateSerialNumber("QRY"),
	// 	MktId:          "12345678901234",
	// 	RfndTrcno:      "RFND20240101120000123",          // 退款流水号
	// 	Vno:            "4",
	// }

	// 5. 调用查询退款结果接口
	// 该方法会自动进行签名、发送请求、验证签名等操作
	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		log.Fatalf("查询退款结果失败: %v", err)
	}

	// 6. 处理响应
	if resp.IsSuccess() {
		fmt.Println("查询退款结果成功!")
		fmt.Printf("客户方退款流水号: %s\n", resp.CustRfndTrcno)
		fmt.Printf("退款流水号: %s\n", resp.RfndTrcno)
		if resp.RfndAmt != nil {
			fmt.Printf("退款金额: %s 元\n", *resp.RfndAmt)
		}
		fmt.Printf("退款状态: %s\n", resp.RefundRspSt)

		// 7. 根据退款状态进行业务处理
		switch resp.RefundRspSt {
		case model.RefundStatusSuccess:
			fmt.Println("退款成功")
		case model.RefundStatusFailed:
			fmt.Println("退款失败")
		case model.RefundStatusDelayed:
			fmt.Println("退款延迟等待,请稍后查询")
		case model.RefundStatusUncertain:
			fmt.Println("退款结果不确定,请稍后查询")
		case model.RefundStatusWaiting:
			fmt.Println("等待确认")
		case model.RefundStatusNotFound:
			fmt.Println("没有查询到符合条件的记录")
		case model.RefundStatusAccepted:
			fmt.Println("退款已受理(异步退款)")
		case model.RefundStatusInterrupted:
			fmt.Println("退款中断(异步退款)")
		}
	} else {
		fmt.Printf("查询退款结果失败: %v\n", resp.GetError())
		fmt.Printf("退款状态: %s\n", resp.RefundRspSt)
	}
}

// QueryRefundWithMinimalParams 演示使用最少参数查询退款结果
//
// 本示例展示了如何使用最少参数查询退款,SDK 会自动填充默认值
func QueryRefundWithMinimalParams(cli *client.Client, custRfndTrcno string) error {
	// 使用最少参数,其他参数会自动填充默认值
	req := &model.QueryRefundRequest{
		CustRfndTrcno: custRfndTrcno, // 必填: 客户方退款流水号
	}

	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	if resp.IsSuccess() {
		if resp.RfndAmt != nil {
			fmt.Printf("退款成功,金额: %s 元\n", *resp.RfndAmt)
		} else {
			fmt.Println("退款成功")
		}
	}

	return nil
}

// QueryRefundWithErrorHandling 演示完整的错误处理
//
// 本示例展示了如何处理各种可能的错误情况
func QueryRefundWithErrorHandling(cli *client.Client, custRfndTrcno string) (*model.QueryRefundResponse, error) {
	req := &model.QueryRefundRequest{
		CustRfndTrcno: custRfndTrcno,
	}

	// 调用 API
	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		// 处理错误
		// 可能的错误类型:
		// 1. 配置错误: 私钥/公钥无效
		// 2. 网络错误: 连接超时、网络不可达
		// 3. 签名错误: 签名生成失败
		// 4. 验签错误: 响应签名验证失败
		// 5. 业务错误: 查询失败
		return nil, fmt.Errorf("查询退款结果失败: %w", err)
	}

	// 检查响应是否成功
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("查询退款结果失败: %w", resp.GetError())
	}

	return resp, nil
}

// QueryRefundWithRefundFlowNo 演示使用退款流水号查询
//
// 本示例展示了如何使用退款流水号查询退款结果
func QueryRefundWithRefundFlowNo(cli *client.Client, rfndTrcno string) error {
	req := &model.QueryRefundRequest{
		RfndTrcno: rfndTrcno, // 使用退款流水号查询
	}

	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	if resp.IsSuccess() {
		if resp.RfndAmt != nil {
			fmt.Printf("退款成功,流水号: %s, 金额: %s 元\n", resp.RfndTrcno, *resp.RfndAmt)
		} else {
			fmt.Printf("退款成功,流水号: %s\n", resp.RfndTrcno)
		}
	}

	return nil
}

// QueryRefundWithRetry 演示带重试机制的退款查询
//
// 本示例展示了如何实现退款查询的重试机制,适用于轮询退款状态
func QueryRefundWithRetry(cli *client.Client, custRfndTrcno string, maxRetries int) (*model.QueryRefundResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := QueryRefundWithErrorHandling(cli, custRfndTrcno)
		if err != nil {
			lastErr = err
			continue
		}

		// 如果退款成功,直接返回
		if resp.IsSuccess() {
			return resp, nil
		}

		// 如果退款失败,直接返回错误
		if resp.RefundRspSt == model.RefundStatusFailed {
			return nil, fmt.Errorf("退款失败")
		}

		// 如果退款延迟等待或结果不确定,继续重试
		fmt.Printf("第 %d 次查询,退款状态: %s\n", i+1, resp.RefundRspSt)

		// 等待一段时间后重试
		// 实际业务中,应该使用更合理的重试间隔
		// 这里仅作示例,实际使用时应该使用 time.Sleep()
	}

	return nil, fmt.Errorf("退款查询失败,已重试 %d 次,最后错误: %w", maxRetries, lastErr)
}

// QueryRefundAndProcess 演示查询退款并进行业务处理
//
// 本示例展示了如何根据退款状态进行不同的业务处理
func QueryRefundAndProcess(cli *client.Client, custRfndTrcno string) error {
	resp, err := QueryRefundWithErrorHandling(cli, custRfndTrcno)
	if err != nil {
		return fmt.Errorf("查询退款失败: %w", err)
	}

	// 根据退款状态进行业务处理
	switch resp.RefundRspSt {
	case model.RefundStatusSuccess:
		// 退款成功状态
		fmt.Println("退款成功!")
		fmt.Printf("退款流水号: %s\n", resp.RfndTrcno)
		if resp.RfndAmt != nil {
			fmt.Printf("退款金额: %s 元\n", *resp.RfndAmt)
		}
		// 可以在这里更新订单状态、通知用户等业务逻辑
		return nil

	case model.RefundStatusFailed:
		// 退款失败状态
		fmt.Println("退款失败")
		// 可以在这里记录失败原因、通知用户等
		return fmt.Errorf("退款失败")

	case model.RefundStatusDelayed:
		// 退款延迟等待状态
		fmt.Println("退款延迟等待,请稍后查询")
		// 可以在这里实现轮询逻辑
		return nil

	case model.RefundStatusUncertain:
		// 退款结果不确定状态
		fmt.Println("退款结果不确定,请稍后查询")
		// 可以在这里实现轮询逻辑
		return nil

	case model.RefundStatusWaiting:
		// 等待确认状态(线下订单类型返回)
		fmt.Println("等待确认")
		return nil

	case model.RefundStatusNotFound:
		// 没有查询到符合条件的记录
		fmt.Println("没有查询到符合条件的记录")
		return fmt.Errorf("未找到退款记录")

	case model.RefundStatusAccepted:
		// 已受理状态(仅异步退款有此状态)
		fmt.Println("退款已受理(异步退款)")
		return nil

	case model.RefundStatusInterrupted:
		// 中断状态(仅异步退款有此状态)
		fmt.Println("退款中断(异步退款)")
		return fmt.Errorf("退款中断")

	default:
		// 未知状态
		return fmt.Errorf("未知退款状态: %s", resp.RefundRspSt)
	}
}

// QueryRefundWithValidation 演示查询前的参数验证
//
// 本示例展示了如何在查询前进行参数验证
func QueryRefundWithValidation(cli *client.Client, custRfndTrcno string) error {
	// 1. 参数验证
	if custRfndTrcno == "" {
		return fmt.Errorf("客户方退款流水号不能为空")
	}

	// 2. 发起查询
	req := &model.QueryRefundRequest{
		CustRfndTrcno: custRfndTrcno,
	}

	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	if resp.IsSuccess() {
		if resp.RfndAmt != nil {
			fmt.Printf("退款成功,金额: %s 元\n", *resp.RfndAmt)
		} else {
			fmt.Println("退款成功")
		}
	}

	return nil
}

// QueryRefundWithCallback 演示查询后的回调处理
//
// 本示例展示了如何在查询成功后执行回调处理
func QueryRefundWithCallback(cli *client.Client, custRfndTrcno string, callback func(*model.QueryRefundResponse)) error {
	req := &model.QueryRefundRequest{
		CustRfndTrcno: custRfndTrcno,
	}

	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	if resp.IsSuccess() {
		// 执行回调处理
		if callback != nil {
			callback(resp)
		}
	}

	return nil
}

// QueryRefundWithBothParams 演示同时使用两个参数查询
//
// 本示例展示了如何同时使用客户方退款流水号和退款流水号查询
// 注意: 实际查询时,两个参数必输其一,同时提供时优先使用客户方退款流水号
func QueryRefundWithBothParams(cli *client.Client, custRfndTrcno string, rfndTrcno string) error {
	req := &model.QueryRefundRequest{
		CustRfndTrcno: custRfndTrcno, // 客户方退款流水号
		RfndTrcno:     rfndTrcno,     // 退款流水号
	}

	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	if resp.IsSuccess() {
		if resp.RfndAmt != nil {
			fmt.Printf("退款成功,流水号: %s, 金额: %s 元\n", resp.RfndTrcno, *resp.RfndAmt)
		} else {
			fmt.Printf("退款成功,流水号: %s\n", resp.RfndTrcno)
		}
	}

	return nil
}

// QueryRefundBatch 演示批量查询退款结果
//
// 本示例展示了如何批量查询多个退款的结果
func QueryRefundBatch(cli *client.Client, custRfndTrcnoList []string) error {
	for i, custRfndTrcno := range custRfndTrcnoList {
		fmt.Printf("查询第 %d 个退款: %s\n", i+1, custRfndTrcno)

		req := &model.QueryRefundRequest{
			CustRfndTrcno: custRfndTrcno,
		}

		resp, err := cli.QueryRefund(context.Background(), req)
		if err != nil {
			fmt.Printf("查询失败: %v\n", err)
			continue
		}

		if resp.IsSuccess() {
			if resp.RfndAmt != nil {
				fmt.Printf("退款成功,金额: %s 元, 状态: %s\n", *resp.RfndAmt, resp.RefundRspSt)
			} else {
				fmt.Printf("退款成功, 状态: %s\n", resp.RefundRspSt)
			}
		} else {
			fmt.Printf("退款失败,状态: %s\n", resp.RefundRspSt)
		}
	}

	return nil
}

// QueryRefundWithDetailedInfo 演示查询并获取详细信息
//
// 本示例展示了如何查询退款并获取详细信息
func QueryRefundWithDetailedInfo(cli *client.Client, custRfndTrcno string) error {
	req := &model.QueryRefundRequest{
		CustRfndTrcno: custRfndTrcno,
	}

	resp, err := cli.QueryRefund(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询退款结果失败: %w", err)
	}

	fmt.Printf("客户方退款流水号: %s\n", resp.CustRfndTrcno)
	fmt.Printf("退款流水号: %s\n", resp.RfndTrcno)
	if resp.RfndAmt != nil {
		fmt.Printf("退款金额: %s 元\n", *resp.RfndAmt)
	}
	fmt.Printf("退款状态: %s\n", resp.RefundRspSt)
	fmt.Printf("发起方时间戳: %s\n", resp.IttpartyTms)
	fmt.Printf("发起方流水号: %s\n", resp.IttpartyJrnlNo)

	return nil
}
