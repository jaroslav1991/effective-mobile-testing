package service

import (
	"database/sql"
	"effective_mobile_testing/internal/model"
	"effective_mobile_testing/internal/service/repository"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/joho/godotenv"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func (s *UserTaskService) CreateUser(passportNumber string, user model.UserFromAPI) (*model.User, error) {
	createdUser, err := s.repo.CreateUser(user.Surname, user.Name, user.Patronymic, user.Address, passportNumber)
	if err != nil {
		if errors.As(err, &repository.ErrUserExists) {
			return nil, repository.ErrUserExists
		}

		slog.Error("can't create user", slog.String("err", err.Error()))
		return nil, err
	}

	return createdUser, nil
}

func (s *UserTaskService) GetUserData(passportSerie, passportNumber string) (model.UserFromAPI, error) {
	baseUrl := os.Getenv("ANOTHER_API_URL")
	if baseUrl == "" {
		slog.Error("can't get ANOTHER_API_URL environment variable")
		return model.UserFromAPI{}, errors.New("can't get ANOTHER_API_URL")
	}

	fullUrl := fmt.Sprintf("%s?passportSerie=%s&passportNumber=%s", baseUrl, passportSerie, passportNumber)

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		slog.Error("error creating new request:", slog.String("err", err.Error()))
		return model.UserFromAPI{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("error making request", slog.String("err", err.Error()))
		return model.UserFromAPI{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response body:", slog.String("err", err.Error()))
		return model.UserFromAPI{}, err
	}

	var user model.UserFromAPI

	if err := json.Unmarshal(body, &user); err != nil {
		slog.Error("error parsing response body", slog.String("err", err.Error()))
		return model.UserFromAPI{}, err
	}

	return user, nil
}

func (s *UserTaskService) StartTracking(req model.RequestStartTracking) error {
	if err := s.repo.StartTask(req.UserID, req.TaskName); err != nil {
		if errors.As(err, &repository.ErrUserNotFound) {
			return fmt.Errorf("%v", repository.ErrUserNotFound)
		}
		slog.Error("can't start tracking:", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *UserTaskService) StopTracking(req model.RequestStopTracking) error {
	if err := s.repo.StopTask(req.UserID, req.TaskName); err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("%v", repository.ErrUserNotFound)
		}
		slog.Error("can't stop tracking:", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *UserTaskService) GetLaborCosts(userID int64) ([]model.ResponseLobarCost, error) {
	resp, err := s.repo.GetLaborCosts(userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("%v", repository.ErrUserNotFound)
		}
		slog.Error("cant' get labor costs", slog.String("err", err.Error()))
		return nil, err
	}

	for i := range resp {
		hours := resp[i].DurationMinutes / 60
		resp[i].DurationHours = hours
		resp[i].DurationMinutes = resp[i].DurationMinutes % 60
	}

	return resp, nil
}

func (s *UserTaskService) GetUserByFilters(limit, offset int, user model.User) (*[]model.User, error) {
	users, err := s.repo.GetUserByFilters(limit, offset, user.ID, user.Surname, user.Name, user.Patronymic, user.Address, user.PassportNumber)
	if err != nil {
		slog.Error("can't get user by filters", slog.String("err", err.Error()))
		return nil, err
	}

	return users, nil
}

func (s *UserTaskService) DeleteUser(id int64) error {
	if err := s.repo.DeleteUser(id); err != nil {
		slog.Error("can't delete user", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *UserTaskService) UpdateUser(id int64, user model.UserUpdateRequest) (*model.User, error) {
	updateUser, err := s.repo.UpdateUser(id, user.Surname, user.Name, user.Patronymic, user.Address, user.PassportNumber)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("%v", repository.ErrUserNotFound)
		}

		slog.Error("can't update user", slog.String("err", err.Error()))
		return nil, err
	}

	return updateUser, nil
}
