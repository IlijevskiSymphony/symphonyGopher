package configuration

var DeveloperDashboardDB = &EnvironmentVariable{"BQ_DD_DB", isNotBlank}
var Port = &EnvironmentVariable{"BQ_DD_PORT", isNumber}
var ElasticSearchURL = &EnvironmentVariable{"BC_ELASTICSEARCHURL", isNotBlank}
var HashSalt = &EnvironmentVariable{"BQ_DD_HASHSALT", isNotBlank}
var MandrillApiKey = &EnvironmentVariable{"BQ_DD_MANDRILLAPIKEY", isNotBlank}
var MandrillHost = &EnvironmentVariable{"BQ_DD_MANDRILLHOST", isNotBlank}
var StaticContentDir = &EnvironmentVariable{"BQ_DD_STATICCONTENTDIR", isNotBlank}

type Configuration struct {
	DeveloperDashboardDB string
	ElasticSearchURL     string
	Port                 string
	HashSalt             string
	MandrillApiKey       string
	MandrillHost         string
	StaticContentDir     string
}

func Read() (*Configuration, []error) {
	r := NewEnvReader()

	c := Configuration{}
	c.DeveloperDashboardDB = r.Read(DeveloperDashboardDB)
	c.Port = r.Read(Port)
	c.ElasticSearchURL = r.Read(ElasticSearchURL)
	c.HashSalt = r.Read(HashSalt)
	c.MandrillApiKey = r.Read(MandrillApiKey)
	c.MandrillHost = r.Read(MandrillHost)
	c.StaticContentDir = r.Read(StaticContentDir)

	return &c, r.Errors
}
