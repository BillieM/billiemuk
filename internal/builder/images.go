package builder

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const maxImageWidth = 1200
const jpegQuality = 85

func processImages(contentDir, distDir string) error {
	imagesDir := filepath.Join(contentDir, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, _ := filepath.Rel(imagesDir, path)
		outPath := filepath.Join(distDir, "images", rel)

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg":
			return compressImage(path, outPath, "jpeg")
		case ".png":
			return compressImage(path, outPath, "png")
		default:
			// Copy non-image files as-is (e.g. SVG, GIF)
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return writeFile(outPath, string(data))
		}
	})
}

func compressImage(srcPath, dstPath, format string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decode image %s: %w", srcPath, err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Resize if wider than max
	if width > maxImageWidth {
		newHeight := int(float64(height) * float64(maxImageWidth) / float64(width))
		resized := image.NewRGBA(image.Rect(0, 0, maxImageWidth, newHeight))
		draw.CatmullRom.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
		img = resized
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: jpegQuality})
	case "png":
		return png.Encode(out, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
