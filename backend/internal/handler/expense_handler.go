package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
	"github.com/luisfucros/expense-tracker-app/internal/middleware"
	appvalidator "github.com/luisfucros/expense-tracker-app/pkg/validator"
)

// ExpenseHandler handles expense-related HTTP routes.
type ExpenseHandler struct {
	*Handler
}

// NewExpenseHandler creates an ExpenseHandler.
func NewExpenseHandler(h *Handler) *ExpenseHandler {
	return &ExpenseHandler{Handler: h}
}

// Create handles POST /api/v1/expenses
//
// @Summary     Create expense
// @Description Creates a new expense for the authenticated user
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body model.CreateExpenseInput true "Expense payload"
// @Success     201 {object} successResponse{data=model.Expense}
// @Failure     400 {object} errorResponse
// @Failure     401 {object} errorResponse
// @Router      /expenses [post]
func (eh *ExpenseHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		eh.Fail(c, apierror.Unauthorized("not authenticated"))
		return
	}

	var input model.CreateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "BAD_REQUEST",
			"message": "invalid request body",
		}})
		return
	}

	if err := appvalidator.Validate(input); err != nil {
		eh.Fail(c, err)
		return
	}

	expense, err := eh.ExpenseService.Create(c.Request.Context(), userID, input)
	if err != nil {
		eh.Fail(c, err)
		return
	}

	eh.Respond(c, http.StatusCreated, expense)
}

// GetByID handles GET /api/v1/expenses/:id
//
// @Summary     Get expense by ID
// @Description Returns a single expense belonging to the authenticated user
// @Tags        expenses
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Expense ID"
// @Success     200 {object} successResponse{data=model.Expense}
// @Failure     400 {object} errorResponse
// @Failure     401 {object} errorResponse
// @Failure     404 {object} errorResponse
// @Router      /expenses/{id} [get]
func (eh *ExpenseHandler) GetByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		eh.Fail(c, apierror.Unauthorized("not authenticated"))
		return
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		eh.Fail(c, apierror.BadRequest(apierror.CodeBadRequest, "invalid expense id"))
		return
	}

	expense, err := eh.ExpenseService.GetByID(c.Request.Context(), userID, id)
	if err != nil {
		eh.Fail(c, err)
		return
	}

	eh.Respond(c, http.StatusOK, expense)
}

// Update handles PUT /api/v1/expenses/:id
//
// @Summary     Update expense
// @Description Updates fields of an existing expense (all fields optional)
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path int                      true "Expense ID"
// @Param       body body model.UpdateExpenseInput true "Fields to update"
// @Success     200 {object} successResponse{data=model.Expense}
// @Failure     400 {object} errorResponse
// @Failure     401 {object} errorResponse
// @Failure     404 {object} errorResponse
// @Router      /expenses/{id} [put]
func (eh *ExpenseHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		eh.Fail(c, apierror.Unauthorized("not authenticated"))
		return
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		eh.Fail(c, apierror.BadRequest(apierror.CodeBadRequest, "invalid expense id"))
		return
	}

	var input model.UpdateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "BAD_REQUEST",
			"message": "invalid request body",
		}})
		return
	}

	if err := appvalidator.Validate(input); err != nil {
		eh.Fail(c, err)
		return
	}

	expense, err := eh.ExpenseService.Update(c.Request.Context(), userID, id, input)
	if err != nil {
		eh.Fail(c, err)
		return
	}

	eh.Respond(c, http.StatusOK, expense)
}

// Delete handles DELETE /api/v1/expenses/:id
//
// @Summary     Delete expense
// @Description Deletes an expense belonging to the authenticated user
// @Tags        expenses
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Expense ID"
// @Success     204
// @Failure     401 {object} errorResponse
// @Failure     404 {object} errorResponse
// @Router      /expenses/{id} [delete]
func (eh *ExpenseHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		eh.Fail(c, apierror.Unauthorized("not authenticated"))
		return
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		eh.Fail(c, apierror.BadRequest(apierror.CodeBadRequest, "invalid expense id"))
		return
	}

	if err := eh.ExpenseService.Delete(c.Request.Context(), userID, id); err != nil {
		eh.Fail(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// List handles GET /api/v1/expenses
//
// @Summary     List expenses
// @Description Returns a paginated, filterable list of the authenticated user's expenses
// @Tags        expenses
// @Produce     json
// @Security    BearerAuth
// @Param       category   query string false "Filter by category" Enums(Groceries,Leisure,Electronics,Utilities,Clothing,Health,Others)
// @Param       start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param       end_date   query string false "Filter to date (YYYY-MM-DD)"
// @Param       page       query int    false "Page number (default 1)"
// @Param       page_size  query int    false "Items per page, max 100 (default 20)"
// @Success     200 {object} successResponse{data=model.ExpenseListResponse}
// @Failure     401 {object} errorResponse
// @Router      /expenses [get]
func (eh *ExpenseHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		eh.Fail(c, apierror.Unauthorized("not authenticated"))
		return
	}

	filter := model.ExpenseFilter{
		Page:     1,
		PageSize: 20,
	}

	if cat := c.Query("category"); cat != "" {
		category := model.Category(cat)
		filter.Category = &category
	}

	if sd := c.Query("start_date"); sd != "" {
		t, err := time.Parse("2006-01-02", sd)
		if err == nil {
			filter.StartDate = &t
		}
	}

	if ed := c.Query("end_date"); ed != "" {
		t, err := time.Parse("2006-01-02", ed)
		if err == nil {
			// Set to end of day
			endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
			filter.EndDate = &endOfDay
		}
	}

	if p := c.Query("page"); p != "" {
		if page, err := strconv.Atoi(p); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if pageSize, err := strconv.Atoi(ps); err == nil && pageSize > 0 && pageSize <= 100 {
			filter.PageSize = pageSize
		}
	}

	result, err := eh.ExpenseService.List(c.Request.Context(), userID, filter)
	if err != nil {
		eh.Fail(c, err)
		return
	}

	eh.Respond(c, http.StatusOK, result)
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	raw := c.Param(param)
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
