// Copyright 2020 Job Stoit. All rights reseved.

// Package form reads the struct
//
// The form package has
package form

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/jobstoit/strct"
)

// Post posts the structure as a form and
func Post(url string, obj interface{}, headers map[string]string) (*http.Response, error) {
	body, ct, err := Marshal(obj)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	req.Header.Add(`Content-Type`, ct)

	client := new(http.Client)
	return client.Do(req)
}

// Marshal returns the form body and contentType
func Marshal(obj interface{}) (body io.Reader, contentType string, xerr error) {
	x := new(bytes.Buffer)
	form := multipart.NewWriter(x)

	err := strct.Scan(obj, func(field reflect.StructField, value *reflect.Value) error {
		name := field.Tag.Get(`form`)
		if name == `` {
			name = strings.ToLower(field.Name)
		}

		if value.Type() == reflect.TypeOf(new(os.File)).Elem() {
			file := value.Interface().(*os.File)
			wr, err := form.CreateFormFile(field.Name, file.Name())
			if err != nil {
				return err
			}

			if _, err := io.Copy(wr, file); err != nil {
				return err
			}
		} else {
			return form.WriteField(field.Name, fmt.Sprint(value.Interface()))
		}
		return nil
	})
	if err != nil {
		xerr = err
		return
	}

	if err = form.Close(); err != nil {
		xerr = err
		return
	}

	return x, form.FormDataContentType(), nil
}

// Unmarshal
