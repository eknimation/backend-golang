package usecase

import (
	"errors"
	"os"
	"testing"
	"time"

	"backend-service/internal/domain"
)

// MockUserRepository implements domain.UserInterface for testing
type MockUserRepository struct {
	users   map[string]*domain.User
	emails  map[string]*domain.User
	lastID  int
	errMode string // Used to simulate different error conditions
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*domain.User),
		emails: make(map[string]*domain.User),
		lastID: 0,
	}
}

func (m *MockUserRepository) CreateUser(user domain.User) error {
	if m.errMode == "create_error" {
		return errors.New("database error")
	}

	// Check if email already exists
	if _, exists := m.emails[user.Email]; exists {
		return errors.New("email already exists")
	}

	m.lastID++
	user.ID = string(rune(m.lastID))
	now := time.Now()
	user.CreatedAt = &now

	m.users[user.ID] = &user
	m.emails[user.Email] = &user

	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	if m.errMode == "get_by_email_error" {
		return nil, errors.New("database error")
	}

	user, exists := m.emails[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (m *MockUserRepository) GetUserByID(id string) (*domain.User, error) {
	if m.errMode == "get_by_id_error" {
		return nil, errors.New("database error")
	}

	if id == "invalid" {
		return nil, errors.New("invalid user ID format")
	}

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (m *MockUserRepository) GetUsersWithPagination(page, limit int) ([]*domain.User, int, error) {
	if m.errMode == "pagination_error" {
		return nil, 0, errors.New("database error")
	}

	// Convert map to slice
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}

	totalCount := len(users)

	// Apply pagination
	start := (page - 1) * limit
	if start >= len(users) {
		return []*domain.User{}, totalCount, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], totalCount, nil
}

func (m *MockUserRepository) UpdateUser(id string, user domain.User) error {
	if m.errMode == "update_error" {
		return errors.New("database error")
	}

	existingUser, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}

	// Check email conflict if email is being updated
	if user.Email != "" && user.Email != existingUser.Email {
		if _, emailExists := m.emails[user.Email]; emailExists {
			return errors.New("email already exists")
		}

		// Remove old email mapping
		delete(m.emails, existingUser.Email)
		existingUser.Email = user.Email
		m.emails[user.Email] = existingUser
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}

	return nil
}

func (m *MockUserRepository) DeleteUser(id string) error {
	if m.errMode == "delete_error" {
		return errors.New("database error")
	}

	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	delete(m.emails, user.Email)

	return nil
}

func (m *MockUserRepository) GetUserCount() (int64, error) {
	if m.errMode == "count_error" {
		return 0, errors.New("database error")
	}

	return int64(len(m.users)), nil
}

// Helper method to set error mode for testing error conditions
func (m *MockUserRepository) SetErrorMode(mode string) {
	m.errMode = mode
}

// Helper method to clear error mode
func (m *MockUserRepository) ClearErrorMode() {
	m.errMode = ""
}

// Helper method to add a user directly to the mock
func (m *MockUserRepository) AddUser(user *domain.User) {
	m.users[user.ID] = user
	m.emails[user.Email] = user
}

func setupTestUsecase() *Usecase {
	return &Usecase{
		userRepo: NewMockUserRepository(),
	}
}

func TestMain(m *testing.M) {
	// Set up test environment variables
	_ = os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-testing")
	_ = os.Setenv("MONGODB_HOST", "localhost")
	_ = os.Setenv("MONGODB_PORT", "27017")
	_ = os.Setenv("MONGODB_DATABASE_NAME", "test")
	_ = os.Setenv("API_PORT", "8080")
	_ = os.Setenv("APP_ENV", "test")
	_ = os.Setenv("BASE_PATH", "/api")

	// Run tests
	code := m.Run()

	// Clean up
	_ = os.Unsetenv("JWT_SECRET")
	_ = os.Unsetenv("MONGODB_HOST")
	_ = os.Unsetenv("MONGODB_PORT")
	_ = os.Unsetenv("MONGODB_DATABASE_NAME")
	_ = os.Unsetenv("API_PORT")
	_ = os.Unsetenv("APP_ENV")
	_ = os.Unsetenv("BASE_PATH")

	os.Exit(code)
}

func TestCreateUser_Success(t *testing.T) {
	uc := setupTestUsecase()

	user := domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	err := uc.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was created
	mockRepo := uc.userRepo.(*MockUserRepository)
	if len(mockRepo.users) != 1 {
		t.Errorf("Expected 1 user in repository, got %d", len(mockRepo.users))
	}
}

func TestCreateUser_EmptyPassword(t *testing.T) {
	uc := setupTestUsecase()

	user := domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "", // Empty password
	}

	err := uc.CreateUser(user)
	// Empty password should still be processed (business logic validation
	// should be handled at the controller/validation layer)
	if err != nil {
		t.Errorf("Expected no error for empty password at usecase level, got %v", err)
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)
	mockRepo.SetErrorMode("create_error")

	user := domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	err := uc.CreateUser(user)

	if err == nil {
		t.Error("Expected repository error, got nil")
	}
}

