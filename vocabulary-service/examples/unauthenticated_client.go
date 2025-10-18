package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/vocal-tracker/vocabulary-service/proto"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func main() {
// 	log.Println("=== Vocabulary Service gRPC Client Example ===")
// 	log.Println("Note: Make sure you have a user in the database first!")
// 	log.Println("You can create a user using the auth service client example.")

// 	// Connect to the gRPC server
// 	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("Failed to connect: %v", err)
// 	}
// 	defer conn.Close()

// 	// Create vocabulary service client
// 	client := proto.NewVocabularyServiceClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Use an existing user ID (you should create a user first using auth service)
// 	// For this example, we'll use ID 1 - make sure this user exists!
// 	userID := uint32(4)

// 	log.Printf("Using user ID: %d", userID)

// 	// Example 1: Create a new vocabulary
// 	log.Println("=== Creating Vocabulary ===")
// 	createResp, err := client.CreateVocabulary(ctx, &proto.CreateVocabularyRequest{
// 		UserId:  userID,
// 		Word:    "serendipity",
// 		Meaning: "A pleasant surprise; the occurrence of events by chance in a happy way",
// 		Example: "Finding my old friend at the coffee shop was pure serendipity.",
// 		Date:    time.Now().Format("2006-01-02"),
// 		Status:  "review_needed",
// 	})
// 	if err != nil {
// 		log.Fatalf("Create vocabulary failed: %v", err)
// 	}
// 	log.Printf("Create response: %v", createResp)

// 	var vocabularyID uint32
// 	if createResp.Success && createResp.Vocabulary != nil {
// 		vocabularyID = createResp.Vocabulary.Id
// 		log.Printf("Created vocabulary with ID: %d", vocabularyID)
// 	}

// 	// Example 2: Get all vocabularies for user
// 	log.Println("\n=== Getting All Vocabularies ===")
// 	getResp, err := client.GetVocabularies(ctx, &proto.GetVocabulariesRequest{
// 		UserId: userID,
// 		Limit:  10,
// 		Offset: 0,
// 	})
// 	if err != nil {
// 		log.Fatalf("Get vocabularies failed: %v", err)
// 	}
// 	log.Printf("Get response: Success=%v, Count=%d, Total=%d",
// 		getResp.Success, getResp.Count, getResp.Total)
// 	for i, vocab := range getResp.Vocabularies {
// 		log.Printf("Vocabulary %d: %s - %s (Status: %s)",
// 			i+1, vocab.Word, vocab.Meaning, vocab.Status)
// 	}

// 	// Example 3: Search vocabularies
// 	log.Println("\n=== Searching Vocabularies ===")
// 	searchResp, err := client.GetVocabularies(ctx, &proto.GetVocabulariesRequest{
// 		UserId: userID,
// 		Search: "serendipity",
// 	})
// 	if err != nil {
// 		log.Fatalf("Search vocabularies failed: %v", err)
// 	}
// 	log.Printf("Search response: Found %d vocabularies", searchResp.Count)

// 	// Example 4: Update vocabulary status (if we created one)
// 	if vocabularyID > 0 {
// 		log.Println("\n=== Updating Vocabulary Status ===")
// 		updateResp, err := client.UpdateVocabulary(ctx, &proto.UpdateVocabularyRequest{
// 			VocabularyId: vocabularyID,
// 			UserId:       userID,
// 			Status:       "learning",
// 		})
// 		if err != nil {
// 			log.Fatalf("Update vocabulary failed: %v", err)
// 		}
// 		log.Printf("Update response: %v", updateResp)
// 	}

// 	// Example 5: Get vocabulary by ID
// 	if vocabularyID > 0 {
// 		log.Println("\n=== Getting Vocabulary by ID ===")
// 		getByIdResp, err := client.GetVocabularyById(ctx, &proto.GetVocabularyByIdRequest{
// 			VocabularyId: vocabularyID,
// 			UserId:       userID,
// 		})
// 		if err != nil {
// 			log.Fatalf("Get vocabulary by ID failed: %v", err)
// 		}
// 		if getByIdResp.Success && getByIdResp.Vocabulary != nil {
// 			vocab := getByIdResp.Vocabulary
// 			log.Printf("Found vocabulary: %s - %s (Status: %s)",
// 				vocab.Word, vocab.Meaning, vocab.Status)
// 		}
// 	}

// 	// Example 6: Get vocabulary statistics
// 	log.Println("\n=== Getting Vocabulary Statistics ===")
// 	statsResp, err := client.GetVocabularyStats(ctx, &proto.GetVocabularyStatsRequest{
// 		UserId: userID,
// 	})
// 	if err != nil {
// 		log.Fatalf("Get vocabulary stats failed: %v", err)
// 	}
// 	if statsResp.Success {
// 		log.Printf("Statistics:")
// 		log.Printf("  Total words: %d", statsResp.TotalWords)
// 		log.Printf("  Words this week: %d", statsResp.WordsThisWeek)
// 		log.Printf("  Words this month: %d", statsResp.WordsThisMonth)
// 		log.Printf("  Status counts: %v", statsResp.StatusCounts)
// 		log.Printf("  Recent daily counts (last 5 days):")
// 		for i := len(statsResp.DailyCounts) - 5; i < len(statsResp.DailyCounts) && i >= 0; i++ {
// 			if i >= 0 {
// 				daily := statsResp.DailyCounts[i]
// 				log.Printf("    %s: %d words", daily.Date, daily.Count)
// 			}
// 		}
// 	}

// 	// Example 7: Create another vocabulary for testing
// 	log.Println("\n=== Creating Another Vocabulary ===")
// 	createResp2, err := client.CreateVocabulary(ctx, &proto.CreateVocabularyRequest{
// 		UserId:  userID,
// 		Word:    "ephemeral",
// 		Meaning: "Lasting for a very short time; transitory",
// 		Example: "The beauty of cherry blossoms is ephemeral, lasting only a few weeks.",
// 		Date:    time.Now().Format("2006-01-02"),
// 		Status:  "review_needed",
// 	})
// 	if err != nil {
// 		log.Printf("Create second vocabulary failed: %v", err)
// 	} else {
// 		log.Printf("Created second vocabulary: %v", createResp2.Success)
// 	}

// 	// Example 8: Get vocabularies with date filter (today's vocabularies)
// 	log.Println("\n=== Getting Today's Vocabularies ===")
// 	todayResp, err := client.GetVocabularies(ctx, &proto.GetVocabulariesRequest{
// 		UserId: userID,
// 		Date:   time.Now().Format("2006-01-02"),
// 	})
// 	if err != nil {
// 		log.Fatalf("Get today's vocabularies failed: %v", err)
// 	}
// 	log.Printf("Today's vocabularies: %d found", todayResp.Count)
// 	for _, vocab := range todayResp.Vocabularies {
// 		log.Printf("  - %s: %s", vocab.Word, vocab.Meaning)
// 	}
// }