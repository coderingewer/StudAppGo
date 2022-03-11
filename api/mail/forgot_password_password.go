package mail

import (
	"net/http"
	"os"

	"github.com/matcornic/hermes/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendMail struct{}

type EmailResponse struct {
	Status   int
	RespBody string
}

type sendMaler interface {
	SendResetPassword(string, string, string, string, string) (*EmailResponse, error)
}

func (s *sendMail) SendResetPassword(ToUser string, FromAdmin string, Token string, Sendgridkey string, AppEnv string) (*EmailResponse, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "StudApp",
			Link: "https://localhost:8000",
		},
	}
	var forgotUrl string

	if os.Getenv("SNGRİ_KEY") == "production" {
		forgotUrl = "https://studapp.com/resetpassword/" + Token
	} else {
		forgotUrl = "https://localhost:8000/resetpassword/" + Token
	}
	email := hermes.Email{
		Body: hermes.Body{
			Name: ToUser,
			Intros: []string{
				"Welcome to StudApp! Good to have you here.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click this link to reset your password",
					Button: hermes.Button{
						Color: "#FFFFFF",
						Text:  "Reset Password",
						Link:  forgotUrl,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return nil, err
	}
	from := mail.NewEmail("StudApp", FromAdmin)
	subject := "Reset Password"
	to := mail.NewEmail("Reset Password", ToUser)
	message := mail.NewSingleEmail(from, subject, to, emailBody, emailBody)
	client := sendgrid.NewSendClient(Sendgridkey)
	_, err = client.Send(message)
	if err != nil {
		return nil, err
	}
	return &EmailResponse{
		Status:   http.StatusOK,
		RespBody: "Success, Please click on the link provided in your email",
	}, nil
}
