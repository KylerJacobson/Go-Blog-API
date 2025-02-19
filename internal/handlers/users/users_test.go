package users

import (
	"bytes"
	"encoding/json"
	"errors"
	userModels "github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUsersRepository struct {
	mock.Mock
}

func (m *mockUsersRepository) UpdateUser(user userModels.UserUpdate) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockUsersRepository) GetUserById(id int) (*userModels.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockUsersRepository) GetUserByEmail(email string) (*userModels.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockUsersRepository) GetAllUsers() (*[]userModels.FrontendUser, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockUsersRepository) DeleteUserById(id int) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockUsersRepository) LoginUser(user userModels.UserLogin) (*userModels.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockUsersRepository) CreateUser(user userModels.UserCreate) (string, error) {
	args := m.Called(user)
	return args.Get(0).(string), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{} // Changed to interface{} to allow invalid JSON
		setupMock      func(*mockUsersRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful_user_creation",
			requestBody: userModels.UserCreate{
				FirstName:         "John",
				LastName:          "Doe",
				Email:             "john@test.com",
				Password:          "password123",
				AccessRequest:     0,
				EmailNotification: true,
			},
			setupMock: func(m *mockUsersRepository) {
				m.On("CreateUser", mock.MatchedBy(func(user userModels.UserCreate) bool {
					return user.FirstName == "John" &&
						user.LastName == "Doe" &&
						user.Email == "john@test.com" &&
						user.Password == "password123" &&
						user.AccessRequest == 0 &&
						user.EmailNotification == true
				})).Return("1", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id": "1",
			},
		},
		{
			name: "missing_required_fields",
			requestBody: userModels.UserCreate{

				// Missing FirstName, LastName and Email
				Password:          "",
				AccessRequest:     0,
				EmailNotification: true,
			},
			setupMock: func(m *mockUsersRepository) {
				// No mock needed as validation should fail before repository call
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"status": float64(400), "message": "Invalid request body", "detail": "first name is required, last name is required, email is required, password is required, password must be at least 8 characters long, invalid email format"},
		},
		{
			name: "invalid_email_format",
			requestBody: userModels.UserCreate{
				FirstName:         "John",
				LastName:          "Doe",
				Email:             "invalid-email", // Invalid email format
				Password:          "password123",
				AccessRequest:     0,
				EmailNotification: true,
			},
			setupMock: func(m *mockUsersRepository) {
				// No mock needed as validation should fail before repository call
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"status": float64(400), "message": "Invalid request body", "detail": "invalid email format"},
		},
		{
			name: "bad access request",
			requestBody: userModels.UserCreate{
				FirstName:         "John",
				LastName:          "Doe",
				Email:             "john@test.com",
				Password:          "password123", // Too short password
				AccessRequest:     4,
				EmailNotification: true,
			},
			setupMock: func(m *mockUsersRepository) {
				// No mock needed as validation should fail before repository call
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"status": float64(400), "message": "Invalid request body", "detail": "access request must be -1, 0 or 2"},
		},
		{
			name:        "bad request body",
			requestBody: `{"firstName": "John", "lastName": "Doe", email: bad-json}`,
			setupMock: func(m *mockUsersRepository) {
				// No mock needed as validation should fail before repository call
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"status": float64(400), "message": "Invalid request body", "detail": "invalid character 'e' looking for beginning of object key string"},
		},
		{
			name: "unsuccessful_user_creation",
			requestBody: userModels.UserCreate{
				FirstName:         "John",
				LastName:          "Doe",
				Email:             "john@test.com",
				Password:          "password123",
				AccessRequest:     0,
				EmailNotification: true,
			},
			setupMock: func(m *mockUsersRepository) {
				m.On("CreateUser", mock.MatchedBy(func(user userModels.UserCreate) bool {
					return user.FirstName == "John" &&
						user.LastName == "Doe" &&
						user.Email == "john@test.com" &&
						user.Password == "password123" &&
						user.AccessRequest == 0 &&
						user.EmailNotification == true
				})).Return("", errors.New("failed to create user"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"detail": "failed to create user", "message": "failed to create user", "status": float64(500)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock repository
			mockRepo := new(mockUsersRepository)

			// Creat test logger
			testLogger := zap.NewNop()

			// Setup mock expectations
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			// Create API instance
			usersApi := New(mockRepo, testLogger)

			// Create request body
			var bodyBytes []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			usersApi.CreateUser(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch for test: %s", tt.name)

			// Assert Content-Type header
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"),
				"Content-Type header mismatch for test: %s", tt.name)

			// Parse response body
			var responseBody map[string]interface{}
			err = json.NewDecoder(rr.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("Failed to decode response body in test %s: %v", tt.name, err)
			}

			// Assert response body
			assert.Equal(t, tt.expectedBody, responseBody,
				"Response body mismatch for test: %s", tt.name)

			// Verify that all mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
