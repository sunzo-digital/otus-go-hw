package hw09structvalidator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		description string
		in          interface{}
		expectedErr error
	}{
		{
			description: "Not a structure was passed",
			in:          0,
			expectedErr: errors.New("received entity is not a structure"),
		},
		{
			description: "User.ID is more than 36",
			in: User{
				ID:     makeStringWithLength(37),
				Age:    18,
				Email:  "w@w.w",
				Role:   UserRole("admin"),
				Phones: []string{makeStringWithLength(11)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   errors.New("length of value is not equal to 36"),
				},
			},
		},
		{
			description: "User.Age is less than 18",
			in: User{
				ID:     makeStringWithLength(36),
				Age:    17,
				Email:  "w@w.w",
				Role:   UserRole("admin"),
				Phones: []string{makeStringWithLength(11)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   errors.New("value is less than 18"),
				},
			},
		},
		{
			description: "User.Age is more than 50",
			in: User{
				ID:     makeStringWithLength(36),
				Age:    51,
				Email:  "w@w.w",
				Role:   UserRole("admin"),
				Phones: []string{makeStringWithLength(11)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   errors.New("value is greater than 50"),
				},
			},
		},
		{
			description: "User.Email isn't match regexp",
			in: User{
				ID:     makeStringWithLength(36),
				Age:    50,
				Email:  "dark_ganjubas_killer_666@mailru", // miss "."
				Role:   UserRole("stuff"),
				Phones: []string{makeStringWithLength(11)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Email",
					Err:   errors.New("value does not match the regular expression"),
				},
			},
		},
		{
			description: "User.Role is not among the valid values",
			in: User{
				ID:     makeStringWithLength(36),
				Age:    50,
				Email:  "w@w.w",
				Role:   UserRole("moderator"),
				Phones: []string{makeStringWithLength(11)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Role",
					Err:   errors.New("value is not among the valid values: [admin stuff]"),
				},
			},
		},
		{
			description: "one of User.Phones is less than 11",
			in: User{
				ID:     makeStringWithLength(36),
				Age:    50,
				Email:  "w@w.w",
				Role:   UserRole("admin"),
				Phones: []string{makeStringWithLength(11), makeStringWithLength(10)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Phones",
					Err:   errors.New("length of value is not equal to 11"),
				},
			},
		},
		{
			description: "Few errors",
			in: User{
				ID:     makeStringWithLength(37),
				Age:    51,
				Email:  "invalid",
				Role:   UserRole("ibragim"),
				Phones: []string{makeStringWithLength(11), makeStringWithLength(10)},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   errors.New("length of value is not equal to 36"),
				},
				{
					Field: "Age",
					Err:   errors.New("value is greater than 50"),
				},
				{
					Field: "Email",
					Err:   errors.New("value does not match the regular expression"),
				},
				{
					Field: "Role",
					Err:   errors.New("value is not among the valid values: [admin stuff]"),
				},
				{
					Field: "Phones",
					Err:   errors.New("length of value is not equal to 11"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}

func makeStringWithLength(strLen int) string {
	return string(make([]byte, strLen))
}
