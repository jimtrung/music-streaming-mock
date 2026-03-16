package generators

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/model"
	"golang.org/x/crypto/bcrypt"
)

// UserGenerator generates mock users
type UserGenerator struct {
	ctx *GeneratorContext
}

// NewUserGenerator creates a new user generator
func NewUserGenerator(ctx *GeneratorContext) *UserGenerator {
	return &UserGenerator{ctx: ctx}
}

// Generate creates mock users
func (ug *UserGenerator) Generate(quantity int) error {
	fmt.Printf("Generating %d users...\n", quantity)

	users := make([]model.UserJSON, quantity)

	for i := 0; i < quantity; i++ {
		username := ug.ctx.GetRandomUsername()
		password := ug.ctx.GetRandomPassword()
		email := fmt.Sprintf("%s@example.com", username)

		// Hash password with bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: Failed to hash password for %s: %v\n", username, err)
			hashedPassword = []byte(password) // Fallback to plain password
		}

		user := model.UserJSON{
			Id:          uuid.New().String(),
			Username:    username,
			Email:       email,
			Password:    string(hashedPassword),
			PlainPasswd: password,
			Role:        "listener",
			Provider:    "local",
			IsVerified:  true,
			IsPremium:   i%10 == 0, // 10% premium users
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		}

		users[i] = user

		if (i+1)%10 == 0 {
			fmt.Printf("  ✓ Generated %d users\n", i+1)
		}
	}

	// Save all users at once
	if err := ug.ctx.AddUsers(users); err != nil {
		return fmt.Errorf("failed to save users: %w", err)
	}

	fmt.Printf("✓ Successfully generated and saved %d users\n", quantity)
	return nil
}

// ProfileGenerator generates mock user profiles
type ProfileGenerator struct {
	ctx *GeneratorContext
}

// NewProfileGenerator creates a new profile generator
func NewProfileGenerator(ctx *GeneratorContext) *ProfileGenerator {
	return &ProfileGenerator{ctx: ctx}
}

// Generate creates profiles for existing users
func (pg *ProfileGenerator) Generate() error {
	fmt.Println("Generating profiles for users...")

	users := pg.ctx.Users
	if len(users) == 0 {
		return fmt.Errorf("no users found in context")
	}

	profiles := make([]model.ProfileJSON, len(users))

	for i, user := range users {
		profile := model.ProfileJSON{
			UserId:    user.Id,
			Name:      fmt.Sprintf("%s Profile", user.Username),
			AvatarUrl: "",
			Bio:       fmt.Sprintf("This is the profile of %s", user.Username),
		}

		profiles[i] = profile
	}

	// Save all profiles
	if err := pg.ctx.AddProfiles(profiles); err != nil {
		return fmt.Errorf("failed to save profiles: %w", err)
	}

	fmt.Printf("✓ Successfully generated and saved %d profiles\n", len(profiles))
	return nil
}

// ArtistGenerator generates mock artists
type ArtistGenerator struct {
	ctx *GeneratorContext
}

// NewArtistGenerator creates a new artist generator
func NewArtistGenerator(ctx *GeneratorContext) *ArtistGenerator {
	return &ArtistGenerator{ctx: ctx}
}

// Generate creates mock artists from users
func (ag *ArtistGenerator) Generate(quantity int) error {
	fmt.Printf("Generating %d artists...\n", quantity)

	users := ag.ctx.Users
	if len(users) == 0 {
		return fmt.Errorf("no users found - generate users first")
	}

	if quantity > len(users) {
		quantity = len(users)
		fmt.Printf("Limited to %d artists (number of users)\n", quantity)
	}

	artists := make([]model.ArtistJSON, quantity)

	for i := 0; i < quantity; i++ {
		user := users[i]
		artistName := ag.ctx.GetRandomArtistName()

		artist := model.ArtistJSON{
			Id:         uuid.New().String(),
			UserId:     user.Id,
			Name:       artistName,
			IsVerified: i%20 == 0, // 5% verified artists
			CreatedAt:  time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
		}

		artists[i] = artist

		if (i+1)%10 == 0 {
			fmt.Printf("  ✓ Generated %d artists\n", i+1)
		}
	}

	// Save all artists
	if err := ag.ctx.AddArtists(artists); err != nil {
		return fmt.Errorf("failed to save artists: %w", err)
	}

	fmt.Printf("✓ Successfully generated and saved %d artists\n", quantity)
	return nil
}

// TrackGenerator generates mock tracks
type TrackGenerator struct {
	ctx *GeneratorContext
}

// NewTrackGenerator creates a new track generator
func NewTrackGenerator(ctx *GeneratorContext) *TrackGenerator {
	return &TrackGenerator{ctx: ctx}
}

