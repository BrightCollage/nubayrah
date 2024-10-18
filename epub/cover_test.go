package epub

import (
	"os"
	"testing"
)

func TestExtractCover(t *testing.T) {

	origCover, err := os.ReadFile("test_data/origMobyDickCover.png")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Opening MobyDick.epub")
	fp := "test_data/MobyDick.epub"

	epub, err := OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	coverPath, err := epub.ExtractCoverImage("test_data")
	defer os.Remove(coverPath)
	if err != nil {
		t.Fatal(err)
	}

	extractedCover, err := os.ReadFile(coverPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(extractedCover) != len(origCover) {
		t.Fatalf("Mismatch cover image sizes. Want: %d Have: %d", len(origCover), len(extractedCover))
	}

	for i, b := range origCover {
		if b != extractedCover[i] {
			t.Fatalf("Mismatch cover image data as position %d", i)
		}
	}
}

func TestGetCoverFile(t *testing.T) {

	t.Log("Opening MobyDick.epub")
	fp := "test_data/MobyDick.epub"

	epub, err := OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	covFile, err := epub.GetCoverFile()
	if err != nil {
		t.Fatal(err)
	}

	if covFile.Name != "OEBPS/8143055649100492814_2701-cover.png" {
		t.Fatalf("Incorrect cover path. Want: `OEBPS\\8143055649100492814_2701-cover.png` Have: `%s`", covFile.Name)
	}
}
