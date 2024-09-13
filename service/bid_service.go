package service

import (
	"database/sql"
	"fmt"
	"tenderservice/logger"
	"tenderservice/model"
	"tenderservice/repository"
)

const bidFields = "id, name, description, status, tender_id, author_type, author_id, version, created_at"

const feedbackFields = "id, bid_id, description, username, created_at"

type BidService struct {
	repository repository.Repository
}

func NewBidService(repository repository.Repository) *BidService {
	return &BidService{
		repository: repository,
	}
}

func (s *BidService) Create(createBidDto *model.CreateBidDto) (*model.Bid, error) {
	query := addReturnFields(`INSERT INTO bid(name, description, tender_id, author_type, author_id) VALUES($1, $2, $3, $4, $5)`, bidFields)
	row := s.repository.GetConnection().QueryRow(query,
		createBidDto.Name,
		createBidDto.Description,
		createBidDto.TenderId,
		createBidDto.AuthorType.String(),
		createBidDto.AuthorId)
	var bid model.Bid
	err := s.scanRow(row, &bid)
	return &bid, err
}

func (s *BidService) GetMy(userId string) ([]*model.Bid, error) {
	query := fmt.Sprintf("SELECT %s FROM bid WHERE author_id = $1", bidFields)
	rows, err := s.repository.GetConnection().Query(query, userId)
	if err != nil {
		return nil, err
	}
	return s.scanRows(rows)
}

func (s *BidService) GetAll(tenderId, username string) ([]*model.Bid, error) {
	query := fmt.Sprintf("SELECT %s FROM bid WHERE tender_id = $1", bidFields)
	rows, err := s.repository.GetConnection().Query(query, tenderId)
	if err != nil {
		return nil, err
	}
	return s.scanRows(rows)
}

func (s *BidService) GetStatus(bidId string) (*model.Status, error) {
	query := `SELECT status FROM bid WHERE tender_id = $1`
	var status model.Status
	err := s.repository.GetConnection().QueryRow(query, bidId).Scan(&status)
	return &status, err
}

func (s *BidService) CreateBidFeedback(bidId, username, description string) (*model.Feedback, error) {
	query := addReturnFields(`INSERT INTO feedback(bid_id, username, description) 
									VALUES($1, $2, $3)`, feedbackFields)
	var feedback model.Feedback
	err := s.repository.GetConnection().QueryRow(query, bidId, username, description).Scan(
		&feedback.Id,
		&feedback.BidId,
		&feedback.Description,
		&feedback.Username,
		&feedback.CreatedAt)
	return &feedback, err
}

func (s *BidService) CheckIfUserBidExists(tenderId, username string) error {
	query := `SELECT id FROM bid WHERE tender_id = $1 and author_id=(SELECT id FROM employee WHERE username = $2)`
	var id string
	return s.repository.GetConnection().QueryRow(query, tenderId, username).Scan(&id)

}

func (s *BidService) GetAllUserReviews(username string) ([]*model.Feedback, error) {
	query := `SELECT id, bid_id, description, username, created_at from feedback WHERE bid_id in (SELECT id FROM bid WHERE author_id=(SELECT id 
              FROM employee WHERE username=$1))`

	rows, err := s.repository.GetConnection().Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := make([]*model.Feedback, 0)
	for rows.Next() {
		var review model.Feedback
		err := rows.Scan(
			&review.Id,
			&review.BidId,
			&review.Description,
			&review.Username,
			&review.CreatedAt)
		if err != nil {
			logger.WarnLogger.Println("error occurred while scanning a query result: ", err.Error())
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (s *BidService) PatchBid(bidId string, patchBidDto *model.PatchBidDto) (*model.Bid, error) {
	versionSaved, err := s.saveVersion(bidId)
	if err != nil {
		return nil, err
	}

	query := addReturnFields(`UPDATE bid SET name=$1, description=$2, version=$3  WHERE id = $4`, bidFields)
	row := s.repository.GetConnection().QueryRow(
		query,
		patchBidDto.Name,
		patchBidDto.Description,
		versionSaved+1,
		bidId)
	var bid model.Bid
	err = s.scanRow(row, &bid)
	return &bid, err
}

func (s *BidService) saveVersion(bidId string) (int, error) {
	query := `SELECT version, name, description FROM bid WHERE id=$1`
	var version int
	var name, description string
	err := s.repository.GetConnection().QueryRow(query, bidId).Scan(&version, &name, &description)
	if err != nil {
		return 0, err
	}
	query = `INSERT INTO bid_version(bid_id, version, name, description)
             VALUES ($1, $2, $3, $4)`
	row := s.repository.GetConnection().QueryRow(query, bidId, version, name, description)
	return version, row.Err()
}

func (s *BidService) RollbackVersion(bidId string, version int) (*model.Bid, error) {
	query := `SELECT name, description FROM bid_version WHERE bid_id=$1 AND version=$2`
	var name, description string
	err := s.repository.GetConnection().QueryRow(query, bidId, version).Scan(&name, &description)
	if err != nil {
		return nil, err
	}

	oldVersion, err := s.saveVersion(bidId)
	if err != nil {
		return nil, err
	}

	query = addReturnFields(`UPDATE bid SET name=$1, description=$2, version=$3 WHERE id=$4`, bidFields)
	row := s.repository.GetConnection().QueryRow(query, name, description, oldVersion+1, bidId)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var bid model.Bid
	err = s.scanRow(row, &bid)
	return &bid, err
}

func (s *BidService) scanRows(rows *sql.Rows) ([]*model.Bid, error) {
	bids := make([]*model.Bid, 0)
	for rows.Next() {
		var bid model.Bid
		err := rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderId,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
			&bid.CreatedAt)
		if err != nil {
			logger.WarnLogger.Println("error occurred while scanning from a query row: ", err)
		}
		bids = append(bids, &bid)
	}
	return bids, nil
}

func (s *BidService) CheckIfTenderExists(tenderId string) error {
	query := `SELECT name FROM tender WHERE id=$1`
	var name string
	row := s.repository.GetConnection().QueryRow(query, tenderId)
	err := row.Scan(&name)
	return err
}

func (s *BidService) scanRow(row *sql.Row, bid *model.Bid) error {
	return row.Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt)
}
