package pokedex

import (
	"errors"
	"strconv"
)

var (
	ErrMissingResources  = errors.New("one of id or name must be provided")
	ErrMultipleResources = errors.New("id and name cannot both be provided")
)

type GetRequest struct {
	ID   int
	Name string
}

func (r *GetRequest) GetResource() (string, error) {
	if r.ID == 0 && r.Name == "" {
		return "", NewError(ErrMissingResources.Error(), CodeInvalidArgs, nil)
	}
	if r.ID != 0 && r.Name != "" {
		return "", NewError(ErrMultipleResources.Error(), CodeInvalidArgs, nil)
	}

	if r.ID != 0 {
		return strconv.Itoa(r.ID), nil
	}
	return r.Name, nil
}

type ListRequest struct {
	PageSize uint
}
