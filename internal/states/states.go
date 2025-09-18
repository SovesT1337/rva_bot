package states

// StateType определяет тип состояния пользователя
type StateType int

const (
	// Основные состояния
	StateStart StateType = iota
	StateError
	StateStartKeyboard
	StateAdminKeyboard

	// Состояния создания тренера
	StateEnterTrainerName
	StateEnterTrainerTgId
	StateEnterTrainerChatId
	StateEnterTrainerInfo
	StateConfirmTrainerCreation

	// Состояния редактирования тренера
	StateSelectTrainerToEdit
	StateEditTrainerName
	StateEditTrainerTgId
	StateEditTrainerInfo
	StateConfirmTrainerEdit
	StateConfirmTrainerDelete

	// Состояния создания трассы
	StateEnterTrackName
	StateEnterTrackInfo
	StateConfirmTrackCreation

	// Состояния редактирования трассы
	StateSelectTrackToEdit
	StateEditTrackName
	StateEditTrackInfo
	StateConfirmTrackEdit
	StateConfirmTrackDelete

	// Состояния создания расписания

	// Состояния регистрации пользователя
	StateEnterUserName
	StateConfirmUserRegistration

	// Состояния создания тренировки
	StateEnterTrainingTrainer
	StateEnterTrainingTrack
	StateEnterTrainingDate
	StateEnterTrainingMaxParticipants
	StateConfirmTrainingCreation

	// Состояния регистрации на тренировку
	StateSelectTrainingToRegister
	StateConfirmTrainingRegistration

	// Новые состояния для пошаговой записи на тренировку
	StateSelectTrackForRegistration
	StateSelectTrainerForRegistration
	StateSelectTrainingTimeForRegistration
)

// State представляет состояние пользователя в боте
type State struct {
	Type StateType
	Data map[string]interface{} // Дополнительные данные состояния
}

// TempTrainerData представляет временные данные тренера при создании
type TempTrainerData struct {
	Name   string
	TgId   string
	ChatId int
	Info   string
}

// TempTrackData представляет временные данные трассы при создании
type TempTrackData struct {
	Name string
	Info string
}

// TempUserData представляет временные данные пользователя при регистрации
type TempUserData struct {
	Name string
}

// TempTrainingData представляет временные данные тренировки при создании
type TempTrainingData struct {
	TrainerID       uint
	TrackID         uint
	Date            string
	MaxParticipants int
}

// TempRegistrationData представляет временные данные для записи на тренировку
type TempRegistrationData struct {
	TrackID   uint
	TrainerID uint
}

// NewState создает новое состояние
func NewState(stateType StateType, data map[string]interface{}) State {
	if data == nil {
		data = make(map[string]interface{})
	}
	return State{Type: stateType, Data: data}
}

// GetID возвращает ID из данных состояния
func (s State) GetID() uint {
	if id, ok := s.Data["id"].(uint); ok {
		return id
	}
	return 0
}

// SetID устанавливает ID в данные состояния
func (s State) SetID(id uint) State {
	s.Data["id"] = id
	return s
}

// GetString возвращает строковое значение из данных состояния
func (s State) GetString(key string) string {
	if val, ok := s.Data[key].(string); ok {
		return val
	}
	return ""
}

// SetString устанавливает строковое значение в данные состояния
func (s State) SetString(key, value string) State {
	s.Data[key] = value
	return s
}

// GetTempTrainerData возвращает временные данные тренера из состояния
func (s State) GetTempTrainerData() *TempTrainerData {
	if data, ok := s.Data["tempTrainer"].(*TempTrainerData); ok {
		return data
	}
	return &TempTrainerData{}
}

// SetTempTrainerData устанавливает временные данные тренера в состояние
func (s State) SetTempTrainerData(data *TempTrainerData) State {
	s.Data["tempTrainer"] = data
	return s
}

// GetTempTrackData возвращает временные данные трассы из состояния
func (s State) GetTempTrackData() *TempTrackData {
	if data, ok := s.Data["tempTrack"].(*TempTrackData); ok {
		return data
	}
	return &TempTrackData{}
}

// SetTempTrackData устанавливает временные данные трассы в состояние
func (s State) SetTempTrackData(data *TempTrackData) State {
	s.Data["tempTrack"] = data
	return s
}

// GetTempUserData возвращает временные данные пользователя из состояния
func (s State) GetTempUserData() *TempUserData {
	if data, ok := s.Data["tempUser"].(*TempUserData); ok {
		return data
	}
	return &TempUserData{}
}

// SetTempUserData устанавливает временные данные пользователя в состояние
func (s State) SetTempUserData(data *TempUserData) State {
	s.Data["tempUser"] = data
	return s
}

// GetTempTrainingData возвращает временные данные тренировки из состояния
func (s State) GetTempTrainingData() *TempTrainingData {
	if data, ok := s.Data["tempTraining"].(*TempTrainingData); ok {
		return data
	}
	return &TempTrainingData{}
}

