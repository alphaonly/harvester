package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserViewHandler(t *testing.T) {

	type want struct {
		code        int
		ID          string
		response    string
		contentType string
	}

	users := make(map[string]User)

	u1 := User{
		ID:        "u1",
		FirstName: "Misha",
		LastName:  "Popov",
	}
	u2 := User{
		ID:        "u2",
		FirstName: "Sasha",
		LastName:  "Popov",
	}
	var u4 User

	users["u1"] = u1
	users["u2"] = u2
	users["u4"] = u4

	tests := []struct {
		name string
		ID   string
		want want
	}{
		{
			name: "test#1 positive",
			ID:   "u1",
			want: want{code: 200,
				response:    `{"status":"ok"}`,
				contentType: "application/json"},
		},
		{
			name: "test#2 negative",
			ID:   "",
			want: want{
				code:        400,
				response:    `{"status":"bad request"}`,
				contentType: "application/json"},
		},
		{
			name: "test#3 negative",
			ID:   "V3",
			want: want{
				code:        404,
				response:    `{"status":"bad request"}`,
				contentType: "application/json"},
		},
		{
			name: "test#4 negative",
			ID:   "u4",
			want: want{
				code:        500,
				response:    ``,
				contentType: "application/json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, "/users?user_id="+tt.ID, nil)

			w := httptest.NewRecorder()
			h := UserViewHandler(users)

			h.ServeHTTP(w, request)

			response := w.Result()

			if response.StatusCode != tt.want.code {
				t.Errorf("error code %v want %v", response.StatusCode, tt.want.code)
			}
			if (response.StatusCode == http.StatusOK) && (response.Header.Get("Content-type") != tt.want.contentType) {
				t.Errorf("error contentType %v want %v", response.Header.Get("Content-type"), tt.want.contentType)
			}

		})

	}
}
