package middleware

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var setupValidationOnce sync.Once

func SetupValidation() {
	setupValidationOnce.Do(func() {
		validatorEngine, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}

		validatorEngine.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.Split(field.Tag.Get("json"), ",")[0]
			if name == "" {
				return field.Name
			}
			if name == "-" {
				return ""
			}

			return name
		})
	})
}
