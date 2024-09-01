package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ginx-contribs/ginx/pkg/resp"
	localeen "github.com/go-playground/locales/en"
	unitrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transen "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

// Validator is wrapper of binding.StructValidator
type Validator interface {
	binding.StructValidator
	// HandleError judge that how to handle params validation errors
	HandleError(ctx *gin.Context, val any, err error)
}

func NewHumanizedValidator(v *validator.Validate, translator unitrans.Translator) *HumanizedValidator {
	return &HumanizedValidator{v: v, translator: translator}
}

// HumanizedValidator return human-readable validation result information
type HumanizedValidator struct {
	translator unitrans.Translator
	v          *validator.Validate
}

func (h *HumanizedValidator) HandleError(ctx *gin.Context, val any, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMsg []string
		for _, validateErr := range validationErrors {
			errorMsg = append(errorMsg, validateErr.Translate(h.translator))
		}
		resp.Fail(ctx).ErrorMsg(strings.Join(errorMsg, ",")).JSON()
		return
	}
	resp.Fail(ctx).ErrorMsg("params validate failed").JSON()
}

func (h *HumanizedValidator) ValidateStruct(a any) error {
	return h.v.Struct(a)
}

func (h *HumanizedValidator) Engine() any {
	return h.v
}

// SetValidator replace the default validator for binding packages
func SetValidator(structValidator binding.StructValidator) {
	binding.Validator = structValidator
}

// EnglishValidator create a validator can return human-readable parameters validation information with language english.
func EnglishValidator(v *validator.Validate) (*HumanizedValidator, error) {
	localeEn := localeen.New()
	universalTranslator := unitrans.New(localeEn)
	enTrans, _ := universalTranslator.GetTranslator(localeEn.Locale())
	err := transen.RegisterDefaultTranslations(v, enTrans)
	if err != nil {
		return nil, err
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		lookupNames := []string{"label", "form", "uri", "json", "yaml", "toml", "xml", "toml", "header"}
		for _, tag := range lookupNames {
			if name, ok := field.Tag.Lookup(tag); ok {
				return name
			}
		}
		return field.Name
	})
	return NewHumanizedValidator(v, enTrans), nil
}

// ValidateHandler will be called if validate failed.
type ValidateHandler func(ctx *gin.Context, val any, err error)

var defaultValidateHandler ValidateHandler

func SetValidateHandler(handler ValidateHandler) {
	defaultValidateHandler = handler
}

func ShouldValidateWith(ctx *gin.Context, val any, binding binding.Binding) error {
	err := ctx.ShouldBindWith(val, binding)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidate(ctx *gin.Context, val any) error {
	err := ctx.ShouldBind(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateJSON(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindJSON(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateQuery(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindQuery(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateURI(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindUri(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateHeader(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindHeader(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateXML(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindXML(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateYAML(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindYAML(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}

func ShouldValidateTOML(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindTOML(val)
	if err != nil {
		if defaultValidateHandler != nil {
			defaultValidateHandler(ctx, val, err)
		}
		return err
	}
	return nil
}
