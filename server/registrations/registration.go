package registrations

const (
	StatusVerifiactionNotSent = status("verification-not-sent")
	StatusVerificationSent    = status("verification-sent")
)

type status string

type registration struct {
	ID        string
	Reference string
	PartnerID string
	Status    status
}

func New(id string, reference string, partnerID string) *registration {
	return &registration{
		ID:        id,
		Reference: reference,
		PartnerID: partnerID,
		Status:    StatusVerifiactionNotSent,
	}
}
