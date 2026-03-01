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

// QueryOrderExample 演示如何使用查询支付结果接口
//
// 本示例展示了完整的支付结果查询流程,包括:
// 1. 配置初始化
// 2. 客户端创建
// 3. 构建查询请求
// 4. 调用 API
// 5. 处理响应和错误
func QueryOrderExample() {
	// 1. 准备 RSA 密钥对
	// 实际使用时,应该从文件加载密钥
	// privateKey, err := utils.LoadPrivateKey("1111")
	// publicKey, err := utils.LoadPublicKey("11111")
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

	// 4. 创建查询支付结果请求
	// 方式1: 使用主订单编号查询
	req := &model.QueryOrderRequest{
		IttpartyStmId:  "00000",                           // 发起渠道编号,默认送5个0
		PyChnlCd:       "0000000000000000000000000",       // 支付渠道代码,默认送25个0
		IttpartyTms:    utils.CurrentTimestamp(),          // 发起方时间戳,格式: yyyyMMddHHmmssfff
		IttpartyJrnlNo: utils.GenerateSerialNumber("QRY"), // 发起方流水号,不允许重复
		MktId:          "12345678901234",                  // 市场编号
		MainOrdrNo:     "ORD20240101120000123001",         // 主订单编号
		Vno:            "4",                               // 版本号,填写版本为4
	}

	// 方式2: 使用支付流水号查询
	// req := &model.QueryOrderRequest{
	// 	IttpartyStmId:  "00000",
	// 	PyChnlCd:       "0000000000000000000000000",
	// 	IttpartyTms:    utils.CurrentTimestamp(),
	// 	IttpartyJrnlNo: utils.GenerateSerialNumber("QRY"),
	// 	MktId:          "12345678901234",
	// 	PyTrnNo:        "PY20240101120000123",            // 支付流水号
	// 	Vno:            "4",
	// }

	// 5. 调用查询支付结果接口
	// 该方法会自动进行签名、发送请求、验证签名等操作
	resp, err := cli.QueryOrder(context.Background(), req)
	if err != nil {
		log.Fatalf("查询支付结果失败: %v", err)
	}

	// 6. 处理响应
	if resp.IsSuccess() {
		fmt.Println("查询支付结果成功!")
		fmt.Printf("主订单编号: %s\n", resp.MainOrdrNo)
		fmt.Printf("支付流水号: %s\n", resp.PyTrnNo)
		fmt.Printf("支付金额: %.2f 元\n", resp.Txnamt)
		fmt.Printf("订单生成时间: %s\n", resp.OrdrGenTm)

		// 7. 检查订单支付状态
		if resp.IsPaid() {
			fmt.Println("订单已支付成功")
		} else {
			fmt.Printf("订单状态: %s\n", resp.OrdrStcd)
			// 订单状态说明:
			// 1-待支付
			// 2-支付成功
			// 3-支付失败
			// 9-待轮询
		}
	} else {
		fmt.Printf("查询支付结果失败: %v\n", resp.GetError())
	}
}

// QueryOrderWithMinimalParams 演示使用最少参数查询支付结果
//
// 本示例展示了如何使用最少参数查询订单,SDK 会自动填充默认值
func QueryOrderWithMinimalParams(cli *client.Client, mainOrderNo string) error {
	// 使用最少参数,其他参数会自动填充默认值
	req := &model.QueryOrderRequest{
		MainOrdrNo: mainOrderNo, // 必填: 主订单编号
	}

	resp, err := cli.QueryOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询支付结果失败: %w", err)
	}

	if resp.IsPaid() {
		fmt.Printf("订单已支付,金额: %.2f 元\n", resp.Txnamt)
	} else {
		fmt.Printf("订单状态: %s\n", resp.OrdrStcd)
	}

	return nil
}

// QueryOrderWithErrorHandling 演示完整的错误处理
//
// 本示例展示了如何处理各种可能的错误情况
func QueryOrderWithErrorHandling(cli *client.Client, mainOrderNo string) (*model.QueryOrderResponse, error) {
	req := &model.QueryOrderRequest{
		MainOrdrNo: mainOrderNo,
	}

	// 调用 API
	resp, err := cli.QueryOrder(context.Background(), req)
	if err != nil {
		// 处理错误
		// 可能的错误类型:
		// 1. 配置错误: 私钥/公钥无效
		// 2. 网络错误: 连接超时、网络不可达
		// 3. 签名错误: 签名生成失败
		// 4. 验签错误: 响应签名验证失败
		// 5. 业务错误: 订单查询失败
		return nil, fmt.Errorf("查询支付结果失败: %w", err)
	}

	// 检查响应是否成功
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("查询支付结果失败: %w", resp.GetError())
	}

	return resp, nil
}

