package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"gorm.io/gorm"

	"github.com/wingc34/mini-commerce-api/pkg/response"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetRecommended returns 4 random products for the homepage.
func (h *ProductHandler) GetRecommended(c *gin.Context) {
	// 第一步：呼叫 service
	products, err := h.service.GetRecommended()

	// 第二步：處理錯誤
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get recommended products")
		return
	}

	// 第三步：回傳結果
	response.Success(c, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.service.GetProductByID(id)

	if err != nil {
		// 找不到是 404，不是 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get product")
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	// 拿 ?page=1，預設是 1
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	products, total, err := h.service.GetProducts(page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  page,
	})

}

type CheckStockRequest struct {
	Attributes map[string]string `json:"attributes" binding:"required,min=1"`
}

func (h *ProductHandler) CheckStock(c *gin.Context) {
	var req CheckStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	stock, err := h.service.CheckStock(id, req.Attributes)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to check stock")
		return
	}

	response.Success(c, gin.H{
		"inStock": stock,
	})
}
