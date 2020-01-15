/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package validator

import (
	"errors"
	"fmt"
	"os"

	schema "github.com/xeipuuv/gojsonschema"
)

const (
	// SchemaDefinitions is the location of the repository
	// of schema definitions on the filepath
	SchemaDefinitions = "./artifacts/json-schema"

	// DefinitionFile is the schema definition file, or the root
	// file from which a definition beings
	DefinitionFile = "schema.json"

	// CurrentVersion represents the default
	CurrentVersion = "1.0.0"
)

var (
	// ErrSchemaVersionUndefined is returned if a validator is requested
	// and no corresponding schema definition that matches the version is present
	// in the SchemaDefinitions artifact repository
	ErrSchemaVersionUndefined = errors.New("schema version undefined")

	// ErrValidationFailed is returned by Validate if the supplied document
	// does not pass validation
	ErrValidationFailed = errors.New("validation failed")
)

// Validator is a validator for JSON documents
type Validator struct {
	Version        string
	Valid          bool
	DefinitionFile string

	definition schema.JSONLoader
	result     *schema.Result
}

// NewValidator will return a new validator configured to
// use the CurrentVersion of the schema definition
func NewValidator() (*Validator, error) {
	return NewValidatorByVersion(CurrentVersion)
}

// NewValidatorByVersion will return a new validator with the definition
// for the supplied version number loaded, as defined in the SchemaDefinitions
// repository and ready for a document to be validated against
func NewValidatorByVersion(version string) (*Validator, error) {
	path, err := schemaPath(version)
	if err != nil {
		return nil, err
	}

	v := Validator{
		Version:        version,
		DefinitionFile: path,
		definition:     schema.NewReferenceLoader(fmt.Sprintf("file://%s", path)),
	}

	return &v, nil
}

// Validate will perform validation of the supplied document
// against the validators supplied defintion, returning the error
// ErrValidationFailed in the case of a failed validation
func (v *Validator) Validate(document []byte) error {
	doc := schema.NewStringLoader(string(document))

	result, err := schema.Validate(v.definition, doc)

	if err != nil {
		return err
	}

	v.result = result
	v.Valid = result.Valid()

	if v.Valid != true {
		return ErrValidationFailed
	}

	return nil
}

// Errors returns any error objects from a failed validation
func (v *Validator) Errors() []schema.ResultError {
	if v != nil {
		return v.result.Errors()
	}

	return []schema.ResultError{}
}

// schemaPath determines the path for the requested definition version
// if defined in the SchemaDefinitions repository. ErrSchemaVersionUndefined
// is returned if a corresponding definition by version is not present
func schemaPath(version string) (string, error) {
	schemaPath := fmt.Sprintf(
		"%s/%s/%s",
		SchemaDefinitions,
		version,
		DefinitionFile,
	)

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return "", ErrSchemaVersionUndefined
	}

	return schemaPath, nil
}
