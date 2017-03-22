package passwordreset

type PasswordReset struct {
	ID        string
	Reference string
	PartnerID string
}

func New(id string, reference string, partnerid string) *PasswordReset {
	return &PasswordReset{
		ID:        id,
		Reference: reference,
		PartnerID: partnerid,
	}
}
