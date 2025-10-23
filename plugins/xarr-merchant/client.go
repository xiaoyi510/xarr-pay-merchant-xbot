package xarrmerchant

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"github.com/xiaoyi510/xbot/storage"
)

const (
	// ConnectType 固定为 xbot
	ConnectType = "xbot"
	// ConfigKey 配置存储key
	ConfigKey = "merchant:config"
)

var (
	// 配置锁
	configMu sync.RWMutex
)

// MerchantClient 商户API客户端
type MerchantClient struct {
	storage storage.Storage
	client  *req.Client
}

// NewMerchantClient 创建商户API客户端
func NewMerchantClient(store storage.Storage) *MerchantClient {
	return &MerchantClient{
		storage: store,
		client:  req.C(),
	}
}

// GetConfig 获取配置
func (c *MerchantClient) GetConfig() (*MerchantConfig, error) {
	configMu.RLock()
	defer configMu.RUnlock()

	data, err := c.storage.Get(ConfigKey)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("未配置商户API，请联系管理员设置")
	}

	var config MerchantConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 保存配置
func (c *MerchantClient) SaveConfig(config *MerchantConfig) error {
	configMu.Lock()
	defer configMu.Unlock()

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return c.storage.Set(ConfigKey, data)
}

// IsGroupAllowed 检查群聊是否在白名单中
func (c *MerchantClient) IsGroupAllowed(groupID int64) bool {
	config, err := c.GetConfig()
	if err != nil {
		return false
	}

	// 如果白名单为空，则不允许任何群
	if len(config.AllowedGroups) == 0 {
		return false
	}

	// 检查群ID是否在白名单中
	for _, allowedGroup := range config.AllowedGroups {
		if allowedGroup == groupID {
			return true
		}
	}

	return false
}

// generateSign 生成签名
// sign = md5(参数按key排序后拼接 + secret)
func generateSign(params map[string]string, secret string) string {
	// 获取所有key并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" { // 排除sign字段本身
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接参数
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(params[k])
	}

	// 添加secret (直接拼接，不加&)
	builder.WriteString(secret)

	// 计算MD5
	hash := md5.Sum([]byte(builder.String()))
	return hex.EncodeToString(hash[:])
}

// addSignature 为请求添加签名
func (c *MerchantClient) addSignature(params map[string]string) (map[string]string, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	// 添加时间戳(10位)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params["timestamp"] = timestamp

	// 生成签名
	sign := generateSign(params, config.Secret)
	params["sign"] = sign

	return params, nil
}

// GetUserInfo 获取用户信息
func (c *MerchantClient) GetUserInfo(openID string) (*UserInfo, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"open_id":      openID,
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/user/info")

	if err != nil {
		return nil, err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return nil, errors.New(message)
	}

	result := gjson.Get(data, "data")
	if !result.Exists() {
		return nil, errors.New("未找到用户信息")
	}

	// 能获取到数据就说明已绑定
	return &UserInfo{
		UID:      result.Get("id").Int(), // 使用id字段
		Username: result.Get("username").String(),
		Balance:  result.Get("balance").Int(),
		Status:   int(result.Get("status").Int()),
		IsBind:   true, // 能查询到数据就是已绑定
	}, nil
}

// GetUserBalance 获取用户余额
func (c *MerchantClient) GetUserBalance(openID string) (int64, error) {
	config, err := c.GetConfig()
	if err != nil {
		return 0, err
	}

	params := map[string]string{
		"open_id":      openID,
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return 0, err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/user/balance")

	if err != nil {
		return 0, err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return 0, errors.New(message)
	}

	balance := gjson.Get(data, "data.balance").Int()
	return balance, nil
}

// GetUserMealInfo 获取用户套餐信息
func (c *MerchantClient) GetUserMealInfo(openID string) (*UserMealInfo, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"open_id":      openID,
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/user/meal-info")

	if err != nil {
		return nil, err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return nil, errors.New(message)
	}

	result := gjson.Get(data, "data")
	if !result.Exists() {
		return nil, errors.New("未找到套餐信息")
	}
	// "{\"code\":200,\"message\":\"获取成功\",\"data\":{\"expire_time\":-1,\"meal_name\":\"终极会员\",\"rate\":0,\"channel_account_count\":-1,\"day_limit\":-1,\"month_limit\":-1},\"redirect\":\"\"}\n"
	return &UserMealInfo{
		MealName:            result.Get("meal_name").String(),
		ExpireTime:          result.Get("expire_time").Int(),
		ChannelAccountCount: int(result.Get("channel_account_count").Int()),
		DayLimit:            result.Get("day_limit").Int(),
		MonthLimit:          result.Get("month_limit").Int(),
		Rate:                int(result.Get("rate").Int()),
	}, nil
}

// GetUserPayStat 获取用户支付统计
func (c *MerchantClient) GetUserPayStat(openID string) (*UserPayStat, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"open_id":      openID,
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/user/pay-stat")

	if err != nil {
		return nil, err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return nil, errors.New(message)
	}

	result := gjson.Get(data, "data")
	if !result.Exists() {
		return nil, errors.New("未找到统计信息")
	}

	return &UserPayStat{
		TodayAmount:     result.Get("today_amount").Int(),
		TodayOrderCount: result.Get("today_order_count").Int(),
		WeekAmount:      result.Get("week_amount").Int(),
		WeekOrderCount:  result.Get("week_order_count").Int(),
		MonthAmount:     result.Get("month_amount").Int(),
		MonthOrderCount: result.Get("month_order_count").Int(),
		TotalAmount:     result.Get("total_amount").Int(),
		TotalOrderCount: result.Get("total_order_count").Int(),
	}, nil
}

// GetChannelAccountList 获取渠道账户列表
func (c *MerchantClient) GetChannelAccountList(openID string) ([]ChannelAccount, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"open_id":      openID,
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/channel-account/list")

	if err != nil {
		return nil, err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return nil, errors.New(message)
	}

	result := gjson.Get(data, "data")
	if !result.Exists() {
		return []ChannelAccount{}, nil
	}

	var accounts []ChannelAccount
	for _, item := range result.Array() {
		accounts = append(accounts, ChannelAccount{
			ID:             item.Get("id").Int(),
			Name:           item.Get("name").String(),
			PayType:        item.Get("pay_type").String(),
			PayTypeName:    item.Get("pay_type_name").String(),
			Status:         int(item.Get("status").Int()),
			Online:         int(item.Get("online").Int()),
			DayAmount:      item.Get("day_amount").Int(),
			DayAmountLimit: item.Get("day_amount_limit").Int(),
		})
	}

	return accounts, nil
}

// BindUser 绑定用户
func (c *MerchantClient) BindUser(ticket, openID string) error {
	config, err := c.GetConfig()
	if err != nil {
		return err
	}

	params := map[string]string{
		"ticket":       ticket,
		"open_id":      openID,
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/user/bind")

	if err != nil {
		return err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return errors.New(message)
	}

	return nil
}

// UnbindUser 解绑用户
func (c *MerchantClient) UnbindUser(openID string) error {
	config, err := c.GetConfig()
	if err != nil {
		return err
	}

	params := map[string]string{
		"connect_type": ConnectType,
	}

	params, err = c.addSignature(params)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetFormData(params).
		Post(config.BaseURL + "/api/system-api/user/unbind")

	if err != nil {
		return err
	}

	data, _ := resp.ToString()
	code := gjson.Get(data, "code").Int()
	if code != 200 {
		message := gjson.Get(data, "message").String()
		return errors.New(message)
	}

	return nil
}
