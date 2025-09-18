package states

type StateType int

const (
	StateStart StateType = iota
	StateError
	StateStartKeyboard
	StateAdminKeyboard

	StateEnterTrainerName
	StateEnterTrainerTgId
	StateEnterTrainerChatId
	StateEnterTrainerInfo
	StateConfirmTrainerCreation

	StateSelectTrainerToEdit
	StateEditTrainerName
	StateEditTrainerTgId
	StateEditTrainerInfo
	StateConfirmTrainerEdit
	StateConfirmTrainerDelete

	StateEnterTrackName
	StateEnterTrackInfo
	StateConfirmTrackCreation

	StateSelectTrackToEdit
	StateEditTrackName
	StateEditTrackInfo
	StateConfirmTrackEdit
	StateConfirmTrackDelete


	StateEnterUserName
	StateConfirmUserRegistration

	StateEnterTrainingTrainer
	StateEnterTrainingTrack
	StateEnterTrainingDate
	StateEnterTrainingMaxParticipants
	StateConfirmTrainingCreation

	StateSelectTrainingToRegister
	StateConfirmTrainingRegistration

	StateSelectTrackForRegistration
	StateSelectTrainerForRegistration
	StateSelectTrainingTimeForRegistration
)

type State struct {
	Type StateType
	Data map[string]interface{}
}

type TempTrainerData struct {
	Name   string
	TgId   string
	ChatId int
	Info   string
}

type TempTrackData struct {
	Name string
	Info string
}

type TempUserData struct {
	Name string
}

type TempTrainingData struct {
	TrainerID       uint
	TrackID         uint
	Date            string
	MaxParticipants int
}

type TempRegistrationData struct {
	TrackID   uint
	TrainerID uint
}

func NewState(stateType StateType, data map[string]interface{}) State {
	if data == nil {
		data = make(map[string]interface{})
	}
	return State{Type: stateType, Data: data}
}

func (s State) GetID() uint {
	if id, ok := s.Data["id"].(uint); ok {
		return id
	}
	return 0
}

func (s State) SetID(id uint) State {
	s.Data["id"] = id
	return s
}

func (s State) GetString(key string) string {
	if val, ok := s.Data[key].(string); ok {
		return val
	}
	return ""
}

func (s State) SetString(key, value string) State {
	s.Data[key] = value
	return s
}

func (s State) GetTempTrainerData() *TempTrainerData {
	if data, ok := s.Data["tempTrainer"].(*TempTrainerData); ok {
		return data
	}
	return &TempTrainerData{}
}

func (s State) SetTempTrainerData(data *TempTrainerData) State {
	s.Data["tempTrainer"] = data
	return s
}

func (s State) GetTempTrackData() *TempTrackData {
	if data, ok := s.Data["tempTrack"].(*TempTrackData); ok {
		return data
	}
	return &TempTrackData{}
}

func (s State) SetTempTrackData(data *TempTrackData) State {
	s.Data["tempTrack"] = data
	return s
}

func (s State) GetTempUserData() *TempUserData {
	if data, ok := s.Data["tempUser"].(*TempUserData); ok {
		return data
	}
	return &TempUserData{}
}

func (s State) SetTempUserData(data *TempUserData) State {
	s.Data["tempUser"] = data
	return s
}

func (s State) GetTempTrainingData() *TempTrainingData {
	if data, ok := s.Data["tempTraining"].(*TempTrainingData); ok {
		return data
	}
	return &TempTrainingData{}
}

func (s State) SetTempTrainingData(data *TempTrainingData) State {
	s.Data["tempTraining"] = data
	return s
}

func (s State) GetTempRegistrationData() *TempRegistrationData {
	if data, ok := s.Data["tempRegistration"].(*TempRegistrationData); ok {
		return data
	}
	return &TempRegistrationData{}
}

func (s State) SetTempRegistrationData(data *TempRegistrationData) State {
	s.Data["tempRegistration"] = data
	return s
}

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

func SetEnterUserName() State {
	return NewState(StateEnterUserName, nil)
}

func SetConfirmUserRegistration() State {
	return NewState(StateConfirmUserRegistration, nil)
}

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

func SetSelectTrainingToRegister() State {
	return NewState(StateSelectTrainingToRegister, nil)
}

func SetConfirmTrainingRegistration(trainingId uint) State {
	return NewState(StateConfirmTrainingRegistration, map[string]interface{}{"trainingId": trainingId})
}

func SetSelectTrackForRegistration() State {
	return NewState(StateSelectTrackForRegistration, nil)
}

func SetSelectTrainerForRegistration() State {
	return NewState(StateSelectTrainerForRegistration, nil)
}

func SetSelectTrainingTimeForRegistration() State {
	return NewState(StateSelectTrainingTimeForRegistration, nil)
}
