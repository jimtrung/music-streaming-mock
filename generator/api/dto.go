package api

import (
	"time"

	"github.com/google/uuid"
)

// Request and Response DTOs for API communication
// These align with the Music Streaming API specifications

// SignInRequest is the DTO for sign in requests
type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignUpRequest is the DTO for sign up requests
type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse is the DTO for token responses
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// RefreshResponse is the DTO for token refresh responses
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// CreateTrackResponse is the DTO for track creation responses
type CreateTrackResponse struct {
	Id        uuid.UUID `json:"id"`
	ArtistId  uuid.UUID `json:"artist_id"`
	Title     string    `json:"title"`
	AudioUrl  string    `json:"audio_url"`
	CoverUrl  string    `json:"cover_url,omitempty"`
	TrackNum  int       `json:"track_number,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateArtistRequest is the DTO for artist creation
type CreateArtistRequest struct {
	Name string `json:"name"`
}

// CreateArtistResponse is the DTO for artist creation responses
type CreateArtistResponse struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UpdateArtistRequest is the DTO for artist updates
type UpdateArtistRequest struct {
	Name string `json:"name"`
}

// CreatePlaylistRequest is the DTO for playlist creation
type CreatePlaylistRequest struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

// CreatePlaylistResponse is the DTO for playlist creation responses
type CreatePlaylistResponse struct {
	Id        uuid.UUID `json:"id"`
	OwnerId   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddTrackRequest is the DTO for adding tracks to playlists
type AddTrackRequest struct {
	PlaylistId uuid.UUID `json:"playlist_id"`
	TrackId    uuid.UUID `json:"track_id"`
	Position   int       `json:"position"`
}

// UpdateProfileRequest is the DTO for profile updates
type UpdateProfileRequest struct {
	Name string `json:"name,omitempty"`
}

// ProfileResponse is the DTO for profile responses
type ProfileResponse struct {
	Id         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email"`
	AvatarUrl  string    `json:"avatar_url,omitempty"`
	IsVerified bool      `json:"is_verified"`
	IsPremium  bool      `json:"is_premium"`
}

// UserResponse is the DTO for user responses
type UserResponse struct {
	Id         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	Provider   string    `json:"provider"`
	IsVerified bool      `json:"is_verified"`
	IsPremium  bool      `json:"is_premium"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
