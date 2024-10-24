package main

import "testing"

func Test_applyProfaneFilter(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no profanity",
			args{"I had something interesting for breakfast"},
			"I had something interesting for breakfast",
		},
		{
			"half body",
			args{"I hear Mastodon is better than Chirpy. sharbert I need to migrate"},
			"I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			"two words",
			args{"I really need a kerfuffle to go to bed sooner, Fornax !"},
			"I really need a **** to go to bed sooner, **** !",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyProfaneFilter(tt.args.body); got != tt.want {
				t.Errorf("applyProfaneFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
