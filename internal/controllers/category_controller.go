package controllers

import (
	"strconv"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/pkg/validator"
	"finance-app-backend/internal/requests"
	"finance-app-backend/internal/resources"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (ctrl *CategoryController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	categories, err := ctrl.categoryService.List(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToCategoryCollection(categories)
	// collection responses wrap data directly
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (ctrl *CategoryController) Store(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	cat, err := ctrl.categoryService.Create(userID, req.Name, req.Icon, req.Color, req.Type)
	if err != nil {
		// Matching Laravel 403 response for Free users trying to add category
		return response.Error(c, fiber.StatusForbidden, err.Error(), fiber.Map{
			"upgrade_required": true,
		})
	}

	res := resources.ToCategoryResource(cat)
	return response.Success(c, "Kategori berhasil dibuat", res)
}

func (ctrl *CategoryController) Show(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	cat, err := ctrl.categoryService.GetByID(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, "Akses ditolak", nil)
	}

	res := resources.ToCategoryResource(cat)
	return response.Success(c, "Kategori ditemukan", res)
}

func (ctrl *CategoryController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	var req requests.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	cat, err := ctrl.categoryService.Update(userID, uint(id), req.Name, req.Icon, req.Color, req.Type)
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	res := resources.ToCategoryResource(cat)
	return response.Success(c, "Kategori berhasil diperbarui", res)
}

func (ctrl *CategoryController) Destroy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	err = ctrl.categoryService.Delete(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Kategori berhasil dihapus")
}
