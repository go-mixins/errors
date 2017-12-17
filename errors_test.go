package errors_test

import (
	errs "errors"
	"testing"

	"github.com/go-mixins/errors"
)

func Test_NewClass(t *testing.T) {
	c1 := errors.NewClass("root")
	c2 := c1.Sub("leaf")
	if c1.Is(c2) {
		t.Errorf("%v should be subclass of %v", c1, c2)
	}
	if !c2.Is(c1) {
		t.Errorf("%v should not be subclass of %v", c1, c2)
	}
}

func TestNew(t *testing.T) {
	err := errors.New("simple error")
	if _, ok := err.(error); !ok {
		t.Errorf("%v should implement error interface", err)
	}
}

func TestClass_Contains(t *testing.T) {
	root := errors.NewClass("root")
	leaf := errors.NewClass("root", "leaf")
	rootError := root.Wrap(errs.New("some error"), "in root")
	leafError := leaf.Wrap(errs.New("some other error"), "in leaf")
	outsideError := errs.New("some outside error")
	if !root.Contains(leafError) {
		t.Errorf("%v should belong to %v", leafError, root)
	}
	if !root.Contains(rootError) {
		t.Errorf("%v should belong to %v", rootError, root)
	}
	if leaf.Contains(rootError) {
		t.Errorf("%v should not belong to %v", rootError, leaf)
	}
	if !leaf.Contains(leafError) {
		t.Errorf("%v should belong to %v", leafError, leaf)
	}
	if root.Contains(outsideError) {
		t.Errorf("%v should not belong to %v", outsideError, root)
	}
	if leaf.Contains(outsideError) {
		t.Errorf("%v should not belong to %v", outsideError, leaf)
	}
}

func TestClass_ErrorsCause(t *testing.T) {
	root := errors.NewClass("root")
	err := errs.New("some error")
	rootError := root.Wrap(err, "in root")
	if errors.Cause(rootError) != err {
		t.Errorf("%v should be the cause of %v", err, rootError)
	}
}

func TestClass_New_Errorf_Wrap_Wrapf(t *testing.T) {
	c1 := errors.NewClass("root")
	err := c1.New("root error 0")
	if !c1.Contains(err) {
		t.Errorf("%v should belong to %v", err, c1)
	}
	err = errs.New("some other error")
	if c1.Contains(err) {
		t.Errorf("%v should not belong to %v", err, c1)
	}
	err = c1.Errorf("root error %d", 1)
	if !c1.Contains(err) {
		t.Errorf("%v should belong to %v", err, c1)
	}
	err = c1.Wrap(errs.New("root error 2"), "wrapped")
	if !c1.Contains(err) {
		t.Errorf("%v should belong to %v", err, c1)
	}
	err = c1.Wrapf(errs.New("root error"), "wrapped %d", 3)
	if !c1.Contains(err) {
		t.Errorf("%v should belong to %v", err, c1)
	}
}
