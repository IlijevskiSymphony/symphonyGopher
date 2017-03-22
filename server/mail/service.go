package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Service struct {
	Host string
}

type response struct {
	Email        string `json:"email"`
	Status       string `json:"status"`
	ID           string `json:"_id"`
	RejectReason string `json:"reject_reason"`
}

func (s *Service) Send(mail Mail) error {
	body, err := json.Marshal(mail)
	if err != nil {
		return errors.Wrap(err, "Cannot serialize mail json")
	}

	req, err := http.NewRequest("POST", s.Host+"/api/1.0/messages/send.json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "Cannot create POST mail request")
	}

	post := &http.Client{}
	resp, err := post.Do(req)
	if err != nil {
		return errors.Wrap(err, "Cannot post to client")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Invalid mandrill response: %+v", resp))
	}

	decoder := json.NewDecoder(resp.Body)
	var result []response
	if err := decoder.Decode(&result); err != nil {
		return errors.Wrap(err, "Cannot decode response.")
	}

	if len(result) != 1 || result[0].Status != "sent" {
		return errors.New(fmt.Sprintf("Invalid response: %+v", result))
	}
	return nil
}
