package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vocal-tracker/broker-service/config"
	"github.com/vocal-tracker/broker-service/middleware"
	pb "github.com/vocal-tracker/broker-service/proto"
)

type VocabHandler struct {
	cfg *config.Config
}

// Request types
type CreateVocabRequest struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
	Example string `json:"example"`
	Date    string `json:"date"`
	Status  string `json:"status,omitempty"`
}

type UpdateVocabRequest struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
	Example string `json:"example"`
	Status  string `json:"status"`
}

// Response types
type VocabResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Vocab   *Vocabulary `json:"vocabulary,omitempty"`
}

type VocabListResponse struct {
	Success      bool         `json:"success"`
	Message      string       `json:"message"`
	Vocabularies []Vocabulary `json:"vocabularies"`
	Count        int32        `json:"count"`
	Total        int32        `json:"total"`
}

type Vocabulary struct {
	ID        uint32 `json:"id"`
	UserID    uint32 `json:"user_id"`
	Word      string `json:"word"`
	Meaning   string `json:"meaning"`
	Example   string `json:"example"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewVocabHandler(cfg *config.Config) *VocabHandler {
	return &VocabHandler{cfg: cfg}
}

// GetVocabularies handles GET /vocab
func (v *VocabHandler) GetVocabularies(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	date := r.URL.Query().Get("date")
	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit, offset int32 = 50, 0 // defaults

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = int32(o)
		}
	}

	// Create gRPC request
	grpcReq := &pb.GetVocabulariesRequest{
		UserId: user.UserID,
		Date:   date,
		Search: search,
		Limit:  limit,
		Offset: offset,
	}

	// Call vocabulary service with authenticated context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Forward JWT token to vocabulary service
	ctx = middleware.CreateAuthenticatedContext(ctx, r)

	resp, err := v.cfg.VocabServiceClient.GetVocabularies(ctx, grpcReq)
	if err != nil {
		middleware.WriteErrorResponse(w, "Failed to get vocabularies", http.StatusInternalServerError)
		return
	}

	// Convert response
	vocabularies := make([]Vocabulary, len(resp.Vocabularies))
	for i, vocab := range resp.Vocabularies {
		vocabularies[i] = Vocabulary{
			ID:        vocab.Id,
			UserID:    vocab.UserId,
			Word:      vocab.Word,
			Meaning:   vocab.Meaning,
			Example:   vocab.Example,
			Date:      vocab.Date,
			Status:    vocab.Status,
			CreatedAt: vocab.CreatedAt,
			UpdatedAt: vocab.UpdatedAt,
		}
	}

	response := VocabListResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		Vocabularies: vocabularies,
		Count:        resp.Count,
		Total:        resp.Total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateVocabulary handles POST /vocab
func (v *VocabHandler) CreateVocabulary(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req CreateVocabRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Word == "" || req.Meaning == "" {
		middleware.WriteErrorResponse(w, "Word and meaning are required", http.StatusBadRequest)
		return
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "review_needed"
	}

	// Create gRPC request
	grpcReq := &pb.CreateVocabularyRequest{
		UserId:  user.UserID,
		Word:    req.Word,
		Meaning: req.Meaning,
		Example: req.Example,
		Date:    req.Date,
		Status:  req.Status,
	}

	// Call vocabulary service with authenticated context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Forward JWT token to vocabulary service
	ctx = middleware.CreateAuthenticatedContext(ctx, r)

	resp, err := v.cfg.VocabServiceClient.CreateVocabulary(ctx, grpcReq)
	if err != nil {
		middleware.WriteErrorResponse(w, "Failed to create vocabulary", http.StatusInternalServerError)
		return
	}

	// Convert response
	var vocab *Vocabulary
	if resp.Vocabulary != nil {
		vocab = &Vocabulary{
			ID:        resp.Vocabulary.Id,
			UserID:    resp.Vocabulary.UserId,
			Word:      resp.Vocabulary.Word,
			Meaning:   resp.Vocabulary.Meaning,
			Example:   resp.Vocabulary.Example,
			Date:      resp.Vocabulary.Date,
			Status:    resp.Vocabulary.Status,
			CreatedAt: resp.Vocabulary.CreatedAt,
			UpdatedAt: resp.Vocabulary.UpdatedAt,
		}
	}

	response := VocabResponse{
		Success: resp.Success,
		Message: resp.Message,
		Vocab:   vocab,
	}

	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateVocabulary handles PUT /vocab/{id}
func (v *VocabHandler) UpdateVocabulary(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Extract vocab ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/vocab/")
	vocabID, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		middleware.WriteErrorResponse(w, "Invalid vocabulary ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req UpdateVocabRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Word == "" || req.Meaning == "" {
		middleware.WriteErrorResponse(w, "Word and meaning are required", http.StatusBadRequest)
		return
	}

	// Create gRPC request
	grpcReq := &pb.UpdateVocabularyRequest{
		VocabularyId: uint32(vocabID),
		UserId:       user.UserID,
		Word:         req.Word,
		Meaning:      req.Meaning,
		Example:      req.Example,
		Status:       req.Status,
	}

	// Call vocabulary service with authenticated context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Forward JWT token to vocabulary service
	ctx = middleware.CreateAuthenticatedContext(ctx, r)

	resp, err := v.cfg.VocabServiceClient.UpdateVocabulary(ctx, grpcReq)
	if err != nil {
		middleware.WriteErrorResponse(w, "Failed to update vocabulary", http.StatusInternalServerError)
		return
	}

	// Convert response
	var vocab *Vocabulary
	if resp.Vocabulary != nil {
		vocab = &Vocabulary{
			ID:        resp.Vocabulary.Id,
			UserID:    resp.Vocabulary.UserId,
			Word:      resp.Vocabulary.Word,
			Meaning:   resp.Vocabulary.Meaning,
			Example:   resp.Vocabulary.Example,
			Date:      resp.Vocabulary.Date,
			Status:    resp.Vocabulary.Status,
			CreatedAt: resp.Vocabulary.CreatedAt,
			UpdatedAt: resp.Vocabulary.UpdatedAt,
		}
	}

	response := VocabResponse{
		Success: resp.Success,
		Message: resp.Message,
		Vocab:   vocab,
	}

	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteVocabulary handles DELETE /vocab/{id}
func (v *VocabHandler) DeleteVocabulary(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Extract vocab ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/vocab/")
	vocabID, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		middleware.WriteErrorResponse(w, "Invalid vocabulary ID", http.StatusBadRequest)
		return
	}

	// Create gRPC request
	grpcReq := &pb.DeleteVocabularyRequest{
		VocabularyId: uint32(vocabID),
		UserId:       user.UserID,
	}

	// Call vocabulary service with authenticated context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Forward JWT token to vocabulary service
	ctx = middleware.CreateAuthenticatedContext(ctx, r)

	resp, err := v.cfg.VocabServiceClient.DeleteVocabulary(ctx, grpcReq)
	if err != nil {
		middleware.WriteErrorResponse(w, "Failed to delete vocabulary", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": resp.Success,
		"message": resp.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}