// QueryOrderWithPaymentFlowNo 演示使用支付流水号查询
//
// 本示例展示了如何使用支付流水号查询订单状态
func QueryOrderWithPaymentFlowNo(cli *client.Client, pyTrnNo string) error {
	req := &model.QueryOrderRequest{
		PyTrnNo: pyTrnNo, // 使用支付流水号查询
	}

	resp, err := cli.QueryOrder(context.Background(), req)
	if err != nil {
		return fmt.Errorf("查询支付结果失败: %w", err)
	}

	if resp.IsPaid() {
		fmt.Printf("订单已支付,流水号: %s, 金额: %.2f 元\n", resp.PyTrnNo, resp.Txnamt)
	} else {
		fmt.Printf("订单状态: %s\n", resp.OrdrStcd)
	}

	return nil
}

// QueryOrderWithRetry 演示带重试机制的订单查询
//
// 本示例展示了如何实现订单查询的重试机制,适用于轮询订单状态
func QueryOrderWithRetry(cli *client.Client, mainOrderNo string, maxRetries int) (*model.QueryOrderResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := QueryOrderWithErrorHandling(cli, mainOrderNo)
		if err != nil {
			lastErr = err
			continue
		}

		// 如果订单已支付,直接返回
		if resp.IsPaid() {
			return resp, nil
		}

		// 如果订单支付失败,直接返回错误
		if resp.OrdrStcd == model.OrderStatusFailed {
			return nil, fmt.Errorf("订单支付失败")
		}

		// 如果订单待支付或待轮询,继续重试
		fmt.Printf("第 %d 次查询,订单状态: %s\n", i+1, resp.OrdrStcd)

		// 等待一段时间后重试
		// 实际业务中,应该使用更合理的重试间隔
		// 这里仅作示例,实际使用时应该使用 time.Sleep()
	}

	return nil, fmt.Errorf("订单查询失败,已重试 %d 次,最后错误: %w", maxRetries, lastErr)
}

// QueryOrderAndProcess 演示查询订单并进行业务处理
//
// 本示例展示了如何根据订单状态进行不同的业务处理
func QueryOrderAndProcess(cli *client.Client, mainOrderNo string) error {
	resp, err := QueryOrderWithErrorHandling(cli, mainOrderNo)
	if err != nil {
		return fmt.Errorf("查询订单失败: %w", err)
	}

	// 根据订单状态进行业务处理
	switch resp.OrdrStcd {
	case model.OrderStatusPending:
		// 待支付状态
		fmt.Println("订单待支付,等待用户支付")
		// 可以在这里实现轮询逻辑
		return nil

	case model.OrderStatusSuccess:
		// 支付成功状态
		fmt.Println("订单支付成功!")
		fmt.Printf("支付流水号: %s\n", resp.PyTrnNo)
		fmt.Printf("支付金额: %.2f 元\n", resp.Txnamt)
		// 可以在这里更新订单状态、发货等业务逻辑
		return nil

	case model.OrderStatusFailed:
		// 支付失败状态
		fmt.Println("订单支付失败")
		// 可以在这里记录失败原因、通知用户等
		return fmt.Errorf("订单支付失败")

	case model.OrderStatusPolling:
		// 待轮询状态
		fmt.Println("订单待轮询,需要继续查询")
		// 可以在这里实现轮询逻辑
		return nil

	default:
		// 未知状态
		return fmt.Errorf("未知订单状态: %s", resp.OrdrStcd)
	}
}

// QueryOrderWithSubOrderInfo 演示查询订单并获取子订单信息
//
// 本示例展示了如何查询订单并解析子订单列表
func QueryOrderWithSubOrderInfo(cli *client.Client, mainOrderNo string) error {
	resp, err := QueryOrderWithErrorHandling(cli, mainOrderNo)
	if err != nil {
		return fmt.Errorf("查询订单失败: %w", err)
	}

	fmt.Printf("主订单编号: %s\n", resp.MainOrdrNo)
	fmt.Printf("支付流水号: %s\n", resp.PyTrnNo)
	fmt.Printf("订单状态: %s\n", resp.OrdrStcd)

	// 解析子订单信息
	if len(resp.Orderlist) > 0 {
		fmt.Printf("子订单数量: %d\n", len(resp.Orderlist))
		for i, subOrder := range resp.Orderlist {
			fmt.Printf("子订单 %d:\n", i+1)
			fmt.Printf("  客户方子订单编号: %s\n", subOrder.CmdtyOrdrNo)
			fmt.Printf("  子订单编号: %s\n", subOrder.SubOrdrId)

			// 解析消费券信息
			if len(subOrder.Cpnlist) > 0 {
				fmt.Printf("  使用消费券数量: %d\n", len(subOrder.Cpnlist))
				for j, coupon := range subOrder.Cpnlist {
					fmt.Printf("    消费券 %d: 订单编号=%s, 金额=%s\n",
						j+1, coupon.CnsmpNoteOrdrId, coupon.Amt)
				}
			}
		}
	}

	return nil
}
