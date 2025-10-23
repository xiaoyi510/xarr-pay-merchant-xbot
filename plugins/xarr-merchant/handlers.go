package xarrmerchant

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xiaoyi510/xbot"
	"github.com/xiaoyi510/xbot/event"
	"github.com/xiaoyi510/xbot/message"
)

// registerAdminHandlers 注册超管命令处理器
func registerAdminHandlers(engine *xbot.Engine) {
	// 超管命令 - 设置商户系统配置
	engine.OnRegex(`^/设置商户系统\s+(\S+)\s+(\S+)`, xbot.OnlyPrivateMessage()).Handle(func(ctx *xbot.Context) {
		if !ctx.IsSuperUser() {
			ctx.Reply("❌ 权限不足，仅超级管理员可操作")
			return
		}

		if ctx.RegexResult == nil || len(ctx.RegexResult.Groups) < 3 {
			ctx.Reply("❌ 参数不完整\n用法: /设置商户系统 <API地址> <Secret密钥>")
			return
		}

		baseURL := ctx.RegexResult.Groups[1]
		secret := ctx.RegexResult.Groups[2]

		// 检查client初始化
		if client == nil {
			ctx.Reply("❌ 系统初始化失败")
			return
		}

		// 创建新配置
		config := &MerchantConfig{
			BaseURL:       baseURL,
			Secret:        secret,
			AllowedGroups: []int64{}, // 初始化为空，需要单独设置
		}

		// 保存配置
		if err := client.SaveConfig(config); err != nil {
			ctx.Reply(fmt.Sprintf("❌ 保存失败: %s", err.Error()))
			return
		}

		msg := fmt.Sprintf("✅ 商户系统配置成功\n\n"+
			"API地址: %s\n"+
			"密钥: %s",
			baseURL,
			maskSecret(secret))

		ctx.Reply(msg)
	})

	// 超管命令 - 设置允许的群聊
	engine.OnRegex(`^/设置商户群聊\s+(.+)`, xbot.OnlyPrivateMessage(), xbot.OnlySuperUsers()).Handle(func(ctx *xbot.Context) {
		if ctx.RegexResult == nil || len(ctx.RegexResult.Groups) < 2 {
			ctx.Reply("❌ 参数不完整\n用法: /设置商户群聊 <群号1,群号2,...>\n示例: /设置商户群聊 123456,789012")
			return
		}

		groupsStr := ctx.RegexResult.Groups[1]

		// 解析群号列表
		groups, err := parseGroupIDs(groupsStr)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 解析失败: %s\n用法: /设置商户群聊 <群号1,群号2,...>", err.Error()))
			return
		}

		if len(groups) == 0 {
			ctx.Reply("❌ 至少需要设置一个群号")
			return
		}

		// 检查client初始化
		if client == nil {
			ctx.Reply("❌ 系统初始化失败")
			return
		}

		// 获取现有配置
		config, err := client.GetConfig()
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 获取配置失败: %s\n请先使用 /设置商户系统 配置API", err.Error()))
			return
		}

		// 更新允许的群聊
		config.AllowedGroups = groups

		// 保存配置
		if err := client.SaveConfig(config); err != nil {
			ctx.Reply(fmt.Sprintf("❌ 保存失败: %s", err.Error()))
			return
		}

		msg := fmt.Sprintf("✅ 允许的群聊设置成功\n\n"+
			"群聊数量: %d个\n"+
			"群号列表: %s",
			len(groups),
			formatGroupIDs(groups))

		ctx.Reply(msg)
	})

	// 超管命令 - 查看配置
	engine.OnCommand("查看商户配置", xbot.OnlyPrivateMessage()).Handle(func(ctx *xbot.Context) {
		if !ctx.IsSuperUser() {
			ctx.Reply("❌ 权限不足，仅超级管理员可操作")
			return
		}

		if client == nil {
			ctx.Reply("❌ 系统初始化失败")
			return
		}

		config, err := client.GetConfig()
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 获取配置失败: %s", err.Error()))
			return
		}

		msg := fmt.Sprintf("⚙️ 商户配置\n\n"+
			"API地址: %s\n"+
			"密钥: %s\n"+
			"允许的群聊: %s (%d个)",
			config.BaseURL,
			maskSecret(config.Secret),
			formatGroupIDs(config.AllowedGroups),
			len(config.AllowedGroups))

		ctx.Reply(msg)
	})
}

