package tests

import (
	"prac/handlers"
	"testing"
)

func TestDelHandler(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		initialData map[string]handlers.CacheItem
		expectError bool
		expectedErr string
	}{
		{
			name:        "Key exists, successful deletion",
			args:        []string{"testKey"},
			initialData: map[string]handlers.CacheItem{"testKey": {Val: "testValue"}},
			expectError: false,
		},
		{
			name:        "Key doesn't exist",
			args:        []string{"nonExistentKey"},
			initialData: map[string]handlers.CacheItem{},
			expectError: true,
			expectedErr: "DEL nonExistentKey : Key doesn't exist !!!",
		},
		{
			name:        "Missing Key",
			args:        []string{},
			initialData: map[string]handlers.CacheItem{},
			expectError: true,
			expectedErr: "DEL : Missing Key",
		},
	}

	for _, test := range tests {
		// Initialize cache state
		handlers.PlainCache.Data = test.initialData

		t.Run(test.name, func(t *testing.T) {
			err := handlers.DelHandler(test.args)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != test.expectedErr {
					t.Errorf("Expected error: %v, but got: %v", test.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}

				if _, exists := handlers.PlainCache.Data[test.args[0]]; exists {
					t.Errorf("Expected key to be deleted, but it still exists")
				}
			}
		})
	}
}

func TestSetHandler(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectError bool
		expectedErr string
	}{
		{
			name:        "Successful Set",
			input:       []string{"testKey", "testValue"},
			expectError: false,
		},
		{
			name:        "Missing Key",
			input:       []string{},
			expectError: true,
			expectedErr: "SET : Missing Key and Value",
		},
		{
			name:        "Missing Value",
			input:       []string{"testKey"},
			expectError: true,
			expectedErr: "SET testKey: Add value as well !!!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := handlers.SetHandler(test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != test.expectedErr {
					t.Errorf("Expected error: %v, but got: %v", test.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if handlers.PlainCache.Data[test.input[0]].Val != test.input[1] {
					t.Errorf("Expected value %v, but got %v", test.input[1], handlers.PlainCache.Data[test.input[0]])
				}
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		initialData map[string]handlers.CacheItem
		expectError bool
		expectedErr string
		expectedVal string
	}{
		{
			name:        "Successful Get",
			input:       []string{"testKey"},
			initialData: map[string]handlers.CacheItem{"testKey": {Val: "testValue"}},
			expectError: false,
			expectedVal: "testValue",
		},
		{
			name:        "Key doesn't exist",
			input:       []string{"nonExistentKey"},
			initialData: map[string]handlers.CacheItem{},
			expectError: true,
			expectedErr: "GET nonExistentKey: Key doesn't exist!!!",
		},
		{
			name:        "Missing Key",
			input:       []string{},
			initialData: map[string]handlers.CacheItem{},
			expectError: true,
			expectedErr: "GET : Missing Key",
		},
	}

	for _, test := range tests {

		handlers.PlainCache.Data = test.initialData

		t.Run(test.name, func(t *testing.T) {
			val, err := handlers.GetHandler(test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != test.expectedErr {
					t.Errorf("Expected error: %v, but got: %v", test.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if val != test.expectedVal {
					t.Errorf("Expected value %v, but got %v", test.expectedVal, val)
				}
			}
		})
	}
}
