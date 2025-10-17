package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abhishek0948/Vocab_Tracker/database"
	"github.com/abhishek0948/Vocab_Tracker/models"

	"github.com/gin-gonic/gin"
)

func GetVocabulary(c *gin.Context) {
	userID, _ := c.Get("userID")
	dateStr := c.Query("date")
	search := c.Query("search")

	var vocabularies []models.Vocabulary
	query := database.DB.Where("user_id = ?", userID)

	// Filter by date if provided
	if dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			query = query.Where("date = ?", date)
		}
	}

	// Search filter
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(word) LIKE ? OR LOWER(meaning) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Order("created_at DESC").Find(&vocabularies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vocabularies": vocabularies,
		"count":       len(vocabularies),
	})
}

func CreateVocabulary(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req models.VocabRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "review_needed"
	}

	vocabulary := models.Vocabulary{
		UserID:  userID.(uint),
		Word:    req.Word,
		Meaning: req.Meaning,
		Example: req.Example,
		Date:    date,
		Status:  status,
	}

	if err := database.DB.Create(&vocabulary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vocabulary"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Vocabulary created successfully",
		"vocabulary": vocabulary,
	})
}

func UpdateVocabulary(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vocabulary ID"})
		return
	}

	var req models.VocabUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var vocabulary models.Vocabulary
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&vocabulary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vocabulary not found"})
		return
	}

	// Update fields if provided
	if req.Word != "" {
		vocabulary.Word = req.Word
	}
	if req.Meaning != "" {
		vocabulary.Meaning = req.Meaning
	}
	if req.Example != "" {
		vocabulary.Example = req.Example
	}
	if req.Status != "" {
		vocabulary.Status = req.Status
	}

	if err := database.DB.Save(&vocabulary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Vocabulary updated successfully",
		"vocabulary": vocabulary,
	})
}

func DeleteVocabulary(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vocabulary ID"})
		return
	}

	var vocabulary models.Vocabulary
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&vocabulary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vocabulary not found"})
		return
	}

	if err := database.DB.Delete(&vocabulary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vocabulary deleted successfully",
	})
}