package validators

import (
	"encoding/base64"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// ValidateBase64Image validates the base64 image field
func ValidateBase64Image(fl validator.FieldLevel) bool {
	base64Image := fl.Field().String()
	if base64Image == "" {
		// If the base64 image is empty, consider it valid
		return true
	}

	// Decode the base64 string
	_, err := base64.StdEncoding.DecodeString(base64Image)
	return err == nil
}

// ValidatePNGImage validates the base64 image field for PNG format
func ValidatePNGImage(fl validator.FieldLevel) bool {
	base64Image := fl.Field().String()
	if base64Image == "" {
		// If the base64 image is empty, consider it valid
		return true
	}

	// Decode the base64 string
	decoded, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return false
	}

	// Check if the decoded data starts with a valid PNG image signature
	return strings.HasPrefix(string(decoded), "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A") // PNG
}

// ValidateJPGImage validates the base64 image field for JPEG format
func ValidateJPGImage(fl validator.FieldLevel) bool {
	base64Image := fl.Field().String()
	if base64Image == "" {
		// If the base64 image is empty, consider it valid
		return true
	}

	// Decode the base64 string
	decoded, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return false
	}

	// Check if the decoded data starts with a valid JPEG image signature
	return strings.HasPrefix(string(decoded), "\xFF\xD8\xFF") // JPEG
}

func ValidateImageFormat(fl validator.FieldLevel) bool {
	imageData := fl.Field().String()
	if imageData == "" {
		// If the field is empty, consider it valid
		return true
	}

	return ValidatePNGImage(fl) || ValidateJPGImage(fl)
}

func validateDateTime(fl validator.FieldLevel) bool {
	// Get the datetime string from the field
	dateTimeStr := fl.Field().String()

	// Parse the datetime string
	_, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)

	if err == nil {
		return true
	}

	return false
}

func validateDate(fl validator.FieldLevel) bool {
	// Get the datetime string from the field
	dateTimeStr := fl.Field().String()

	// Parse the datetime string
	_, err := time.Parse("2006-01-02", dateTimeStr)

	if err != nil {
		return false
	}

	return true
}

func ValidateDateTimeRequiredIf(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		// If the date string is empty, it's considered valid (since it's not required)
		return true
	}

	tag := fl.Param()
	parts := strings.Split(tag, " ")

	// Access individual parts
	tagName := parts[0]
	tagValue := parts[1]

	// Ensure fl.Parent() is a reflect.Value of type reflect.Struct
	if fl.Parent().Kind() != reflect.Struct {
		// If fl.Parent() is not a struct, skip the validation
		return true
	}

	// Get the parent struct as a reflect.Value
	parent := fl.Parent()

	// Find the field named "RegistrationType"
	registrationTypeField, found := parent.Type().FieldByName(tagName)
	if !found {
		// If "RegistrationType" field not found, skip the validation
		return true
	}

	// Ensure the field type is what we expect
	if registrationTypeField.Type.Kind() != reflect.String {
		// If the field type is not a string, skip the validation
		return true
	}

	// Get the value of RegistrationType field
	fieldValue := parent.FieldByName(tagName)
	if !fieldValue.IsValid() {
		// If field value is not valid, skip the validation
		return true
	}
	if fieldValue.String() != tagValue {
		// If RegistrationType doesn't match the required value, skip the validation
		return true
	}

	// Perform datetime validatio
	return validateDateTime(fl)
}

func ValidateDateRequiredIf(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		// If the date string is empty, it's considered valid (since it's not required)
		return true
	}

	tag := fl.Param()
	parts := strings.Split(tag, " ")

	// Access individual parts
	tagName := parts[0]
	tagValue := parts[1]

	// Ensure fl.Parent() is a reflect.Value of type reflect.Struct
	if fl.Parent().Kind() != reflect.Struct {
		// If fl.Parent() is not a struct, skip the validation
		return true
	}

	// Get the parent struct as a reflect.Value
	parent := fl.Parent()

	// Find the field named "RegistrationType"
	registrationTypeField, found := parent.Type().FieldByName(tagName)
	if !found {
		// If "RegistrationType" field not found, skip the validation
		return true
	}

	// Ensure the field type is what we expect
	if registrationTypeField.Type.Kind() != reflect.String {
		// If the field type is not a string, skip the validation
		return true
	}

	// Get the value of RegistrationType field
	fieldValue := parent.FieldByName(tagName)
	if !fieldValue.IsValid() {
		// If field value is not valid, skip the validation
		return true
	}
	if fieldValue.String() != tagValue {
		// If RegistrationType doesn't match the required value, skip the validation
		return true
	}

	// Perform datetime validatio
	return validateDate(fl)
}

func validateBirthday(fl validator.FieldLevel) bool {
	birthday := fl.Field().String()
	// Check if the date matches the format YYYY-MM-DD
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, birthday)
	if !matched {
		return false
	}
	// Try to parse the date to ensure it's a valid date
	_, err := time.Parse("2006-01-02", birthday)
	return err == nil
}

func validateMonth(fl validator.FieldLevel) bool {
	month := fl.Field().String()
	matched, _ := regexp.MatchString("^(0[1-9]|1[0-2])$", month)
	return matched
}

func RegisterCustomValidators(validate *validator.Validate) {

	validate.RegisterValidation("base64image", ValidateBase64Image)
	validate.RegisterValidation("pngimage", ValidatePNGImage)
	validate.RegisterValidation("jpgimage", ValidateJPGImage)
	validate.RegisterValidation("imageformat", ValidateImageFormat)
	validate.RegisterValidation("datetime", validateDateTime)
	validate.RegisterValidation("datetime_if", ValidateDateTimeRequiredIf)
	validate.RegisterValidation("date", validateDate)
	validate.RegisterValidation("date_if", ValidateDateRequiredIf)
	validate.RegisterValidation("birthday", validateBirthday)
	validate.RegisterValidation("month", validateMonth)
}
