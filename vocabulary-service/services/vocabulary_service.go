package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/vocal-tracker/vocabulary-service/database"
	"github.com/vocal-tracker/vocabulary-service/middleware"
	"github.com/vocal-tracker/vocabulary-service/models"
	"github.com/vocal-tracker/vocabulary-service/proto"

	"gorm.io/gorm"
)

type VocabularyServiceImpl struct {
	proto.UnimplementedVocabularyServiceServer
}

func NewVocabularyService() *VocabularyServiceImpl {
	return &VocabularyServiceImpl{}
}

// GetVocabularies implements the GetVocabularies RPC method
func (s *VocabularyServiceImpl) GetVocabularies(ctx context.Context, req *proto.GetVocabulariesRequest) (*proto.GetVocabulariesResponse, error) {
	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return &proto.GetVocabulariesResponse{
			Success: false,
			Message: "Authentication required",
		}, nil
	}

	// Verify the request is for the authenticated user
	if req.UserId != authenticatedUserID {
		return &proto.GetVocabulariesResponse{
			Success: false,
			Message: "Access denied: can only access your own vocabularies",
		}, nil
	}

	var vocabularies []models.Vocabulary
	query := database.DB.Where("user_id = ?", authenticatedUserID)

	// Filter by date if provided
	if req.Date != "" {
		if date, err := time.Parse("2006-01-02", req.Date); err == nil {
			query = query.Where("date = ?", date)
		}
	}

	// Search filter
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(word) LIKE ? OR LOWER(meaning) LIKE ?", searchTerm, searchTerm)
	}

	// Apply limit and offset
	if req.Limit > 0 {
		query = query.Limit(int(req.Limit))
	}
	if req.Offset > 0 {
		query = query.Offset(int(req.Offset))
	}

	// Get total count for pagination
	var total int64
	countQuery := database.DB.Model(&models.Vocabulary{}).Where("user_id = ?", authenticatedUserID)
	if req.Date != "" {
		if date, err := time.Parse("2006-01-02", req.Date); err == nil {
			countQuery = countQuery.Where("date = ?", date)
		}
	}
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		countQuery = countQuery.Where("LOWER(word) LIKE ? OR LOWER(meaning) LIKE ?", searchTerm, searchTerm)
	}
	countQuery.Count(&total)

	if err := query.Order("created_at DESC").Find(&vocabularies).Error; err != nil {
		return &proto.GetVocabulariesResponse{
			Success: false,
			Message: "Failed to fetch vocabularies",
		}, err
	}

	// Convert to proto format
	protoVocabs := make([]*proto.Vocabulary, len(vocabularies))
	for i, vocab := range vocabularies {
		protoVocabs[i] = &proto.Vocabulary{
			Id:        uint32(vocab.ID),
			UserId:    uint32(vocab.UserID),
			Word:      vocab.Word,
			Meaning:   vocab.Meaning,
			Example:   vocab.Example,
			Date:      vocab.Date.Format("2006-01-02"),
			Status:    vocab.Status,
			CreatedAt: vocab.CreatedAt.Format(time.RFC3339),
			UpdatedAt: vocab.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetVocabulariesResponse{
		Success:      true,
		Message:      "Vocabularies retrieved successfully",
		Vocabularies: protoVocabs,
		Count:        int32(len(vocabularies)),
		Total:        int32(total),
	}, nil
}

// CreateVocabulary implements the CreateVocabulary RPC method
func (s *VocabularyServiceImpl) CreateVocabulary(ctx context.Context, req *proto.CreateVocabularyRequest) (*proto.VocabularyResponse, error) {
	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Authentication required",
		}, nil
	}

	// Verify the request is for the authenticated user
	if req.UserId != authenticatedUserID {
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Access denied: can only create vocabularies for your own account",
		}, nil
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Invalid date format. Use YYYY-MM-DD",
		}, nil
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "review_needed"
	}

	// Create vocabulary using authenticated user ID
	vocab := models.Vocabulary{
		UserID:  uint(authenticatedUserID),
		Word:    req.Word,
		Meaning: req.Meaning,
		Example: req.Example,
		Date:    date,
		Status:  status,
	}

	if err := database.DB.Create(&vocab).Error; err != nil {
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Failed to create vocabulary",
		}, err
	}

	return &proto.VocabularyResponse{
		Success: true,
		Message: "Vocabulary created successfully",
		Vocabulary: &proto.Vocabulary{
			Id:        uint32(vocab.ID),
			UserId:    uint32(vocab.UserID),
			Word:      vocab.Word,
			Meaning:   vocab.Meaning,
			Example:   vocab.Example,
			Date:      vocab.Date.Format("2006-01-02"),
			Status:    vocab.Status,
			CreatedAt: vocab.CreatedAt.Format(time.RFC3339),
			UpdatedAt: vocab.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// UpdateVocabulary implements the UpdateVocabulary RPC method
func (s *VocabularyServiceImpl) UpdateVocabulary(ctx context.Context, req *proto.UpdateVocabularyRequest) (*proto.VocabularyResponse, error) {
	var vocab models.Vocabulary
	if err := database.DB.Where("id = ? AND user_id = ?", req.VocabularyId, req.UserId).First(&vocab).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &proto.VocabularyResponse{
				Success: false,
				Message: "Vocabulary not found",
			}, nil
		}
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Database error",
		}, err
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Word != "" {
		updates["word"] = req.Word
	}
	if req.Meaning != "" {
		updates["meaning"] = req.Meaning
	}
	if req.Example != "" {
		updates["example"] = req.Example
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	log.Println(updates);

	if len(updates) > 0 {
		if err := database.DB.Model(&vocab).Updates(updates).Error; err != nil {
			return &proto.VocabularyResponse{
				Success: false,
				Message: "Failed to update vocabulary",
			}, err
		}
	}

	// Reload the updated vocabulary
	if err := database.DB.First(&vocab, vocab.ID).Error; err != nil {
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Failed to reload vocabulary",
		}, err
	}

	return &proto.VocabularyResponse{
		Success: true,
		Message: "Vocabulary updated successfully",
		Vocabulary: &proto.Vocabulary{
			Id:        uint32(vocab.ID),
			UserId:    uint32(vocab.UserID),
			Word:      vocab.Word,
			Meaning:   vocab.Meaning,
			Example:   vocab.Example,
			Date:      vocab.Date.Format("2006-01-02"),
			Status:    vocab.Status,
			CreatedAt: vocab.CreatedAt.Format(time.RFC3339),
			UpdatedAt: vocab.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// DeleteVocabulary implements the DeleteVocabulary RPC method
func (s *VocabularyServiceImpl) DeleteVocabulary(ctx context.Context, req *proto.DeleteVocabularyRequest) (*proto.DeleteVocabularyResponse, error) {
	result := database.DB.Where("id = ? AND user_id = ?", req.VocabularyId, req.UserId).Delete(&models.Vocabulary{})
	if result.Error != nil {
		return &proto.DeleteVocabularyResponse{
			Success: false,
			Message: "Failed to delete vocabulary",
		}, result.Error
	}

	if result.RowsAffected == 0 {
		return &proto.DeleteVocabularyResponse{
			Success: false,
			Message: "Vocabulary not found",
		}, nil
	}

	return &proto.DeleteVocabularyResponse{
		Success: true,
		Message: "Vocabulary deleted successfully",
	}, nil
}

// GetVocabularyById implements the GetVocabularyById RPC method
func (s *VocabularyServiceImpl) GetVocabularyById(ctx context.Context, req *proto.GetVocabularyByIdRequest) (*proto.VocabularyResponse, error) {
	var vocab models.Vocabulary
	if err := database.DB.Where("id = ? AND user_id = ?", req.VocabularyId, req.UserId).First(&vocab).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &proto.VocabularyResponse{
				Success: false,
				Message: "Vocabulary not found",
			}, nil
		}
		return &proto.VocabularyResponse{
			Success: false,
			Message: "Database error",
		}, err
	}

	return &proto.VocabularyResponse{
		Success: true,
		Message: "Vocabulary retrieved successfully",
		Vocabulary: &proto.Vocabulary{
			Id:        uint32(vocab.ID),
			UserId:    uint32(vocab.UserID),
			Word:      vocab.Word,
			Meaning:   vocab.Meaning,
			Example:   vocab.Example,
			Date:      vocab.Date.Format("2006-01-02"),
			Status:    vocab.Status,
			CreatedAt: vocab.CreatedAt.Format(time.RFC3339),
			UpdatedAt: vocab.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetVocabularyStats implements the GetVocabularyStats RPC method
func (s *VocabularyServiceImpl) GetVocabularyStats(ctx context.Context, req *proto.GetVocabularyStatsRequest) (*proto.VocabularyStatsResponse, error) {
	// Total words count
	var totalWords int64
	database.DB.Model(&models.Vocabulary{}).Where("user_id = ?", req.UserId).Count(&totalWords)

	// Words this week
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	var wordsThisWeek int64
	database.DB.Model(&models.Vocabulary{}).Where("user_id = ? AND created_at >= ?", req.UserId, weekStart).Count(&wordsThisWeek)

	// Words this month
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	var wordsThisMonth int64
	database.DB.Model(&models.Vocabulary{}).Where("user_id = ? AND created_at >= ?", req.UserId, monthStart).Count(&wordsThisMonth)

	// Status counts
	statusCounts := make(map[string]int32)
	var statusResults []struct {
		Status string
		Count  int64
	}
	database.DB.Model(&models.Vocabulary{}).
		Select("status, count(*) as count").
		Where("user_id = ?", req.UserId).
		Group("status").
		Find(&statusResults)

	for _, result := range statusResults {
		statusCounts[result.Status] = int32(result.Count)
	}

	// Daily counts for the last 30 days
	var dailyCounts []*proto.DailyCount
	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		var count int64
		database.DB.Model(&models.Vocabulary{}).
			Where("user_id = ? AND DATE(created_at) = ?", req.UserId, dateStr).
			Count(&count)

		dailyCounts = append(dailyCounts, &proto.DailyCount{
			Date:  dateStr,
			Count: int32(count),
		})
	}

	return &proto.VocabularyStatsResponse{
		Success:        true,
		Message:        "Statistics retrieved successfully",
		TotalWords:     int32(totalWords),
		WordsThisWeek:  int32(wordsThisWeek),
		WordsThisMonth: int32(wordsThisMonth),
		StatusCounts:   statusCounts,
		DailyCounts:    dailyCounts,
	}, nil
}
