// 特殊 Hook 使用示例（Telegram、Trace 等）

package main

import (
	"fmt"
	"os"

	"github.com/ergoapi/util/log/hooks/tg"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("=== 特殊 Hook 使用示例 ===")
	fmt.Println()

	// 检查环境变量
	token := os.Getenv("TG_BOT_TOKEN")
	chatID := os.Getenv("TG_CHAT_ID")

	if token == "" || chatID == "" {
		fmt.Println(">>> Telegram Hook 示例（需要配置环境变量）")
		fmt.Println("请设置以下环境变量：")
		fmt.Println("  export TG_BOT_TOKEN='你的Bot Token'")
		fmt.Println("  export TG_CHAT_ID='你的Chat ID'")
		fmt.Println()
		telegramHookExample()
	} else {
		fmt.Println(">>> 使用真实的 Telegram Hook")
		realTelegramExample(token, chatID)
	}
}

// telegramHookExample 展示 Telegram Hook 的配置方法
func telegramHookExample() {
	fmt.Println("Telegram Hook 配置示例：")
	fmt.Println(`
	tghook, err := tg.NewTGHook(tg.TGConfig{
		Token:  "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",  // Bot Token
		ChatID: -1001234567890,                          // Chat ID（负数表示群组）
		Level:  logrus.WarnLevel,                        // 只发送 Warn 及以上级别
	})
	if err != nil {
		logrus.Fatal(err)
	}
	
	logrus.AddHook(tghook)
	
	// 这些日志会发送到 Telegram
	logrus.Warn("⚠️ 系统警告：内存使用率超过 80%")
	logrus.Error("❌ 系统错误：数据库连接失败")
	
	// Info 级别不会发送（低于 WarnLevel）
	logrus.Info("普通信息不会发送到 Telegram")
	`)

	fmt.Println("\n使用场景：")
	fmt.Println("  - 生产环境的紧急告警")
	fmt.Println("  - 关键错误的即时通知")
	fmt.Println("  - 系统监控和运维告警")

	fmt.Println("\n配置步骤：")
	fmt.Println("  1. 在 Telegram 中找 @BotFather 创建 Bot")
	fmt.Println("  2. 获取 Bot Token")
	fmt.Println("  3. 将 Bot 加入群组或对话")
	fmt.Println("  4. 获取 Chat ID（可以用 @userinfobot）")
}

// realTelegramExample 使用真实配置发送消息
func realTelegramExample(token, chatID string) {
	// 转换 chatID
	var id int64
	fmt.Sscanf(chatID, "%d", &id)

	tghook, err := tg.NewTGHook(tg.TGConfig{
		Token:  token,
		ChatID: id,
		Level:  logrus.WarnLevel,
	})
	if err != nil {
		fmt.Printf("创建 Telegram Hook 失败: %v\n", err)
		return
	}

	// 创建独立的 logger 避免影响全局
	log := logrus.New()
	log.AddHook(tghook)

	// 发送测试消息
	log.Warn("⚠️ 测试警告：这是一条来自 ergoapi/util 的测试消息")
	log.Error("❌ 测试错误：请忽略这条测试错误")

	fmt.Println("✓ 已发送测试消息到 Telegram")
	fmt.Printf("  Chat ID: %s\n", chatID)
}
