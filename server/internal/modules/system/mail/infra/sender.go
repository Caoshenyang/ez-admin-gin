package infra

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	mailapp "ez-admin-gin/server/internal/modules/system/mail/application"
	"ez-admin-gin/server/internal/platform/model"
)

// SMTPSender 使用标准 SMTP 协议发送邮件。
type SMTPSender struct{}

func NewSMTPSender() *SMTPSender {
	return &SMTPSender{}
}

func (s *SMTPSender) Send(account model.MailAccount, message mailapp.Message) error {
	address := net.JoinHostPort(account.Host, fmt.Sprintf("%d", account.Port))
	client, err := s.connect(address, account)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Quit()
	}()

	if account.Encryption == model.MailEncryptionSTARTTLS {
		tlsConfig := &tls.Config{ServerName: account.Host, MinVersion: tls.VersionTLS12}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	username := strings.TrimSpace(account.Username)
	if username != "" {
		auth := smtp.PlainAuth("", username, account.Password, account.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(account.FromEmail); err != nil {
		return fmt.Errorf("set mail from: %w", err)
	}
	recipients := append(append([]string{}, message.To...), message.Cc...)
	recipients = append(recipients, message.Bcc...)
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("set recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("open smtp data: %w", err)
	}
	if _, err := writer.Write(buildMimeMessage(account, message)); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write smtp data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close smtp data: %w", err)
	}
	return nil
}

func (s *SMTPSender) connect(address string, account model.MailAccount) (*smtp.Client, error) {
	if account.Encryption != model.MailEncryptionSSL {
		conn, err := net.DialTimeout("tcp", address, 10*time.Second)
		if err != nil {
			return nil, fmt.Errorf("connect smtp: %w", err)
		}
		client, err := smtp.NewClient(conn, account.Host)
		if err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("create smtp client: %w", err)
		}
		return client, nil
	}

	tlsConfig := &tls.Config{ServerName: account.Host, MinVersion: tls.VersionTLS12}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", address, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("connect smtp ssl: %w", err)
	}
	client, err := smtp.NewClient(conn, account.Host)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("create smtp client: %w", err)
	}
	return client, nil
}

func buildMimeMessage(account model.MailAccount, message mailapp.Message) []byte {
	var buffer bytes.Buffer
	from := mail.Address{Name: account.FromName, Address: account.FromEmail}

	headers := map[string]string{
		"From":         from.String(),
		"To":           strings.Join(message.To, ", "),
		"Subject":      mime.QEncoding.Encode("utf-8", message.Subject),
		"MIME-Version": "1.0",
		"Date":         time.Now().Format(time.RFC1123Z),
	}
	if len(message.Cc) > 0 {
		headers["Cc"] = strings.Join(message.Cc, ", ")
	}
	if message.IsHTML {
		headers["Content-Type"] = `text/html; charset="utf-8"`
	} else {
		headers["Content-Type"] = `text/plain; charset="utf-8"`
	}
	headers["Content-Transfer-Encoding"] = "8bit"

	headerOrder := []string{"From", "To", "Cc", "Subject", "MIME-Version", "Date", "Content-Type", "Content-Transfer-Encoding"}
	for _, key := range headerOrder {
		value := headers[key]
		if value == "" {
			continue
		}
		buffer.WriteString(key)
		buffer.WriteString(": ")
		buffer.WriteString(value)
		buffer.WriteString("\r\n")
	}
	buffer.WriteString("\r\n")
	buffer.WriteString(message.Content)
	return buffer.Bytes()
}
