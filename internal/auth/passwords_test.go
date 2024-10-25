package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "correctPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			hash:     hash1,
			password: password1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			hash:     hash1,
			password: "wrongPassword",
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			hash:     hash2,
			password: password1,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			hash:     hash1,
			password: "",
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			hash:     "invalidhash",
			password: password1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
