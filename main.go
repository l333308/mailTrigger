package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var recipient string = "304754343@qq.com"

type GmailAccount struct {
	Email    string
	Password string
}

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// 从环境变量中读取账号列表,账号格式： 邮箱:密码
	accountsStr := os.Getenv("ACCOUNTS")
	accounts := strings.Split(accountsStr, ",")

	for _, senderRaw := range accounts {
		senderArr := strings.Split(senderRaw, ":")
		email := senderArr[0]
		password := senderArr[1]
		sender := GmailAccount{Email: email, Password: password}

		sendEmail(sender)
	}
}

func sendEmail(sender GmailAccount) {
	// 设置邮件服务器配置
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 设置邮件内容
	subject := "Weekly Update"
	body := "This is a weekly update email."
	message := []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// 认证
	auth := smtp.PlainAuth("", sender.Email, sender.Password, smtpHost)

	// 发送邮件
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender.Email, []string{recipient}, message)
	if err != nil {
		log.Printf("发送邮件失败: %v", err)
		if strings.Contains(err.Error(), "534 5.7.9") {
			log.Println("这可能是因为需要使用应用专用密码。请检查您的 Google 账户设置。")
		}
	} else {
		fmt.Printf("从 %s 成功发送邮件到 %s\n", sender.Email, recipient)
	}
}
