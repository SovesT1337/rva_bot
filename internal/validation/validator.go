package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("поле '%s': %s", e.Field, e.Message)
}

// ValidationResult содержит результат валидации
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// AddError добавляет ошибку валидации
func (vr *ValidationResult) AddError(field, message string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, ValidationError{Field: field, Message: message})
}

// HasErrors проверяет наличие ошибок
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// GetErrorMessages возвращает все сообщения об ошибках
func (vr *ValidationResult) GetErrorMessages() []string {
	messages := make([]string, len(vr.Errors))
	for i, err := range vr.Errors {
		messages[i] = err.Error()
	}
	return messages
}

// Validator содержит методы валидации
type Validator struct{}

// validateStringLength проверяет длину строки
func (v *Validator) validateStringLength(value, fieldName string, minLen, maxLen int) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	if len(value) < minLen {
		result.AddError(fieldName, fmt.Sprintf("должно содержать минимум %d символов", minLen))
	}

	if len(value) > maxLen {
		result.AddError(fieldName, fmt.Sprintf("не должно превышать %d символов", maxLen))
	}

	return result
}

// validateRequired проверяет, что поле не пустое
func (v *Validator) validateRequired(value, fieldName string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	if strings.TrimSpace(value) == "" {
		result.AddError(fieldName, "не может быть пустым")
	}

	return result
}

// validateInvalidChars проверяет на недопустимые символы
func (v *Validator) validateInvalidChars(value, fieldName string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	invalidChars := regexp.MustCompile(`[<>{}[\]\\|` + "`" + `~!@#$%^&*()+=]`)
	if invalidChars.MatchString(value) {
		result.AddError(fieldName, "содержит недопустимые символы")
	}

	return result
}

// NewValidator создает новый валидатор
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateTrainerName валидирует имя тренера
func (v *Validator) ValidateTrainerName(name string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(name, "name"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем длину
	if lengthResult := v.validateStringLength(name, "name", 2, 100); !lengthResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, lengthResult.Errors...)
	}

	// Проверяем недопустимые символы
	if charsResult := v.validateInvalidChars(name, "name"); !charsResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, charsResult.Errors...)
	}

	return result
}

// ValidateTelegramID валидирует Telegram ID
func (v *Validator) ValidateTelegramID(tgID string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(tgID, "tg_id"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем формат: @username или числовой ID
	if strings.HasPrefix(tgID, "@") {
		username := strings.TrimPrefix(tgID, "@")
		if len(username) < 5 {
			result.AddError("tg_id", "username должен содержать минимум 5 символов")
		}
		if len(username) > 32 {
			result.AddError("tg_id", "username не должен превышать 32 символа")
		}
		// Проверяем допустимые символы для username
		usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		if !usernameRegex.MatchString(username) {
			result.AddError("tg_id", "username содержит недопустимые символы")
		}
	} else {
		// Проверяем числовой ID
		if _, err := strconv.ParseInt(tgID, 10, 64); err != nil {
			result.AddError("tg_id", "неверный формат Telegram ID")
		}
	}

	return result
}

// ValidateChatID валидирует Chat ID
func (v *Validator) ValidateChatID(chatID string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(chatID, "chat_id"); !requiredResult.IsValid {
		return requiredResult
	}

	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		result.AddError("chat_id", "Chat ID должен быть числом")
		return result
	}

	if id <= 0 {
		result.AddError("chat_id", "Chat ID должен быть положительным числом")
	}

	return result
}

// ValidateTrainerInfo валидирует информацию о тренере
func (v *Validator) ValidateTrainerInfo(info string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(info, "info"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем длину
	if lengthResult := v.validateStringLength(info, "info", 10, 500); !lengthResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, lengthResult.Errors...)
	}

	return result
}

// ValidateTrackName валидирует название трассы
func (v *Validator) ValidateTrackName(name string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(name, "name"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем длину
	if lengthResult := v.validateStringLength(name, "name", 2, 100); !lengthResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, lengthResult.Errors...)
	}

	return result
}

// ValidateTrackInfo валидирует информацию о трассе
func (v *Validator) ValidateTrackInfo(info string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(info, "info"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем длину
	if lengthResult := v.validateStringLength(info, "info", 10, 500); !lengthResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, lengthResult.Errors...)
	}

	return result
}

// ValidateUserName валидирует имя пользователя
func (v *Validator) ValidateUserName(name string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(name, "name"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем длину
	if lengthResult := v.validateStringLength(name, "name", 2, 50); !lengthResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, lengthResult.Errors...)
	}

	// Проверяем недопустимые символы
	if charsResult := v.validateInvalidChars(name, "name"); !charsResult.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, charsResult.Errors...)
	}

	return result
}

// ValidateTime валидирует время в формате HH:MM
func (v *Validator) ValidateTime(timeStr string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(timeStr, "time"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем формат HH:MM
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timeRegex.MatchString(timeStr) {
		result.AddError("time", "неверный формат времени. Используйте HH:MM")
		return result
	}

	// Парсим время для дополнительной проверки
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		result.AddError("time", "неверное время")
	}

	return result
}

// ValidateDateTime валидирует дату и время в формате 2006-01-02 15:04
func (v *Validator) ValidateDateTime(dateTimeStr string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(dateTimeStr, "datetime"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем формат даты и времени
	dateTimeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`)
	if !dateTimeRegex.MatchString(dateTimeStr) {
		result.AddError("datetime", "неверный формат даты и времени. Используйте YYYY-MM-DD HH:MM")
		return result
	}

	// Парсим дату и время для дополнительной проверки
	_, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		result.AddError("datetime", "неверная дата или время")
	}

	return result
}

// ValidateMaxParticipants валидирует максимальное количество участников
func (v *Validator) ValidateMaxParticipants(participantsStr string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(participantsStr, "max_participants"); !requiredResult.IsValid {
		return requiredResult
	}

	participants, err := strconv.Atoi(participantsStr)
	if err != nil {
		result.AddError("max_participants", "количество участников должно быть числом")
		return result
	}

	if participants < 1 {
		result.AddError("max_participants", "количество участников должно быть больше 0")
	}

	if participants > 100 {
		result.AddError("max_participants", "количество участников не должно превышать 100")
	}

	return result
}

// ValidateID валидирует ID
func (v *Validator) ValidateID(idStr string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(idStr, "id"); !requiredResult.IsValid {
		return requiredResult
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.AddError("id", "ID должен быть числом")
		return result
	}

	if id == 0 {
		result.AddError("id", "ID должен быть больше 0")
	}

	return result
}

// ValidateMessageText валидирует текст сообщения
func (v *Validator) ValidateMessageText(text string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Проверяем обязательность
	if requiredResult := v.validateRequired(text, "text"); !requiredResult.IsValid {
		return requiredResult
	}

	// Проверяем максимальную длину
	if len(text) > 4096 {
		result.AddError("text", "текст сообщения не должен превышать 4096 символов")
	}

	return result
}
