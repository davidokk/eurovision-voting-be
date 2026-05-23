package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewSMTPFromEnv() *Client {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "smtp.yandex.ru"
	}
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "465"
	}
	user := strings.TrimSpace(os.Getenv("SMTP_USER"))
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if from == "" {
		from = user
	}
	return &Client{
		host:     host,
		port:     port,
		user:     user,
		password: os.Getenv("SMTP_PASSWORD"),
		from:     from,
	}
}

func (c *Client) Enabled() bool {
	return c.user != "" && c.password != "" && c.from != ""
}

func (c *Client) SendVerificationCode(to, code, subjectLine string) error {
	if !c.Enabled() {
		return fmt.Errorf("email service not configured (SMTP_USER, SMTP_PASSWORD)")
	}

	subject := subjectLine
	if subject == "" {
		subject = "Код подтверждения — Eurovision Voting"
	}

	d0, d1, d2, d3, d4, d5 := code[0:1], code[1:2], code[2:3], code[3:4], code[4:5], code[5:6]
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#0b1220;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#0b1220;padding:32px 16px;">
    <tr><td align="center">
      <table width="100%%" style="max-width:480px;background:linear-gradient(160deg,#0f172a 0%%,#1e1b4b 100%%);border-radius:20px;border:1px solid rgba(79,124,255,0.25);overflow:hidden;">
        <tr><td style="padding:28px 28px 8px;text-align:center;">
          <div style="font-size:28px;line-height:1">🎤</div>
          <div style="margin-top:12px;font-size:11px;font-weight:700;letter-spacing:0.12em;text-transform:uppercase;color:#7aa2ff;">Eurovision Voting</div>
        </td></tr>
        <tr><td style="padding:8px 28px 4px;">
          <p style="margin:0;font-size:22px;font-weight:800;color:#f8fafc;">Привет!</p>
          <p style="margin:12px 0 0;font-size:15px;line-height:1.55;color:#94a3b8;">Вот ваш код для подтверждения. Введите его в приложении — и можно продолжать.</p>
        </td></tr>
        <tr><td style="padding:20px 28px;">
          <table cellpadding="0" cellspacing="0" align="center"><tr>
            <td style="width:44px;height:52px;background:rgba(15,23,42,0.9);border:2px solid rgba(79,124,255,0.4);border-radius:12px;text-align:center;font-size:24px;font-weight:800;color:#e2e8f0;">%s</td>
            <td width="8"></td>
            <td style="width:44px;height:52px;background:rgba(15,23,42,0.9);border:2px solid rgba(79,124,255,0.4);border-radius:12px;text-align:center;font-size:24px;font-weight:800;color:#e2e8f0;">%s</td>
            <td width="8"></td>
            <td style="width:44px;height:52px;background:rgba(15,23,42,0.9);border:2px solid rgba(79,124,255,0.4);border-radius:12px;text-align:center;font-size:24px;font-weight:800;color:#e2e8f0;">%s</td>
            <td width="8"></td>
            <td style="width:44px;height:52px;background:rgba(15,23,42,0.9);border:2px solid rgba(79,124,255,0.4);border-radius:12px;text-align:center;font-size:24px;font-weight:800;color:#e2e8f0;">%s</td>
            <td width="8"></td>
            <td style="width:44px;height:52px;background:rgba(15,23,42,0.9);border:2px solid rgba(79,124,255,0.4);border-radius:12px;text-align:center;font-size:24px;font-weight:800;color:#e2e8f0;">%s</td>
            <td width="8"></td>
            <td style="width:44px;height:52px;background:rgba(15,23,42,0.9);border:2px solid rgba(124,77,255,0.55);border-radius:12px;text-align:center;font-size:24px;font-weight:800;color:#c4b5fd;">%s</td>
          </tr></table>
        </td></tr>
        <tr><td style="padding:4px 28px 28px;">
          <p style="margin:0;font-size:13px;line-height:1.5;color:#64748b;text-align:center;">Код действует <strong style="color:#94a3b8;">15 минут</strong>. Если вы не запрашивали письмо — просто проигнорируйте его.</p>
        </td></tr>
      </table>
      <p style="margin:20px 0 0;font-size:12px;color:#475569;">© Eurovision Voting</p>
    </td></tr>
  </table>
</body>
</html>`, d0, d1, d2, d3, d4, d5)

	return c.send(strings.TrimSpace(to), subject, html)
}

func (c *Client) send(to, subject, htmlBody string) error {
	msg := buildMessage(c.from, to, subject, htmlBody)
	addr := net.JoinHostPort(c.host, c.port)

	port, _ := strconv.Atoi(c.port)
	if port == 465 {
		return c.sendSMTPS(addr, to, msg)
	}
	auth := smtp.PlainAuth("", c.user, c.password, c.host)
	return smtp.SendMail(addr, auth, c.from, []string{to}, msg)
}

func (c *Client) sendSMTPS(addr, to string, msg []byte) error {
	tlsCfg := &tls.Config{ServerName: c.host}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return fmt.Errorf("smtp tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", c.user, c.password, c.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(c.from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp data close: %w", err)
	}
	return client.Quit()
}

func buildMessage(from, to, subject, htmlBody string) []byte {
	var buf bytes.Buffer
	buf.WriteString("From: " + from + "\r\n")
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)
	return buf.Bytes()
}
