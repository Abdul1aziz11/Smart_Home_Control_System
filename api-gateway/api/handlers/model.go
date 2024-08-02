package handlers

import "api-gateway/genproto"

type Handlers struct {
	Auth          genproto.AuthServiceClient
	DeviceControl genproto.DeviceServiceClient
}

// Register

// @UserRegister godoc
type UserRegister struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Profile  Profile `json:"profile"`
}

// @Profile godoc
type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
}

type RespStatus struct {
	Status string `json:"status"`
}

// Verify
type UserVerify struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RespStatusUserID struct {
	UserId string `json:"user_id"`
	Status string `json:"status"`
}

// SignIn
type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRespFull struct {
	UserId string `json:"user_id"`
	Status string `json:"status"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

// Profile
type UserProfile struct {
	Email string `json:"email"`
}

type UserProfileResp struct {
	UserId       string `json:"user_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// CreateDevice
type DeviceRegister struct {
	UserId             string `json:"user_id"`
	DeviceType         string `json:"device_type"`
	DeviceName         string `json:"device_name"`
	DeviceStatus       string `json:"device_status"`
	Location           string `json:"location"`
	FirmwareVersion    string `json:"firmware_version"`
	ConnectivityStatus string `json:"connectivity_status"`
}

type Device struct {
	DeviceId              string            `json:"device_id"`
	UserId                string            `json:"user_id"`
	DeviceType            string            `json:"device_type"`
	DeviceName            string            `json:"device_name"`
	DeviceStatus          string            `json:"device_status"`
	ConfigurationSettings map[string]string `json:"configuration_settings"`
	LastUpdated           string            `json:"last_updated"`
	FirmwareVersion       string            `json:"firmware_version"`
	ConnectivityStatus    string            `json:"connectivity_status"`
}

// DeleteDevice
type DeleteResp struct {
	Success string `json:"success"`
}

// ControlDevice
type Command struct {
	DeviceId       string            `json:"device_id"`
	UserId         string            `json:"user_id"`
	CommandType    string            `json:"command_type"`
	CommandPayload map[string]string `json:"command_payload"`
	Status         string            `json:"status"`
}
