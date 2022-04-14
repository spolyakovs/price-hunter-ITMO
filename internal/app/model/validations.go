package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var ValidationRulesPassword = []validation.Rule{
	validation.Required,
	validation.Length(8, 0),
	validation.Match(regexp.MustCompile(".*[A-Z].*")),
	validation.Match(regexp.MustCompile(".*[0-9].*")),
	validation.Match(regexp.MustCompile(".*[-!#$%&*+\\/=?^_{|}~].*")),
}

var ValidationRulesEmail = []validation.Rule{
	validation.Required,
	is.Email,
}

var ValidationRulesUsername = []validation.Rule{
	validation.Required,
	validation.Length(6, 30),
	validation.Match(regexp.MustCompile("^[a-zA-Z0-9_-]{6,30}$")),
}
