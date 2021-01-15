package rpc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateRequest(r *http.Request, schema gojsonschema.JSONLoader) error {
	compiledSchema, err := gojsonschema.NewSchemaLoader().Compile(schema)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return io.ErrUnexpectedEOF
		}

		return fmt.Errorf("crpc failed to read request body: %w", err)
	}

	ld := gojsonschema.NewBytesLoader(body)

	result, err := compiledSchema.Validate(ld)
	if err != nil {
		return fmt.Errorf("crpc schema validation failed: %w", err)
	}

	err = coerceJSONSchemaError(result)
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	return nil
}

func coerceJSONSchemaError(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	var reasons []cher.E

	errs := result.Errors()
	for _, err := range errs {
		reasons = append(reasons, cher.E{
			Code: "schema_failure",
			Meta: cher.M{
				"field":   err.Field(),
				"type":    err.Type(),
				"message": err.Description(),
			},
		})
	}

	return cher.New(cher.BadRequest, nil, reasons...)
}
