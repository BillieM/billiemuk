package builder

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessImagesResizesLargeJPEG(t *testing.T) {
	root := t.TempDir()
	imagesDir := filepath.Join(root, "content", "images")
	distDir := filepath.Join(root, "dist")

	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a 2400x1600 JPEG
	img := image.NewRGBA(image.Rect(0, 0, 2400, 1600))
	f, err := os.Create(filepath.Join(imagesDir, "large.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if err := processImages(filepath.Join(root, "content"), distDir); err != nil {
		t.Fatal(err)
	}

	// Verify output exists
	outPath := filepath.Join(distDir, "images", "large.jpg")
	outF, err := os.Open(outPath)
	if err != nil {
		t.Fatal("output image not created")
	}
	defer outF.Close()

	outImg, _, err := image.DecodeConfig(outF)
	if err != nil {
		t.Fatal(err)
	}
	if outImg.Width > maxImageWidth {
		t.Errorf("output width = %d, want <= %d", outImg.Width, maxImageWidth)
	}
}

func TestProcessImagesKeepsSmallPNG(t *testing.T) {
	root := t.TempDir()
	imagesDir := filepath.Join(root, "content", "images")
	distDir := filepath.Join(root, "dist")

	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a small 400x300 PNG
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	f, err := os.Create(filepath.Join(imagesDir, "small.png"))
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if err := processImages(filepath.Join(root, "content"), distDir); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(distDir, "images", "small.png")
	outF, err := os.Open(outPath)
	if err != nil {
		t.Fatal("output image not created")
	}
	defer outF.Close()

	outImg, _, err := image.DecodeConfig(outF)
	if err != nil {
		t.Fatal(err)
	}
	if outImg.Width != 400 {
		t.Errorf("output width = %d, want 400 (should not resize)", outImg.Width)
	}
}
