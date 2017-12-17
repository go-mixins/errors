package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Class represents error hierarchy path
type Class []string

// Sub creates subpath in the class hierarchy
func (c Class) Sub(path ...string) (res Class) {
	res = make([]string, len(c)+len(path))
	copy(res, c)
	copy(res[len(c):], path)
	return
}

func (c Class) String() string {
	return strings.Join(c, "/")
}

// NewClass creates independent new path class hierarchy
func NewClass(path ...string) Class {
	return Class(nil).Sub(path...)
}

// Cause simply calls original errors.Cause function
func Cause(err error) error {
	return errors.Cause(err)
}

type classer interface {
	Class() Class
}

type causer interface {
	Cause() error
}

// Is returns true if the class belongs to specific parent class
func (c Class) Is(parent Class) bool {
	if len(parent) > len(c) {
		return false
	}
	for i := range parent {
		if parent[i] != c[i] {
			return false
		}
	}
	return true
}

// Contains is true if the error belongs to certain class
func (c Class) Contains(err error) bool {
	if e, ok := err.(classer); ok {
		return e.Class().Is(c)
	}
	if e, ok := err.(causer); ok {
		return c.Contains(e.Cause())
	}
	return false
}

type simpleError string

func (se simpleError) Error() string {
	return string(se)
}

type errorOfClass struct {
	class Class
	error
}

func (ec errorOfClass) Error() string {
	return strings.Join(ec.class, ": ") + ": " + ec.error.Error()
}

func (ec errorOfClass) Cause() error {
	return ec.error
}

func (ec errorOfClass) Class() Class {
	return ec.class
}

// Wrap marks the error with certain class and wraps it using errors package
func (c Class) Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(errorOfClass{c, err}, message)
}

// Wrapf marks the error with certain class and wraps it using errors package
func (c Class) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(errorOfClass{c, err}, format, args...)
}

// New returns an error with the supplied message and class
func (c Class) New(message string) error {
	return errors.WithStack(errorOfClass{c, simpleError(message)})
}

// Errorf returns an error formatted against supplied format
func (c Class) Errorf(format string, args ...interface{}) error {
	return errors.WithStack(errorOfClass{c, fmt.Errorf(format, args...)})
}
