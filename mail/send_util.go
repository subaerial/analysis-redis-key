package mail

import (
	"analysis.redis/config"
	"analysis.redis/util"
	"fmt"
)

// SendErrorEmail 发送运行异常邮件
func SendErrorEmail(err error) {
	subject := "redis全量分析脚本运行异常"
	body := fmt.Sprintf(`运行是发生异常: %s`, err)
	sendEmail(config.Properties.Mail.Dialer.Username, config.Properties.Mail.Receiver,
		subject, body)
}

// SendEndEmail 发送运行结束邮件
func SendEndEmail(runtime int64) {
	ip, _ := util.GetLocalIP()
	subject := "redis全量分析脚本执行通知"
	body := fmt.Sprintf(`redis key全量分析脚本一执行结束, 服务器: %s , 共耗时: %ds`, ip, runtime)
	sendEmail(config.Properties.Mail.Dialer.Username, config.Properties.Mail.Receiver,
		subject, body)
}
