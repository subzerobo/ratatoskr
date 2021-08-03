package mailer

type SendConfig struct {
	SenderName string
	SenderEmail string
	SenderTemplate string
}

type SendOption func(config *SendConfig)

func WithSenderName(senderName string) SendOption {
	return func(args *SendConfig) {
		if senderName != "" {
			args.SenderName = senderName
		}
	}
}

func WithSenderEmail(senderEmail string) SendOption {
	return func(args *SendConfig) {
		if senderEmail != "" {
			args.SenderName = senderEmail
		}
	}
}

func WithSenderTemplate(template string) SendOption {
	return func(args *SendConfig) {
		if template != "" {
			args.SenderTemplate = template
		}
	}
}


