// @title        Customer Manager API
// @version      1.0
// @description  API para gestión de clientes
// @host         localhost:8080
// @BasePath     /api/v1
package inbound

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	mwr "github.com/devpablocristo/tech-house/pkg/rest/middlewares/gin"
	ginserver "github.com/devpablocristo/tech-house/pkg/rest/servers/gin"
	gindefs "github.com/devpablocristo/tech-house/pkg/rest/servers/gin/defs"
	swagdoc "github.com/devpablocristo/tech-house/pkg/swagger"
	swagin "github.com/devpablocristo/tech-house/pkg/swagger/adapters"
	swagdefs "github.com/devpablocristo/tech-house/pkg/swagger/defs"
	types "github.com/devpablocristo/tech-house/pkg/types"
	utils "github.com/devpablocristo/tech-house/pkg/utils"

	config "github.com/devpablocristo/tech-house/projects/customers-manager/internal/config"
	transport "github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/adapters/inbound/transport"
	ports "github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/core/ports"
)

type Handler struct {
	Ucs ports.UseCases
	Svr gindefs.Server
	Swg swagdefs.Service
}

func NewHandler(u ports.UseCases) (*Handler, error) {
	s, err := ginserver.Bootstrap(false)
	if err != nil {
		return nil, err
	}

	g, err := swagdoc.Bootstrap()
	if err != nil {
		return nil, err
	}

	return &Handler{
		Ucs: u,
		Svr: s,
		Swg: g,
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

	// Configurar Swagger
	if err := swagin.SetupSwagger(router, h.Swg); err != nil {
		panic(err)
	}
}

// @Summary     Health check
// @Description Verifica el estado del servicio
// @Tags        system
// @Produce     json
// @Success     200 {object} gin.H
// @Router      /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
	})
}

// @Summary     Ping
// @Description Simple ping para verificar conectividad
// @Tags        system
// @Produce     json
// @Success     200 {object} gin.H
// @Router      /ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// @Summary     Protected ping
// @Description Ping protegido que requiere autenticación
// @Tags        system
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} gin.H
// @Failure     401 {object} types.APIError
// @Router      /protected/ping [get]
func (h *Handler) ProtectedPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "protected pong"})
}

// @Summary     Get list of customers
// @Description Obtiene la lista de todos los clientes
// @Tags        customers
// @Produce     json
// @Success     200 {object} transport.GetCustomersResponse
// @Failure     500 {object} types.APIError
// @Router      /customers [get]
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

// @Summary     Get customer by ID
// @Description Obtiene un cliente por su ID
// @Tags        customers
// @Produce     json
// @Param       id path int true "Customer ID"
// @Success     200 {object} transport.GetCustomerResponse
// @Failure     400 {object} types.APIError
// @Failure     404 {object} types.APIError
// @Failure     500 {object} types.APIError
// @Router      /customers/{id} [get]
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

// @Summary     Create customer
// @Description Crea un nuevo cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Param       customer body transport.CustomerJson true "Customer Data"
// @Success     201
// @Failure     400 {object} types.APIError
// @Failure     500 {object} types.APIError
// @Router      /customers [post]
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

// @Summary     Update customer
// @Description Actualiza un cliente existente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Param       id path int true "Customer ID"
// @Param       customer body transport.CustomerJson true "Customer Data"
// @Success     200
// @Failure     400 {object} types.APIError
// @Failure     404 {object} types.APIError
// @Failure     500 {object} types.APIError
// @Router      /customers/{id} [put]
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

// @Summary     Delete customer
// @Description Elimina un cliente
// @Tags        customers
// @Param       id path int true "Customer ID"
// @Success     204
// @Failure     400 {object} types.APIError
// @Failure     404 {object} types.APIError
// @Failure     500 {object} types.APIError
// @Router      /customers/{id} [delete]
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

// @Summary     Get KPIs
// @Description Obtiene los KPIs de clientes
// @Tags        customers
// @Produce     json
// @Success     200 {object} transport.GetKPIJson
// @Failure     500 {object} types.APIError
// @Router      /customers/kpi [get]
func (h *Handler) GetKPI(c *gin.Context) {
	kpi, err := h.Ucs.GetKPI(c.Request.Context())
	if err != nil {
		apiErr, status := types.NewAPIError(err)
		c.JSON(status, apiErr)
		return
	}
	c.JSON(http.StatusOK, transport.ToGetKPIJson(kpi))
}
