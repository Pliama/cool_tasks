package users_test

import (
	"bytes"
	"github.com/Nastya-Kruglikova/cool_tasks/src/models"
	"github.com/Nastya-Kruglikova/cool_tasks/src/services"
	"github.com/Nastya-Kruglikova/cool_tasks/src/services/users"
	"github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/Nastya-Kruglikova/cool_tasks/src/services/auth"
)

var router = services.NewRouter()

type usersCRUDTestCase struct {
	name              string
	url               string
	want              int
	mockedGetUser     models.User
	mockedCreateUser  models.User
	mockedGetUsers    []models.User
	mockedUserError   error
	mockedDeleteUsers uuid.UUID
	permission bool
	mock              func()
	error             string
	testUser          models.User
}

func TestGetUsers(t *testing.T) {
	tests := []usersCRUDTestCase{
		{
			name:            "Get_Users_200",
			url:             "/v1/users",
			want:            200,
			mockedGetUsers:  []models.User{},
			mockedUserError: nil,
		},
		{
			name:            "Get_Users_404",
			url:             "/v1/users",
			want:            404,
			mockedGetUsers:  []models.User{},
			mockedUserError: http.ErrBodyNotAllowed,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models.MockedGetUsers(tc.mockedGetUsers, tc.mockedUserError)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []usersCRUDTestCase{
		{
			name:            "Get_Users_200",
			url:             "/v1/users/a7264252-6ef4-11e8-9982-0242ac110002",
			want:            200,
			mockedGetUser:   models.User{},
			mockedUserError: nil,
		},
		{
			name:            "Get_Users_400",
			url:             "/v1/users/asdad",
			want:            400,
			mockedGetUser:   models.User{},
			mockedUserError: nil,
		},
		{
			name:            "Get_Users_404",
			url:             "/v1/users/a7264252-6ef4-11e8-9982-0242ac110002",
			want:            404,
			mockedGetUser:   models.User{},
			mockedUserError: http.ErrNoLocation,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models.MockedGetUser(tc.mockedGetUser, tc.mockedUserError)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {

	userId, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	tests := []usersCRUDTestCase{
		{
			name:              "Delete_Users_200",
			url:               "/v1/users/00000000-0000-0000-0000-000000000001",
			want:              200,
			mockedDeleteUsers: userId,
			mockedUserError:   nil,
			permission: true,
			mock: func() {
			},
		},
		{
			name:              "Delete_Users_404",
			url:               "/v1/users/00000000-0000-0000-0000-000000000001",
			want:              404,
			mockedDeleteUsers: userId,
			mockedUserError:   nil,
			permission: true,
			mock: func() {
				var err = http.ErrBodyNotAllowed
				models.DeleteUser = func(id uuid.UUID) error {
					return err
				}
			},
		},
		{
			name:              "Delete_Users_400",
			url:               "/v1/users/sadsad",
			want:              400,
			mockedDeleteUsers: userId,
			mockedUserError:   nil,
			permission: true,
			mock: func() {
			},
		},
		{
			name:              "Delete_Users_403",
			url:               "/v1/users/00000000-0000-0000-0000-000000000001",
			want:              403,
			mockedDeleteUsers: userId,
			mockedUserError:   nil,
			permission: false,
			mock: func() {
			},
		},
	}
	defer func(){auth.CheckPermission=auth.CheckPermission}()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models.MockedDeleteUser(userId, nil)

			auth.MockedCheckPermission(tc.permission)
			tc.mock()
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, tc.url, nil)

			router.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Errorf("Expected: %d , got %d", tc.want, rec.Code)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []usersCRUDTestCase{
		{
			name:             "Add_Users_200",
			url:              "/v1/users",
			want:             200,
			mockedCreateUser: models.User{},
			mockedUserError:  nil,
			permission: true,
		},
		{
			name:             "Add_Users_403",
			url:              "/v1/users",
			want:             403,
			mockedCreateUser: models.User{},
			mockedUserError:  nil,
			permission: false,
		},
	}
	data := url.Values{}
	data.Add("name", "Karim")
	data.Add("login", "Karim123")
	data.Add("password", "1324qwer")
	defer func(){auth.CheckPermission=auth.CheckPermission}()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			auth.MockedCheckPermission(tc.permission)
			models.MockedCreateUser(tc.mockedCreateUser)
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

func TestIsValid(t *testing.T) {
	id, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	tests := []usersCRUDTestCase{
		{
			name:  "Valid data",
			error: "",
			testUser: models.User{
				ID:       id,
				Name:     "Validname",
				Login:    "Validlogin",
				Password: "Validpassword",
			},
		},
		{
			name:  "Invalid Password",
			error: "Invalid Password",
			testUser: models.User{
				ID:       id,
				Name:     "Validname",
				Login:    "Validlogin",
				Password: "1234",
			},
		},
		{
			name:  "Invalid username",
			error: " Invalid Name",
			testUser: models.User{
				ID:       id,
				Name:     "invalidname",
				Login:    "Validlogin",
				Password: "Validpassword",
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			_, err := users.IsValid(tc.testUser)
			if err != tc.error {
				t.Errorf("Expected: %s , got %s", tc.error, err)
			}
		})
	}
}