package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/mock"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
)

func TestSignUP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := mock.NewMockService(ctrl)

	h := New(config.Config{Port: "8080"}, userServiceMock)

	testUser := models.UserInfo{
		Name:     "John Doe",
		Email:    "johndoe@example.com",
		PhoneNO:  "+1-245-213-2222",
		Password: "securepassword",
	}

	userJSON, _ := json.Marshal(testUser)
	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Expect the SignUP function to be called with the correct arguments
	userServiceMock.EXPECT().SignUP(gomock.Any(), testUser).Return(nil)

	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	h.SignUP(ctx)

	assert.Equal(t, http.StatusCreated, rr.Code)

	expectedBody := `{"message": "User created successfully"}`

	var expected, actual map[string]interface{}
	json.Unmarshal([]byte(expectedBody), &expected)
	json.Unmarshal([]byte(rr.Body.String()), &actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestSignIN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := mock.NewMockService(ctrl)

	h := New(config.Config{Port: "8080"}, userServiceMock)

	userSignInfo := SignINInfo{
		Phone:    "+1-245-213-2222",
		Password: "dfbhsjdkflgdb",
	}

	sigINJson, _ := json.Marshal(userSignInfo)

	req, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(sigINJson))
	if err != nil {
		t.Errorf("failed to create requests: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	userInfo := models.UserInfo{
		PhoneNO:  "+1-245-213-2222",
		Password: "dfbhsjdkflgdb",
	}
	rr := httptest.NewRecorder()
	tokenString := "token"
	userServiceMock.EXPECT().SignIN(gomock.Any(), userInfo).Return(tokenString, nil)

	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	h.SignIN(ctx)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody := `{"message": "user successfully authorized"}`

	var expected, actual map[string]interface{}
	json.Unmarshal([]byte(expectedBody), &expected)
	json.Unmarshal([]byte(rr.Body.String()), &actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v but got %v", expected, actual)
	}

}

func TestLogOut(t *testing.T) {
	ctrl := gomock.NewController(t)

	userServiceMock := mock.NewMockService(ctrl)

	h := New(config.Config{Port: "8080"}, userServiceMock)

	req, err := http.NewRequest(http.MethodPost, "/logout", nil)
	if err != nil {
		t.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	h.LogOut(ctx)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody := `{"message": "Successfully logged out"}`

	var expected, actual map[string]interface{}
	json.Unmarshal([]byte(expectedBody), &expected)
	json.Unmarshal([]byte(rr.Body.String()), &actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestAuthMiddleWare(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	NewServiceMock := mock.NewMockService(ctrl)

	h := New(config.Config{Port: "8088"}, NewServiceMock)

	req, err := http.NewRequest(http.MethodGet, "/profile", nil)
	if err != nil {
		t.Errorf("failed to create request: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "token"})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	NewServiceMock.EXPECT().VerifyToken(gomock.Any(), "token").Return("UserID", nil)
	NewServiceMock.EXPECT().CheckTokenInRedis(gomock.Any(), "token").Return(nil)

	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	h.AuthMiddleWare(ctx)
	h.Profile(ctx)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody := `{"message": "Welcome to the protected profile area!"}`

	var expected, actual map[string]interface{}
	json.Unmarshal([]byte(expectedBody), &expected)
	json.Unmarshal([]byte(rr.Body.String()), &actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}
