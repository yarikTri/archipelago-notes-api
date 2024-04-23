package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type InvitationType uint8

const (
	NoteInvitationType InvitationType = iota
	DirInvitationType
)

func (vt *InvitationType) String() string {
	switch *vt {
	case NoteInvitationType:
		return "note"
	case DirInvitationType:
		return "directory"
	}

	return "note"
}

func (vt *InvitationType) toLink(id string) string {
	var subpath string
	switch *vt {
	case NoteInvitationType:
		subpath = "notes"
	case DirInvitationType:
		subpath = "dirs"
	}

	return fmt.Sprintf("%s/%s/%s", os.Getenv("SCHEME_AND_HOST"), subpath, id)
}

type IEmailInvitationClient interface {
	SendInvitation(to string, visitType InvitationType, id string) error
}

type SmtpInvitationClient struct {
	From     string
	Endpoint string
	Auth     smtp.Auth
}

func NewSmtpInvitationClient() *SmtpInvitationClient {
	host := os.Getenv("EMAIL_HOST")
	inbox := os.Getenv("EMAIL_INBOX")

	return &SmtpInvitationClient{
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

func (s *SmtpInvitationClient) SendInvitation(to string, visitType InvitationType, id string) error {
	msg := fmt.Sprintf(
		"Subject: Invitation to the %s\n"+
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
			"Hi! Archipelago team here!<br/>"+
			"You received an invitation to the %s<br/>"+
			"Tap the link below:<br/>"+
			"<br/>%s<br/><br/>"+
			"Any questions? Contact us in telegram: @yarik_tri",
		visitType.String(), visitType.String(), visitType.toLink(id),
	)

	return smtp.SendMail(s.Endpoint, s.Auth, s.From, []string{to}, []byte(msg))
}
