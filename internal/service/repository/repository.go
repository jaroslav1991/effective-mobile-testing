package repository

import (
	"database/sql"
	"effective_mobile_testing/internal/db"
	"effective_mobile_testing/internal/model"
	"effective_mobile_testing/internal/validators"
	"errors"
	"fmt"
	"time"
)

const (
	createUserQuery = `insert into person (surname, name, patronymic, address, passport_number)
						values ($1, $2, $3, $4, $5) returning surname, name,patronymic,address`
	startTaskQuery         = `insert into task (name, start_tracking, user_id) values ($1, $2, $3)`
	checkUserIDTaskQuery   = `select user_id from task where user_id = $1`
	checkUserIDPersonQuery = `select id from person where id = $1`
	stopTaskQuery          = `update task set stop_tracking=$1 where user_id=$2 and name=$3`
	getLaborCosts          = `select name, floor(EXTRACT(EPOCH from (stop_tracking - start_tracking)) / 60) as duration from task 
						where user_id=$1 and stop_tracking is not null  order by duration desc `
	deleteFromTaskQuery   = `delete from task where user_id = $1`
	deleteFromPersonQuery = `delete from person where id = $1`
	selectForUpdateQuery  = `select surname, name, patronymic, address, passport_number from person where id =$1 for update`
	updatePersonQuery     = `update person set surname=$1, name=$2, address=$3, patronymic=$4, passport_number=$5 where id = $6
								returning id, surname, name, patronymic, address, passport_number`
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UserTaskRepo struct {
	DB db.DB
}

func NewUserTaskRepo(db db.DB) *UserTaskRepo {
	return &UserTaskRepo{DB: db}
}

func (repo *UserTaskRepo) CreateUser(surname, name, patronymic, address, passportNumber string) (*model.User, error) {
	var user model.User

	if err := repo.DB.QueryRowx(createUserQuery, surname, name, patronymic, address, passportNumber).Scan(
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.Address,
	); err != nil {
		err, ok := validators.IsConstrainError(err)
		if ok {
			return nil, fmt.Errorf("%v:%w", ErrUserExists, err)
		}

		return nil, err
	}

	return &user, nil
}

func (repo *UserTaskRepo) StartTask(userId int64, taskName string) error {
	startTime := time.Now()

	_, err := repo.DB.Exec(startTaskQuery, taskName, startTime, userId)
	if err != nil {
		err, ok := validators.IsConstrainError(err)
		if ok {
			return fmt.Errorf("%v:%w", ErrUserNotFound, err)
		}

		return err
	}

	return nil
}

func (repo *UserTaskRepo) StopTask(userId int64, taskName string) error {
	stopTime := time.Now()

	if err := repo.CheckUserIDTask(userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%v", ErrUserNotFound)
		}
	}

	_, err := repo.DB.Exec(stopTaskQuery, stopTime, userId, taskName)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserTaskRepo) GetLaborCosts(userId int64) ([]model.ResponseLobarCost, error) {

	if err := repo.CheckUserIDPerson(userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%v", ErrUserNotFound)
		}
	}

	var lobarCosts []model.ResponseLobarCost

	rows, err := repo.DB.Query(getLaborCosts, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var lobarCost model.ResponseLobarCost
		if err := rows.Scan(&lobarCost.TaskName, &lobarCost.DurationMinutes); err != nil {
			return nil, err
		}

		lobarCosts = append(lobarCosts, lobarCost)
	}

	return lobarCosts, nil
}

func (repo *UserTaskRepo) CheckUserIDPerson(userID int64) error {
	var userId int
	err := repo.DB.QueryRowx(checkUserIDPersonQuery, userID).Scan(&userId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserTaskRepo) CheckUserIDTask(id int64) error {
	var userId int
	err := repo.DB.QueryRowx(checkUserIDTaskQuery, id).Scan(&userId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserTaskRepo) GetUserByFilters(limit, offset int, id int64, surname, name, patronymic, address, passportNumber string) (*[]model.User, error) {
	query := "SELECT id, surname, name, patronymic, address, passport_number FROM person WHERE 1=1 "
	var args []interface{}
	paramIndex := 1

	if id != 0 {
		query += fmt.Sprintf(" and id = $%d", paramIndex)
		args = append(args, id)
		paramIndex++
	}

	if surname != "" {
		query += fmt.Sprintf(" and surname = $%d", paramIndex)
		args = append(args, surname)
		paramIndex++
	}

	if name != "" {
		query += fmt.Sprintf(" and name = $%d", paramIndex)
		args = append(args, name)
		paramIndex++
	}

	if patronymic != "" {
		query += fmt.Sprintf(" and patronymic = $%d", paramIndex)
		args = append(args, patronymic)
		paramIndex++
	}

	if address != "" {
		query += fmt.Sprintf(" and address = $%d", paramIndex)
		args = append(args, address)
		paramIndex++
	}

	if passportNumber != "" {
		query += fmt.Sprintf(" and passport_number = $%d", paramIndex)
		args = append(args, address)
		paramIndex++
	}

	if limit != 0 {
		query += fmt.Sprintf(" limit $%d", paramIndex)
		args = append(args, limit)
		paramIndex++
	}

	if offset != 0 {
		query += fmt.Sprintf(" offset $%d", paramIndex)
		args = append(args, offset-1)
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.ID, &user.Surname, &user.Name, &user.Patronymic, &user.Address, &user.PassportNumber); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (repo *UserTaskRepo) DeleteUser(id int64) error {
	if _, err := repo.DB.Exec(deleteFromTaskQuery, id); err != nil {
		return err
	}

	if _, err := repo.DB.Exec(deleteFromPersonQuery, id); err != nil {
		return err
	}

	return nil
}

func (repo *UserTaskRepo) UpdateUser(id int64, surname, name, patronymic, address, passportNumber string) (*model.User, error) {
	if err := repo.CheckUserIDPerson(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%v", ErrUserNotFound)
		}
	}

	tx := repo.DB.MustBegin()
	defer tx.Rollback()

	if _, err := tx.Exec(selectForUpdateQuery, id); err != nil {
		return nil, err
	}

	var user model.User

	if err := tx.QueryRowx(
		updatePersonQuery,
		surname,
		name,
		patronymic,
		address,
		passportNumber,
		id,
	).Scan(
		&user.ID,
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.Address,
		&user.PassportNumber,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}
