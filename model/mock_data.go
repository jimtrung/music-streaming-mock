package model

import (
	"time"

	"github.com/google/uuid"
)

// MockDataJSON represents the complete mock data structure
// aligned with MOCK_DATA_SCHEMA.json
type MockDataJSON struct {
	Metadata struct {
		Version        string `json:"version"`
		GeneratedAt    string `json:"generatedAt"`
		DatasetName    string `json:"datasetName"`
		Description    string `json:"description,omitempty"`
		TotalUsers     int    `json:"totalUsers"`
		TotalArtists   int    `json:"totalArtists"`
		TotalTracks    int    `json:"totalTracks"`
		TotalPlaylists int    `json:"totalPlaylists"`
	} `json:"metadata"`

	Users          []UserJSON          `json:"users"`
	Profiles       []ProfileJSON       `json:"profiles"`
	Artists        []ArtistJSON        `json:"artists"`
	Tracks         []TrackJSON         `json:"tracks"`
	Playlists      []PlaylistJSON      `json:"playlists"`
	PlaylistTracks []PlaylistTrackJSON `json:"playlistTracks"`
	RefreshTokens  []RefreshTokenJSON  `json:"refreshTokens,omitempty"`
}

// UserJSON represents a user entity
type UserJSON struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PlainPasswd string `json:"plainPassword,omitempty"`
	Role        string `json:"role"`
	Provider    string `json:"provider"`
	IsVerified  bool   `json:"isVerified"`
	IsPremium   bool   `json:"isPremium"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ProfileJSON represents a user profile
type ProfileJSON struct {
	UserId    string `json:"userId"`
	Name      string `json:"name,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

// ArtistJSON represents an artist entity
type ArtistJSON struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	Name       string `json:"name"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// TrackJSON represents a track/song entity
type TrackJSON struct {
	Id        string `json:"id"`
	ArtistId  string `json:"artistId"`
	AlbumId   string `json:"albumId,omitempty"`
	Title     string `json:"title"`
	AudioUrl  string `json:"audioUrl"`
	CoverUrl  string `json:"coverUrl,omitempty"`
	TrackNum  int    `json:"trackNumber,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Genre     string `json:"genre,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// PlaylistJSON represents a playlist entity
type PlaylistJSON struct {
	Id          string   `json:"id"`
	OwnerId     string   `json:"ownerId"`
	Name        string   `json:"name"`
	IsPublic    bool     `json:"isPublic"`
	Description string   `json:"description,omitempty"`
	CoverUrl    string   `json:"coverUrl,omitempty"`
	TrackIds    []string `json:"trackIds"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// PlaylistTrackJSON represents the junction table entry
type PlaylistTrackJSON struct {
	PlaylistId string `json:"playlistId"`
	TrackId    string `json:"trackId"`
	Position   int    `json:"position"`
	AddedAt    string `json:"addedAt"`
}

// RefreshTokenJSON represents a refresh token
type RefreshTokenJSON struct {
	UserId    string `json:"userId"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
}

// Helper functions for data generation

// NewUser creates a new user JSON object
func NewUser(username, email, hashedPassword, plainPassword string, isPremium bool) UserJSON {
	now := time.Now().UTC().Format(time.RFC3339)
	return UserJSON{
		Id:          uuid.New().String(),
		Username:    username,
		Email:       email,
		Password:    hashedPassword,
		PlainPasswd: plainPassword,
		Role:        "listener",
		Provider:    "local",
		IsVerified:  true,
		IsPremium:   isPremium,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewArtist creates a new artist JSON object
func NewArtist(userId, name string) ArtistJSON {
	now := time.Now().UTC().Format(time.RFC3339)
	return ArtistJSON{
		Id:         uuid.New().String(),
		UserId:     userId,
		Name:       name,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// NewTrack creates a new track JSON object
func NewTrack(artistId, title, audioUrl string) TrackJSON {
	now := time.Now().UTC().Format(time.RFC3339)
	return TrackJSON{
		Id:        uuid.New().String(),
		ArtistId:  artistId,
		Title:     title,
		AudioUrl:  audioUrl,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewPlaylist creates a new playlist JSON object
func NewPlaylist(ownerId, name string, isPublic bool) PlaylistJSON {
	now := time.Now().UTC().Format(time.RFC3339)
	return PlaylistJSON{
		Id:        uuid.New().String(),
		OwnerId:   ownerId,
		Name:      name,
		IsPublic:  isPublic,
		TrackIds:  []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewMockData creates a new empty MockDataJSON
func NewMockData(name, description string) *MockDataJSON {
	return &MockDataJSON{
		Users:          []UserJSON{},
		Profiles:       []ProfileJSON{},
		Artists:        []ArtistJSON{},
		Tracks:         []TrackJSON{},
		Playlists:      []PlaylistJSON{},
		PlaylistTracks: []PlaylistTrackJSON{},
		RefreshTokens:  []RefreshTokenJSON{},
	}
}

// UpdateMetadata sets the metadata for the mock data
func (m *MockDataJSON) UpdateMetadata(version, datasetName, description string) {
	m.Metadata.Version = version
	m.Metadata.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	m.Metadata.DatasetName = datasetName
	m.Metadata.Description = description
	m.Metadata.TotalUsers = len(m.Users)
	m.Metadata.TotalArtists = len(m.Artists)
	m.Metadata.TotalTracks = len(m.Tracks)
	m.Metadata.TotalPlaylists = len(m.Playlists)
}
