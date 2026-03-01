# 更新响应体结构 Spec

## Why
根据建行提供的实际响应示例，需要更新 SDK 中的响应体结构，确保字段与实际 API 返回的数据完全匹配，避免因字段缺失导致的解析错误。

## What Changes
- **BREAKING**: 为 `QueryOrderResponse` 添加 `Txnamt` 字段（交易金额）
- **BREAKING**: 为 `RefundOrderResponse` 添加 `RefundFundsSource` 字段（退款资金来源）
- **BREAKING**: 为 `RefundOrderResponse` 添加 `RefundRspInf` 字段（退款响应信息）

## Impact
- Affected specs: 查询支付结果接口、订单退款接口
- Affected code: `pkg/model/order.go`

## ADDED Requirements

### Requirement: QueryOrderResponse 添加交易金额字段
The system SHALL provide `Txnamt` field in QueryOrderResponse to represent transaction amount.

#### Scenario: Success case
- **WHEN** querying payment result and API returns transaction amount
- **THEN** the response should correctly parse `Txnamt` field from JSON

### Requirement: RefundOrderResponse 添加退款资金来源字段
The system SHALL provide `RefundFundsSource` field in RefundOrderResponse to represent refund funds source.

#### Scenario: Success case
- **WHEN** processing refund and API returns refund funds source
- **THEN** the response should correctly parse `Refund_Funds_Source` field from JSON

### Requirement: RefundOrderResponse 添加退款响应信息字段
The system SHALL provide `RefundRspInf` field in RefundOrderResponse to represent refund response message.

#### Scenario: Success case
- **WHEN** processing refund and API returns refund response message
- **THEN** the response should correctly parse `Refund_Rsp_Inf` field from JSON

## MODIFIED Requirements

### Requirement: QueryOrderResponse 结构更新
Update QueryOrderResponse to include transaction amount field for accurate payment result tracking.

### Requirement: RefundOrderResponse 结构更新
Update RefundOrderResponse to include refund funds source and response message fields for complete refund information.

## REMOVED Requirements
None
