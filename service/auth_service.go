package service

import (
	"github.com/google/uuid"
	"tenderservice/model"
	"tenderservice/myerrors"
	"tenderservice/repository"
)

type AuthService struct {
	repository repository.Repository
}

func NewAuthService(repository repository.Repository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) CheckRightForOrg(username, organizationId string) (err error) {
	query := `SELECT id from employee WHERE username=$1`
	var userId uuid.UUID
	s.repository.GetConnection().QueryRow(query, username).Scan(&userId)
	if userId == uuid.Nil {
		return &myerrors.NotFoundError{Text: "employee not found"}
	}

	query = `SELECT id FROM organization_responsible WHERE user_id = $1 and organization_id=$2`
	var relId uuid.UUID
	s.repository.GetConnection().QueryRow(query, userId, organizationId).Scan(&relId)
	if relId == uuid.Nil {
		return &myerrors.NoRightsError{Text: "employee does not have enough rights"}
	}

	return err
}

func (s *AuthService) CheckRightsForTender(username, tenderId string) (err error) {
	query := `SELECT organization_id FROM tender WHERE id = $1`
	var orgId string
	err = s.repository.GetConnection().QueryRow(query, tenderId).Scan(&orgId)
	if err != nil {
		return &myerrors.NotFoundError{Text: "tender not found"}
	}

	return s.CheckRightForOrg(username, orgId)
}

func (s *AuthService) CheckRightsForBid(username, bidId string) (err error) {
	query := `SELECT author_id FROM bid WHERE id = $1`
	var authorId string
	err = s.repository.GetConnection().QueryRow(query, bidId).Scan(&authorId)
	if err != nil {
		return err
	}

	var userId string
	query = `SELECT id FROM employee WHERE username=$1`
	err = s.repository.GetConnection().QueryRow(query, username).Scan(&userId)
	if err != nil {
		return err
	}

	if userId != authorId {
		return &myerrors.NoRightsError{Text: "employee does not have enough rights"}
	}
	return nil
}

func (s *AuthService) CheckIfAuthorExists(authorType model.AuthorType, authorId string) error {
	var query string
	if authorType == model.User {
		query = `SELECT username FROM employee WHERE id=$1`
	} else {
		query = `SELECT name FROM organization WHERE id=$1`
	}
	var name string
	row := s.repository.GetConnection().QueryRow(query, authorId)
	err := row.Scan(&name)
	return err
}

func (s *AuthService) GetEmployeeId(username string) (string, error) {
	query := `SELECT id FROM employee WHERE username=$1`
	var id string
	err := s.repository.GetConnection().QueryRow(query, username).Scan(&id)
	return id, err
}
