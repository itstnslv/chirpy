package auth

import (
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestValidateJWT(t *testing.T) {
	randomUUID := uuid.New()
	tokenSecret := "tokenSecret"
	tokenDuration := 2 * time.Second
	token, _ := MakeJWT(randomUUID, tokenSecret, tokenDuration)
	expiredToken, _ := MakeJWT(randomUUID, tokenSecret, 1)
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
			randomUUID,
			false,
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
	token := "ajfbdslkgbjrbegoernbvijsdfnvkjdfnsvkjdfsnv"
	testHeaders := http.Header{}
	testHeaders.Add("Authorization", "Bearer "+token)
	falseHeaders := http.Header{}
	falseHeaders.Add("Authorization", "Bearer ")
	type args struct {
		headers http.Header
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"valid token",
			args{
				headers: testHeaders,
			},
			token,
			false,
		},
		{
			"no header",
			args{
				headers: nil,
			},
			"",
			true,
		},
		{
			"no token",
			args{
				headers: falseHeaders,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.args.headers)
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
