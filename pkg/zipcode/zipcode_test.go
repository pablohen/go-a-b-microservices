package zipcode

import (
	"testing"

	"go-a-b-microservices/pkg/apperror"
)

func TestZipCodeRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		zipCode string
		wantErr error
	}{
		{
			name:    "valid zipcode",
			zipCode: "13484000",
			wantErr: nil,
		},
		{
			name:    "empty zipcode",
			zipCode: "",
			wantErr: apperror.ErrZipCodeRequired,
		},
		{
			name:    "invalid zipcode - too short",
			zipCode: "1234567",
			wantErr: apperror.ErrZipCodeInvalid,
		},
		{
			name:    "invalid zipcode - too long",
			zipCode: "123456789",
			wantErr: apperror.ErrZipCodeInvalid,
		},
		{
			name:    "invalid zipcode - contains letters",
			zipCode: "1234567a",
			wantErr: apperror.ErrZipCodeInvalid,
		},
		{
			name:    "invalid zipcode - contains special characters",
			zipCode: "1234-567",
			wantErr: apperror.ErrZipCodeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ZipCodeRequest{CEP: tt.zipCode}

			err := req.Validate()

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ZipCodeRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ZipCodeRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
