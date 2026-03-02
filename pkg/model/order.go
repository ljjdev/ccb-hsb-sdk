// Package model 定义了与建行对公专业结算综合服务平台交互的数据模型。
//
// 该包包含了请求参数、响应数据、错误信息等结构体定义。
package model

import (
	"encoding/json"
	"fmt"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	// OrderStatusPending 待支付
	OrderStatusPending OrderStatus = "1"
	// OrderStatusSuccess 支付成功
	OrderStatusSuccess OrderStatus = "2"
	// OrderStatusFailed 支付失败
	OrderStatusFailed OrderStatus = "3"
	// OrderStatusPolling 待轮询
	OrderStatusPolling OrderStatus = "9"
)

// PaymentMethod 支付方式代码
type PaymentMethod string

const (
	// PaymentMethodPC PC端收银台
	PaymentMethodPC PaymentMethod = "01"
	// PaymentMethodOffline 线下支付
	PaymentMethodOffline PaymentMethod = "02"
	// PaymentMethodMobileH5 移动端H5页面
	PaymentMethodMobileH5 PaymentMethod = "03"
	// PaymentMethodWechatMini 微信小程序
	PaymentMethodWechatMini PaymentMethod = "05"
	// PaymentMethodOnlineBank 对私网银
	PaymentMethodOnlineBank PaymentMethod = "06"
	// PaymentMethodQRCode 聚合二维码
	PaymentMethodQRCode PaymentMethod = "07"
	// PaymentMethodDragonPay 龙支付
	PaymentMethodDragonPay PaymentMethod = "08"
	// PaymentMethodScan 被扫
	PaymentMethodScan PaymentMethod = "09"
	// PaymentMethodDigitalWallet 数字电子钱包
	PaymentMethodDigitalWallet PaymentMethod = "11"
	// PaymentMethodContactless 无感支付
	PaymentMethodContactless PaymentMethod = "12"
	// PaymentMethodSharedWallet 共享钱包
	PaymentMethodSharedWallet PaymentMethod = "13"
	// PaymentMethodAlipayMini 支付宝小程序
	PaymentMethodAlipayMini PaymentMethod = "14"
	// PaymentMethodSilentPay 免密支付
	PaymentMethodSilentPay PaymentMethod = "15"
)

// OrderType 订单类型
type OrderType string

