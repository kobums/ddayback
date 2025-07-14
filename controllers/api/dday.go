package api

import (
	"dday-backend/controllers"
	"dday-backend/models"
	"dday-backend/models/dday"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DdayController struct {
	*controllers.Controller
	manager *models.DdayManager
}

func NewDdayController(c *fiber.Ctx) *DdayController {
	return &DdayController{
		Controller: controllers.NewController(c),
		manager:    models.NewDdayManager(),
	}
}

func (ctrl *DdayController) GetDdays(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	page, pageSize := ctrl.GetPagination()
	orderBy := ctrl.GetOrderBy()
	search := ctrl.GetSearch()
	category := ctrl.GetCategory()
	isImportant := ctrl.GetIsImportant()

	var args []interface{}

	if search != "" {
		args = append(args, models.NewCustom("(d_title LIKE ? OR d_memo LIKE ?)", "%"+search+"%", "%"+search+"%"))
	}

	if category != "" && dday.IsValidCategory(category) {
		args = append(args, models.NewWhere("d_category", category, "="))
	}

	if isImportant != nil {
		args = append(args, models.NewWhere("d_is_important", *isImportant, "="))
	}

	if orderBy != "" {
		args = append(args, models.NewOrdering(orderBy))
	}

	args = append(args, models.NewPaging(page, pageSize))

	ddays, err := ctrl.manager.GetAll(args...)
	if err != nil {
		return ctrl.InternalServerError("Failed to fetch D-Days")
	}

	totalCount, err := ctrl.manager.Count(args[:len(args)-1]...)
	if err != nil {
		return ctrl.InternalServerError("Failed to count D-Days")
	}

	response := fiber.Map{
		"data":       ddays,
		"pagination": fiber.Map{
			"page":       page,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPages": (totalCount + pageSize - 1) / pageSize,
		},
	}

	return ctrl.Success(response)
}

func (ctrl *DdayController) CreateDday(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	var req struct {
		Title       string `json:"title"`
		TargetDate  string `json:"target_date"`
		Category    string `json:"category"`
		Memo        string `json:"memo"`
		IsImportant bool   `json:"is_important"`
	}

	if err := ctrl.Body(&req); err != nil {
		return ctrl.BadRequest("Invalid request body")
	}

	if strings.TrimSpace(req.Title) == "" {
		return ctrl.BadRequest("Title is required")
	}

	if req.TargetDate == "" {
		return ctrl.BadRequest("Target date is required")
	}

	if _, err := time.Parse("2006-01-02", req.TargetDate); err != nil {
		return ctrl.BadRequest("Invalid target date format. Use YYYY-MM-DD")
	}

	if req.Category == "" {
		req.Category = dday.GetDefaultCategory()
	} else if !dday.IsValidCategory(req.Category) {
		return ctrl.BadRequest("Invalid category")
	}

	newDday := &models.DDay{
		ID:          uuid.New().String(),
		Title:       strings.TrimSpace(req.Title),
		TargetDate:  req.TargetDate,
		Category:    req.Category,
		Memo:        strings.TrimSpace(req.Memo),
		IsImportant: req.IsImportant,
		CreatedAt:   time.Now(),
	}

	if err := ctrl.manager.Create(newDday); err != nil {
		return ctrl.InternalServerError("Failed to create D-Day")
	}

	return ctrl.Created(newDday)
}

func (ctrl *DdayController) GetDday(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	id := ctrl.Params("id")
	if id == "" {
		return ctrl.BadRequest("ID is required")
	}

	dday, err := ctrl.manager.GetByID(id)
	if err != nil {
		return ctrl.NotFound("D-Day not found")
	}

	return ctrl.Success(dday)
}

func (ctrl *DdayController) UpdateDday(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	id := ctrl.Params("id")
	if id == "" {
		return ctrl.BadRequest("ID is required")
	}

	existingDday, err := ctrl.manager.GetByID(id)
	if err != nil {
		return ctrl.NotFound("D-Day not found")
	}

	var req struct {
		Title       string `json:"title"`
		TargetDate  string `json:"target_date"`
		Category    string `json:"category"`
		Memo        string `json:"memo"`
		IsImportant bool   `json:"is_important"`
	}

	if err := ctrl.Body(&req); err != nil {
		return ctrl.BadRequest("Invalid request body")
	}

	if strings.TrimSpace(req.Title) == "" {
		return ctrl.BadRequest("Title is required")
	}

	if req.TargetDate == "" {
		return ctrl.BadRequest("Target date is required")
	}

	if _, err := time.Parse("2006-01-02", req.TargetDate); err != nil {
		return ctrl.BadRequest("Invalid target date format. Use YYYY-MM-DD")
	}

	if req.Category == "" {
		req.Category = dday.GetDefaultCategory()
	} else if !dday.IsValidCategory(req.Category) {
		return ctrl.BadRequest("Invalid category")
	}

	updatedDday := &models.DDay{
		ID:          id,
		Title:       strings.TrimSpace(req.Title),
		TargetDate:  req.TargetDate,
		Category:    req.Category,
		Memo:        strings.TrimSpace(req.Memo),
		IsImportant: req.IsImportant,
		CreatedAt:   existingDday.CreatedAt,
	}

	if err := ctrl.manager.Update(id, updatedDday); err != nil {
		return ctrl.InternalServerError("Failed to update D-Day")
	}

	return ctrl.Success(updatedDday)
}

func (ctrl *DdayController) DeleteDday(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	id := ctrl.Params("id")
	if id == "" {
		return ctrl.BadRequest("ID is required")
	}

	_, err := ctrl.manager.GetByID(id)
	if err != nil {
		return ctrl.NotFound("D-Day not found")
	}

	if err := ctrl.manager.Delete(id); err != nil {
		return ctrl.InternalServerError("Failed to delete D-Day")
	}

	return ctrl.Success(fiber.Map{
		"message": "D-Day deleted successfully",
	})
}