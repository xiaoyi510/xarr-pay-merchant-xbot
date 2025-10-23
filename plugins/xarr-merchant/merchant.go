package xarrmerchant

import (
	"github.com/xiaoyi510/xbot"
	"github.com/xiaoyi510/xbot/logger"
	"github.com/xiaoyi510/xbot/storage"
)

var (
	// API客户端
	client *MerchantClient
	// 插件专属storage
	storageDB storage.Storage
)

func init() {
	engine := xbot.NewEngine()
	engine.UseRecovery().UseLogger()

	// 初始化插件专属storage
	storageDB = xbot.GetStorage("xarr_merchant")

	// 初始化商户客户端
	client = NewMerchantClient(storageDB)
	logger.Info("商户机器人客户端初始化成功")

	// 注册超管命令
	registerAdminHandlers(engine)

	// 注册用户命令
	registerUserHandlers(engine)

	// 注册群消息处理
	registerGroupMessageHandler(engine)

	logger.Info("商户机器人插件已加载")
}