const (
	// OrderTypeCoupon 消费券购买订单
	OrderTypeCoupon OrderType = "02"
	// OrderTypeTransit 在途订单
	OrderTypeTransit OrderType = "03"
	// OrderTypeNormal 普通订单
	OrderTypeNormal OrderType = "04"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	// IttpartyStmId 发起渠道编号,默认送5个0
	IttpartyStmId string `json:"Ittparty_Stm_Id"`

	// PyChnlCd 支付渠道代码,默认送25个0
	PyChnlCd string `json:"Py_Chnl_Cd"`

	// IttpartyTms 发起方时间戳,格式: yyyyMMddHHmmssfff
	IttpartyTms string `json:"Ittparty_Tms"`

	// IttpartyJrnlNo 发起方流水号,不允许重复
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No"`

	// MktId 市场编号
	MktId string `json:"Mkt_Id"`

	// MainOrdrNo 主订单编号,不允许重复
	MainOrdrNo string `json:"Main_Ordr_No"`

	// PymdCd 支付方式代码
	PymdCd PaymentMethod `json:"Pymd_Cd"`

	// QRCODE 码信息(一维码、二维码)
	QRCODE string `json:"QRCODE,omitempty"`

	// SVCID 建行钱包合约号
	SVCID string `json:"SVCID,omitempty"`

	// APPID 第三方APP平台号
	APPID string `json:"APPID,omitempty"`

	// WALLETNAME 钱包名称
	WALLETNAME string `json:"WALLETNAME,omitempty"`

	// PyOrdrTpcd 订单类型
	PyOrdrTpcd OrderType `json:"Py_Ordr_Tpcd"`

	// PyRsltNtcSn 支付结果通知序号
	PyRsltNtcSn string `json:"Py_Rslt_Ntc_Sn,omitempty"`

	// BnkCd 银行编码
	BnkCd string `json:"Bnk_Cd,omitempty"`

	// OprNo 操作员号
	OprNo string `json:"Opr_No,omitempty"`

	// UsrId 用户ID
	UsrId string `json:"Usr_Id,omitempty"`

	// Ccy 币种,默认156(人民币)
	Ccy string `json:"Ccy"`

	// PgfcRetUrlAdr 页面返回URL地址
	PgfcRetUrlAdr string `json:"Pgfc_Ret_Url_Adr,omitempty"`

	// OrdrTamt 订单总金额
	OrdrTamt float64 `json:"Ordr_Tamt"`

	// TxnTamt 交易总金额
	TxnTamt float64 `json:"Txn_Tamt"`

	// SubAppid 小程序的APPID
	SubAppid string `json:"Sub_Appid,omitempty"`

	// SubOpenid 用户子标识
	SubOpenid string `json:"Sub_Openid,omitempty"`

	// InstallNum 分期期数
	InstallNum string `json:"Install_Num,omitempty"`

	// HdcgBrsId 手续费承担方编号
	HdcgBrsId string `json:"Hdcg_Brs_Id,omitempty"`

	// ClrgDt 确认收货日期
	ClrgDt string `json:"Clrg_Dt,omitempty"`

	// PayDsc 支付描述
	PayDsc string `json:"Pay_Dsc,omitempty"`

	// OrderTimeOut 订单超时时间(秒)
	OrderTimeOut string `json:"Order_Time_Out,omitempty"`

	// PltNo 车牌号
	PltNo string `json:"Plt_No,omitempty"`

	// ApntCnsmpAmt 是否指定消费券核销金额
	ApntCnsmpAmt float64 `json:"Apnt_Cnsmp_Amt,omitempty"`

	// CustomerIdr 消费者唯一标识
	CustomerIdr string `json:"Customer_Idr,omitempty"`

	// Orderlist 子订单列表
	Orderlist []SubOrder `json:"Orderlist,omitempty"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// SubOrder 子订单
type SubOrder struct {
	// MktMrchId 商家编号
	MktMrchId string `json:"Mkt_Mrch_Id"`

	// UdfId 商家自定义编号
	UdfId string `json:"Udf_Id,omitempty"`

	// CmdtyOrdrNo 客户方子订单编号
	CmdtyOrdrNo string `json:"Cmdty_Ordr_No"`

	// OrdrAmt 订单金额
	OrdrAmt float64 `json:"Ordr_Amt"`

	// Txnamt 交易金额
	Txnamt float64 `json:"Txnamt"`

	// ApdTamt 附加项总金额
	ApdTamt float64 `json:"Apd_Tamt,omitempty"`

	// CmdtyDsc 商品描述
	CmdtyDsc string `json:"Cmdty_Dsc,omitempty"`

	// CmdtyTp 商品类型
	CmdtyTp string `json:"Cmdty_Tp,omitempty"`

	// ClrgRuleId 分账规则编号
	ClrgRuleId string `json:"Clrg_Rule_Id,omitempty"`

	// Parlist 分账方列表
	Parlist []Participant `json:"Parlist,omitempty"`

	// Cpnlist 消费券列表
	Cpnlist []Coupon `json:"Cpnlist,omitempty"`
}

// Participant 分账方
type Participant struct {
	// SeqNo 顺序号
	SeqNo int `json:"Seq_No"`

	// MktMrchId 商家编号
	MktMrchId string `json:"Mkt_Mrch_Id"`

	// Amt 金额
	Amt float64 `json:"Amt,omitempty"`
}

// Coupon 消费券
type Coupon struct {
	// CnsmpNoteOrdrId 消费券订单编号
	CnsmpNoteOrdrId string `json:"Cnsmp_Note_Ordr_Id"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	// IttpartyTms 发起方时间戳
	IttpartyTms string `json:"Ittparty_Tms"`

	// IttpartyJrnlNo 发起方流水号
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No"`

	// MainOrdrNo 主订单编号
	MainOrdrNo string `json:"Main_Ordr_No"`

	// PyTrnNo 支付流水号
	PyTrnNo string `json:"Py_Trn_No"`

	// PrimOrdrNo 订单编号
	PrimOrdrNo string `json:"Prim_Ordr_No"`

	// OrdrGenTm 订单生成时间
	OrdrGenTm string `json:"Ordr_Gen_Tm"`

	// OrdrOvtmTm 订单超时时间
	OrdrOvtmTm string `json:"Ordr_Ovtm_Tm"`

	// CshdkUrl 收银台URL
	CshdkUrl string `json:"Cshdk_Url"`

	// PayUrl 支付URL
	PayUrl string `json:"Pay_Url,omitempty"`

	// PayQrCode 支付二维码串
	PayQrCode string `json:"Pay_Qr_Code,omitempty"`

	// RtnParData 返回参数数据
	RtnParData string `json:"Rtn_Par_Data,omitempty"`

	// WaitTime 等待时间(秒)
	WaitTime string `json:"Wait_Time,omitempty"`

	// OrdrStcd 订单状态代码
	OrdrStcd OrderStatus `json:"Ordr_Stcd"`

	// MatchPayerAcctRslt 付款人账号匹配结果
	MatchPayerAcctRslt string `json:"Match_Payer_Acct_Rslt,omitempty"`

	// MatchPayerNameRslt 付款人户名匹配结果
	MatchPayerNameRslt string `json:"Match_Payer_Name_Rslt,omitempty"`

	// Orderlist 子订单列表
	Orderlist []SubOrderResponse `json:"Orderlist,omitempty"`

	// SvcRspSt 服务响应状态
	SvcRspSt string `json:"Svc_Rsp_St"`

	// SvcRspCd 服务响应码
	SvcRspCd string `json:"Svc_Rsp_Cd,omitempty"`

	// RspInf 响应信息
	RspInf string `json:"Rsp_Inf,omitempty"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// SubOrderResponse 子订单响应
type SubOrderResponse struct {
	// CmdtyOrdrNo 客户方子订单编号
	CmdtyOrdrNo string `json:"Cmdty_Ordr_No"`

	// SubOrdrId 子订单编号
	SubOrdrId string `json:"Sub_Ordr_Id"`

	// Cpnlist 使用消费券列表
	Cpnlist []UsedCoupon `json:"Cpnlist,omitempty"`
}

// UsedCoupon 使用消费券
type UsedCoupon struct {
	// CnsmpNoteOrdrId 消费券订单编号
	CnsmpNoteOrdrId string `json:"Cnsmp_Note_Ordr_Id"`

	// Amt 金额
	Amt float64 `json:"Amt"`

	// BalAmt 余额
	BalAmt float64 `json:"Bal_Amt,omitempty"`
}

// ToMap 将请求转换为 map,用于签名
func (r *CreateOrderRequest) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}

// ToMap 将响应转换为 map,用于签名验证
func (r *CreateOrderResponse) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}

// IsSuccess 判断响应是否成功
func (r *CreateOrderResponse) IsSuccess() bool {
	return r.SvcRspSt == "00"
}

// GetError 获取错误信息
func (r *CreateOrderResponse) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return fmt.Errorf("service error: code=%s, message=%s", r.SvcRspCd, r.RspInf)
}

// RefundStatus 退款响应状态
type RefundStatus string

const (
	// RefundStatusSuccess 退款成功
	RefundStatusSuccess RefundStatus = "00"
	// RefundStatusFailed 退款失败
	RefundStatusFailed RefundStatus = "01"
	// RefundStatusDelayed 退款延迟等待
	RefundStatusDelayed RefundStatus = "02"
	// RefundStatusUncertain 退款结果不确定
	RefundStatusUncertain RefundStatus = "03"
	// RefundStatusWaiting 等待确认(线下订单类型返回)
	RefundStatusWaiting RefundStatus = "04"
	// RefundStatusNotFound 没有查询到符合条件的记录
	RefundStatusNotFound RefundStatus = "05"
	// RefundStatusAccepted 已受理(仅异步退款有此状态)
	RefundStatusAccepted RefundStatus = "0a"
	// RefundStatusInterrupted 中断(仅异步退款有此状态)
	RefundStatusInterrupted RefundStatus = "0b"
)

// QueryRefundRequest 查询退款结果请求
type QueryRefundRequest struct {
	// IttpartyStmId 发起渠道编号,默认送5个0
	IttpartyStmId string `json:"Ittparty_Stm_Id"`

	// PyChnlCd 支付渠道代码,默认送25个0
	PyChnlCd string `json:"Py_Chnl_Cd"`

	// IttpartyTms 发起方时间戳,格式: yyyyMMddHHmmssfff
	IttpartyTms string `json:"Ittparty_Tms"`

	// IttpartyJrnlNo 发起方流水号,不允许重复
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No"`

	// MktId 市场编号
	MktId string `json:"Mkt_Id"`

	// CustRfndTrcno 客户方退款流水号
	// 客户方退款流水号与退款流水号必输其一
	CustRfndTrcno string `json:"Cust_Rfnd_Trcno,omitempty"`

	// RfndTrcno 退款流水号
	// 客户方退款流水号与退款流水号必输其一
	RfndTrcno string `json:"Rfnd_Trcno,omitempty"`

	// Vno 版本号,填写版本为4
	Vno string `json:"Vno"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// QueryRefundResponse 查询退款结果响应
type QueryRefundResponse struct {
	// IttpartyTms 发起方时间戳
	IttpartyTms string `json:"Ittparty_Tms"`

	// IttpartyJrnlNo 发起方流水号
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No"`

	// CustRfndTrcno 客户方退款流水号
	CustRfndTrcno string `json:"Cust_Rfnd_Trcno,omitempty"`

	// RfndTrcno 退款流水号
	RfndTrcno string `json:"Rfnd_Trcno"`

	// RfndAmt 退款金额
	RfndAmt *float64 `json:"Rfnd_Amt,omitempty"`

	// RefundRspSt 退款响应状态
	RefundRspSt RefundStatus `json:"Refund_Rsp_St"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// ToMap 将请求转换为 map,用于签名
func (r *QueryRefundRequest) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}

// IsSuccess 判断响应是否成功
func (r *QueryRefundResponse) IsSuccess() bool {
	return r.RefundRspSt == RefundStatusSuccess
}

// GetError 获取错误信息
func (r *QueryRefundResponse) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return fmt.Errorf("refund query failed: status=%s", r.RefundRspSt)
}

// QueryOrderRequest 查询支付结果请求
type QueryOrderRequest struct {
	// IttpartyStmId 发起渠道编号,默认送5个0
	IttpartyStmId string `json:"Ittparty_Stm_Id"`

	// PyChnlCd 支付渠道代码,默认送25个0
	PyChnlCd string `json:"Py_Chnl_Cd"`

	// IttpartyTms 发起方时间戳,格式: yyyyMMddHHmmssfff
	IttpartyTms string `json:"Ittparty_Tms"`

	// IttpartyJrnlNo 发起方流水号,不允许重复
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No"`

	// MktId 市场编号
	MktId string `json:"Mkt_Id"`

	// MainOrdrNo 主订单编号
	// 主订单号与支付流水号必输其一
	MainOrdrNo string `json:"Main_Ordr_No,omitempty"`

	// PyTrnNo 支付流水号
	PyTrnNo string `json:"Py_Trn_No,omitempty"`

	// Vno 版本号,填写版本为4
	Vno string `json:"Vno"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// QueryOrderResponse 查询支付结果响应
type QueryOrderResponse struct {
	// MainOrdrNo 主订单编号
	MainOrdrNo string `json:"Main_Ordr_No,omitempty"`

	// PyTrnNo 支付流水号
	PyTrnNo string `json:"Py_Trn_No"`

	// Txnamt 交易金额
	Txnamt float64 `json:"Txnamt,omitempty"`

	// OrdrGenTm 订单生成时间
	OrdrGenTm string `json:"Ordr_Gen_Tm,omitempty"`

	// OrdrOvtmTm 订单超时时间
	OrdrOvtmTm string `json:"Ordr_Ovtm_Tm,omitempty"`

	// OrdrStcd 订单状态代码
	OrdrStcd OrderStatus `json:"Ordr_Stcd"`

	// MatchPayerAcctRslt 付款人账号匹配结果
	MatchPayerAcctRslt string `json:"Match_Payer_Acct_Rslt,omitempty"`

	// MatchPayerNameRslt 付款人户名匹配结果
	MatchPayerNameRslt string `json:"Match_Payer_Name_Rslt,omitempty"`

	// PrimOrdrNo 订单编号
	PrimOrdrNo string `json:"Prim_Ordr_No,omitempty"`

	// Orderlist 子订单列表
	Orderlist []SubOrderResponse `json:"Orderlist,omitempty"`

	// SvcRspSt 服务响应状态
	SvcRspSt string `json:"Svc_Rsp_St"`

	// SvcRspCd 服务响应码
	SvcRspCd string `json:"Svc_Rsp_Cd,omitempty"`

	// RspInf 响应信息
	RspInf string `json:"Rsp_Inf,omitempty"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// ToMap 将请求转换为 map,用于签名
func (r *QueryOrderRequest) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}

