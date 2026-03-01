// Package signature 提供了 JSON 签名字符串拼接功能。
//
// 该包实现了建行对公专业结算综合服务平台所需的 JSON 签名字符串拼接算法,
// 支持递归处理嵌套的 JSON 对象和数组。
package signature

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// BuildSignatureStringFromJSON 从 JSON 字符串构建待签名字符串
//
// 该函数按照建行规范将 JSON 数据拼接成待签名字符串。
//
// 拼接规则:
// 1. 递归处理 JSON 对象和数组
// 2. 将所有参数按字典序排序
// 3. 使用 key=value&key=value 的格式拼接
// 4. 忽略值为空的参数
// 5. 忽略以下字段: SIGN_INF, Svc_Rsp_St, Svc_Rsp_Cd, Rsp_Inf
// 6. 数组类型的值直接拼接所有元素的签名字符串
// 7. 嵌套对象类型的值直接拼接其签名字符串（不包含 key=）
//
// 该方法参考 Java 的 SplicingUtil.createSign 实现
func BuildSignatureStringFromJSON(jsonStr string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal json: %w", err)
	}

	result := splicingSign(data)
	if len(result) > 0 && result[len(result)-1] == '&' {
		result = result[:len(result)-1]
	}

	return result, nil
}

// splicingSign 递归拼接签名字符串
func splicingSign(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		return splicingObject(v)
	case []interface{}:
		return splicingArray(v)
	default:
		return ""
	}
}

// splicingObject 拼接 JSON 对象
func splicingObject(obj map[string]interface{}) string {
	// 获取所有键并排序
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 构建排序后的 map
	sortedMap := make(map[string]interface{})
	for _, key := range keys {
		sortedMap[key] = obj[key]
	}

	// 处理每个字段
	jsonMap := make(map[string]interface{})
	for key, value := range sortedMap {
		switch val := value.(type) {
		case []interface{}:
			// 数组类型: 递归处理每个元素
			strList := make([]string, 0, len(val))
			for _, item := range val {
				strList = append(strList, splicingSign(item))
			}
			jsonMap[key] = strList
		case map[string]interface{}:
			// 嵌套对象类型: 递归处理，放入列表
			rstStr := splicingSign(val)
			strList := []string{rstStr}
			jsonMap[key] = strList
		default:
			// 基本类型: 转换为字符串
			strValue := fmt.Sprintf("%v", val)
			if strValue != "" {
				jsonMap[key] = strValue
			}
		}
	}

	// 构建签名字符串
	// 获取所有键并排序
	sortedKeys := make([]string, 0, len(jsonMap))
	for key := range jsonMap {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	var builder strings.Builder
	for _, key := range sortedKeys {
		value := jsonMap[key]

		// 跳过特定字段
		if key == "SIGN_INF" || key == "Sign_Inf" ||
			key == "Svc_Rsp_St" || key == "Svc_Rsp_Cd" || key == "Rsp_Inf" {
			continue
		}

		switch val := value.(type) {
		case []string:
			// 数组类型: 直接拼接所有元素
			for _, v := range val {
				builder.WriteString(v)
			}
		default:
			// 其他类型: 使用 key=value& 格式
			builder.WriteString(key)
			builder.WriteByte('=')
			builder.WriteString(fmt.Sprintf("%v", val))
			builder.WriteByte('&')
		}
	}

	return builder.String()
}

// splicingArray 拼接 JSON 数组
func splicingArray(arr []interface{}) string {
	var builder strings.Builder
	for _, item := range arr {
		builder.WriteString(splicingSign(item))
	}
	return builder.String()
}
