package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"tigerhallKittens/app/lib/logger"
	"tigerhallKittens/app/utils"
)

type Request struct {
	*http.Request
	pathParams map[string]string
	params     map[string]string
}

type ValidationErrorInterface interface {
	Type() string
	Error() string
	Unwrap() error
}

type ValidationError struct {
	errorType string
	message   string
	Err       error
}

func (e *ValidationError) Error() string { return e.message }

func (e *ValidationError) Type() string { return e.errorType }

func (e *ValidationError) Unwrap() error { return e.Err }

var (
	UnexpectedErr  = "Unexpected"
	ErrInvalidType = func(field string, expectedType interface{}, err error) ValidationErrorInterface {
		return &ValidationError{
			errorType: "InvalidType",
			message:   fmt.Sprintf("InvalidType for field: %s. Expected: %s", field, expectedType),
			Err:       err,
		}
	}
	ErrInvalidJson = func(err error) ValidationErrorInterface {
		return &ValidationError{
			errorType: "InvalidJson",
			message:   fmt.Sprintf("InvalidJson: %s", err.Error()),
			Err:       err,
		}
	}
	ErrInvalidValue = func(message string, err error) ValidationErrorInterface {
		if message == "" {
			message = err.Error()
		}
		return &ValidationError{
			errorType: "InvalidValue",
			message:   fmt.Sprintf("InvalidValue: %s", message),
			Err:       err,
		}
	}
)

func NewRequest(r *http.Request) Request {
	return Request{Request: r}
}

func (r *Request) SetPathParam(key, value string) {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}
	r.pathParams[key] = value
}

func (r *Request) GetPathParam(key string) string {
	if value, ok := r.pathParams[key]; ok {
		return value
	}
	return ""
}

func (r *Request) QueryParams() map[string]string {
	if r.params != nil {
		return r.params
	}
	r.params = map[string]string{}
	for key, val := range r.URL.Query() {
		r.params[key] = strings.Join(val, " | ")
	}
	return r.params
}

func (r *Request) Bind(v interface{}) error {
	return Bind(r.Context(), r.Request.Body, v)
}

func Bind(ctx context.Context, body io.ReadCloser, v interface{}) error {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.E(ctx, err, "error closing body")
		}
	}(body)

	decoder := json.NewDecoder(body)
	decodeErr := decoder.Decode(v)
	return decodeErr
}

func (r *Request) ParseAndValidateBody(s interface{}) error {
	if err := r.Bind(s); err != nil {
		return handleValidationErrors(err)
	}

	return validateStruct(s)
}

func validateStruct(s interface{}, structValidations ...validator.StructLevelFunc) ValidationErrorInterface {
	var validate = validator.New()

	_ = validate.RegisterValidation("notblank", validators.NotBlank)
	_ = validate.RegisterValidation("date", utils.ValidateDate)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	for _, sValidation := range structValidations {
		validate.RegisterStructValidation(sValidation, s)
	}

	if err := validate.Struct(s); err != nil {
		return handleValidationErrors(err)
	}

	return nil
}

func handleValidationErrors(err error) ValidationErrorInterface {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		return ErrInvalidType(e.Field, e.Type, e)
	case validator.ValidationErrors:
		var msgs []string
		for _, fe := range e {
			switch fe.Tag() {
			case "required":
				msgs = append(msgs, fmt.Sprintf("%s is a required field", fe.Field()))
			case "notblank":
				msgs = append(msgs, fmt.Sprintf("%s should not be empty", fe.Field()))
			case "max":
				msgs = append(msgs, fmt.Sprintf("%s must be a maximum of %s in length", fe.Field(), fe.Param()))
			case "url":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", fe.Field()))
			case "uuid":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid uuid", fe.Field()))
			case "date":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid date", fe.Field()))
			default:
				msgs = append(msgs, fmt.Sprintf("validation failed for %s on %s", fe.Field(), fe.Tag()))
			}
		}
		return ErrInvalidValue(strings.Join(msgs, ", "), e)
	case *json.SyntaxError:
		return ErrInvalidJson(e)
	default:
		return &ValidationError{
			errorType: UnexpectedErr,
			message:   e.Error(),
			Err:       e,
		}
	}
}
