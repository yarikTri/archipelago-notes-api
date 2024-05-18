package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type IEmailInvitationClient interface {
	SendInvitation(to string, visitType InvitationType, id string) error
}

type IEmailConfirmationClient interface {
	SendConfirmation(to string, userID string) error
}

type EmailClient struct {
	From     string
	Endpoint string
	Auth     smtp.Auth
}

func NewEmailClient() *EmailClient {
	host := os.Getenv("EMAIL_HOST")
	inbox := os.Getenv("EMAIL_INBOX")

	return &EmailClient{
		From:     inbox,
		Endpoint: fmt.Sprintf("%s:%s", host, os.Getenv("EMAIL_PORT")),
		Auth: smtp.PlainAuth(
			"",
			inbox,
			os.Getenv("EMAIL_PASSWORD"),
			host,
		),
	}
}

type InvitationType uint8

const (
	NoteInvitationType InvitationType = iota
	DirInvitationType
)

func (vt *InvitationType) stringRussianTo() string {
	switch *vt {
	case NoteInvitationType:
		return "заметку"
	case DirInvitationType:
		return "папку"
	}

	return "заметку"
}

func (vt *InvitationType) toLink(resoureID string) string {
	var subpath string
	switch *vt {
	case NoteInvitationType:
		subpath = "notes"
	case DirInvitationType:
		subpath = "dirs"
	}

	return fmt.Sprintf("%s/%s/%s", os.Getenv("SCHEME_AND_HOST"), subpath, resoureID)
}

func toConfirmationLink(userID string) string {
	return fmt.Sprintf("%s?confirm_email_user_id=%s", os.Getenv("SCHEME_AND_HOST"), userID)
}

func (s *EmailClient) SendInvitation(to string, visitType InvitationType, resourceID string) error {
	msg := fmt.Sprintf(
		"Subject: Приглашение в %s\n"+
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
			"Привет! Это команда Archipelago!<br/>"+
			"Ты получил приглашение в %s<br/>"+
			"Перейти в неё можно по ссылке снизу:<br/>"+
			"<br/>%s<br/><br/>"+
			"Есть вопросы? Свяжись с нами в телеграме: @yarik_tri или @rbeketov",
		visitType.stringRussianTo(), visitType.stringRussianTo(), visitType.toLink(resourceID),
	)

	return smtp.SendMail(s.Endpoint, s.Auth, s.From, []string{to}, []byte(msg))
}

func (s *EmailClient) SendConfirmation(to string, userID string) error {
	msg := fmt.Sprintf(
		"Subject: Подтверждение почты\n"+
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
			"Привет! Это команда Archipelago!<br/>"+
			"Чтобы полноценно пользоваться нашим сервисом, необходимо подтвердить свой почтовый ящик<br/>"+
			"Для этого достаточно перейти по ссылке снизу:"+
			"<br/><br/>"+
			"<a href=\"%s\">Ссылка для подтверждения почты</a>"+
			"<br/><br/>"+
			"Есть вопросы? Свяжись с нами в телеграме: @yarik_tri или @rbeketov",
		toConfirmationLink(userID),
	)

	return smtp.SendMail(s.Endpoint, s.Auth, s.From, []string{to}, []byte(msg))
}
