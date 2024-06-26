// Code generated by ent, DO NOT EDIT.

package ent

import (
	"Backend/ent/schema"
	"Backend/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.NameValidator is a validator for the "name" field. It is called by the builders before save.
	user.NameValidator = userDescName.Validators[0].(func(string) error)
	// userDescAge is the schema descriptor for age field.
	userDescAge := userFields[1].Descriptor()
	// user.AgeValidator is a validator for the "age" field. It is called by the builders before save.
	user.AgeValidator = userDescAge.Validators[0].(func(int) error)
	// userDescSalary is the schema descriptor for salary field.
	userDescSalary := userFields[2].Descriptor()
	// user.SalaryValidator is a validator for the "salary" field. It is called by the builders before save.
	user.SalaryValidator = userDescSalary.Validators[0].(func(float64) error)
}
