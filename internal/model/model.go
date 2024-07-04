package model

import (
	"time"
)

//func (u User) Value() (driver.Value, error) {
//	data, err := json.Marshal(u)
//	if err != nil {
//		return nil, err
//	}
//	return string(data), nil
//}
//
//func (u *User) Scan(src interface{}) error {
//	if data, ok := src.([]byte); ok {
//		return json.Unmarshal(data, &u)
//	}
//	return nil
//}

type User struct {
	ID             int64  `json:"id,omitempty"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address"`
	PassportNumber string `json:"passport_number,omitempty"`
}

type Task struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	StartTracking time.Time `json:"start_tracking"`
	StopTracking  time.Time `json:"stop_tracking"`
	UserID        int64     `json:"user_id"`
}

type UserFromAPI struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic,omitempty"`
	Address    string `json:"address"`
}

type RequestLaborCost struct {
	UserID int64 `json:"user_id"`
}

type ResponseLobarCost struct {
	TaskName        string `json:"task_name"`
	DurationHours   int    `json:"duration_hours,omitempty"`
	DurationMinutes int    `json:"duration_minutes"`
}

type RequestStartTracking struct {
	TaskName string `json:"task_name"`
	UserID   int64  `json:"user_id"`
}

type RequestStopTracking struct {
	TaskName string `json:"task_name"`
	UserID   int64  `json:"user_id"`
}

type CreateUserRequest struct {
	PassportNumber string `json:"passportNumber"`
}

type Args struct {
	ID             int64  `json:"id,omitempty"`
	Surname        string `json:"surname,omitempty"`
	Name           string `json:"name,omitempty"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address,omitempty"`
	PassportNumber string `json:"passport_number	,omitempty"`
}

type UserUpdateRequest struct {
	Surname        string `json:"surname,omitempty"`
	Name           string `json:"name,omitempty"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
}
