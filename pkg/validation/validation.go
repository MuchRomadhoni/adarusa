package validation

import (
	"fmt"
	"fp_pinjaman_online/model/dto/json"
	adminEntity "fp_pinjaman_online/model/entity/admin"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/stoewer/go-strcase"
	"golang.org/x/crypto/bcrypt"
)

func GetValidationError(err error) []json.ValidationField {
	var validationFields []json.ValidationField
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range ve {
			log.Debug().Msg(fmt.Sprintf("validationError : %v", validationError))
			myField := convertFieldRequired(validationError.Namespace())
			validationFields = append(validationFields, json.ValidationField{
				FieldName: myField,
				Message:   formatMessage(validationError),
			})
		}
	}
	return validationFields
}

func convertFieldRequired(myValue string) string {
	log.Debug().Msg("convertFieldRequired: " + myValue)
	fieldSegmen := strings.Split(myValue, ".")
	myField := ""
	length := len(fieldSegmen)
	i := 1
	for _, val := range fieldSegmen {
		if i == 1 {
			i++
			continue
		}

		if i == length {
			myField += strcase.SnakeCase(val)
			break
		}

		myField += strcase.LowerCamelCase(val) + `/`
		i++
	}

	return myField
}

func formatMessage(err validator.FieldError) string {
	var message string

	switch err.Tag() {
	case "required":
		message = "required"
	case "number":
		message = "must be number"
	case "email":
		message = "invalid format email"
	case "DateOnly":
		message = "invalid format date"
	case "min":
		message = "minimum value is not exceed"
	case "max":
		message = "max value is exceed"
	case "oneof":
		message = "value must be one of " + err.Param()
	case "datetime":
		message = "datetime format is invalid"
	case "password":
		message = "password must 1 lowercase, 1 uppercase, 1 numeric, 1 special character"
	case "numeric":
		message = "must be numeric"
	case "len":
		message = "length must be " + err.Param()
	}

	return message
}

func ValidationPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashAndPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidateUserComplete(user adminEntity.UserCompleteInfo) bool {
	return user.AccountNumber.Valid && user.AccountNumber.String != "" &&
		user.BankName.Valid && user.BankName.String != "" &&
		user.EmergencyContact.Valid && user.EmergencyContact.String != "" &&
		user.EmergencyPhone.Valid && user.EmergencyPhone.String != "" &&
		user.JobName.Valid && user.JobName.String != "" &&
		user.OfficeName.Valid && user.OfficeName.String != "" &&
		user.NIK.Valid && user.NIK.String != "" &&
		user.FullName.Valid && user.FullName.String != "" &&
		user.PersonalPhoneNumber.Valid && user.PersonalPhoneNumber.String != "" &&
		user.PersonalAddress.Valid && user.PersonalAddress.String != "" &&
		user.City.Valid && user.City.String != "" &&
		user.FotoKTP.Valid && user.FotoKTP.String != "" &&
		user.FotoSelfie.Valid && user.FotoSelfie.String != "" &&
		user.Email != "" && user.Gaji.Float64 != 0
}