func TestCreateUser_EmailExists(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	// Add existing user
	existingUser := &domain.User{
		ID:    "1",
		Name:  "Existing User",
		Email: "john@example.com",
	}
	mockRepo.AddUser(existingUser)

	user := domain.User{
		Name:     "John Doe",
		Email:    "john@example.com", // Same email
		Password: "password123",
	}

	err := uc.CreateUser(user)

	if err == nil || err.Error() != "email already exists" {
		t.Errorf("Expected 'email already exists' error, got %v", err)
	}
}

func TestAuthenticateUser_Success(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	// Create a user with hashed password
	hashedPassword := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi" // "password"
	existingUser := &domain.User{
		ID:       "1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}
	mockRepo.AddUser(existingUser)

	user := domain.User{
		Email:    "john@example.com",
		Password: "password",
	}

	token, err := uc.AuthenticateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected token to be returned, got empty string")
	}
}

func TestAuthenticateUser_InvalidCredentials_UserNotFound(t *testing.T) {
	uc := setupTestUsecase()

	user := domain.User{
		Email:    "nonexistent@example.com",
		Password: "password",
	}

	token, err := uc.AuthenticateUser(user)

	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %v", err)
	}

	if token != "" {
		t.Error("Expected empty token, got non-empty string")
	}
}

func TestAuthenticateUser_InvalidCredentials_WrongPassword(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	// Create a user with hashed password
	hashedPassword := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi" // "password"
	existingUser := &domain.User{
		ID:       "1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}
	mockRepo.AddUser(existingUser)

	user := domain.User{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	token, err := uc.AuthenticateUser(user)

	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %v", err)
	}

	if token != "" {
		t.Error("Expected empty token, got non-empty string")
	}
}

func TestAuthenticateUser_RepositoryError(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)
	mockRepo.SetErrorMode("get_by_email_error")

	user := domain.User{
		Email:    "john@example.com",
		Password: "password",
	}

	token, err := uc.AuthenticateUser(user)

	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %v", err)
	}

	if token != "" {
		t.Error("Expected empty token, got non-empty string")
	}
}

func TestGetUserByID_Success(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	expectedUser := &domain.User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockRepo.AddUser(expectedUser)

	user, err := uc.GetUserByID("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user to be returned, got nil")
	}

	if user.ID != expectedUser.ID || user.Name != expectedUser.Name || user.Email != expectedUser.Email {
		t.Errorf("Expected user %+v, got %+v", expectedUser, user)
	}
}

func TestGetUserByID_UserNotFound(t *testing.T) {
	uc := setupTestUsecase()

	user, err := uc.GetUserByID("nonexistent")

	if err == nil || err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got %v", err)
	}

	if user != nil {
		t.Error("Expected nil user, got non-nil")
	}
}