// registerUserHandlers 注册用户命令处理器
func registerUserHandlers(engine *xbot.Engine) {
	// 绑定商户账号
	engine.OnRegex(`^/绑定\s+(\S+)`).Handle(func(ctx *xbot.Context) {
		if ctx.RegexResult == nil || len(ctx.RegexResult.Groups) < 2 {
			ctx.Reply("❌ 请提供绑定ticket\n用法: /绑定 <ticket>")
			return
		}

		ticket := ctx.RegexResult.Groups[1]
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		err := client.BindUser(ticket, openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 绑定失败: %s", err.Error()))
			return
		}

		ctx.Reply("✅ 绑定成功!")
	})

	// 解绑商户账号
	engine.OnCommand("解绑").Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		err := client.UnbindUser(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 解绑失败: %s", err.Error()))
			return
		}

		ctx.Reply("✅ 解绑成功!")
	})

	// 查询用户信息
	engine.OnCommandGroup([]string{"我的信息", "个人信息"}).Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		userInfo, err := client.GetUserInfo(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 查询失败: %s", err.Error()))
			return
		}

		// 能查询到数据就说明已绑定，不需要额外判断
		statusText := "正常"
		if userInfo.Status != 1 {
			statusText = "已禁用"
		}

		// 判断是否是私聊，决定是否掩码处理
		var uidText, usernameText string
		if ctx.IsPrivateMessage() {
			// 私聊显示完整信息
			uidText = strconv.FormatInt(userInfo.UID, 10)
			usernameText = userInfo.Username
		} else {
			// 群聊掩码处理
			uidText = maskUserID(userInfo.UID)
			usernameText = maskUsername(userInfo.Username)
		}

		msg := fmt.Sprintf("📋 个人信息\n\n"+
			"用户ID: %s\n"+
			"用户名: %s\n"+
			"余额: ¥%s\n"+
			"状态: %s",
			uidText,
			usernameText,
			formatAmount(userInfo.Balance),
			statusText)

		ctx.Reply(msg)
	})

	// 查询余额
	engine.OnCommandGroup([]string{"余额", "查询余额"}).Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		balance, err := client.GetUserBalance(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 查询失败: %s", err.Error()))
			return
		}

		ctx.Reply(fmt.Sprintf("💰 当前余额: ¥%s", formatAmount(balance)))
	})

	// 查询套餐信息
	engine.OnCommandGroup([]string{"套餐信息", "我的套餐"}).Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		mealInfo, err := client.GetUserMealInfo(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 查询失败: %s", err.Error()))
			return
		}

		// 处理到期时间显示
		var expireTimeText, expireStatus string
		if mealInfo.ExpireTime == -1 {
			expireTimeText = "永久"
			expireStatus = "永久有效"
		} else if mealInfo.ExpireTime < time.Now().Unix() {
			expireTimeText = formatTime(mealInfo.ExpireTime)
			expireStatus = "已过期"
		} else {
			expireTimeText = formatTime(mealInfo.ExpireTime)
			expireStatus = "正常"
		}

		// 处理通道账号数
		var channelCountText string
		if mealInfo.ChannelAccountCount == -1 {
			channelCountText = "不限制"
		} else {
			channelCountText = fmt.Sprintf("%d", mealInfo.ChannelAccountCount)
		}

		// 处理日限额
		var dayLimitText string
		if mealInfo.DayLimit == -1 {
			dayLimitText = "不限制"
		} else {
			dayLimitText = "¥" + formatAmount(mealInfo.DayLimit)
		}

		// 处理月限额
		var monthLimitText string
		if mealInfo.MonthLimit == -1 {
			monthLimitText = "不限制"
		} else {
			monthLimitText = "¥" + formatAmount(mealInfo.MonthLimit)
		}

		msg := fmt.Sprintf("📦 套餐信息\n\n"+
			"套餐名称: %s\n"+
			"到期时间: %s\n"+
			"状态: %s\n"+
			"通道账号数: %s\n"+
			"日限额: %s\n"+
			"月限额: %s\n"+
			"费率: %.2f%%",
			mealInfo.MealName,
			expireTimeText,
			expireStatus,
			channelCountText,
			dayLimitText,
			monthLimitText,
			float64(mealInfo.Rate)/100)

		ctx.Reply(msg)
	})

	// 今日统计
	engine.OnCommandGroup([]string{"今日统计", "今日"}).Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		stat, err := client.GetUserPayStat(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 查询失败: %s", err.Error()))
			return
		}

		msg := fmt.Sprintf("📊 今日统计\n\n"+
			"💰 今日收款: ¥%s\n"+
			"📦 订单数量: %d 笔\n"+
			"📈 平均订单: ¥%s",
			formatAmount(stat.TodayAmount),
			stat.TodayOrderCount,
			formatAvgAmount(stat.TodayAmount, stat.TodayOrderCount))

		ctx.Reply(msg)
	})

	// 查询支付统计（全部）
	engine.OnCommandGroup([]string{"统计", "支付统计"}).Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		stat, err := client.GetUserPayStat(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 查询失败: %s", err.Error()))
			return
		}

		msg := fmt.Sprintf("📊 支付统计\n\n"+
			"【今日】\n"+
			"金额: ¥%s\n"+
			"订单: %d 笔\n\n"+
			"【本周】\n"+
			"金额: ¥%s\n"+
			"订单: %d 笔\n\n"+
			"【本月】\n"+
			"金额: ¥%s\n"+
			"订单: %d 笔\n\n"+
			"【总计】\n"+
			"金额: ¥%s\n"+
			"订单: %d 笔",
			formatAmount(stat.TodayAmount), stat.TodayOrderCount,
			formatAmount(stat.WeekAmount), stat.WeekOrderCount,
			formatAmount(stat.MonthAmount), stat.MonthOrderCount,
			formatAmount(stat.TotalAmount), stat.TotalOrderCount)

		ctx.Reply(msg)
	})

	// 查询渠道账户列表
	engine.OnCommandGroup([]string{"渠道列表", "账户列表"}).Handle(func(ctx *xbot.Context) {
		openID := strconv.FormatInt(ctx.GetUserID(), 10)

		accounts, err := client.GetChannelAccountList(openID)
		if err != nil {
			ctx.Reply(fmt.Sprintf("❌ 查询失败: %s", err.Error()))
			return
		}

		if len(accounts) == 0 {
			ctx.Reply("📋 暂无渠道账户")
			return
		}

		isPrivate := ctx.IsPrivateMessage()
		var msg strings.Builder
		msg.WriteString(fmt.Sprintf("📋 渠道账户列表 (共%d个)\n\n", len(accounts)))

		for i, acc := range accounts {
			statusText := "禁用"
			if acc.Status == 1 {
				statusText = "启用"
			}

			onlineText := "离线"
			if acc.Online == 1 {
				onlineText = "在线"
			}

			// 账户名称处理
			var accountName string
			if isPrivate {
				accountName = acc.Name
			} else {
				accountName = maskAccountName(acc.Name)
			}

			msg.WriteString(fmt.Sprintf("%d. %s\n", i+1, accountName))
			msg.WriteString(fmt.Sprintf("   支付方式: %s\n", acc.PayTypeName))
			msg.WriteString(fmt.Sprintf("   状态: %s | %s\n", statusText, onlineText))
			msg.WriteString(fmt.Sprintf("   今日: ¥%s / ¥%s\n",
				formatAmount(acc.DayAmount),
				formatAmount(acc.DayAmountLimit)))
			if i < len(accounts)-1 {
				msg.WriteString("\n")
			}
		}

		ctx.Reply(msg.String())
	})

	// 帮助菜单
	engine.OnCommandGroup([]string{"商户帮助"}).Handle(func(ctx *xbot.Context) {
		// 检查是否是超管
		isSuperUser := ctx.IsSuperUser()

		msg := `🤖 商户机器人使用指南

📝 账户管理
/绑定 <ticket> - 绑定商户账号
/解绑 - 解绑商户账号
/我的信息 - 查看个人信息
/余额 - 查看账户余额

📦 套餐相关
/套餐信息 - 查看套餐详情

📊 统计查询
/今日统计 - 查看今日数据
/统计 - 查看完整统计
/渠道列表 - 查看渠道账户`

		// 超管显示额外命令
		if isSuperUser {
			msg += `

⚙️ 超管命令
/设置商户系统 <API地址> <Secret> - 设置系统配置
/设置商户群聊 <群号1,群号2,...> - 设置允许的群聊
/查看商户配置 - 查看当前配置`
		}

		msg += `

💡 提示: 私聊机器人使用，部分已开通的群聊也可使用`

		ctx.Reply(msg)
	})
}