// IsSuccess 判断响应是否成功
func (r *QueryOrderResponse) IsSuccess() bool {
	return r.SvcRspSt == "00"
}

// IsPaid 判断订单是否已支付
func (r *QueryOrderResponse) IsPaid() bool {
	return r.OrdrStcd == OrderStatusSuccess
}

// GetError 获取错误信息
func (r *QueryOrderResponse) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return fmt.Errorf("order query failed: code=%s, message=%s", r.SvcRspCd, r.RspInf)
}

// ToMap 将响应转换为 map,用于签名验证
func (r *QueryOrderResponse) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}

// RefundOrderRequest 退款订单请求
type RefundOrderRequest struct {
	// MktId 市场编号
	MktId string `json:"Mkt_Id"`

	// MainOrdrNo 主订单编号
	MainOrdrNo string `json:"Main_Ordr_No"`

	// RefundOrdrNo 退款订单编号,不允许重复
	RefundOrdrNo string `json:"Refund_Ordr_No"`

	// RefundAmt 退款金额
	RefundAmt float64 `json:"Refund_Amt"`

	// RefundRsn 退款原因
	RefundRsn string `json:"Refund_Rsn,omitempty"`

	// IttpartyJrnlNo 发起方流水号,不允许重复
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No,omitempty"`

	// IttpartyTms 发起方时间戳,格式: yyyyMMddHHmmssfff
	IttpartyTms string `json:"Ittparty_Tms,omitempty"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// RefundOrderResponse 退款订单响应
type RefundOrderResponse struct {
	// IttpartyTms 发起方时间戳
	IttpartyTms string `json:"Ittparty_Tms"`

	// IttpartyJrnlNo 发起方流水号
	IttpartyJrnlNo string `json:"Ittparty_Jrnl_No"`

	// MainOrdrNo 主订单编号
	MainOrdrNo string `json:"Main_Ordr_No"`

	// RefundOrdrNo 退款订单编号
	RefundOrdrNo string `json:"Refund_Ordr_No"`

	// RefundAmt 退款金额
	RefundAmt float64 `json:"Refund_Amt"`

	// RefundTrnNo 退款流水号
	RefundTrnNo string `json:"Refund_Trn_No,omitempty"`

	// RefundTm 退款时间
	RefundTm string `json:"Refund_Tm,omitempty"`

	// RefundStcd 退款状态代码
	RefundStcd string `json:"Refund_Stcd,omitempty"`

	// RefundFundsSource 退款资金来源
	RefundFundsSource string `json:"Refund_Funds_Source,omitempty"`

	// RefundRspInf 退款响应信息
	RefundRspInf string `json:"Refund_Rsp_Inf,omitempty"`

	// SvcRspSt 服务响应状态
	SvcRspSt string `json:"Svc_Rsp_St"`

	// SvcRspCd 服务响应码
	SvcRspCd string `json:"Svc_Rsp_Cd,omitempty"`

	// RspInf 响应信息
	RspInf string `json:"Rsp_Inf,omitempty"`

	// SignInf 签名信息
	SignInf string `json:"Sign_Inf,omitempty"`
}

// ToMap 将退款请求转换为 map,用于签名
func (r *RefundOrderRequest) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}

// IsSuccess 判断响应是否成功
func (r *RefundOrderResponse) IsSuccess() bool {
	return r.SvcRspSt == "00"
}

// GetError 获取错误信息
func (r *RefundOrderResponse) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return fmt.Errorf("service error: code=%s, message=%s", r.SvcRspCd, r.RspInf)
}

// ToMap 将退款响应转换为 map,用于签名验证
func (r *RefundOrderResponse) ToMap() (map[string]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 转换为 map[string]string
	stringMap := make(map[string]string)
	for k, v := range result {
		if v != nil {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	return stringMap, nil
}
