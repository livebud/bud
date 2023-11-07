package req

import "net/http"

func Form(r *http.Request, in any) error {
	return nil
}

func Parse(r *http.Request, in any) error {
	return nil
}

func Json(r *http.Request, in any) error {
	return nil
}

func Unmarshal(r *http.Request, in any) error {
	return nil
}

func Accept(r *http.Request, mimes ...string) bool {
	return true
}

type Sessions struct {
}

func (s *Sessions) Load(r *http.Request, session any) error {
	return nil
}

func (s *Sessions) Save(w http.ResponseWriter, session any) error {
	return nil
}
