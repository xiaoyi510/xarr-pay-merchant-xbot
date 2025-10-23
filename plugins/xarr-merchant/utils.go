package xarrmerchant

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// formatAmount 格式化金额(分转元)
func formatAmount(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}

// formatTime 格式化时间
func formatTime(timestamp int64) string {
	if timestamp == 0 {
		return "未设置"
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// maskSecret 密钥脱敏
func maskSecret(secret string) string {
	if secret == "" {
		return "未设置"
	}
	if len(secret) > 8 {
		return secret[:4] + "****" + secret[len(secret)-4:]
	}
	return "****"
}

// parseGroupIDs 解析群号列表（逗号分隔）
func parseGroupIDs(groupStr string) ([]int64, error) {
	if groupStr == "" {
		return []int64{}, nil
	}

	parts := strings.Split(groupStr, ",")
	groups := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		groupID, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("无效的群号: %s", part)
		}

		groups = append(groups, groupID)
	}

	return groups, nil
}

// formatGroupIDs 格式化群号列表
func formatGroupIDs(groups []int64) string {
	if len(groups) == 0 {
		return "无"
	}

	strs := make([]string, len(groups))
	for i, g := range groups {
		strs[i] = strconv.FormatInt(g, 10)
	}

	return strings.Join(strs, ", ")
}

// maskUserID 掩码用户ID
func maskUserID(uid int64) string {
	uidStr := strconv.FormatInt(uid, 10)
	if len(uidStr) <= 4 {
		return "****"
	}
	// 显示前2位，后面用****代替
	return uidStr[:2] + "****"
}

// maskUsername 掩码用户名
func maskUsername(username string) string {
	if username == "" {
		return "****"
	}

	// 转换为rune数组以正确处理中文
	runes := []rune(username)
	if len(runes) == 0 {
		return "****"
	}

	if len(runes) == 1 {
		return string(runes[0]) + "**"
	}

	// 显示第一个字符，后面用**代替
	return string(runes[0]) + "**"
}

// formatAvgAmount 格式化平均金额
func formatAvgAmount(totalAmount int64, orderCount int64) string {
	if orderCount == 0 {
		return "0.00"
	}
	avgAmount := float64(totalAmount) / float64(orderCount) / 100
	return fmt.Sprintf("%.2f", avgAmount)
}

// maskAmount 掩码金额（仅在群聊中使用）
func maskAmount(amount int64) string {
	if amount == 0 {
		return "0.00"
	}
	// 显示金额区间而不是具体数值
	amountFloat := float64(amount) / 100
	switch {
	case amountFloat < 100:
		return "< 100"
	case amountFloat < 500:
		return "100 - 500"
	case amountFloat < 1000:
		return "500 - 1000"
	case amountFloat < 5000:
		return "1K - 5K"
	case amountFloat < 10000:
		return "5K - 10K"
	case amountFloat < 50000:
		return "10K - 50K"
	case amountFloat < 100000:
		return "50K - 100K"
	default:
		return "> 100K"
	}
}

// maskAccountName 掩码账户名称
func maskAccountName(name string) string {
	if name == "" {
		return "****"
	}

	// 转换为rune数组以正确处理中文
	runes := []rune(name)
	if len(runes) == 0 {
		return "****"
	}

	if len(runes) <= 2 {
		return string(runes[0]) + "**"
	}

	// 显示首尾字符，中间用**代替
	return string(runes[0]) + "**" + string(runes[len(runes)-1])
}
