package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strconv"
	"tenderservice/model"
	"tenderservice/myerrors"
	"tenderservice/service"
)

type bidHandler struct {
	bidService  *service.BidService
	authService *service.AuthService
	validator   *validator.Validate
}

func NewBidHandler(bidService *service.BidService, authService *service.AuthService) Handler {
	return &bidHandler{
		bidService:  bidService,
		authService: authService,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *bidHandler) Setup(router *gin.Engine) {
	router.POST("/api/bids/new", h.create)
	router.GET("/api/bids/my", h.getMy)
	router.GET("/api/bids/:tenderId/list", h.get)
	router.PUT("/api/bids/:bidId/feedback", h.createFeedback)
	router.GET("/api/bids/:tenderId/reviews", h.getReviews)
	router.PATCH("/api/bids/:bidId/edit", h.patchBid)
	router.PUT("/api/bids/:bidId/rollback/:version", h.rollbackBid)

}

func (h *bidHandler) create(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createBidDto model.CreateBidDto
	err = json.Unmarshal(data, &createBidDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.validator.Struct(createBidDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.authService.CheckIfAuthorExists(createBidDto.AuthorType, createBidDto.AuthorId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "this user does not exist"})
		return
	}

	err = h.bidService.CheckIfTenderExists(createBidDto.TenderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "this tender does not exist"})
		return
	}

	bid, err := h.bidService.Create(&createBidDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bid)
}

func (h *bidHandler) getMy(c *gin.Context) {
	username := c.Query("username")
	limit := -1
	offset := 0

	var err error
	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if limit < 0 || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if c.Query("offset") != "" {
		offset, err = strconv.Atoi(c.Query("offset"))
		if offset < 0 || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	userId, err := h.authService.GetEmployeeId(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	bids, err := h.bidService.GetMy(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if limit > -1 && limit < len(bids) {
		bids = bids[:limit]
	}
	if offset > 0 && offset < len(bids) {
		bids = bids[offset:]
	}

	c.JSON(http.StatusOK, bids)
}

func (h *bidHandler) get(c *gin.Context) {
	username := c.Query("username")
	tenderId := c.Param("tenderId")
	limit := -1
	offset := 0

	var err error
	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if limit < 0 || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if c.Query("offset") != "" {
		offset, err = strconv.Atoi(c.Query("offset"))
		if offset < 0 || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	bids, err := h.bidService.GetAll(tenderId, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if limit > -1 && limit < len(bids) {
		bids = bids[:limit]
	}
	if offset > 0 && offset < len(bids) {
		bids = bids[offset:]
	}

	c.JSON(http.StatusOK, bids)
}

func (h *bidHandler) createFeedback(c *gin.Context) {
	description := c.Query("bidFeedback")
	username := c.Query("username")
	bidId := c.Param("bidId")

	_, err := h.authService.GetEmployeeId(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	feedback, err := h.bidService.CreateBidFeedback(bidId, username, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feedback)
}

func (h *bidHandler) getReviews(c *gin.Context) {
	limit := -1
	offset := 0

	var err error
	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if limit < 0 || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if c.Query("offset") != "" {
		offset, err = strconv.Atoi(c.Query("offset"))
		if offset < 0 || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	tenderId := c.Param("tenderId")
	authorUsername := c.Query("authorUsername")
	reqUsername := c.Query("requesterUsername")

	if err := h.bidService.CheckIfUserBidExists(tenderId, authorUsername); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tender does not have a user's bid"})
		return
	}

	if err := h.authService.CheckRightsForTender(reqUsername, tenderId); err != nil {
		switch {
		case errors.Is(err, &myerrors.NotFoundError{}):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, &myerrors.NoRightsError{}):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	reviews, err := h.bidService.GetAllUserReviews(authorUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if limit > -1 && limit < len(reviews) {
		reviews = reviews[:limit]
	}
	if offset > 0 && offset < len(reviews) {
		reviews = reviews[offset:]
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *bidHandler) patchBid(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var patchBidDto model.PatchBidDto
	err = json.Unmarshal(data, &patchBidDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.validator.Struct(patchBidDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.Query("username")
	bidId := c.Param("bidId")
	err = h.authService.CheckRightsForBid(username, bidId)
	if err != nil {
		switch {
		case errors.Is(err, &myerrors.NotFoundError{}):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, &myerrors.NoRightsError{}):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	bid, err := h.bidService.PatchBid(bidId, &patchBidDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bid)

}

func (h *bidHandler) rollbackBid(c *gin.Context) {
	username := c.Query("username")
	bidId := c.Param("bidId")
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.authService.CheckRightsForBid(username, bidId)
	if err != nil {
		switch {
		case errors.Is(err, &myerrors.NotFoundError{}):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, &myerrors.NoRightsError{}):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	bid, err := h.bidService.RollbackVersion(bidId, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bid)
}
