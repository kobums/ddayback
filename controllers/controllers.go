package controllers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	c *fiber.Ctx
}

func NewController(c *fiber.Ctx) *Controller {
	return &Controller{c: c}
}

func (ctrl *Controller) Get(key string) string {
	return ctrl.c.Get(key)
}

func (ctrl *Controller) Geti(key string) int {
	value := ctrl.c.Get(key)
	if value == "" {
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func (ctrl *Controller) Query(key string) string {
	return ctrl.c.Query(key)
}

func (ctrl *Controller) Queryi(key string) int {
	value := ctrl.c.Query(key)
	if value == "" {
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func (ctrl *Controller) Params(key string) string {
	return ctrl.c.Params(key)
}

func (ctrl *Controller) ParamsInt(key string) int {
	value := ctrl.c.Params(key)
	if value == "" {
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func (ctrl *Controller) Body(out interface{}) error {
	return ctrl.c.BodyParser(out)
}

func (ctrl *Controller) JSON(data interface{}) error {
	return ctrl.c.JSON(data)
}

func (ctrl *Controller) Status(code int) *fiber.Ctx {
	return ctrl.c.Status(code)
}

func (ctrl *Controller) Error(code int, message string) error {
	return ctrl.c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}

func (ctrl *Controller) Success(data interface{}) error {
	return ctrl.c.JSON(data)
}

func (ctrl *Controller) Created(data interface{}) error {
	return ctrl.c.Status(201).JSON(data)
}

func (ctrl *Controller) NoContent() error {
	return ctrl.c.SendStatus(204)
}

func (ctrl *Controller) BadRequest(message string) error {
	return ctrl.Error(400, message)
}

func (ctrl *Controller) NotFound(message string) error {
	return ctrl.Error(404, message)
}

func (ctrl *Controller) InternalServerError(message string) error {
	return ctrl.Error(500, message)
}

func (ctrl *Controller) GetPagination() (page int, pageSize int) {
	page = ctrl.Queryi("page")
	if page <= 0 {
		page = 1
	}

	pageSize = ctrl.Queryi("pageSize")
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func (ctrl *Controller) GetOrderBy() string {
	orderBy := ctrl.Query("orderBy")
	if orderBy == "" {
		return ""
	}

	// Map API column names to database column names
	columnMapping := map[string]string{
		"id":           "d_id",
		"title":        "d_title",
		"target_date":  "d_target_date",
		"category":     "d_category",
		"is_important": "d_is_important",
		"created_at":   "d_created_at",
		"updated_at":   "d_updated_at",
	}

	direction := ctrl.Query("direction")
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}

	if dbColumn, exists := columnMapping[strings.ToLower(orderBy)]; exists {
		return dbColumn + " " + direction
	}

	return ""
}

func (ctrl *Controller) GetSearch() string {
	return strings.TrimSpace(ctrl.Query("search"))
}

func (ctrl *Controller) GetCategory() string {
	return strings.TrimSpace(ctrl.Query("category"))
}

func (ctrl *Controller) GetIsImportant() *bool {
	value := ctrl.Query("isImportant")
	if value == "" {
		return nil
	}

	if value == "true" || value == "1" {
		result := true
		return &result
	}
	if value == "false" || value == "0" {
		result := false
		return &result
	}

	return nil
}