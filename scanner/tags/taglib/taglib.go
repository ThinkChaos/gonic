package taglib

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sentriz/audiotags"
	"go.senan.xyz/gonic/scanner/tags/tagcommon"
)

type TagLib struct{}

func (TagLib) CanRead(absPath string) bool {
	switch ext := filepath.Ext(absPath); ext {
	case ".mp3", ".flac", ".aac", ".m4a", ".m4b", ".ogg", ".opus", ".wma", ".wav", ".wv":
		return true
	}
	return false
}

func (TagLib) Read(absPath string) (tagcommon.Info, error) {
	raw, props, err := audiotags.Read(absPath)
	return &info{raw, props}, err
}

type info struct {
	raw   map[string][]string
	props *audiotags.AudioProperties
}

// https://picard-docs.musicbrainz.org/downloads/MusicBrainz_Picard_Tag_Map.html

func (i *info) Title() string          { return first(find(i.raw, "title")) }
func (i *info) BrainzID() string       { return first(find(i.raw, "musicbrainz_trackid")) } // musicbrainz recording ID
func (i *info) Artist() string         { return first(find(i.raw, "artist")) }
func (i *info) Album() string          { return first(find(i.raw, "album")) }
func (i *info) AlbumArtist() string    { return first(find(i.raw, "albumartist", "album artist")) }
func (i *info) AlbumArtists() []string { return find(i.raw, "albumartists", "album_artists") }
func (i *info) AlbumBrainzID() string  { return first(find(i.raw, "musicbrainz_albumid")) } // musicbrainz release ID
func (i *info) Genre() string          { return first(find(i.raw, "genre")) }
func (i *info) Genres() []string       { return find(i.raw, "genres") }
func (i *info) TrackNumber() int       { return intSep("/", first(find(i.raw, "tracknumber"))) }                  // eg. 5/12
func (i *info) DiscNumber() int        { return intSep("/", first(find(i.raw, "discnumber"))) }                   // eg. 1/2
func (i *info) Year() int              { return intSep("-", first(find(i.raw, "originaldate", "date", "year"))) } // eg. 2023-12-01
func (i *info) Length() int            { return i.props.Length }
func (i *info) Bitrate() int           { return i.props.Bitrate }

func first[T comparable](is []T) T {
	var z T
	for _, i := range is {
		if i != z {
			return i
		}
	}
	return z
}

func find(m map[string][]string, keys ...string) []string {
	for _, k := range keys {
		if r := filterStr(m[k]); len(r) > 0 {
			return r
		}
	}
	return nil
}

func filterStr(ss []string) []string {
	var r []string
	for _, s := range ss {
		if strings.TrimSpace(s) != "" {
			r = append(r, s)
		}
	}
	return r
}

func intSep(sep, in string) int {
	start, _, _ := strings.Cut(in, sep)
	out, _ := strconv.Atoi(start)
	return out
}
