package partners

import (
	"crypto/sha512"
	"encoding/base64"
	"regexp"
	"unicode"
)

type Status string

const (
	StatusNotVerified = Status("not-verified")
	StatusVerified    = Status("verified")
	StatusActive      = Status("active")
)

type companySize string

const (
	CompanySizeOneToFife         = companySize("1-5")
	CompanySizeSixToTen          = companySize("6-10")
	CompanySizeElevenToFifty     = companySize("11-50")
	CompanySizeFiftyoneToHunderd = companySize("51-100")
	CompanySizeInvalid           = companySize("invalid")
)

type language string

const (
	LanguageNodeJs  = language("node-js")
	LanguageJava    = language("java")
	LanguageGo      = language("go-lang")
	LanguageInvalid = language("invalid")
)

const EmailFormatRegexp = `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
const (
	PasswordMinLength = 8
	PasswordMaxLength = 15
	UpperCaseLettersMinCount = 1
	LowerCaseLettersMinCount = 1
	NumeralsMinCount = 1
)

type Partner struct {
	ID           string
	Email        string
	PasswordHash string
	Status       Status
	Name         string
	Company      string
	CompanySize  companySize
	Language     language
}

func New(id string, email string, password string, salt string) *Partner {
	return &Partner{
		ID:           id,
		Email:        email,
		PasswordHash: Hash(password + salt),
		Status:       StatusNotVerified,
	}
}

func Hash(password string) string {
	var result = sha512.Sum512([]byte(password))
	return base64.URLEncoding.EncodeToString(result[:64])
}

func CompanySize(value string) (companySize, bool) {
	switch companySize(value) {
	case CompanySizeOneToFife:
		return CompanySizeOneToFife, true
	case CompanySizeSixToTen:
		return CompanySizeSixToTen, true
	case CompanySizeElevenToFifty:
		return CompanySizeElevenToFifty, true
	case CompanySizeFiftyoneToHunderd:
		return CompanySizeFiftyoneToHunderd, true
	}
	return CompanySizeInvalid, false
}

func Language(value string) (language, bool) {
	switch language(value) {
	case LanguageNodeJs:
		return LanguageNodeJs, true
	case LanguageJava:
		return LanguageJava, true
	case LanguageGo:
		return LanguageGo, true
	}
	return LanguageInvalid, false
}

func ValidateEmailFormat(value string) bool {
	var validEmailRegexp = regexp.MustCompile(EmailFormatRegexp)
	return validEmailRegexp.MatchString(value)
}

func ValidatePasswordFormat(value string) bool {
	if len(value) < PasswordMinLength || len(value) > PasswordMaxLength {
		return false
	}

	upperCaseCount := 0
	lowerCaseCount := 0
	numeralsCount := 0
	specialsCount := 0

	for _, char := range value {
		switch {
			case unicode.IsUpper(char):
				upperCaseCount++
			case unicode.IsLower(char):
				lowerCaseCount++
			case unicode.IsNumber(char):
				numeralsCount++
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				specialsCount++
			case char == ' ':
				return false
			default:

		}
	}

	if upperCaseCount < UpperCaseLettersMinCount || lowerCaseCount < LowerCaseLettersMinCount || numeralsCount < NumeralsMinCount {
		return false
	}

	return true
}