// Generate creates mock tracks from artists
func (tg *TrackGenerator) Generate(quantity int) error {
	fmt.Printf("Generating %d tracks...\n", quantity)

	artists := tg.ctx.Artists
	if len(artists) == 0 {
		return fmt.Errorf("no artists found - generate artists first")
	}

	tracks := make([]model.TrackJSON, quantity)

	genres := []string{"Pop", "Rock", "Hip-Hop", "Jazz", "Classical", "Electronic", "Country", "Reggae"}

	for i := 0; i < quantity; i++ {
		// Distribute tracks evenly among artists
		artistIdx := i % len(artists)
		artist := artists[artistIdx]

		title := tg.ctx.GetRandomTrackTitle()
		genre := genres[i%len(genres)]

		track := model.TrackJSON{
			Id:        uuid.New().String(),
			ArtistId:  artist.Id,
			Title:     title,
			AudioUrl:  fmt.Sprintf("https://storage.example.com/tracks/%s.mp3", uuid.New().String()),
			CoverUrl:  fmt.Sprintf("https://storage.example.com/covers/%s.jpg", uuid.New().String()),
			TrackNum:  (i % 20) + 1,
			Duration:  180 + (i * 5 % 120), // 3-5 minutes
			Genre:     genre,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		tracks[i] = track

		if (i+1)%50 == 0 {
			fmt.Printf("  ✓ Generated %d tracks\n", i+1)
		}
	}

	// Save all tracks
	if err := tg.ctx.AddTracks(tracks); err != nil {
		return fmt.Errorf("failed to save tracks: %w", err)
	}

	fmt.Printf("✓ Successfully generated and saved %d tracks\n", quantity)
	return nil
}

// PlaylistGenerator generates mock playlists
type PlaylistGenerator struct {
	ctx *GeneratorContext
}

// NewPlaylistGenerator creates a new playlist generator
func NewPlaylistGenerator(ctx *GeneratorContext) *PlaylistGenerator {
	return &PlaylistGenerator{ctx: ctx}
}

// Generate creates mock playlists with tracks
func (pg *PlaylistGenerator) Generate(quantity int, tracksPerPlaylist int) error {
	fmt.Printf("Generating %d playlists with ~%d tracks each...\n", quantity, tracksPerPlaylist)

	users := pg.ctx.Users
	tracks := pg.ctx.Tracks

	if len(users) == 0 {
		return fmt.Errorf("no users found - generate users first")
	}

	if len(tracks) == 0 {
		return fmt.Errorf("no tracks found - generate tracks first")
	}

	playlists := make([]model.PlaylistJSON, quantity)
	var playlistTracks []model.PlaylistTrackJSON

	for i := 0; i < quantity; i++ {
		// Select a random user to own the playlist
		userIdx := i % len(users)
		user := users[userIdx]

		playlistName := pg.ctx.GetRandomPlaylistName()

		// Select random tracks for this playlist
		trackIDsForPlaylist := make([]string, 0, tracksPerPlaylist)
		selectedTrackIndices := make(map[int]bool)

		for j := 0; j < tracksPerPlaylist && j < len(tracks); j++ {
			trackIdx := i * j % len(tracks)

			// Avoid duplicate tracks in same playlist
			if !selectedTrackIndices[trackIdx] {
				selectedTrackIndices[trackIdx] = true
				trackIDsForPlaylist = append(trackIDsForPlaylist, tracks[trackIdx].Id)

				// Add to playlist tracks junction
				playlistTracks = append(playlistTracks, model.PlaylistTrackJSON{
					PlaylistId: uuid.New().String(), // Will be set from playlist ID
					TrackId:    tracks[trackIdx].Id,
					Position:   j + 1,
					AddedAt:    time.Now().UTC().Format(time.RFC3339),
				})
			}
		}

		playlist := model.PlaylistJSON{
			Id:          uuid.New().String(),
			OwnerId:     user.Id,
			Name:        playlistName,
			IsPublic:    i%2 == 0, // 50% public
			Description: fmt.Sprintf("A great playlist with %d tracks", len(trackIDsForPlaylist)),
			CoverUrl:    fmt.Sprintf("https://storage.example.com/covers/%s.jpg", uuid.New().String()),
			TrackIds:    trackIDsForPlaylist,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		}

		playlists[i] = playlist

		if (i+1)%10 == 0 {
			fmt.Printf("  ✓ Generated %d playlists\n", i+1)
		}
	}

	// Save all playlists
	if err := pg.ctx.AddPlaylists(playlists); err != nil {
		return fmt.Errorf("failed to save playlists: %w", err)
	}

	// Save playlist track relationships
	if err := pg.ctx.AddPlaylistTracks(playlistTracks); err != nil {
		return fmt.Errorf("failed to save playlist tracks: %w", err)
	}

	fmt.Printf("✓ Successfully generated and saved %d playlists\n", quantity)
	return nil
}
