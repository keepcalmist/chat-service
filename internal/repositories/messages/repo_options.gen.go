// Code generated by options-gen. DO NOT EDIT.
package messagesrepo

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"github.com/keepcalmist/chat-service/internal/store"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	db *store.Database,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.db = db

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("db", _validate_Options_db(o)))
	return errs.AsError()
}

func _validate_Options_db(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.db, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `db` did not pass the test: %w", err)
	}
	return nil
}
