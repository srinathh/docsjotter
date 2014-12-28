package utils

import (
	"os"
	"path/filepath"
	"strings"
)

var extfiletype map[string]string

const (
	HTML_FILE     = "html"
	IMAGE_FILE    = "image"
	ARCHIVE_FILE  = "zip"
	AUDIO_FILE    = "audio"
	CODE_FILE     = "text"
	EXCEL_FILE    = "vnd.ms-word"
	VIDEO_FILE    = "video"
	PPT_FILE      = "vnd.ms-powerpoint"
	PDF_FILE      = "pdf"
	TEXT_FILE     = "text"
	WORD_FILE     = "vnd.ms-word"
	MARKDOWN_FILE = "text"
	DIRECTORY     = "directory"
	DEFAULT       = "default"
)

func init() {

	extfiletype = map[string]string{
		".html": HTML_FILE,
		".htm":  HTML_FILE,
		".jpg":  IMAGE_FILE,
		".gif":  IMAGE_FILE,
		".bmp":  IMAGE_FILE,
		".png":  IMAGE_FILE,
		".tiff": IMAGE_FILE,
		".zip":  ARCHIVE_FILE,
		".tar":  ARCHIVE_FILE,
		".bz2":  ARCHIVE_FILE,
		".gz":   ARCHIVE_FILE,
		".7z":   ARCHIVE_FILE,
		".mp3":  AUDIO_FILE,
		".wma":  AUDIO_FILE,
		".wav":  AUDIO_FILE,
		".ra":   AUDIO_FILE,
		".aiff": AUDIO_FILE,
		".flac": AUDIO_FILE,
		".c":    CODE_FILE,
		".go":   CODE_FILE,
		".cpp":  CODE_FILE,
		".py":   CODE_FILE,
		".js":   CODE_FILE,
		".json": CODE_FILE,
		".pl":   CODE_FILE,
		".bat":  CODE_FILE,
		".xls":  EXCEL_FILE,
		".xlsx": EXCEL_FILE,
		".xlm":  EXCEL_FILE,
		".xlsm": EXCEL_FILE,
		".mov":  VIDEO_FILE,
		".mp4":  VIDEO_FILE,
		".mpg":  VIDEO_FILE,
		".mpeg": VIDEO_FILE,
		".wmv":  VIDEO_FILE,
		".flv":  VIDEO_FILE,
		".3gp":  VIDEO_FILE,
		".webm": VIDEO_FILE,
		".pdf":  PDF_FILE,
		".ppt":  PPT_FILE,
		".pptx": PPT_FILE,
		".txt":  TEXT_FILE,
		".csv":  EXCEL_FILE,
		".doc":  WORD_FILE,
		".docx": WORD_FILE,
		".dot":  WORD_FILE,
		".md":   MARKDOWN_FILE,
	}
}

func GetFileType(info os.FileInfo) string {

	if info.IsDir() {
		return DIRECTORY
	} else {
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if val, ok := extfiletype[ext]; ok {
			return val
		} else {
			return DEFAULT
		}
	}
}
