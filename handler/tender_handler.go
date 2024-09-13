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

type tenderHandler struct {
	tenderService *service.TenderService
	authService   *service.AuthService
	validator     *validator.Validate
}

func NewTenderHandler(tenderService *service.TenderService, authService *service.AuthService) Handler {
	return &tenderHandler{
		tenderService: tenderService,
		authService:   authService,
		validator:     validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *tenderHandler) Setup(router *gin.Engine) {
	router.GET("/api/tenders", h.getAll)
	router.POST("/api/tenders/new", h.create)
	router.GET("/api/tenders/my", h.getMy)
	router.GET("/api/tenders/:tenderId/status", h.getStatus)
	router.PUT("/api/tenders/:tenderId/status", h.changeStatus)
	router.PATCH("/api/tenders/:tenderId/edit", h.patchById)
	router.PUT("/api/tenders/:tenderId/rollback/:version", h.putToVersion)
}

func (h *tenderHandler) getAll(c *gin.Context) {
	limit := -1
	offset := 0
	service_type := c.QueryArray("service_type")

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

	tenders, err := h.tenderService.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if service_type != nil {
		types := make(map[string]bool, 3)
		for _, t := range service_type {
			types[t] = true
		}

		filteredTenders := make([]*model.Tender, 0, len(tenders))
		for _, t := range tenders {
			if types[t.ServiceType.String()] {
				filteredTenders = append(filteredTenders, t)
			}
		}
		tenders = filteredTenders
	}

	if limit > -1 && limit < len(tenders) {
		tenders = tenders[:limit]
	}

	if offset > 0 && offset < len(tenders) {
		tenders = tenders[offset:]
	}

	c.JSON(http.StatusOK, tenders)
}

func (h *tenderHandler) getMy(c *gin.Context) {
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

	tenders, err := h.tenderService.GetMy(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if limit > -1 && limit < len(tenders) {
		tenders = tenders[:limit]
	}

	if offset > 0 && offset < len(tenders) {
		tenders = tenders[offset:]
	}
	c.JSON(http.StatusOK, tenders)
}

func (h *tenderHandler) create(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createTenderDto model.CreateTenderDto
	err = json.Unmarshal(data, &createTenderDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.validator.Struct(createTenderDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.authService.CheckRightForOrg(createTenderDto.CreatorUsername, createTenderDto.OrganizationId)
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

	tender, err := h.tenderService.Create(createTenderDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tender)
}

func (h *tenderHandler) getStatus(c *gin.Context) {
	tenderId := c.Param("tenderId")
	username := c.Query("username")

	if err := h.authService.CheckRightsForTender(username, tenderId); err != nil {
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

	status, err := h.tenderService.GetStatus(tenderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *tenderHandler) changeStatus(c *gin.Context) {
	tenderId := c.Param("tenderId")
	username := c.Query("username")
	status, err := model.ParseStatus(c.Query("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.CheckRightsForTender(username, tenderId); err != nil {
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

	tender, err := h.tenderService.ChangeStatus(tenderId, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tender)
}

func (h *tenderHandler) patchById(c *gin.Context) {
	tenderId := c.Param("tenderId")
	username := c.Query("username")

	if err := h.authService.CheckRightsForTender(username, tenderId); err != nil {
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

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var patchTenderDto model.PatchTenderDto
	err = json.Unmarshal(data, &patchTenderDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.validator.Struct(patchTenderDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tender, err := h.tenderService.Patch(tenderId, patchTenderDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tender)

}

func (h *tenderHandler) putToVersion(c *gin.Context) {
	tenderId := c.Param("tenderId")
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.Query("username")

	if err := h.authService.CheckRightsForTender(username, tenderId); err != nil {
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

	tender, err := h.tenderService.RollbackVersion(tenderId, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tender)
}
