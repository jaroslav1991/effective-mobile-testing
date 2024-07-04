package model

import (
	"time"
)

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
