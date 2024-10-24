package epub

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetRootDoc(t *testing.T) {

	orig, err := os.ReadFile("../test_data/MobyDickContent.opf")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Opening MobyDick.epub")
	fp := "../test_data/MobyDick.epub"

	epub, err := OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	b, err := epub.RootFile.WriteToString()
	if err != nil {
		t.Fatal(err)
	}

	if len(orig) != len(b) {
		t.Fatalf("File content length mismatch. Want: %d Have: %d", len(orig), len(b))
	}

	for i, a := range orig {
		if a != b[i] {
			t.Fatalf("File content mismatch at position %d", i)
		}
	}
}

func TestReadMetadata(t *testing.T) {

	// Test MobyDick
	fp := "../test_data/MobyDick.epub"
	epub, err := OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	mdataWant := &Metadata{
		Title:      "Moby Dick; Or, The Whale",
		TitleSort:  "",
		Author:     "Herman Melville",
		AuthorSort: "Melville, Herman",
		Language:   "en",
		Series:     "",
		SeriesNum:  -1,
		Subjects: []string{
			"Whaling -- Fiction",
			"Sea stories",
			"Psychological fiction",
			"Ship captains -- Fiction",
			"Adventure stories",
			"Mentally ill -- Fiction",
			"Ahab, Captain (Fictitious character) -- Fiction",
			"Whales -- Fiction",
			"Whaling ships -- Fiction",
		},
		Isbn:         "",
		Publisher:    "",
		PubDate:      "2001-07-01",
		Rights:       "Public domain in the USA.",
		Contributors: []Contributor{},
		Description:  "",
		Uid:          "http://www.gutenberg.org/2701",
	}

	if !assert.Equal(t, mdataWant, epub.Metadata) {
		t.Fatal()
	}

	// Test The Stone Age
	fp = "../test_data/TheStoneAgeInNorthAmericaVol2.epub"
	epub, err = OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	mdataWant = &Metadata{
		Title:        "The stone age in North America, vol. II",
		TitleSort:    "",
		Author:       "Warren K. Moorehead",
		AuthorSort:   "Moorehead, Warren K. (Warren King)",
		Language:     "en",
		Series:       "The Stone Age In North America",
		SeriesNum:    2,
		Subjects:     []string{},
		Isbn:         "",
		Publisher:    "",
		PubDate:      "2024-09-07",
		Rights:       "Public domain in the USA.",
		Contributors: []Contributor{},
		Description:  "",
		Uid:          "http://www.gutenberg.org/74390",
	}
	if !assert.Equal(t, mdataWant, epub.Metadata) {
		t.Fatal()
	}

	// Test The Brothers Karamazov
	fp = "../test_data/TheBrothersKaramazov.epub"
	epub, err = OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	mdataWant = &Metadata{
		Title:      "The Brothers Karamazov",
		TitleSort:  "Brothers Karamazov, The",
		Author:     "Fyodor Dostoyevsky",
		AuthorSort: "Dostoyevsky, Fyodor",
		Language:   "en",
		Series:     "",
		SeriesNum:  -1,
		Subjects: []string{
			"Didactic fiction",
			"Fathers and sons -- Fiction",
			"Russia -- Social life and customs -- 1533-1917 -- Fiction",
			"Brothers -- Fiction",
		},
		Isbn:      "0374528373",
		Publisher: "",
		PubDate:   "2009-02-12",
		Rights:    "Public domain in the USA.",
		Contributors: []Contributor{
			{Name: "Constance Garnett", Role: "trl"},
		},
		Description: "",
		Uid:         "http://www.gutenberg.org/28054",
	}

	if !assert.Equal(t, mdataWant, epub.Metadata) {
		t.Fatal()
	}

	// Test The Stones of Venice
	fp = "../test_data/TheStonesOfVeniceVol2.epub"
	epub, err = OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	mdataWant = &Metadata{
		Title:      "The Stones of Venice",
		TitleSort:  "Stones of Venice, The",
		Author:     "John Ruskin",
		AuthorSort: "Ruskin, John",
		Language:   "en",
		Series:     "The Stones of Venice",
		SeriesNum:  2,
		Subjects: []string{
			"Architecture -- Italy -- Venice",
		},
		Isbn:      "",
		Publisher: "",
		PubDate:   "2009-12-31",
		Rights:    "Public domain in the USA.",
		Contributors: []Contributor{
			{Name: "calibre (7.12.0) [https://calibre-ebook.com]", Role: "bkp"},
		},
		Description: "",
		Uid:         "http://www.gutenberg.org/30755",
	}

	if !assert.Equal(t, mdataWant, epub.Metadata) {
		t.Fatal()
	}
}

