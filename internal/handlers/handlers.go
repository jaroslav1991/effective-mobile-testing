package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"

	"effective_mobile_testing/internal/model"
	"effective_mobile_testing/internal/service/repository"
	"effective_mobile_testing/internal/validators"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	service HandlerInterface
}

type HandlerInterface interface {
	CreateUser(passportNumber string, user model.UserFromAPI) (*model.User, error)
	GetUserData(passportSerie, passportNumber string) (model.UserFromAPI, error)
	StartTracking(req model.RequestStartTracking) error
	StopTracking(req model.RequestStopTracking) error
	GetLaborCosts(userID int64) ([]model.ResponseLobarCost, error)
	GetUserByFilters(limit, offset int, user model.User) (*[]model.User, error)
	DeleteUser(id int64) error
	UpdateUser(id int64, user model.UserUpdateRequest) (*model.User, error)
}

func NewHandlers(service HandlerInterface) *Handlers {
	return &Handlers{service}
}

// @Summary      Create User
// @Description  create user by passportNumber
// @Tags         users
// @Accept       json
// @Produce      json
// @Param  		 input body model.CreateUserRequest true "passport number" example("1234 123456")
// @Router		 /user/create/ [post]
func (h *Handlers) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "createUser"

		var req model.CreateUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Error(fmt.Sprintf("%s error with binding request:%v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		passport := strings.Split(req.PassportNumber, " ")

		if !validators.CheckLenPassport(passport) {
			slog.Error(fmt.Sprintf("%s passport not valid", handler))
			c.JSON(http.StatusBadRequest, gin.H{"error": "passport not valid"})
			return
		}

		userFromAPI, err := h.service.GetUserData(passport[0], passport[1])
		if err != nil {
			slog.Error("Error getting user from API:", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		totalPassport := passport[0] + " " + passport[1]

		user, err := h.service.CreateUser(totalPassport, userFromAPI)
		if err != nil {
			if errors.As(err, &repository.ErrUserExists) {
				slog.Error(fmt.Sprintf("%s error create user: %v", handler, err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
				return
			}
			slog.Error("error creating user:", slog.String("error", err.Error()))
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
			return
		}

		c.JSON(http.StatusOK, user)
		slog.Debug(fmt.Sprintf("%s created user", handler))
	}
}

// @Summary      Start tracking
// @Description  start time for task
// @Tags         users
// @Accept       json
// @Produce      json
// @Param  		 input body model.RequestStartTracking true "choose task and user"
// @Router		 /user/start-tracking/ [patch]
func (h *Handlers) StartTracking() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "startTracking"

		var req model.RequestStartTracking

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Error(fmt.Sprintf("%s error with binding request:%v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if err := h.service.StartTracking(req); err != nil {
			if errors.As(err, &repository.ErrUserNotFound) {
				slog.Error(fmt.Sprintf("%s error start tracking: %v", handler, err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
				return
			}
			slog.Error("error start tracking:", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "start tracking"})
		slog.Debug(fmt.Sprintf("%s started tracking", handler))
	}
}

// @Summary      Stop tracking
// @Description  stop time for task
// @Tags         users
// @Accept       json
// @Produce      json
// @Param  		 input body model.RequestStopTracking true "choose task and user"
// @Router		 /user/stop-tracking/ [patch]
func (h *Handlers) StopTracking() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "stopTracking"

		var req model.RequestStopTracking

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Error(fmt.Sprintf("%s error with binding request:%v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if err := h.service.StopTracking(req); err != nil {
			if err != sql.ErrNoRows {
				slog.Error("error stop tracking:", slog.String("error", err.Error()))
				c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
				return
			}

			slog.Error("error stop tracking:", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "stop tracking"})
		slog.Debug(fmt.Sprintf("%s stop tracking", handler))
	}
}

// @Summary      Get labor cost
// @Description  get info about working
// @Tags         users
// @Accept       json
// @Produce      json
// @Param  		 id query string true "User ID"
// @Router		 /user/get-costs/ [get]
func (h *Handlers) GetLaborCosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "getLaborCosts"

		userID := c.Query("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error with binding request:%v", handler, err))
			//log.Printf("%s error convertion: %v", handler, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
			return
		}

		costs, err := h.service.GetLaborCosts(int64(id))
		if err != nil {
			if err != sql.ErrNoRows {
				slog.Error(fmt.Sprintf("%s error get labor costs: %v", handler, err))
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}

			slog.Error("error get labor costs:", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"costs": costs})
		slog.Debug(fmt.Sprintf("%s labor costs finished", handler))
	}
}

