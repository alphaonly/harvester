package main

import (
	"testing"
)

func TestFullName(t *testing.T) {
	tests := []struct {
		name      string
		FirstName string
		LastName  string
		want      string
	}{
		{
			name:      "test#1",
			FirstName: "Misha",
			LastName:  "Popov",
			want:      "Misha Popov",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(tst *testing.T) {
			result := User{tt.FirstName, tt.LastName}.FullName()
			if result != tt.want {
				t.Errorf("Abs()=%v, want %v", result, tt.want)

			}
		})
	}
}