// SetTempTrainingData устанавливает временные данные тренировки в состояние
func (s State) SetTempTrainingData(data *TempTrainingData) State {
	s.Data["tempTraining"] = data
	return s
}

// GetTempRegistrationData возвращает временные данные регистрации из состояния
func (s State) GetTempRegistrationData() *TempRegistrationData {
	if data, ok := s.Data["tempRegistration"].(*TempRegistrationData); ok {
		return data
	}
	return &TempRegistrationData{}
}

// SetTempRegistrationData устанавливает временные данные регистрации в состояние
func (s State) SetTempRegistrationData(data *TempRegistrationData) State {
	s.Data["tempRegistration"] = data
	return s
}

// Конструкторы состояний
func SetStart() State {
	return NewState(StateStart, nil)
}

func SetError() State {
	return NewState(StateError, nil)
}

func SetStartKeyboard() State {
	return NewState(StateStartKeyboard, nil)
}

func SetAdminKeyboard() State {
	return NewState(StateAdminKeyboard, nil)
}

func SetEnterTrainerName(id uint) State {
	return NewState(StateEnterTrainerName, map[string]interface{}{"id": id})
}

func SetEnterTrainerTgId(id uint) State {
	return NewState(StateEnterTrainerTgId, map[string]interface{}{"id": id})
}

func SetEnterTrainerChatId(id uint) State {
	return NewState(StateEnterTrainerChatId, map[string]interface{}{"id": id})
}

func SetEnterTrainerInfo(id uint) State {
	return NewState(StateEnterTrainerInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrainerCreation() State {
	return NewState(StateConfirmTrainerCreation, nil)
}

// Конструкторы для редактирования тренеров
func SetSelectTrainerToEdit() State {
	return NewState(StateSelectTrainerToEdit, nil)
}

func SetEditTrainerName(id uint) State {
	return NewState(StateEditTrainerName, map[string]interface{}{"id": id})
}

func SetEditTrainerTgId(id uint) State {
	return NewState(StateEditTrainerTgId, map[string]interface{}{"id": id})
}

func SetEditTrainerInfo(id uint) State {
	return NewState(StateEditTrainerInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrainerEdit(id uint) State {
	return NewState(StateConfirmTrainerEdit, map[string]interface{}{"id": id})
}

func SetConfirmTrainerDelete(id uint) State {
	return NewState(StateConfirmTrainerDelete, map[string]interface{}{"id": id})
}

func SetEnterTrackName(id uint) State {
	return NewState(StateEnterTrackName, map[string]interface{}{"id": id})
}

func SetEnterTrackInfo(id uint) State {
	return NewState(StateEnterTrackInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrackCreation() State {
	return NewState(StateConfirmTrackCreation, nil)
}

// Конструкторы для редактирования трасс
func SetSelectTrackToEdit() State {
	return NewState(StateSelectTrackToEdit, nil)
}

func SetEditTrackName(id uint) State {
	return NewState(StateEditTrackName, map[string]interface{}{"id": id})
}

func SetEditTrackInfo(id uint) State {
	return NewState(StateEditTrackInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrackEdit(id uint) State {
	return NewState(StateConfirmTrackEdit, map[string]interface{}{"id": id})
}

func SetConfirmTrackDelete(id uint) State {
	return NewState(StateConfirmTrackDelete, map[string]interface{}{"id": id})
}

// Конструкторы для регистрации пользователя
func SetEnterUserName() State {
	return NewState(StateEnterUserName, nil)
}

func SetConfirmUserRegistration() State {
	return NewState(StateConfirmUserRegistration, nil)
}

// Конструкторы для создания тренировки
func SetEnterTrainingTrainer(id uint) State {
	return NewState(StateEnterTrainingTrainer, map[string]interface{}{"id": id})
}

func SetEnterTrainingTrack(id uint) State {
	return NewState(StateEnterTrainingTrack, map[string]interface{}{"id": id})
}

func SetEnterTrainingDate(id uint) State {
	return NewState(StateEnterTrainingDate, map[string]interface{}{"id": id})
}

func SetEnterTrainingMaxParticipants(id uint) State {
	return NewState(StateEnterTrainingMaxParticipants, map[string]interface{}{"id": id})
}

func SetConfirmTrainingCreation() State {
	return NewState(StateConfirmTrainingCreation, nil)
}

// Конструкторы для регистрации на тренировку
func SetSelectTrainingToRegister() State {
	return NewState(StateSelectTrainingToRegister, nil)
}

func SetConfirmTrainingRegistration(trainingId uint) State {
	return NewState(StateConfirmTrainingRegistration, map[string]interface{}{"trainingId": trainingId})
}

// Конструкторы для пошаговой записи на тренировку
func SetSelectTrackForRegistration() State {
	return NewState(StateSelectTrackForRegistration, nil)
}

func SetSelectTrainerForRegistration() State {
	return NewState(StateSelectTrainerForRegistration, nil)
}

func SetSelectTrainingTimeForRegistration() State {
	return NewState(StateSelectTrainingTimeForRegistration, nil)
}