func TestWriteMetadata(t *testing.T) {
	tmpFp := "../test_data/TestEpub.epub"

	og, err := os.ReadFile("../test_data/TheBrothersKaramazov.epub")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFp, og, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFp)

	epub, err := OpenEpub(tmpFp)
	if err != nil {
		t.Fatal(err)
	}

	newMetadata := &Metadata{
		Title:      "newTitle",
		TitleSort:  "titleNew",
		Author:     "newAuthor",
		AuthorSort: "authorNew",
		Language:   "klingon",
		Series:     "newSeries",
		SeriesNum:  42,
		Subjects: []string{
			"subject1",
			"subject2",
		},
		Isbn:      "8675309",
		Publisher: "newPub",
		PubDate:   "1999-12-31",
		Rights:    "",
		Contributors: []Contributor{
			{Name: "Bob Ross", Role: "art"},
		},
		Description: "fakeMetadata",
		Uid:         "notauid",
	}

	epub.Metadata = newMetadata

	err = epub.WriteChanges()
	if err != nil {
		t.Fatal(err)
	}

	epub.Close()

	epub, err = OpenEpub(tmpFp)
	if err != nil {
		t.Fatal(err)
	}
	defer epub.Close()

	if !assert.Equal(t, newMetadata, epub.Metadata) {
		t.Fatal()
	}
}

func TestWriteCoverImage(t *testing.T) {
	tmpFp := "../test_data/TestEpub.epub"

	og, err := os.ReadFile("../test_data/TheBrothersKaramazov.epub")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFp, og, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFp)

	epub, err := OpenEpub(tmpFp)
	if err != nil {
		t.Fatal(err)
	}

	newImage, err := os.ReadFile("../test_data/miniCoverImg.png")
	if err != nil {
		t.Fatal(err)
	}

	err = epub.SetCoverImage(newImage)
	if err != nil {
		t.Fatal(err)
	}

	err = epub.WriteChanges()
	if err != nil {
		t.Fatal(err)
	}

	epub.Close()

	epub, err = OpenEpub(tmpFp)
	if err != nil {
		t.Fatal(err)
	}
	defer epub.Close()

	coverPath, err := epub.GetCoverPath()
	if err != nil {
		t.Fatal(err)
	}

	cv, err := epub.ReadFile(coverPath)
	if err != nil {
		t.Fatal(err)
	}

	// When go converts the test cover image it gets re-encoded and compressed
	// it is of no use to match file contents or size. Since nobody in their
	// right mind would use a 2x2 pixel image as a book cover it is probably
	// safe to assume that an image under 1kb means the write was successful
	if len(cv) > 1000 {
		t.Fatal()
	}
}

func TestImport(t *testing.T) {
	libRoot := filepath.Join("../test_data", "test_library_root")
	viper.SetDefault("library_path", libRoot)

	file, err := os.Open("../test_data/MobyDick.epub")
	if err != nil {
		t.Fatal(err)
	}

	e, err := Import(file)
	if err != nil {
		t.Fatal(err)
	}

	fp := e.FilePath
	if fp != filepath.Join(libRoot, "Herman Melville", "Moby Dick; Or, The Whale.epub") {
		t.Fatalf("Imported book filepath is incorrect. Want: %s Have: %s", filepath.Join(libRoot, "Herman Melville", "Moby Dick.epub"), fp)
	}

	if _, err = os.Stat(fp); err != nil {
		t.Fatal(err)
	}

	e.Close()
	os.RemoveAll(libRoot)
}

// Tests importing a file that already exists in the filesystem
func TestImportDuplicate(t *testing.T) {
	libRoot := filepath.Join("../test_data", "test_library_root")
	viper.SetDefault("library_path", libRoot)

	var e *Epub

	t.Cleanup(func() {
		if e != nil {
			e.Close()
		}
		os.RemoveAll(libRoot)
	})

	// import once
	file, err := os.Open("../test_data/MobyDick.epub")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Import(file)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	// import twice
	file, err = os.Open("../test_data/MobyDick.epub")
	if err != nil {
		t.Fatal(err)
	}

	e, err = Import(file)
	if err != nil {
		t.Fatal(err)
	}

	fp := e.FilePath
	if fp != filepath.Join(libRoot, "Herman Melville", "Moby Dick; Or, The Whale_1.epub") {
		t.Fatalf("Imported book filepath is incorrect. Want: %s Have: %s", filepath.Join(libRoot, "Herman Melville", "Moby Dick.epub"), fp)
	}

	if _, err = os.Stat(fp); err != nil {
		t.Fatal(err)
	}
}
