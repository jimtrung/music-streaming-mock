package generators

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/jimtrung/music-streaming-mock-data/generator/repository"
	"github.com/jimtrung/music-streaming-mock-data/model"
	"github.com/jimtrung/music-streaming-mock-data/util"
)

// GeneratorContext holds shared context for all generators
type GeneratorContext struct {
	Repository  *repository.DataRepository
	DatasetName string
	Mutex       sync.RWMutex

	// In-memory collections
	Users          []model.UserJSON
	Profiles       []model.ProfileJSON
	Artists        []model.ArtistJSON
	Tracks         []model.TrackJSON
	Playlists      []model.PlaylistJSON
	PlaylistTracks []model.PlaylistTrackJSON

	// Asset pools
	Usernames     []string
	ArtistNames   []string
	TrackTitles   []string
	PlaylistNames []string
	Passwords     []string
}

// NewGeneratorContext creates a new generator context
func NewGeneratorContext(repo *repository.DataRepository, datasetName string) (*GeneratorContext, error) {
	ctx := &GeneratorContext{
		Repository:     repo,
		DatasetName:    datasetName,
		Users:          []model.UserJSON{},
		Profiles:       []model.ProfileJSON{},
		Artists:        []model.ArtistJSON{},
		Tracks:         []model.TrackJSON{},
		Playlists:      []model.PlaylistJSON{},
		PlaylistTracks: []model.PlaylistTrackJSON{},
	}

	// Load asset pools
	if err := ctx.loadAssets(); err != nil {
		return nil, fmt.Errorf("failed to load assets: %w", err)
	}

	return ctx, nil
}

// loadAssets loads all available asset files
func (ctx *GeneratorContext) loadAssets() error {
	var err error

	ctx.Usernames = util.ReadFileTxt("./assets/usernames.txt")
	if err != nil && len(ctx.Usernames) == 0 {
		return fmt.Errorf("failed to load usernames: %w", err)
	}

	ctx.Passwords = util.ReadFileTxt("./assets/passwords.txt")
	if err != nil && len(ctx.Passwords) == 0 {
		return fmt.Errorf("failed to load passwords: %w", err)
	}

	ctx.ArtistNames = util.ReadFileTxt("./assets/artist_names.txt")
	// Artist names are optional
	if len(ctx.ArtistNames) == 0 {
		ctx.ArtistNames = ctx.generateDefaultArtistNames(100)
	}

	ctx.TrackTitles = util.ReadFileTxt("./assets/track_titles.txt")
	if err != nil && len(ctx.TrackTitles) == 0 {
		log.Printf("Warning: track_titles.txt not found, using defaults\n")
		ctx.TrackTitles = ctx.generateDefaultTrackTitles(200)
	}

	ctx.PlaylistNames = util.ReadFileTxt("./assets/playlist_names.txt")
	if err != nil && len(ctx.PlaylistNames) == 0 {
		log.Printf("Warning: playlist_names.txt not found, using defaults\n")
		ctx.PlaylistNames = ctx.generateDefaultPlaylistNames(100)
	}

	return nil
}

// generateDefaultArtistNames generates placeholder artist names
func (ctx *GeneratorContext) generateDefaultArtistNames(count int) []string {
	names := make([]string, count)
	for i := 0; i < count; i++ {
		names[i] = fmt.Sprintf("Artist_%d", i+1)
	}
	return names
}

// generateDefaultTrackTitles generates placeholder track titles
func (ctx *GeneratorContext) generateDefaultTrackTitles(count int) []string {
	titles := make([]string, count)
	for i := 0; i < count; i++ {
		titles[i] = fmt.Sprintf("Track_%d", i+1)
	}
	return titles
}

// generateDefaultPlaylistNames generates placeholder playlist names
func (ctx *GeneratorContext) generateDefaultPlaylistNames(count int) []string {
	names := make([]string, count)
	for i := 0; i < count; i++ {
		names[i] = fmt.Sprintf("Playlist_%d", i+1)
	}
	return names
}

// AddUser adds a user to the collection and saves
func (ctx *GeneratorContext) AddUser(user model.UserJSON) error {
	ctx.Mutex.Lock()
	ctx.Users = append(ctx.Users, user)
	ctx.Mutex.Unlock()

	// Save incrementally
	return ctx.Repository.SaveUsers(ctx.DatasetName, ctx.Users)
}

// AddUsers adds multiple users and saves
func (ctx *GeneratorContext) AddUsers(users []model.UserJSON) error {
	ctx.Mutex.Lock()
	ctx.Users = append(ctx.Users, users...)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveUsers(ctx.DatasetName, ctx.Users)
}

// AddProfile adds a profile and saves
func (ctx *GeneratorContext) AddProfile(profile model.ProfileJSON) error {
	ctx.Mutex.Lock()
	ctx.Profiles = append(ctx.Profiles, profile)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveProfiles(ctx.DatasetName, ctx.Profiles)
}

// AddProfiles adds multiple profiles and saves
func (ctx *GeneratorContext) AddProfiles(profiles []model.ProfileJSON) error {
	ctx.Mutex.Lock()
	ctx.Profiles = append(ctx.Profiles, profiles...)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveProfiles(ctx.DatasetName, ctx.Profiles)
}

// AddArtist adds an artist and saves
func (ctx *GeneratorContext) AddArtist(artist model.ArtistJSON) error {
	ctx.Mutex.Lock()
	ctx.Artists = append(ctx.Artists, artist)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveArtists(ctx.DatasetName, ctx.Artists)
}

