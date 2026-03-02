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

// PlaceOrderExample 演示如何使用支付订单生成接口
//
// 本示例展示了完整的支付订单创建流程,包括:
// 1. 配置初始化
// 2. 客户端创建
// 3. 构建订单请求
// 4. 调用 API
// 5. 处理响应和错误
func PlaceOrderExample() {
	// 1. 准备 RSA 密钥对
	// 实际使用时,加载私钥的base64字符串
	privateKey, err := utils.LoadPrivateKey("1111")
	// 实际使用时,加载公钥的base64字符串
	publicKey, err := utils.LoadPublicKey("11111")

	// 2. 创建配置
	// 使用配置选项模式初始化 SDK 配置
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

	// 4. 生成订单相关编号
	// 主订单编号: 不允许重复,建议使用 UUID 或时间戳+随机数
	mainOrderNo := utils.GenerateSerialNumber("ORD")
	// 发起方流水号: 不允许重复,用于追踪请求
	ittpartyJrnlNo := utils.GenerateSerialNumber("JNL")
	// 商品订单号: 不允许重复,用于标识具体商品
	subOrderNo := mainOrderNo + "01"

	// 5. 构建支付订单请求
	// 支付金额: 100.01 元
	amount := 100.01

	req := &model.CreateOrderRequest{
		// 基础信息
		IttpartyStmId:  "00000",                     // 发起渠道编号,默认送5个0
		PyChnlCd:       "0000000000000000000000000", // 支付渠道代码,默认送25个0
		IttpartyTms:    utils.CurrentTimestamp(),    // 发起方时间戳,格式: yyyyMMddHHmmssfff
		IttpartyJrnlNo: ittpartyJrnlNo,              // 发起方流水号,不允许重复
		MktId:          "12345678901234",            // 市场编号

		// 订单信息
		MainOrdrNo:   mainOrderNo,                 // 主订单编号,不允许重复
		PymdCd:       model.PaymentMethodMobileH5, // 支付方式: 03-移动端H5页面
		PyOrdrTpcd:   model.OrderTypeNormal,       // 订单类型: 04-普通订单
		Ccy:          "156",                       // 币种: 156-人民币
		OrdrTamt:     amount,                      // 订单总金额
		TxnTamt:      amount,                      // 交易总金额
		PayDsc:       "商品购买",                      // 支付描述
		OrderTimeOut: "1800",                      // 订单超时时间(秒): 30分钟

		// 子订单信息
		Orderlist: []model.SubOrder{
			{
				MktMrchId:   "12345678901234567890", // 商家编号
				CmdtyOrdrNo: subOrderNo,             // 客户方子订单编号
				OrdrAmt:     amount,                 // 订单金额
				Txnamt:      amount,                 // 交易金额
				CmdtyDsc:    "商品",                   // 商品描述
				ClrgRuleId:  "123456",               // 分账规则编号

				// 分账方列表
				Parlist: []model.Participant{
					{
						SeqNo:     1,                      // 顺序号
						MktMrchId: "12345678901234567890", // 商家编号
					},
				},
			},
		},
	}

	// 6. 调用支付订单生成接口
	// 该方法会自动进行签名、发送请求、验证签名等操作
	payURL, err := cli.PlaceOrder(context.Background(), req)
	if err != nil {
		log.Fatalf("创建支付订单失败: %v", err)
	}

	// 7. 处理响应
	fmt.Println("创建支付订单成功!")
	fmt.Printf("主订单编号: %s\n", mainOrderNo)
	fmt.Printf("发起方流水号: %s\n", ittpartyJrnlNo)
	fmt.Printf("支付金额: %.2f 元\n", amount)
	fmt.Printf("支付URL: %s\n", payURL)

	// 8. 将支付URL返回给前端,引导用户完成支付
	// 实际业务中,应该将 payURL 返回给前端,前端跳转到该URL进行支付
}
