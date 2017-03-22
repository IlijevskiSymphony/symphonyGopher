package configuration

var SymphonyGopherDB = &EnvironmentVariable{"SG_DB", isNotBlank}
var SymphonyGopherPort = &EnvironmentVariable{"SG_PORT", isNumber}
var HashSalt = &EnvironmentVariable{"SG_HASHSALT", isNotBlank}

//mandril api variables
// var MandrillApiKey = &EnvironmentVariable{"SG_MANDRILLAPIKEY", isNotBlank}
// var MandrillHost = &EnvironmentVariable{"SG_MANDRILLHOST", isNotBlank}

type Configuration struct {
	SymphonyGopherDB   string
	SymphonyGopherPort string
	HashSalt           string
	// MandrillApiKey     string
	// MandrillHost       string
}

func Read() (*Configuration, []error) {
	r := NewEnvReader()

	c := Configuration{}
	c.SymphonyGopherDB = r.Read(SymphonyGopherDB)
	c.SymphonyGopherPort = r.Read(SymphonyGopherPort)
	c.HashSalt = r.Read(HashSalt)

	//mandril api variables
	// c.MandrillApiKey = r.Read(MandrillApiKey)
	// c.MandrillHost = r.Read(MandrillHost)

	return &c, r.Errors
}
