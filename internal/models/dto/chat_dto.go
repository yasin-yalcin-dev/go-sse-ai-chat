/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package dto

// Chat request and response DTOs

// CreateChatRequest represents the request to create a new chat
type CreateChatRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateChatRequest represents the request to update a chat
type UpdateChatRequest struct {
	Title string `json:"title" binding:"required"`
}

// ChatResponse represents the response for a chat
type ChatResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	LastMessageAt string `json:"last_message_at,omitempty"`
	MessageCount  int    `json:"message_count"`
}

// ChatListResponse represents the response for a list of chats
type ChatListResponse struct {
	Chats      []ChatResponse `json:"chats"`
	Pagination PaginationInfo `json:"pagination"`
}
