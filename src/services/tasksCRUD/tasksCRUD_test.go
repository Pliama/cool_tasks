package tasksCRUD_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nastya-Kruglikova/cool_tasks/src/services"
	"net/url"
	"bytes"
)

var router = services.NewRouter()

type tasksCRUDTestCase struct {
	name string
	url  string
	want int
}

func TestGetTasks(t *testing.T) {
	tests := []tasksCRUDTestCase{
		{
			name: "Get_Tasks_200",
			url:  "/v1/tasks",
			want: 200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}
func TestGetTasksByID(t *testing.T) {
	tests := []tasksCRUDTestCase{
		{
			name: "Get_TaskById_400",
			url:  "/v1/tasks/wrongID",
			want: 400,
		},
		{
			name: "Get_TaskById_404",
			url:  "/v1/tasks/-1",
			want: 404,
		},
		{
			name: "Get_TaskById_200",
			url:  "/v1/tasks/1",
			want: 200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}
func TestDeleteTasks(t *testing.T) {
	tests := []tasksCRUDTestCase{
		{
			name: "Delete_Task_415",
			url:  "/v1/tasks/wrongID",
			want: 415,
		},
		{
			name: "Delete_Task_404",
			url:  "/v1/tasks/-1",
			want: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, tc.url, nil)

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}
func TestAddTasks(t *testing.T) {
	tests := []tasksCRUDTestCase{
		{
			name: "Add_Task_400",
			url:  "/v1/tasks",
			want: 400,
		},
	}

	data := url.Values{}
	data.Add("user_id", "error") //bad value
	data.Add("name", "JustUser")
	data.Add("time", "Mon Jan 2 15:04:05 MST 2006")
	data.Add("desc", "Desc of my task")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, tc.url, bytes.NewBufferString(data.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}
