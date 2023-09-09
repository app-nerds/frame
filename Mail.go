package frame

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

type emailDynamicTemplateRequest struct {
	From             emailAddress            `json:"from"`
	Personalizations []emailPersonalizations `json:"personalizations"`
	TemplateID       string                  `json:"template_id"`
}

type emailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type emailPersonalizations struct {
	To           []emailAddress         `json:"to"`
	TemplateData map[string]interface{} `json:"dynamic_template_data"`
}

type EmailServicer interface {
	Clear() EmailServicer
	From(email, name string) EmailServicer
	Send(templateID string) error
	TemplateData(to string, data map[string]interface{}) EmailServicer
	To(email, name string) EmailServicer
	ToMultipleAddresses(emails []string) EmailServicer
}

type emailServiceConfig struct {
	ApiKey string
}

type emailService struct {
	apiKey       string
	from         emailAddress
	to           []emailAddress
	templateData map[string]map[string]interface{}
}

func NewEmailService(config emailServiceConfig) *emailService {
	return &emailService{
		apiKey:       config.ApiKey,
		templateData: map[string]map[string]interface{}{},
		to:           []emailAddress{},
	}
}

func (s *emailService) Clear() EmailServicer {
	s.from = emailAddress{}
	s.to = []emailAddress{}
	s.templateData = map[string]map[string]interface{}{}

	return s
}

func (s *emailService) From(email, name string) EmailServicer {
	s.from = emailAddress{
		Email: email,
		Name:  name,
	}

	return s
}

func (s *emailService) Send(templateID string) error {
	var (
		err         error
		request     rest.Request
		requestBody []byte
		response    *rest.Response
	)

	if s.from.Email == "" {
		return fmt.Errorf("from email address is required")
	}

	if len(s.to) == 0 {
		return fmt.Errorf("to email address is required")
	}

	templateRequest := emailDynamicTemplateRequest{
		From:             s.from,
		Personalizations: []emailPersonalizations{},
		TemplateID:       templateID,
	}

	if len(s.templateData) > 0 {
		for to, data := range s.templateData {
			templateRequest.Personalizations = append(templateRequest.Personalizations, emailPersonalizations{
				To: []emailAddress{
					{
						Email: to,
					},
				},
				TemplateData: data,
			})
		}
	}

	if requestBody, err = json.Marshal(templateRequest); err != nil {
		return fmt.Errorf("error parsing Send Grid dynamic template request body: %w", err)
	}

	request = sendgrid.GetRequest(s.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = http.MethodPost
	request.Body = requestBody

	if response, err = sendgrid.API(request); err != nil {
		return fmt.Errorf("error sending mail request to SendGrid: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 399 {
		return fmt.Errorf("error sending mail: %s", response.Body)
	}

	return nil
}

func (s *emailService) TemplateData(to string, data map[string]interface{}) EmailServicer {
	s.templateData[to] = data
	return s
}

func (s *emailService) To(email, name string) EmailServicer {
	s.to = append(s.to, emailAddress{
		Email: email,
		Name:  name,
	})

	return s
}

func (s *emailService) ToMultipleAddresses(emails []string) EmailServicer {
	for _, email := range emails {
		s.to = append(s.to, emailAddress{
			Email: email,
		})
	}

	return s
}
