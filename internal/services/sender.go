package services

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type MessageSender interface {
	Send(to string, message string) error
}

type TwilioSender struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewTwilioSender(accountSID string, authToken string, fromNumber string) MessageSender {
	if accountSID == "" || authToken == "" || fromNumber == "" {
		log.Fatal("Twilio 설정 오류: 필수 환경 변수가 설정되지 않았습니다")
	}
	return &TwilioSender{
		AccountSID: accountSID,
		AuthToken:  authToken,
		FromNumber: fromNumber,
	}
}

func (t *TwilioSender) Send(to string, message string) error {
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSID)

	//twilio 번호 형식 변환
	if len(to) > 0 && to[0] == '0' {
		to = "+82" + to[1:]
	}

	data := url.Values{}
	data.Set("To", to)
	data.Set("From", t.FromNumber)
	data.Set("Body", message)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.AccountSID, t.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("문자 전송 실패: status %s", resp.Status)
	}

	return nil
}
