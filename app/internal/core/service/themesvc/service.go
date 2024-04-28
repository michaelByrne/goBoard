package themesvc

import (
	"bytes"
	"context"
	php "github.com/leeqvip/gophp"
	"go.uber.org/zap"
	"goBoard/internal/core/ports"
	"strconv"
	"strings"
	"text/template"
)

type ThemeService struct {
	themeRepo ports.ThemeRepo

	logger *zap.SugaredLogger
}

func NewThemeService(themeRepo ports.ThemeRepo, logger *zap.SugaredLogger) ThemeService {
	return ThemeService{
		themeRepo: themeRepo,
		logger:    logger,
	}
}

func (s ThemeService) GetTheme(ctx context.Context, name string) (string, error) {
	phpTheme, err := s.themeRepo.GetTheme(ctx, name)
	if err != nil {
		s.logger.Errorf("error getting theme: %v", err)
		return "", err
	}

	themeIface, err := php.Unserialize([]byte(phpTheme))
	if err != nil {
		s.logger.Errorf("error unserializing theme: %v", err)
		return "", err
	}

	theme, ok := themeIface.(map[string]interface{})
	if !ok {
		s.logger.Errorf("error casting theme to map[string]interface{}")
		return "", err
	}

	themeOut := make(map[string]string)
	for k, v := range theme {
		switch v := v.(type) {
		case string:
			themeOut[k] = v

			if strings.HasPrefix(v, "#") {
				themeOut[k+"_font"] = calcColor(v)
			}
		default:
			s.logger.Errorf("error casting theme value %v to string", v)
		}
	}

	temp, err := template.ParseFiles("../internal/core/service/themesvc/styles.tmpl")
	if err != nil {
		s.logger.Errorf("error parsing template: %v", err)
		return "", err
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, themeOut)
	if err != nil {
		s.logger.Errorf("error executing template: %v", err)
		return "", err
	}

	return buf.String(), nil
}

func calcColor(hex string) string {
	hex = hex[1:]
	rgb := strings.Split(wordWrap(hex, 2, ":"), ":")
	r, _ := strconv.ParseInt(rgb[0], 16, 0)
	g, _ := strconv.ParseInt(rgb[1], 16, 0)
	b, _ := strconv.ParseInt(rgb[2], 16, 0)

	minRGB := minInt(int(r), int(g), int(b))
	maxRGB := maxInt(int(r), int(g), int(b))
	lum := float64((minRGB + maxRGB) / 510.0)

	if lum < 0.45 {
		return "#ffffff"
	} else {
		return "#000000"
	}
}

func wordWrap(s string, width int, breakChar string) string {
	var sb strings.Builder
	for i, r := range s {
		if i > 0 && i%width == 0 {
			sb.WriteString(breakChar)
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func minInt(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

func maxInt(a, b, c int) int {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}
