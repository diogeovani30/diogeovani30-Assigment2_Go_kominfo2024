package controllers

import (
	"assignment2/database"
	"assignment2/dto"
	"assignment2/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateItem(ctx *gin.Context) {
	var req dto.CreateItem

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	item := models.Items{
		ItemCode:    req.ItemCode,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create item"})
		return
	}

	response := dto.CreateItemResponse{
		ID:          item.ID,
		ItemCode:    item.ItemCode,
		Description: item.Description,
		Quantity:    item.Quantity,
	}

	ctx.JSON(201, response)
}

func CreateOrder(ctx *gin.Context) {
	req := dto.CreateOrderRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		CustomerName: req.CustomerName,
		OrderedAt:    time.Now(),
	}

	if err := database.DB.Create(&order).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create order"})
		return
	}

	for _, item := range req.Items {
		orderItem := models.Items{
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
		}

		// Create the item
		if err := database.DB.Create(&orderItem).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create order item"})
			return
		}

		// Associate the item with the order in the order_items table
		if err := database.DB.Model(&order).Association("Items").Append(&orderItem); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to associate order item with order"})
			return
		}
	}

	// Construct the response
	response := dto.CreateOrderResponse{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		OrderAt:      order.OrderedAt,
		Items:        make([]dto.CreateItem, len(req.Items)),
	}

	for i, item := range req.Items {
		response.Items[i] = dto.CreateItem{
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
		}
	}

	ctx.JSON(201, response)
}

func GetOrders(ctx *gin.Context) {
	var orders []models.Order
	database.DB.Preload("Items").Find(&orders)

	ctx.JSON(200, orders)
}

func UpdateOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	req := dto.UpdateOrderRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := database.DB.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	// Update order details
	order.CustomerName = req.CustomerName

	if err := database.DB.Save(&order).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update order"})
		return
	}

	// Clear existing items associated with the order
	if err := database.DB.Model(&order).Association("Items").Clear(); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to clear existing order items"})
		return
	}

	// Update items and associate them with the order
	for _, itemReq := range req.Items {
		orderItem := models.Items{
			ItemCode:    itemReq.ItemCode,
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
		}

		// Check if the item already exists
		var existingItem models.Items
		if err := database.DB.Where("order_id = ? AND item_code = ?", order.ID, itemReq.ItemCode).First(&existingItem).Error; err == nil {
			// Update existing item
			existingItem.Description = itemReq.Description
			existingItem.Quantity = itemReq.Quantity

			if err := database.DB.Save(&existingItem).Error; err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to update existing order item"})
				return
			}

			// Associate the updated item with the order
			if err := database.DB.Model(&order).Association("Items").Append(&existingItem); err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to associate updated order item with order"})
				return
			}
		} else {
			// Create a new item
			if err := database.DB.Create(&orderItem).Error; err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to create new order item"})
				return
			}

			// Associate the new item with the order
			if err := database.DB.Model(&order).Association("Items").Append(&orderItem); err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to associate new order item with order"})
				return
			}
		}
	}

	ctx.JSON(200, gin.H{"order_id": order.ID})
}

func DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	var order models.Order
	if err := database.DB.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	// Delete associated items from the order_items table
	if err := database.DB.Model(&order).Association("Items").Clear(); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete associated items"})
		return
	}

	// Delete the order
	if err := database.DB.Delete(&order).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete order"})
		return
	}

	ctx.JSON(200, gin.H{"order_id": order.ID})
}
