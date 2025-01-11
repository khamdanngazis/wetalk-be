package helper

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func RegisterTranslator(validate *validator.Validate) ut.Translator {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	// Register default translations
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Register custom translations
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email address", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	return trans
}

func GetMessageValidator(validate *validator.Validate, err error) string {
	trans := RegisterTranslator(validate)
	erroMessage := ""
	for _, err := range err.(validator.ValidationErrors) {
		erroMessage = err.Translate(trans)
	}
	return erroMessage
}
