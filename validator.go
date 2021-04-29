package semanticid

import "github.com/go-playground/validator/v10"

func SemanticIDCollectionValidation(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(SemanticID)
	if value.IsNil() {
		return false
	}

	collection := fl.Param()
	return value.Collection == collection
}

func RegisterValidation(v *validator.Validate) {
	v.RegisterValidation("sidcollection", SemanticIDCollectionValidation)
}
