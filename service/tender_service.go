package service

import (
	"database/sql"
	"fmt"
	"strings"
	"tenderservice/logger"
	"tenderservice/model"
	"tenderservice/repository"
)

const tenderFields = "id, name, description, service_type, status, version, organization_id, creator_username, created_at"

type TenderService struct {
	repository repository.Repository
}

func NewTenderService(repository repository.Repository) *TenderService {
	return &TenderService{
		repository: repository,
	}
}

func (s *TenderService) Get() ([]*model.Tender, error) {
	query := fmt.Sprintf("SELECT %s FROM tender", tenderFields)
	rows, err := s.repository.GetConnection().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanRows(rows)
}

func (s *TenderService) GetMy(username string) ([]*model.Tender, error) {
	query := fmt.Sprintf("SELECT %s FROM tender WHERE creator_username=$1", tenderFields)
	rows, err := s.repository.GetConnection().Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanRows(rows)
}

func (s *TenderService) scanRows(rows *sql.Rows) ([]*model.Tender, error) {
	tenders := make([]*model.Tender, 0)
	for rows.Next() {
		var tender model.Tender
		err := rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.Version,
			&tender.OrganizationId,
			&tender.CreatorUsername,
			&tender.CreatedAt)
		if err != nil {
			logger.WarnLogger.Println("error occurred while scanning from a query row: ", err)
		}
		tenders = append(tenders, &tender)
	}
	return tenders, nil
}

func (s *TenderService) Create(tenderDto model.CreateTenderDto) (*model.Tender, error) {
	query := addReturnFields(`INSERT INTO tender(name, description, service_type, organization_id, creator_username) 
			VALUES ($1, $2, $3, $4, $5)`, tenderFields)
	fmt.Println(query)
	var tender model.Tender
	row := s.repository.GetConnection().QueryRow(
		query,
		tenderDto.Name,
		tenderDto.Description,
		tenderDto.ServiceType.String(),
		tenderDto.OrganizationId,
		tenderDto.CreatorUsername)

	err := s.scanRow(row, &tender)
	if err != nil {
		return nil, err
	}

	return &tender, err
}

func (s *TenderService) GetStatus(tenderId string) (*model.Status, error) {
	query := `SELECT status FROM tender WHERE id=$1`
	var status model.Status
	err := s.repository.GetConnection().QueryRow(query, tenderId).Scan(&status)
	return &status, err
}

func (s *TenderService) ChangeStatus(tenderId string, status model.Status) (*model.Tender, error) {
	query := addReturnFields(`UPDATE tender SET status = $1 WHERE id = $2`, tenderFields)
	fmt.Println(status.String(), tenderId)
	var tender model.Tender
	row := s.repository.GetConnection().QueryRow(query, status.String(), tenderId)
	err := s.scanRow(row, &tender)
	return &tender, err
}

func (s *TenderService) Patch(tenderId string, patchTenderDto model.PatchTenderDto) (*model.Tender, error) {
	versionSaved, err := s.saveVersion(tenderId)
	if err != nil {
		return nil, err
	}

	query := addReturnFields(`UPDATE tender SET name=$1, description=$2,service_type=$3, version=$4  WHERE id = $5`, tenderFields)
	var tender model.Tender
	row := s.repository.GetConnection().QueryRow(
		query,
		patchTenderDto.Name,
		patchTenderDto.Description,
		patchTenderDto.ServiceType.String(),
		versionSaved+1,
		tenderId)
	err = s.scanRow(row, &tender)
	return &tender, err
}

func (s *TenderService) saveVersion(tenderId string) (int, error) {
	query := `SELECT version, name, description, service_type FROM tender WHERE id=$1`
	var version int
	var name, description string
	var serviceType model.ServiceType
	err := s.repository.GetConnection().QueryRow(query, tenderId).Scan(&version, &name, &description, &serviceType)
	if err != nil {
		return 0, err
	}
	query = `INSERT INTO tender_version(tender_id, version, name, description, service_type)
             VALUES ($1, $2, $3, $4, $5)`
	row := s.repository.GetConnection().QueryRow(query, tenderId, version, name, description, serviceType.String())
	return version, row.Err()
}

func (s *TenderService) RollbackVersion(tenderId string, version int) (*model.Tender, error) {
	query := `SELECT name, description, service_type FROM tender_version WHERE tender_id=$1 AND version=$2`
	var name, description string
	var serviceType model.ServiceType
	err := s.repository.GetConnection().QueryRow(query, tenderId, version).Scan(&name, &description, &serviceType)
	if err != nil {
		return nil, err
	}

	oldVersion, err := s.saveVersion(tenderId)
	if err != nil {
		return nil, err
	}

	query = addReturnFields(`UPDATE tender SET name=$1, description=$2, service_type=$3, version=$4 WHERE id=$5`, tenderFields)
	row := s.repository.GetConnection().QueryRow(query, name, description, serviceType.String(), oldVersion+1, tenderId)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var tender model.Tender
	err = s.scanRow(row, &tender)
	return &tender, err
}

func (s *TenderService) scanRow(row *sql.Row, tender *model.Tender) error {
	return row.Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&tender.OrganizationId,
		&tender.CreatorUsername,
		&tender.CreatedAt)
}

func addReturnFields(query string, fields string) string {
	var builder strings.Builder
	builder.WriteString(query)
	builder.WriteString(" RETURNING ")
	builder.WriteString(fields)
	return builder.String()
}
