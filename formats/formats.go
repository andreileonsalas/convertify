package formats

import (
	"path/filepath"
	"strings"
)

// Group represents the type of a media file
type Group int

const (
	Unknown Group = iota
	Photo
	Video
	Audio
)

func (g Group) String() string {
	switch g {
	case Photo:
		return "Imagen"
	case Video:
		return "Video"
	case Audio:
		return "Audio"
	}
	return "Desconocido"
}

// OutputFormat describes a conversion target
type OutputFormat struct {
	Ext       string   // output extension, e.g. "mp4"
	Label     string   // display name, e.g. "MP4"
	Args      []string // extra ffmpeg args
	IsExtract bool     // true = extract audio from video
}

// ---- Input extension maps ----

var photoExts = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "webp": true,
	"gif": true, "bmp": true, "tiff": true, "tif": true,
	"avif": true, "heic": true, "ico": true,
}

var videoExts = map[string]bool{
	"mp4": true, "avi": true, "mkv": true, "mov": true,
	"webm": true, "flv": true, "wmv": true, "m4v": true,
	"3gp": true, "ts": true, "mpeg": true, "mpg": true,
	"vob": true, "ogv": true,
}

var audioExts = map[string]bool{
	"mp3": true, "wav": true, "flac": true, "aac": true,
	"ogg": true, "m4a": true, "wma": true, "opus": true,
	"aiff": true, "aif": true, "ape": true,
}

// DetectGroup returns the Group for a file path based on its extension
func DetectGroup(path string) Group {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if photoExts[ext] {
		return Photo
	}
	if videoExts[ext] {
		return Video
	}
	if audioExts[ext] {
		return Audio
	}
	return Unknown
}

// AllExtensions returns all supported input extensions (with dot)
func AllExtensions() []string {
	all := []string{}
	for ext := range photoExts {
		all = append(all, "."+ext)
	}
	for ext := range videoExts {
		all = append(all, "."+ext)
	}
	for ext := range audioExts {
		all = append(all, "."+ext)
	}
	return all
}

func PhotoExtensions() []string {
	r := []string{}
	for ext := range photoExts {
		r = append(r, ext)
	}
	return r
}

func VideoExtensions() []string {
	r := []string{}
	for ext := range videoExts {
		r = append(r, ext)
	}
	return r
}

func AudioExtensions() []string {
	r := []string{}
	for ext := range audioExts {
		r = append(r, ext)
	}
	return r
}

// ---- Output format lists ----

// PhotoOutputs - conversion targets when input is an image
var PhotoOutputs = []OutputFormat{
	{Ext: "jpg",  Label: "JPG",  Args: []string{"-q:v", "2"}},
	{Ext: "png",  Label: "PNG",  Args: []string{}},
	{Ext: "webp", Label: "WebP", Args: []string{"-quality", "80"}},
	{Ext: "bmp",  Label: "BMP",  Args: []string{}},
	{Ext: "tiff", Label: "TIFF", Args: []string{}},
	{Ext: "gif",  Label: "GIF",  Args: []string{}},
}

// VideoOutputs - video-to-video conversions
var VideoOutputs = []OutputFormat{
	{
		Ext:   "mp4",
		Label: "MP4",
		Args:  []string{"-c:v", "libx264", "-crf", "23", "-preset", "medium", "-c:a", "aac", "-b:a", "192k", "-movflags", "+faststart"},
	},
	{
		Ext:   "avi",
		Label: "AVI",
		Args:  []string{"-c:v", "libxvid", "-q:v", "5", "-c:a", "libmp3lame", "-q:a", "3"},
	},
	{
		Ext:   "mkv",
		Label: "MKV",
		Args:  []string{"-c:v", "libx264", "-crf", "23", "-preset", "medium", "-c:a", "aac"},
	},
	{
		Ext:   "mov",
		Label: "MOV",
		Args:  []string{"-c:v", "libx264", "-crf", "23", "-c:a", "aac"},
	},
	{
		Ext:   "webm",
		Label: "WebM",
		Args:  []string{"-c:v", "libvpx-vp9", "-crf", "30", "-b:v", "0", "-c:a", "libopus"},
	},
	{
		Ext:   "gif",
		Label: "GIF (animado)",
		Args:  []string{"-vf", "fps=15,scale=480:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse", "-loop", "0"},
	},
}

// VideoExtractOutputs - audio extraction from video
var VideoExtractOutputs = []OutputFormat{
	{
		Ext:       "mp3",
		Label:     "Extraer audio → MP3",
		IsExtract: true,
		Args:      []string{"-vn", "-c:a", "libmp3lame", "-q:a", "2"},
	},
	{
		Ext:       "aac",
		Label:     "Extraer audio → AAC",
		IsExtract: true,
		Args:      []string{"-vn", "-c:a", "aac", "-b:a", "192k"},
	},
	{
		Ext:       "wav",
		Label:     "Extraer audio → WAV",
		IsExtract: true,
		Args:      []string{"-vn", "-c:a", "pcm_s16le"},
	},
	{
		Ext:       "flac",
		Label:     "Extraer audio → FLAC",
		IsExtract: true,
		Args:      []string{"-vn", "-c:a", "flac"},
	},
	{
		Ext:       "ogg",
		Label:     "Extraer audio → OGG",
		IsExtract: true,
		Args:      []string{"-vn", "-c:a", "libvorbis", "-q:a", "4"},
	},
}

// AudioOutputs - audio-to-audio conversions
var AudioOutputs = []OutputFormat{
	{Ext: "mp3",  Label: "MP3",  Args: []string{"-c:a", "libmp3lame", "-q:a", "2"}},
	{Ext: "wav",  Label: "WAV",  Args: []string{"-c:a", "pcm_s16le"}},
	{Ext: "flac", Label: "FLAC", Args: []string{"-c:a", "flac"}},
	{Ext: "aac",  Label: "AAC",  Args: []string{"-c:a", "aac", "-b:a", "192k"}},
	{Ext: "ogg",  Label: "OGG",  Args: []string{"-c:a", "libvorbis", "-q:a", "4"}},
	{Ext: "opus", Label: "OPUS", Args: []string{"-c:a", "libopus", "-b:a", "128k"}},
	{Ext: "m4a",  Label: "M4A",  Args: []string{"-c:a", "aac", "-b:a", "192k"}},
}

// OutputsFor returns the available output formats for a given input group
func OutputsFor(g Group) []OutputFormat {
	switch g {
	case Photo:
		return PhotoOutputs
	case Video:
		return append(VideoOutputs, VideoExtractOutputs...)
	case Audio:
		return AudioOutputs
	}
	return nil
}

// FindFormat searches for an OutputFormat by extension in the given group
func FindFormat(g Group, ext string) (OutputFormat, bool) {
	ext = strings.ToLower(ext)
	for _, f := range OutputsFor(g) {
		if f.Ext == ext {
			return f, true
		}
	}
	return OutputFormat{}, false
}
