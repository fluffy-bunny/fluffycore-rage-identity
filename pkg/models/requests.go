package models

import "encoding/json"

type LoginGetRequest struct {
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
}

func ConvertFromInterface[T any](value interface{}, out *T) error {
	// ...
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var result T
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	*out = result
	return nil
}
