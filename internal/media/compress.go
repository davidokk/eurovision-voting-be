package media

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"

	"golang.org/x/image/draw"
)

const (
	MaxImageWidth   = 1280
	JPEGQuality     = 82
	MaxAvatarWidth  = 512
	MaxVoiceBytes   = 2 * 1024 * 1024
	MaxVideoBytes   = 8 * 1024 * 1024
	MaxImageBytes   = 5 * 1024 * 1024
)

func CompressImage(r io.Reader, maxWidth int) ([]byte, string, error) {
	if maxWidth <= 0 {
		maxWidth = MaxImageWidth
	}
	raw, err := io.ReadAll(io.LimitReader(r, MaxImageBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(raw) > MaxImageBytes {
		return nil, "", fmt.Errorf("image too large")
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w > maxWidth {
		h = h * maxWidth / w
		w = maxWidth
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Over, nil)
	var out bytes.Buffer
	if err := jpeg.Encode(&out, dst, &jpeg.Options{Quality: JPEGQuality}); err != nil {
		return nil, "", err
	}
	return out.Bytes(), "image/jpeg", nil
}

func ValidateVoice(r io.Reader) ([]byte, string, error) {
	return readLimited(r, MaxVoiceBytes, "audio/webm")
}

func ValidateVideoNote(r io.Reader) ([]byte, string, error) {
	return readLimited(r, MaxVideoBytes, "video/webm")
}

func readLimited(r io.Reader, max int, contentType string) ([]byte, string, error) {
	data, err := io.ReadAll(io.LimitReader(r, int64(max)+1))
	if err != nil {
		return nil, "", err
	}
	if len(data) > max {
		return nil, "", fmt.Errorf("file too large")
	}
	if len(data) == 0 {
		return nil, "", fmt.Errorf("empty file")
	}
	return data, contentType, nil
}
