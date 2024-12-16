package inbound

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	mwr "github.com/devpablocristo/golang-monorepo/pkg/rest/middlewares/gin"
	ginserver "github.com/devpablocristo/golang-monorepo/pkg/rest/servers/gin"
	gindefs "github.com/devpablocristo/golang-monorepo/pkg/rest/servers/gin/defs"
	types "github.com/devpablocristo/golang-monorepo/pkg/types"
	utils "github.com/devpablocristo/golang-monorepo/pkg/utils"

	config "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/config"
	transport "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/adapters/inbound/transport"
	ports "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/core/ports"
)

type Handler struct {
	Ucs ports.UseCases
	Svr gindefs.Server
}

func NewHandler(u ports.UseCases) (*Handler, error) {
	s, err := ginserver.Bootstrap(false)
	if err != nil {
		return nil, err
	}
	return &Handler{
		Ucs: u,
		Svr: s,
	}, nil
}

func (h *Handler) GetRouter() *gin.Engine {
	return h.Svr.GetRouter()
}

func (h *Handler) Start(ctx context.Context) error {
	h.Routes()
	return h.Svr.RunServer(ctx)
}

func (h *Handler) Routes() {
	router := h.Svr.GetRouter()

	router.GET("/health", h.Health)

	apiVersion := h.Svr.GetApiVersion()
	apiBase := "/api/" + apiVersion

	customers := router.Group(apiBase + "/customers")
	{
		customers.GET("", h.GetCustomers)
		customers.GET("/:id", h.GetCustomer)
		customers.POST("", h.CreateCustomer)
		customers.PUT("/:id", h.UpdateCustomer)
		customers.DELETE("/:id", h.DeleteCustomer)
		customers.GET("/kpi", h.GetKPI)
	}

	router.GET(apiBase+"/ping", h.Ping)

	protected := router.Group(apiBase + "/protected")
	protected.Use(mwr.Validate(config.Auth()))
	{
		protected.GET("/ping", h.ProtectedPing)
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
	})
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) ProtectedPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "protected pong"})
}

func (h *Handler) GetCustomers(c *gin.Context) {
	customers, err := h.Ucs.GetCustomers(c.Request.Context())
	if err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}
	c.JSON(http.StatusOK, transport.GetCustomersResponse{
		Customers: transport.DomainListToCustomerJsonList(customers),
	})
}

func (h *Handler) GetCustomer(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr, status := types.NewAPIError(
			types.NewError(
				types.ErrInvalidInput,
				"invalid customer ID format",
				err,
			),
		)
		c.JSON(status, apiErr)
		return
	}

	if err := utils.ValidateID(ID); err != nil {
		apiErr, status := types.NewAPIError(
			types.NewError(
				types.ErrInvalidInput,
				"invalid customer ID format",
				err,
			),
		)
		c.JSON(status, apiErr)
		return
	}

	customer, err := h.Ucs.GetCustomerByID(c.Request.Context(), ID)
	if err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}

	c.JSON(http.StatusOK, transport.GetCustomerResponse{
		Customers: *transport.DomainToCustomerJson(customer),
	})
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var req transport.CustomerJson
	if err := c.ShouldBindJSON(&req); err != nil {
		errStr := err.Error()
		var message string
		switch {
		case strings.Contains(errStr, "Email' failed on the 'required' tag"):
			message = "invalid email format"
		case strings.Contains(errStr, "Age' failed on the 'required' tag"):
			message = "invalid age"
		case strings.Contains(errStr, "failed on the 'required' tag"):
			message = "missing required field"
		case strings.Contains(errStr, "cannot unmarshal"):
			message = "invalid data type"
		default:
			message = "request cannot be nil"
		}

		apiErr, status := types.NewAPIError(
			types.NewError(
				types.ErrValidation,
				message,
				err,
			),
		)
		c.JSON(status, apiErr)
		return
	}

	if err := validateRequest(&req); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}

	if err := h.Ucs.CreateCustomer(c.Request.Context(), transport.CustomerJsonToDomain(&req)); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr, status := types.NewAPIError(
			types.NewError(
				types.ErrInvalidInput,
				"invalid customer ID format",
				err,
			),
		)
		c.JSON(status, apiErr)
		return
	}

	if err := utils.ValidateID(ID); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}

	var req transport.CustomerJson
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr, status := types.NewAPIError(
			types.NewError(
				types.ErrValidation,
				"invalid request body",
				err,
			),
		)
		c.JSON(status, apiErr)
		return
	}

	if err := validateRequest(&req); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}

	customer := transport.CustomerJsonToDomain(&req)
	customer.ID = ID

	if err := h.Ucs.UpdateCustomer(c.Request.Context(), customer); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr, status := types.NewAPIError(
			types.NewError(
				types.ErrInvalidInput,
				"invalid customer ID format",
				err,
			),
		)
		c.JSON(status, apiErr)
		return
	}

	if err := utils.ValidateID(ID); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}

	if err := h.Ucs.DeleteCustomer(c.Request.Context(), ID); err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetKPI(c *gin.Context) {
	kpi, err := h.Ucs.GetKPI(c.Request.Context())
	if err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}
	c.JSON(http.StatusOK, transport.ToGetKPIJson(kpi))
}