// registerGroupMessageHandler 注册群消息处理器
func registerGroupMessageHandler(engine *xbot.Engine) {
	// 群聊白名单中间件
	engine.Use(func(next func(*xbot.Context)) func(*xbot.Context) {
		return func(ctx *xbot.Context) {
			// 只处理群消息
			if evt, ok := ctx.Event.(*event.GroupMessageEvent); ok {
				// 检查是否是商户命令
				text := ctx.GetPlainText()
				if strings.HasPrefix(text, "/") {
					re := regexp.MustCompile(`^/(\S+)`)
					if matches := re.FindStringSubmatch(text); len(matches) > 1 {
						cmd := matches[1]
						// 商户相关命令列表
						merchantCmds := []string{"绑定", "解绑", "我的信息", "个人信息", "余额", "查询余额", "套餐信息", "我的套餐", "今日统计", "今日", "统计", "支付统计", "渠道列表", "账户列表", "商户帮助", "商户菜单"}

						isMerchantCmd := false
						for _, mc := range merchantCmds {
							if cmd == mc {
								isMerchantCmd = true
								break
							}
						}

						if isMerchantCmd {
							// 检查群聊是否在白名单中
							if client != nil && !client.IsGroupAllowed(evt.GroupID) {
								msg := message.NewBuilder().
									Reply(evt.MessageID).
									Text("⚠️ 该群聊未开通商户功能\n如需开通，请联系超级管理员").
									Build()
								ctx.Reply(msg)
								ctx.Abort()
								return
							}
						}
					}
				}
			}
			next(ctx)
		}
	})
}
