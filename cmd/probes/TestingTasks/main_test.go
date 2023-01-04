package main

import (
	"testing"
)

func TestAbs(test *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{
			name:  "test#1",
			value: 1.2,
			want:  999,
		},
	}
	for _, tt := range tests {

		test.Run(tt.name, func(tst *testing.T) {

			if abs := Abs(tt.value); abs != tt.want {
				test.Errorf("Abs()=%v, want %v", abs, tt.want)
				test.Errorf("Abs()=%v, want %v", abs, tt.want)

			}
		})
	}
}
