package internalerror

import "testing"

func TestCEPInvalidError(t *testing.T) {
	err := CEPInvalidError()

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Code != 422 {
		t.Errorf("Expected code 422, got %d", err.Code)
	}
	if err.MSG != "Invalid zipcode" {
		t.Errorf("Expected message 'Invalid zipcode', got '%s'", err.MSG)
	}
}

func TestCEPNotFoundError(t *testing.T) {
	err := CEPNotFoundError()

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Code != 404 {
		t.Errorf("Expected code 404, got %d", err.Code)
	}
	if err.MSG != "Can not find zipcode" {
		t.Errorf("Expected message 'Can not find zipcode', got '%s'", err.MSG)
	}
}

func TestInternalError_Consistency(t *testing.T) {
	err1 := CEPInvalidError()
	err2 := CEPInvalidError()

	if err1 == err2 {
		t.Error("Expected different pointers, got same pointer")
	}
	if err1.Code != err2.Code {
		t.Errorf("Expected consistent codes, got %d and %d", err1.Code, err2.Code)
	}
	if err1.MSG != err2.MSG {
		t.Errorf("Expected consistent messages, got '%s' and '%s'", err1.MSG, err2.MSG)
	}
}

func TestInternalError_HTTPSemantics(t *testing.T) {
	t.Run("422 means unprocessable entity", func(t *testing.T) {
		err := CEPInvalidError()
		if err.Code != 422 {
			t.Errorf("CEPInvalidError should return 422, got %d", err.Code)
		}
	})

	t.Run("404 means not found", func(t *testing.T) {
		err := CEPNotFoundError()
		if err.Code != 404 {
			t.Errorf("CEPNotFoundError should return 404, got %d", err.Code)
		}
	})
}

func TestInternalError_NonNilReturn(t *testing.T) {
	errorFuncs := []func() *InternalError{
		CEPInvalidError,
		CEPNotFoundError,
	}

	for i, fn := range errorFuncs {
		err := fn()
		if err == nil {
			t.Errorf("Error function %d should never return nil", i)
		}
	}
}