// @Summary      Get Users
// @Description  Get info by any filters
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id             query string false "User ID"
// @Param        surname        query string false "Surname"
// @Param        name           query string false "Name"
// @Param        patronymic     query string false "Patronymic"
// @Param        address        query string false "Address"
// @Param        passport_number query string false "Passport Number"
// @Param        limit          query int    false "Limit"
// @Param        offset         query int    false "Offset"
// @Router       /users/ [get]
func (h *Handlers) GetUserByFilters() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "geUserByFilters"

		var (
			idInt     int
			limitInt  int
			offsetInt int
		)

		userID := c.Query("id")
		if userID != "" {
			id, err := strconv.Atoi(userID)
			if err != nil {
				slog.Error(fmt.Sprintf("%s error convertion:%v", handler, err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
				return
			}
			idInt = id
		}

		surname := c.Query("surname")
		name := c.Query("name")
		patronymic := c.Query("patronymic")
		address := c.Query("address")
		passportNumber := c.Query("passport_number")

		l := c.Query("limit")
		if l != "" {
			limit, err := strconv.Atoi(l)
			if err != nil {
				slog.Error(fmt.Sprintf("%s error covertion: %v", handler, err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
				return
			}
			limitInt = limit
		}
		o := c.Query("offset")
		if o != "" {
			offset, err := strconv.Atoi(o)
			if err != nil {
				slog.Error(fmt.Sprintf("%s error covertion: %v", handler, err))
				log.Printf("%s error convertion: %v", handler, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
				return
			}
			offsetInt = offset
		}

		var user model.User

		user.ID = int64(idInt)
		user.Surname = surname
		user.Name = name
		user.Patronymic = patronymic
		user.Address = address
		user.PassportNumber = passportNumber

		users, err := h.service.GetUserByFilters(limitInt, offsetInt, user)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error get user by filter: %v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
		slog.Debug(fmt.Sprintf("%s get user by filter finished", handler))

	}
}

// @Summary      Delete User
// @Description  delete user data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param  		 id query string true "ID"
// @Router		 /user/ [delete]
func (h *Handlers) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "deleteUser"

		userID := c.Query("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error convertion: %v", handler, err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
			return
		}

		if err := h.service.DeleteUser(int64(id)); err != nil {
			slog.Error(fmt.Sprintf("%s error deleting user: %v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "status bad request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "user deleted"})
		slog.Debug(fmt.Sprintf("%s user deleted success", handler))

	}
}

// @Summary      Update User
// @Description  update user data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id             query string false "User ID"
// @Param  		 input body model.UserUpdateRequest false "user udpate data"
// @Router       /user/ [patch]
func (h *Handlers) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		const handler = "updateUser"

		userID := c.Query("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error convertion: %v", handler, err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "status bad request"})
			return
		}

		var user model.UserUpdateRequest

		if err := c.ShouldBindJSON(&user); err != nil {
			slog.Error(fmt.Sprintf("%s error binding json request: %v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		updateUser, err := h.service.UpdateUser(int64(id), user)
		if err != nil {
			if err != sql.ErrNoRows {
				slog.Error(fmt.Sprintf("%s error update user: %v", handler, err))
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}

			slog.Error(fmt.Sprintf("%s error update user: %v", handler, err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "status bad request"})
			return
		}

		c.JSON(http.StatusOK, updateUser)
		slog.Debug(fmt.Sprintf("%s user updated success", handler))

	}
}
