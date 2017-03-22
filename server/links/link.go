package links

type Link struct {
	ID        string
	PartnerID string
	Link      string
	Tags      []string
}

func New(id string, partnerID string, link string, tags []string) *Link {
	return &Link{
		ID:        id,
		PartnerID: partnerID,
		Link:      link,
		Tags:      tags,
	}
}
