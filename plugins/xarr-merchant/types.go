package xarrmerchant

// MerchantConfig 商户配置
type MerchantConfig struct {
	BaseURL       string  `json:"base_url"`       // API基础地址
	Secret        string  `json:"secret"`         // API密钥
	AllowedGroups []int64 `json:"allowed_groups"` // 允许使用的群聊列表
}

// UserInfo 用户信息
type UserInfo struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"` // 余额(分)
	Status   int    `json:"status"`  // 状态
	IsBind   bool   `json:"is_bind"` // 是否绑定
}

// UserMealInfo 用户套餐信息
type UserMealInfo struct {
	MealName            string `json:"meal_name"`             // 套餐名称
	ExpireTime          int64  `json:"expire_time"`           // 到期时间
	ChannelAccountCount int    `json:"channel_account_count"` // 通道账号可添加数
	DayLimit            int64  `json:"day_limit"`             // 每日收款限额(分)
	MonthLimit          int64  `json:"month_limit"`           // 每月收款限额(分)
	Rate                int    `json:"rate"`                  // 费率(分)
}

// UserPayStat 用户支付统计
type UserPayStat struct {
	TodayAmount     int64 `json:"today_amount"`      // 今日支付金额(分)
	TodayOrderCount int64 `json:"today_order_count"` // 今日订单数量
	WeekAmount      int64 `json:"week_amount"`       // 本周支付金额(分)
	WeekOrderCount  int64 `json:"week_order_count"`  // 本周订单数量
	MonthAmount     int64 `json:"month_amount"`      // 本月支付金额(分)
	MonthOrderCount int64 `json:"month_order_count"` // 本月订单数量
	TotalAmount     int64 `json:"total_amount"`      // 总支付金额(分)
	TotalOrderCount int64 `json:"total_order_count"` // 总订单数量
}

// ChannelAccount 渠道账户信息
type ChannelAccount struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`             // 账号名称
	PayType        string `json:"pay_type"`         // 支付类型
	PayTypeName    string `json:"pay_type_name"`    // 支付方式名称
	Status         int    `json:"status"`           // 状态
	Online         int    `json:"online"`           // 在线状态
	DayAmount      int64  `json:"day_amount"`       // 今日已支付额度(分)
	DayAmountLimit int64  `json:"day_amount_limit"` // 单日限额(分)
}
