# 修正响应体字段类型 Spec

## Why
对照建行响应示例检查发现，部分字段类型与实际 API 返回的数据类型不匹配，可能导致 JSON 解析失败或数据类型错误。

## What Changes
- **BREAKING**: 修改 `QueryOrderResponse` 的 `Txnamt` 字段类型从 `string` 改为 `float64`
- **BREAKING**: 修改 `QueryRefundResponse` 的 `RfndAmt` 字段类型从 `string` 改为 `*string`（支持 null 值）

## Impact
- Affected specs: 查询支付结果接口、查询退款结果接口
- Affected code: `pkg/model/order.go`

## ADDED Requirements

### Requirement: QueryOrderResponse Txnamt 字段类型修正
The system SHALL use `float64` type for `Txnamt` field in QueryOrderResponse to match API response.

#### Scenario: Success case
- **WHEN** querying payment result and API returns numeric transaction amount
- **THEN** the response should correctly parse `Txnamt` as float64 (e.g., 20, 20.55)

### Requirement: QueryRefundResponse RfndAmt 字段类型修正
The system SHALL use `*string` type for `RfndAmt` field in QueryRefundResponse to support null values.

#### Scenario: Success case with null amount
- **WHEN** querying refund result and API returns null for refund amount
- **THEN** the response should correctly parse `Rfnd_Amt` as null (nil pointer)

#### Scenario: Success case with amount
- **WHEN** querying refund result and API returns refund amount
- **THEN** the response should correctly parse `Rfnd_Amt` as string

## MODIFIED Requirements

### Requirement: QueryOrderResponse 字段类型更新
Update QueryOrderResponse to use correct data types matching API responses.

### Requirement: QueryRefundResponse 字段类型更新
Update QueryRefundResponse to support null values for optional fields.

## REMOVED Requirements
None
