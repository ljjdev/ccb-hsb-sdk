// Package examples 提供了建行对公专业结算综合服务平台 SDK 的使用示例。
//
// 本包包含了各种常见场景的使用示例,帮助开发者快速上手。
package examples

import (
	"context"
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
	// 实际使用时,加载私钥的base64字符串
	privateKey, err := utils.LoadPrivateKey("1111")
	// 实际使用时,加载公钥的base64字符串
	publicKey, err := utils.LoadPublicKey("11111")

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
			fmt.Printf("退款金额: %.2f 元\n", *resp.RfndAmt)
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
