package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorType определяет тип ошибки
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeDatabase   ErrorType = "database"
	ErrorTypeTelegram   ErrorType = "telegram"
	ErrorTypeNetwork    ErrorType = "network"
	ErrorTypeInternal   ErrorType = "internal"
	ErrorTypeUser       ErrorType = "user"
)

// AppError представляет структурированную ошибку приложения
type AppError struct {
	Type    ErrorType
	Message string
	Details string
	UserMsg string
	Code    string
	Context map[string]interface{}
	Inner   error
	Stack   string
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Inner)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap возвращает внутреннюю ошибку
func (e *AppError) Unwrap() error {
	return e.Inner
}

// GetUserMessage возвращает сообщение для пользователя
func (e *AppError) GetUserMessage() string {
	if e.UserMsg != "" {
		return e.UserMsg
	}
	return e.Message
}

// NewAppError создает новую ошибку приложения
func NewAppError(errorType ErrorType, message string, inner error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Inner:   inner,
		Stack:   getStackTrace(),
	}
}

// NewValidationError создает ошибку валидации
func NewValidationError(message string, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Details: details,
		UserMsg: "Проверьте введенные данные",
		Stack:   getStackTrace(),
	}
}

// NewDatabaseError создает ошибку базы данных
func NewDatabaseError(message string, inner error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Inner:   inner,
		UserMsg: "Ошибка работы с базой данных",
		Stack:   getStackTrace(),
	}
}

// NewTelegramError создает ошибку Telegram API
func NewTelegramError(message string, inner error) *AppError {
	return &AppError{
		Type:    ErrorTypeTelegram,
		Message: message,
		Inner:   inner,
		UserMsg: "Ошибка связи с Telegram",
		Stack:   getStackTrace(),
	}
}

// NewNetworkError создает сетевую ошибку
func NewNetworkError(message string, inner error) *AppError {
	return &AppError{
		Type:    ErrorTypeNetwork,
		Message: message,
		Inner:   inner,
		UserMsg: "Ошибка сети",
		Stack:   getStackTrace(),
	}
}

// NewInternalError создает внутреннюю ошибку
func NewInternalError(message string, inner error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Inner:   inner,
		UserMsg: "Внутренняя ошибка системы",
		Stack:   getStackTrace(),
	}
}

// NewUserError создает ошибку пользователя
func NewUserError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUser,
		Message: message,
		UserMsg: message,
		Stack:   getStackTrace(),
	}
}

// WithContext добавляет контекст к ошибке
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithUserMessage устанавливает сообщение для пользователя
func (e *AppError) WithUserMessage(userMsg string) *AppError {
	e.UserMsg = userMsg
	return e
}

// WithCode устанавливает код ошибки
func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

// getStackTrace получает стек вызовов
func getStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	// Убираем первые несколько строк (сама функция getStackTrace и NewAppError)
	lines := strings.Split(stack, "\n")
	if len(lines) > 6 {
		lines = lines[6:]
	}

	return strings.Join(lines, "\n")
}

// IsValidationError проверяет, является ли ошибка ошибкой валидации
func IsValidationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeValidation
	}
	return false
}

// IsDatabaseError проверяет, является ли ошибка ошибкой базы данных
func IsDatabaseError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeDatabase
	}
	return false
}

// IsTelegramError проверяет, является ли ошибка ошибкой Telegram
func IsTelegramError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeTelegram
	}
	return false
}

// IsNetworkError проверяет, является ли ошибка сетевой ошибкой
func IsNetworkError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeNetwork
	}
	return false
}

// IsInternalError проверяет, является ли ошибка внутренней ошибкой
func IsInternalError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeInternal
	}
	return false
}

// IsUserError проверяет, является ли ошибка ошибкой пользователя
func IsUserError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeUser
	}
	return false
}

// WrapError оборачивает обычную ошибку в AppError
func WrapError(err error, errorType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return NewAppError(errorType, message, err)
}

// HandleError обрабатывает ошибку и возвращает сообщение для пользователя
func HandleError(err error) string {
	if err == nil {
		return ""
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.GetUserMessage()
	}

	// Для обычных ошибок возвращаем общее сообщение
	return "Произошла ошибка, попробуйте позже"
}
