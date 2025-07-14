package rest

import (
	"dday-backend/controllers"
	"dday-backend/models"
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

func (ctrl *DdayController) List(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	page, pageSize := ctrl.GetPagination()
	orderBy := ctrl.GetOrderBy()

	var args []interface{}

	if orderBy != "" {
		args = append(args, models.NewOrdering(orderBy))
	}

	args = append(args, models.NewPaging(page, pageSize))

	ddays, err := ctrl.manager.GetAll(args...)
	if err != nil {
		return ctrl.InternalServerError("Failed to fetch D-Days")
	}

	return ctrl.Success(ddays)
}

func (ctrl *DdayController) Get(c *fiber.Ctx) error {
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

func (ctrl *DdayController) Create(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	var dday models.DDay
	if err := ctrl.Body(&dday); err != nil {
		return ctrl.BadRequest("Invalid request body")
	}

	dday.ID = uuid.New().String()
	dday.CreatedAt = time.Now()

	if err := ctrl.manager.Create(&dday); err != nil {
		return ctrl.InternalServerError("Failed to create D-Day")
	}

	return ctrl.Created(dday)
}

func (ctrl *DdayController) Update(c *fiber.Ctx) error {
	ctrl.Controller = controllers.NewController(c)

	id := ctrl.Params("id")
	if id == "" {
		return ctrl.BadRequest("ID is required")
	}

	existingDday, err := ctrl.manager.GetByID(id)
	if err != nil {
		return ctrl.NotFound("D-Day not found")
	}

	var updatedDday models.DDay
	if err := ctrl.Body(&updatedDday); err != nil {
		return ctrl.BadRequest("Invalid request body")
	}

	updatedDday.ID = id
	updatedDday.CreatedAt = existingDday.CreatedAt

	if err := ctrl.manager.Update(id, &updatedDday); err != nil {
		return ctrl.InternalServerError("Failed to update D-Day")
	}

	return ctrl.Success(updatedDday)
}

func (ctrl *DdayController) Delete(c *fiber.Ctx) error {
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

	return ctrl.NoContent()
}