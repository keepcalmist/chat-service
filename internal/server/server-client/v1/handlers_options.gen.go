// Code generated by options-gen. DO NOT EDIT.
package clientv1

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	getHistory getHistoryUseCase,
	sendMsg sendMessageUseCase,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.getHistory = getHistory
	o.sendMsg = sendMsg

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("getHistory", _validate_Options_getHistory(o)))
	errs.Add(errors461e464ebed9.NewValidationError("sendMsg", _validate_Options_sendMsg(o)))
	return errs.AsError()
}

func _validate_Options_getHistory(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.getHistory, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `getHistory` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_sendMsg(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.sendMsg, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `sendMsg` did not pass the test: %w", err)
	}
	return nil
}