// AddArtists adds multiple artists and saves
func (ctx *GeneratorContext) AddArtists(artists []model.ArtistJSON) error {
	ctx.Mutex.Lock()
	ctx.Artists = append(ctx.Artists, artists...)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveArtists(ctx.DatasetName, ctx.Artists)
}

// AddTrack adds a track and saves
func (ctx *GeneratorContext) AddTrack(track model.TrackJSON) error {
	ctx.Mutex.Lock()
	ctx.Tracks = append(ctx.Tracks, track)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveTracks(ctx.DatasetName, ctx.Tracks)
}

// AddTracks adds multiple tracks and saves
func (ctx *GeneratorContext) AddTracks(tracks []model.TrackJSON) error {
	ctx.Mutex.Lock()
	ctx.Tracks = append(ctx.Tracks, tracks...)
	ctx.Mutex.Unlock()

	return ctx.Repository.SaveTracks(ctx.DatasetName, ctx.Tracks)
}

// AddPlaylist adds a playlist and saves
func (ctx *GeneratorContext) AddPlaylist(playlist model.PlaylistJSON) error {
	ctx.Mutex.Lock()
	ctx.Playlists = append(ctx.Playlists, playlist)
	ctx.Mutex.Unlock()

	return ctx.Repository.SavePlaylists(ctx.DatasetName, ctx.Playlists)
}

// AddPlaylists adds multiple playlists and saves
func (ctx *GeneratorContext) AddPlaylists(playlists []model.PlaylistJSON) error {
	ctx.Mutex.Lock()
	ctx.Playlists = append(ctx.Playlists, playlists...)
	ctx.Mutex.Unlock()

	return ctx.Repository.SavePlaylists(ctx.DatasetName, ctx.Playlists)
}

// AddPlaylistTracks adds playlist-track relationships
func (ctx *GeneratorContext) AddPlaylistTracks(playlistTracks []model.PlaylistTrackJSON) error {
	ctx.Mutex.Lock()
	ctx.PlaylistTracks = append(ctx.PlaylistTracks, playlistTracks...)
	ctx.Mutex.Unlock()

	return nil
}

// GetRandomUsername returns a random username and removes it from the pool
func (ctx *GeneratorContext) GetRandomUsername() string {
	if len(ctx.Usernames) == 0 {
		return fmt.Sprintf("user_%d", time.Now().UnixNano())
	}

	j := rand.Intn(len(ctx.Usernames))
	username := ctx.Usernames[j]
	ctx.Usernames = append(ctx.Usernames[:j], ctx.Usernames[j+1:]...)

	return username
}

// GetRandomPassword returns a random password
func (ctx *GeneratorContext) GetRandomPassword() string {
	if len(ctx.Passwords) == 0 {
		return "DefaultPass@123"
	}

	return ctx.Passwords[rand.Intn(len(ctx.Passwords))]
}

// GetRandomArtistName returns a random artist name and removes it
func (ctx *GeneratorContext) GetRandomArtistName() string {
	if len(ctx.ArtistNames) == 0 {
		return fmt.Sprintf("Artist_%d", time.Now().UnixNano())
	}

	j := rand.Intn(len(ctx.ArtistNames))
	name := ctx.ArtistNames[j]
	ctx.ArtistNames = append(ctx.ArtistNames[:j], ctx.ArtistNames[j+1:]...)

	return name
}

// GetRandomTrackTitle returns a random track title and removes it
func (ctx *GeneratorContext) GetRandomTrackTitle() string {
	if len(ctx.TrackTitles) == 0 {
		return fmt.Sprintf("Track_%d", time.Now().UnixNano())
	}

	j := rand.Intn(len(ctx.TrackTitles))
	title := ctx.TrackTitles[j]
	ctx.TrackTitles = append(ctx.TrackTitles[:j], ctx.TrackTitles[j+1:]...)

	return title
}

// GetRandomPlaylistName returns a random playlist name and removes it
func (ctx *GeneratorContext) GetRandomPlaylistName() string {
	if len(ctx.PlaylistNames) == 0 {
		return fmt.Sprintf("Playlist_%d", time.Now().UnixNano())
	}

	j := rand.Intn(len(ctx.PlaylistNames))
	name := ctx.PlaylistNames[j]
	ctx.PlaylistNames = append(ctx.PlaylistNames[:j], ctx.PlaylistNames[j+1:]...)

	return name
}

// GetStats returns statistics about generated data
func (ctx *GeneratorContext) GetStats() map[string]int {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	return map[string]int{
		"users":          len(ctx.Users),
		"profiles":       len(ctx.Profiles),
		"artists":        len(ctx.Artists),
		"tracks":         len(ctx.Tracks),
		"playlists":      len(ctx.Playlists),
		"playlistTracks": len(ctx.PlaylistTracks),
	}
}

// CompileToMockData creates a complete MockDataJSON from all collections
func (ctx *GeneratorContext) CompileToMockData() *model.MockDataJSON {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	return &model.MockDataJSON{
		Users:          ctx.Users,
		Profiles:       ctx.Profiles,
		Artists:        ctx.Artists,
		Tracks:         ctx.Tracks,
		Playlists:      ctx.Playlists,
		PlaylistTracks: ctx.PlaylistTracks,
		RefreshTokens:  []model.RefreshTokenJSON{},
		// Metadata would be set separately
	}
}
