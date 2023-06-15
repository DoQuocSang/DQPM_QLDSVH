package controllers

import (
	"heritage-management/api/db"
	"heritage-management/api/models"
	"heritage-management/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPagedManagementUnits trả về danh sách tất cả các đơn vị quản lý với phân trang
func GetPagedManagementUnits(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	columnName := c.DefaultQuery("columnName", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	var total int64
	var ManagementUnits []models.Management_Unit

	// Đếm tổng số lượng
	if err := db.GetDB().Model(&models.Management_Unit{}).Count(&total).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not get data")
		return
	}

	// Đếm tổng số trang
	// Chia % vì nếu chia có dư thì đồng nghĩa vẫn còn trang sau nên phải tăng thêm 1
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	// Phân trang
	offset := (page - 1) * limit
	orderClause := columnName + " " + sortOrder
	if err := db.GetDB().Order(orderClause).Offset(offset).Limit(limit).Find(&ManagementUnits).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not get data")
		return
	}

	// Kiểm tra dữ liệu trả về rỗng
	if len(ManagementUnits) == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "No data available")
		return
	}

	// Tạo đối tượng phản hồi phân trang
	pagination := utils.Pagination{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       ManagementUnits,
	}

	utils.SuccessResponse(c, http.StatusOK, pagination)
}

// GetManagementUnitByID trả về thông tin của một đơn vị quản lý dựa trên ID
func GetManagementUnitByID(c *gin.Context) {
	id := c.Param("id")

	var ManagementUnit models.Management_Unit

	if err := db.GetDB().Where("id = ?", id).First(&ManagementUnit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Management unit not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, ManagementUnit)
}

// CreateManagementUnit tạo mới một đơn vị quản lý
func CreateManagementUnit(c *gin.Context) {
	var ManagementUnit models.Management_Unit

	if err := c.ShouldBindJSON(&ManagementUnit); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := db.GetDB().Create(&ManagementUnit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not create management unit")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, ManagementUnit)
}

// UpdateManagementUnit cập nhật thông tin của một đơn vị quản lý dựa trên ID
func UpdateManagementUnit(c *gin.Context) {
	id := c.Param("id")

	var ManagementUnit models.Management_Unit

	// Lấy thông tin về đơn vị quản lý dựa trên ID từ cơ sở dữ liệu
	if err := db.GetDB().Where("id = ?", id).First(&ManagementUnit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Management unit not found")
		return
	}

	// Parse thông tin cập nhật từ request body
	if err := c.ShouldBindJSON(&ManagementUnit); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Lưu thông tin cập nhật vào cơ sở dữ liệu
	if err := db.GetDB().Save(&ManagementUnit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not update management unit")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, ManagementUnit)
}

// DeleteManagementUnit xóa một đơn vị quản lý dựa trên ID
func DeleteManagementUnit(c *gin.Context) {
	id := c.Param("id")

	var ManagementUnit models.Management_Unit

	// Lấy thông tin về đơn vị quản lý dựa trên ID từ cơ sở dữ liệu
	if err := db.GetDB().Where("id = ?", id).First(&ManagementUnit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Management unit not found")
		return
	}

	// Xóa đơn vị quản lý khỏi cơ sở dữ liệu
	if err := db.GetDB().Delete(&ManagementUnit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not delete management unit")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Management unit deleted successfully"})
}

// GetHeritageByUnitSlug trả về danh sách di sản văn hóa dựa trên URL slug của đơn vị
func GetHeritageByUnitSlug(c *gin.Context) {
	unitSlug := c.Param("urlSlug")

	var unit models.Management_Unit
	if err := db.GetDB().Where("urlslug = ?", unitSlug).First(&unit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Unit not found")
		return
	}

	var heritage []models.Heritage
	if err := db.GetDB().Where("management_unit_id = ?", unit.ID).Find(&heritage).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Heritage not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, heritage)
}

// GetPagedHeritageByUnitSlug trả về danh sách di sản văn hóa dựa trên URL slug của đơn vị có phân trang
func GetPagedHeritageByUnitSlug(c *gin.Context) {
	unitSlug := c.Param("urlSlug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Lấy thông tin của địa điểm dựa trên URL slug
	var unit models.Management_Unit
	if err := db.GetDB().Where("urlslug = ?", unitSlug).First(&unit).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Unit not found")
		return
	}

	// Tìm tổng số lượng di sản văn hóa dựa trên ID của đơn vị
	var total int64
	if err := db.GetDB().Model(&models.Heritage{}).Where("management_unit_id = ?", unit.ID).Count(&total).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not get heritage")
		return
	}

	// Tính toán số trang và offset
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	offset := (page - 1) * limit

	// Truy vấn di sản văn hóa dựa trên ID của đơn vị và phân trang
	var heritage []models.Heritage
	if err := db.GetDB().Where("management_unit_id = ?", unit.ID).Offset(offset).Limit(limit).Find(&heritage).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not get heritage")
		return
	}

	// Kiểm tra dữ liệu trả về rỗng
	if len(heritage) == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "No heritage available")
		return
	}

	// Tạo đối tượng phản hồi phân trang
	pagination := utils.Pagination{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       heritage,
	}

	utils.SuccessResponse(c, http.StatusOK, pagination)
}

// SearchUnit trả về thông tin đơn vị theo tên
func SearchUnit(c *gin.Context) {
	hq := models.HeritageQuery{}
	if err := c.ShouldBindQuery(&hq); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid search parameters")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	columnName := c.DefaultQuery("columnName", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	query := db.GetDB().Model(&models.Management_Unit{})

	if hq.Key != "" {
		query = query.Where("name LIKE ?", "%"+hq.Key+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not get data")
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	offset := (page - 1) * limit
	orderClause := columnName + " " + sortOrder

	var units []models.Management_Unit
	if err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&units).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not get data")
		return
	}

	pagination := utils.Pagination{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       units,
	}

	utils.SuccessResponse(c, http.StatusOK, pagination)
}
