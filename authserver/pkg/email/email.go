package email

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/resend/resend-go/v2"
	"github.com/simple-auth-roles/internal/config"
)

// EmailService defines email operations
type EmailService interface {
	SendLoginCodeEmail(ctx context.Context, email, code string) error
	SendWelcomeEmail(ctx context.Context, email, name string) error
}

type emailService struct {
	config       *config.Config
	logger       *slog.Logger
	resendClient *resend.Client
}

func NewEmailService(cfg *config.Config, logger *slog.Logger) EmailService {
	var resendClient *resend.Client
	if cfg.Email.ResendAPIKey != "" {
		resendClient = resend.NewClient(cfg.Email.ResendAPIKey)
	}

	return &emailService{
		config:       cfg,
		logger:       logger.With("service", "email"),
		resendClient: resendClient,
	}
}

func (e *emailService) SendLoginCodeEmail(ctx context.Context, email, code string) error {
	subject := "Your Login Code"
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your Login Code</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .code { background: #f4f4f4; padding: 15px; font-size: 24px; font-weight: bold; text-align: center; margin: 20px 0; border-radius: 8px; }
        .footer { margin-top: 30px; font-size: 14px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Your Login Code</h2>
        <p>Hello!</p>
        <p>Your login code is:</p>
        <div class="code">%s</div>
        <p>This code will expire in 10 minutes.</p>
        <p>If you didn't request this code, please ignore this email.</p>
        <div class="footer">
            <p>Best regards,<br>%s Team</p>
        </div>
    </div>
</body>
</html>
`, code, e.config.Email.FromName)

	return e.sendEmail(email, subject, htmlBody)
}

func (e *emailService) SendWelcomeEmail(ctx context.Context, email, name string) error {
	displayName := name
	if displayName == "" {
		displayName = "there"
	}

	subject := fmt.Sprintf("Welcome to %s!", e.config.Email.FromName)
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .welcome { background: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0; }
        .footer { margin-top: 30px; font-size: 14px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="welcome">
            <h2>Welcome to %s!</h2>
        </div>
        <p>Hello %s!</p>
        <p>We're excited to have you on board!</p>
        <p>You can now log in using your email address. We'll send you a secure login code each time you sign in.</p>
        <p>If you have any questions, feel free to reach out to our support team.</p>
        <div class="footer">
            <p>Best regards,<br>%s Team</p>
        </div>
    </div>
</body>
</html>
`, e.config.Email.FromName, displayName, e.config.Email.FromName)

	return e.sendEmail(email, subject, htmlBody)
}

func (e *emailService) sendEmail(to, subject, htmlBody string) error {
	// If Resend is not configured, just log the email
	if e.resendClient == nil {
		e.logger.Info("Email would be sent (Resend not configured)",
			"to", to,
			"subject", subject,
			"body", htmlBody,
		)
		return nil
	}

	// Send email using Resend
	params := &resend.SendEmailRequest{
		From:    e.config.Email.FromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := e.resendClient.Emails.Send(params)
	if err != nil {
		e.logger.Error("Failed to send email", "error", err, "to", to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.Info("Email sent successfully", "to", to, "subject", subject)
	return nil
}
