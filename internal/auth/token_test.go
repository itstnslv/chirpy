package auth

import (
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "tokenSecret"
	tokenDuration := time.Hour
	token, _ := MakeJWT(userId, tokenSecret, tokenDuration)
	expiredToken, _ := MakeJWT(userId, tokenSecret, 1)
	type args struct {
		tokenString string
		tokenSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			"Valid token",
			args{
				token,
				tokenSecret,
			},
			userId,
			false,
		},
		{
			"Invalid token",
			args{
				"token.invalid.string",
				tokenSecret,
			},
			uuid.Nil,
			true,
		},
		{
			"Invalid secret",
			args{
				token,
				"invalidSecret",
			},
			uuid.Nil,
			true,
		},
		{
			"Expired token",
			args{
				expiredToken,
				tokenSecret,
			},
			uuid.Nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.args.tokenString, tt.args.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			"valid token",
			http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			"valid_token",
			false,
		},
		{
			"no header",
			http.Header{},
			"",
			true,
		},
		{
			"malformed header",
			http.Header{
				"Authorization": []string{"InvalidBearer invalid_token"},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
