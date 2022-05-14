package model_test

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

func TestValidationRulesPassword(t *testing.T) {
	passwordEmpty := ""
	passwordShort := "qwer"
	passwordSimple := "qwertyqweqwe"
	passwordCorrect := "Qwert_y_1"

	if err := validation.Validate(&passwordEmpty, model.ValidationRulesPassword...); err == nil {
		t.Error("Empty password was accepted")
	}
	if err := validation.Validate(&passwordShort, model.ValidationRulesPassword...); err == nil {
		t.Errorf("Short password (%s) was accepted", passwordShort)
	}
	if err := validation.Validate(&passwordSimple, model.ValidationRulesPassword...); err == nil {
		t.Errorf("Simple password (%s) was accepted", passwordSimple)
	}
	if err := validation.Validate(&passwordCorrect, model.ValidationRulesPassword...); err != nil {
		t.Errorf("Correct password (%s) wasn't accepted:\n\t%s", passwordCorrect, err.Error())
	}
}

func TestValidationRulesEmail(t *testing.T) {
	emailEmpty := ""
	emailShort := "qwer"
	emailCorrect := "test@example.com"

	if err := validation.Validate(&emailEmpty, model.ValidationRulesEmail...); err == nil {
		t.Error("Empty email was accepted")
	}
	if err := validation.Validate(&emailShort, model.ValidationRulesEmail...); err == nil {
		t.Errorf("Short email (%s) was accepted", emailShort)
	}
	if err := validation.Validate(&emailCorrect, model.ValidationRulesEmail...); err != nil {
		t.Errorf("Correct email (%s) wasn't accepted:\n\t%s", emailCorrect, err.Error())
	}
}

func TestValidationRulesUsername(t *testing.T) {
	usernameEmpty := ""
	usernameShort := "qwer"
	usernameLong := "qwerashgdfhsagfdysatfdhgsafdhgsafdhgsafdhgfsahdgfashgdfahsfdhsagfdhgsa"
	usernameIncorrect := "qywterwq^:1asdsa"
	usernameCorrect := "test_username"

	if err := validation.Validate(&usernameEmpty, model.ValidationRulesUsername...); err == nil {
		t.Error("Empty username was accepted")
	}
	if err := validation.Validate(&usernameShort, model.ValidationRulesUsername...); err == nil {
		t.Errorf("Short username (%s) was accepted", usernameShort)
	}
	if err := validation.Validate(&usernameLong, model.ValidationRulesUsername...); err == nil {
		t.Errorf("Long username (%s) was accepted", usernameLong)
	}
	if err := validation.Validate(&usernameShort, model.ValidationRulesUsername...); err == nil {
		t.Errorf("Incorrect username (%s) was accepted", usernameIncorrect)
	}
	if err := validation.Validate(&usernameCorrect, model.ValidationRulesUsername...); err != nil {
		t.Errorf("Correct username (%s) wasn't accepted:\n\t%s", usernameCorrect, err.Error())
	}
}
