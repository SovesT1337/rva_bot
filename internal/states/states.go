package states

type StateType string

const (
	StateStart         = "StateStart"
	StateError         = "StateError"
	StateStartKeyboard = "StateStartKeyboard"
	StateAdminKeyboard = "StateAdminKeyboard"

	StateSetTrainerName         = "StateEnterTrainerName"
	StateSetTrainerTgId         = "StateSetTrainerTgId"
	StateSetTrainerChatId       = "StateSetTrainerChatId"
	StateSetTrainerInfo         = "StateSetTrainerInfo"
	StateConfirmTrainerCreation = "StateConfirmTrainerCreation"

	StateEditTrainerName      = "StateEditTrainerName"
	StateEditTrainerTgId      = "StateEditTrainerTgId"
	StateEditTrainerInfo      = "StateEditTrainerInfo"
	StateConfirmTrainerDelete = "StateConfirmTrainerDelete"

	StateSetTrackName         = "StateSetTrackName"
	StateSetTrackInfo         = "StateSetTrackInfo"
	StateConfirmTrackCreation = "StateConfirmTrackCreation"

	StateEditTrackName      = "StateEditTrackName"
	StateEditTrackInfo      = "StateEditTrackInfo"
	StateConfirmTrackDelete = "StateConfirmTrackDelete"

	StateSetUserName             = "StateSetUserName"
	StateConfirmUserRegistration = "StateConfirmUserRegistration"

	StateSetTrainingTrack           = "StateSetTrainingTrack"
	StateSetTrainingTrainer         = "StateSetTrainingTrainer"
	StateSetTrainingStartTime       = "StateSetTrainingStartTime"
	StateSetTrainingEndTime         = "StateSetTrainingEndTime"
	StateSetTrainingMaxParticipants = "StateSetTrainingMaxParticipants"
	StateConfirmTrainingCreation    = "StateConfirmTrainingCreation"

	StateConfirmTrainingRegistration = "StateConfirmTrainingRegistration"
	StateConfirmTrainingDelete       = "StateConfirmTrainingDelete"

	StateSelectTrackForRegistration        = "StateSelectTrackForRegistration"
	StateSelectTrainerForRegistration      = "StateSelectTrainerForRegistration"
	StateSelectTrainingTimeForRegistration = "StateSelectTrainingTimeForRegistration"
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
	StartTime       string
	EndTime         string
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
	return NewState(StateSetTrainerName, map[string]interface{}{"id": id})
}

func SetEnterTrainerTgId(id uint) State {
	return NewState(StateSetTrainerTgId, map[string]interface{}{"id": id})
}

func SetEnterTrainerChatId(id uint) State {
	return NewState(StateSetTrainerChatId, map[string]interface{}{"id": id})
}

func SetEnterTrainerInfo(id uint) State {
	return NewState(StateSetTrainerInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrainerCreation() State {
	return NewState(StateConfirmTrainerCreation, nil)
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

func SetConfirmTrainerDelete(id uint) State {
	return NewState(StateConfirmTrainerDelete, map[string]interface{}{"id": id})
}

func SetEnterTrackName(id uint) State {
	return NewState(StateSetTrackName, map[string]interface{}{"id": id})
}

func SetEnterTrackInfo(id uint) State {
	return NewState(StateSetTrackInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrackCreation() State {
	return NewState(StateConfirmTrackCreation, nil)
}

func SetEditTrackName(id uint) State {
	return NewState(StateEditTrackName, map[string]interface{}{"id": id})
}

func SetEditTrackInfo(id uint) State {
	return NewState(StateEditTrackInfo, map[string]interface{}{"id": id})
}

func SetConfirmTrackDelete(id uint) State {
	return NewState(StateConfirmTrackDelete, map[string]interface{}{"id": id})
}

func SetEnterUserName() State {
	return NewState(StateSetUserName, nil)
}

func SetConfirmUserRegistration() State {
	return NewState(StateConfirmUserRegistration, nil)
}

func SetSetTrainingTrainer(id uint) State {
	return NewState(StateSetTrainingTrainer, map[string]interface{}{"id": id})
}

func SetSetTrainingTrack(id uint) State {
	return NewState(StateSetTrainingTrack, map[string]interface{}{"id": id})
}

func SetSetTrainingStartTime(id uint) State {
	return NewState(StateSetTrainingStartTime, map[string]interface{}{"id": id})
}

func SetSetTrainingEndTime(id uint) State {
	return NewState(StateSetTrainingEndTime, map[string]interface{}{"id": id})
}

func SetSetTrainingMaxParticipants(id uint) State {
	return NewState(StateSetTrainingMaxParticipants, map[string]interface{}{"id": id})
}

func SetConfirmTrainingCreation() State {
	return NewState(StateConfirmTrainingCreation, nil)
}

func SetConfirmTrainingRegistration(trainingId uint) State {
	return NewState(StateConfirmTrainingRegistration, map[string]interface{}{"trainingId": trainingId})
}

func SetConfirmTrainingDelete(trainingId uint) State {
	return NewState(StateConfirmTrainingDelete, map[string]interface{}{"id": trainingId})
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