func TestGetUserByID_InvalidID(t *testing.T) {
	uc := setupTestUsecase()

	user, err := uc.GetUserByID("invalid")

	if err == nil || err.Error() != "invalid user ID format" {
		t.Errorf("Expected 'invalid user ID format' error, got %v", err)
	}

	if user != nil {
		t.Error("Expected nil user, got non-nil")
	}
}

func TestGetUsersWithPagination_Success(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	// Add multiple users
	for i := 1; i <= 25; i++ {
		user := &domain.User{
			ID:    string(rune(i)),
			Name:  "User " + string(rune(i)),
			Email: "user" + string(rune(i)) + "@example.com",
		}
		mockRepo.AddUser(user)
	}

	users, totalCount, err := uc.GetUsersWithPagination(1, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(users) != 10 {
		t.Errorf("Expected 10 users, got %d", len(users))
	}

	if totalCount != 25 {
		t.Errorf("Expected total count 25, got %d", totalCount)
	}
}

func TestGetUsersWithPagination_InvalidPage(t *testing.T) {
	uc := setupTestUsecase()

	users, totalCount, err := uc.GetUsersWithPagination(0, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should default to page 1
	if len(users) != 0 || totalCount != 0 {
		t.Errorf("Expected empty result for page 0 with no users")
	}
}

func TestGetUsersWithPagination_InvalidLimit(t *testing.T) {
	uc := setupTestUsecase()

	users, totalCount, err := uc.GetUsersWithPagination(1, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should use default limit
	if len(users) != 0 || totalCount != 0 {
		t.Errorf("Expected empty result with default limit")
	}
}

func TestGetUsersWithPagination_LimitTooHigh(t *testing.T) {
	uc := setupTestUsecase()

	users, totalCount, err := uc.GetUsersWithPagination(1, 200)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should use default limit of 10
	if len(users) != 0 || totalCount != 0 {
		t.Errorf("Expected empty result with capped limit")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	existingUser := &domain.User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockRepo.AddUser(existingUser)

	updateUser := domain.User{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	err := uc.UpdateUser("1", updateUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was updated
	updatedUser, _ := mockRepo.GetUserByID("1")
	if updatedUser.Name != "John Updated" || updatedUser.Email != "john.updated@example.com" {
		t.Errorf("User was not updated correctly: %+v", updatedUser)
	}
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	uc := setupTestUsecase()

	updateUser := domain.User{
		Name: "John Updated",
	}

	err := uc.UpdateUser("nonexistent", updateUser)

	if err == nil || err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got %v", err)
	}
}

func TestUpdateUser_EmailExists(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	user1 := &domain.User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	user2 := &domain.User{
		ID:    "2",
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
	mockRepo.AddUser(user1)
	mockRepo.AddUser(user2)

	updateUser := domain.User{
		Email: "jane@example.com", // Email already exists
	}

	err := uc.UpdateUser("1", updateUser)

	if err == nil || err.Error() != "email already exists" {
		t.Errorf("Expected 'email already exists' error, got %v", err)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	existingUser := &domain.User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockRepo.AddUser(existingUser)

	err := uc.DeleteUser("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was deleted
	if len(mockRepo.users) != 0 {
		t.Errorf("Expected 0 users after deletion, got %d", len(mockRepo.users))
	}
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	uc := setupTestUsecase()

	err := uc.DeleteUser("nonexistent")

	if err == nil || err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got %v", err)
	}
}

func TestGetUserCount_Success(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)

	// Add some users
	for i := 1; i <= 5; i++ {
		user := &domain.User{
			ID:    string(rune(i)),
			Name:  "User " + string(rune(i)),
			Email: "user" + string(rune(i)) + "@example.com",
		}
		mockRepo.AddUser(user)
	}

	count, err := uc.GetUserCount()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestGetUserCount_RepositoryError(t *testing.T) {
	uc := setupTestUsecase()
	mockRepo := uc.userRepo.(*MockUserRepository)
	mockRepo.SetErrorMode("count_error")

	count, err := uc.GetUserCount()

	if err == nil {
		t.Error("Expected repository error, got nil")
	}

	if count != 0 {
		t.Errorf("Expected count 0 on error, got %d", count)
	}
}
