package repository

import (
	"context"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

func TestAIConversationRepo_Lifecycle(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAIConversationRepo(db)
	ctx := context.Background()

	// Query an existing user from the database
	var testUserID int
	err := db.QueryRow(ctx, "SELECT id FROM users LIMIT 1").Scan(&testUserID)
	if err != nil {
		t.Skip("Skipping test: no users exist in database")
		return
	}

	// 1. Create Conversation
	conv, err := repo.CreateConversation(ctx, testUserID, "Test Conversation")
	if err != nil {
		t.Fatalf("CreateConversation failed: %v", err)
	}
	if conv.ID == "" {
		t.Fatalf("expected non-empty conversation ID")
	}
	if conv.Title != "Test Conversation" {
		t.Errorf("expected title 'Test Conversation', got '%s'", conv.Title)
	}

	// 2. Get Conversation with ownership check
	fetched, err := repo.GetConversationByID(ctx, conv.ID, testUserID)
	if err != nil {
		t.Fatalf("GetConversationByID failed: %v", err)
	}
	if fetched.ID != conv.ID {
		t.Errorf("expected conv ID %s, got %s", conv.ID, fetched.ID)
	}

	// Ownership enforcement: another user cannot access it
	_, err = repo.GetConversationByID(ctx, conv.ID, 999999)
	if err == nil {
		t.Errorf("expected unauthorized error for user 999999, got nil")
	}

	// 3. Add Messages
	msg1 := &models.AIMessage{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        "Find me a phone in Addis",
	}
	if err := repo.CreateMessage(ctx, msg1); err != nil {
		t.Fatalf("CreateMessage (user) failed: %v", err)
	}

	msg2 := &models.AIMessage{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        "Here is what I found",
		Metadata: map[string]interface{}{
			"model": "huntiq-assistant:latest",
		},
	}
	if err := repo.CreateMessage(ctx, msg2); err != nil {
		t.Fatalf("CreateMessage (assistant) failed: %v", err)
	}

	// 4. Retrieve Messages
	messages, err := repo.GetMessagesByConversationID(ctx, conv.ID)
	if err != nil {
		t.Fatalf("GetMessagesByConversationID failed: %v", err)
	}
	if len(messages) < 2 {
		t.Fatalf("expected at least 2 messages, got %d", len(messages))
	}
	if messages[0].Role != "user" || messages[1].Role != "assistant" {
		t.Errorf("unexpected message sequence: %s, %s", messages[0].Role, messages[1].Role)
	}

	// 5. Update Title
	err = repo.UpdateConversationTitle(ctx, conv.ID, testUserID, "Updated Phone Search")
	if err != nil {
		t.Fatalf("UpdateConversationTitle failed: %v", err)
	}

	// 6. List Conversations
	list, err := repo.ListConversationsByUserID(ctx, testUserID, 10, 0)
	if err != nil {
		t.Fatalf("ListConversationsByUserID failed: %v", err)
	}
	found := false
	for _, c := range list {
		if c.ID == conv.ID {
			found = true
			if c.Title != "Updated Phone Search" {
				t.Errorf("expected updated title, got '%s'", c.Title)
			}
			break
		}
	}
	if !found {
		t.Errorf("expected created conversation to appear in list")
	}

	// 7. Delete Conversation
	err = repo.DeleteConversation(ctx, conv.ID, testUserID)
	if err != nil {
		t.Fatalf("DeleteConversation failed: %v", err)
	}

	// Verify conversation is gone
	_, err = repo.GetConversationByID(ctx, conv.ID, testUserID)
	if err == nil {
		t.Errorf("expected conversation to be deleted")
	}

	// Verify messages are deleted (cascade)
	remainingMsgs, err := repo.GetMessagesByConversationID(ctx, conv.ID)
	if err != nil {
		t.Fatalf("GetMessagesByConversationID after delete failed: %v", err)
	}
	if len(remainingMsgs) != 0 {
		t.Errorf("expected 0 remaining messages after conversation delete, got %d", len(remainingMsgs))
	}
}
