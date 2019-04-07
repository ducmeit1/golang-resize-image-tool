package model

import (
	"bytes"
	"fmt"
	"github.com/ducmeit1/imaging"
	"image"
	"regexp"
	"strconv"
	"strings"
)

const MatchOptionalPattern = "^(?P<optional>(?P<dimension>[\\d]{1,4}x[\\d]{1,4})\\_?(?P<crop>top|bottom|center|left|right)*\\_)(?P<file>[a-zA-Z-_0-9]+\\.(?P<ext>jpg|png|gif|jpeg))$"

const (
	JPG  Extension = ".jpg"
	JPEG Extension = ".jpeg"
	PNG  Extension = ".png"
	GIF  Extension = ".gif"
)

const (
	Top    CropOption = "top"
	Bottom CropOption = "bottom"
	Center CropOption = "center"
	Left   CropOption = "left"
	Right  CropOption = "right"
)

type Image struct {
	Optional  string
	FileName  string
	Dimension string
	Extension Extension
	Crop      CropOption
	Height    int
	Width     int
}

type CropOption string
type Extension string

type DownloadOutput struct {
	Data *bytes.Buffer
}

func (s *Image) IsMatchFormat() bool {
	rgx := regexp.MustCompile(MatchOptionalPattern)
	if !rgx.MatchString(s.Optional) {
		return false
	}
	optionals := make(map[string]string)
	matches := rgx.FindStringSubmatch(s.Optional)
	for i, name := range rgx.SubexpNames() {
		if i != 0 && name != "" {
			optionals[name] = matches[i]
		}
	}
	defer func() {
		s.Extension = Extension(optionals["ext"])
		s.FileName = optionals["file"]
	}()
	if optionals["dimension"] == "" {
		return false
	}

	s.Dimension = optionals["dimension"]
	configs := strings.Split(optionals["dimension"], "x")
	height, err := strconv.ParseInt(configs[0], 10, 64)
	if err != nil {
		return false
	}
	width, err := strconv.ParseInt(configs[1], 10, 64)
	if err != nil {
		return false
	}
	s.Height = int(height)
	s.Width = int(width)

	if s.Height < 0 || s.Width < 0 {
		return false
	}

	if s.Height == 0 && s.Width == 0 {
		return false
	}

	if optionals["crop"] != "" && s.Height > 0 && s.Width > 0 {
		s.Crop = CropOption(optionals["crop"])
	}
	return true
}

func (s *Image) GetOutputFileName() string {
	names := make([]string, 0)
	names = append(names, s.Dimension)
	if s.Crop != "" {
		names = append(names, string(s.Crop))
	}
	names = append(names, s.FileName)
	return strings.Join(names, "_")
}

func ParseExtension(contentType string) (format imaging.Format) {
	switch contentType {
	case "image/png":
		format = imaging.PNG
		break
	case "image/gif":
		format = imaging.GIF
		break
	case "image/jpeg":
		format = imaging.JPEG
		break
	case "image/jpg":
		format = imaging.JPEG
		break
	}
	return
}

func ParseCropOption(option CropOption) (anchor imaging.Anchor) {
	switch option {
	case Top:
		anchor = imaging.Top
		break
	case Bottom:
		anchor = imaging.Bottom
		break
	case Left:
		anchor = imaging.Left
		break
	case Right:
		anchor = imaging.Right
		break
	case Center:
		anchor = imaging.Center
		break
	}
	return
}

func ParseContentType(ext Extension) (contentType string) {
	switch ext {
	case JPG:
		contentType = "image/jpeg"
		break
	case JPEG:
		contentType = "image/jpeg"
		break
	case PNG:
		contentType = "image/png"
		break
	case GIF:
		contentType = "image/gif"
		break
	}
	return
}

// ResizeOrCrop
func (s *Image) ResizeOrCrop(img image.Image) (residedImage *image.NRGBA) {
	if s.Crop == "" {
		residedImage = imaging.Fit(img, s.Width, s.Height, imaging.Lanczos)
	} else {
		residedImage = imaging.CropAnchor(img, s.Width, s.Height, ParseCropOption(s.Crop))
	}
	return
}

func (s *Image) GetS3Key(folder string, key string) string {
	return fmt.Sprintf("%s/%s", folder, key)
}
