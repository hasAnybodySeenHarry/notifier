package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(payload []byte, dest any) error {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&dest)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var destinationErr *json.InvalidUnmarshalError

		switch {
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed payload")
		case errors.Is(err, io.EOF):
			return errors.New("payload not found")
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("syntax error at char %d", syntaxErr.Offset)
		case err.Error() == "http: request body too large":
			return errors.New("payload too large")
		case errors.As(err, &typeErr):
			return errors.New("unmatched type in payload")
		case errors.As(err, &destinationErr):
			return errors.New("destination is not a pointer")
		case strings.HasPrefix(err.Error(), "json: unknown fields "):
			return fmt.Errorf("unknown field %s", strings.TrimPrefix(err.Error(), "json: unknown fields "))
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("multiple payloads detected")
	}

	return nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Println(fmt.Errorf("%s", err))
			}
		}()

		fn()
	}()
}
