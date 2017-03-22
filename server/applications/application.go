package applications

type applicationType string

const (
	Native                   = applicationType("native")
	SinglePageWebApplication = applicationType("single-page-web-application")
	RegularWebApplication    = applicationType("regular-web-application")
	NonInteractiveClient     = applicationType("non-interactive-client")
  TypeInvalid              = applicationType("invalid")
)

type Status string

const (
  Active   = Status("active")
  Inactive = Status("inactive")
)

type Application struct {
  ID          string
  Name        string
  Description string
  Type        applicationType
  PartnerID   string
  Status      Status
}

func New (id string, name string, description string, appType applicationType, partnerId string) *Application {
  return &Application{
    ID:              id,
    Name:            name,
    Description:     description,
    Type:            appType,
    PartnerID:       partnerId,
    Status:          Inactive,
  }
}

func ApplicationType(value string) (applicationType, bool) {
  switch applicationType(value) {
  case Native:
    return Native, true
  case SinglePageWebApplication:
    return SinglePageWebApplication, true
  case RegularWebApplication:
    return RegularWebApplication, true
  case NonInteractiveClient:
    return NonInteractiveClient, true
  }
  return TypeInvalid, false
}
