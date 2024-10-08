package epub

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRootDoc(t *testing.T) {

	orig, err := os.ReadFile("test_data/MobyDickContent.opf")
	assert.Nil(t, err)

	t.Log("Opening MobyDick.epub")
	fp := "test_data/MobyDick.epub"

	epub, err := OpenEpub(fp)
	assert.Nil(t, err)

	assert.Nil(t, err)

	b, err := epub.rootFile.WriteToString()
	assert.Nil(t, err)

	assert.Equal(t, len(orig), len(b), "File content length mismatch. Want: %d Have: %d", len(orig), len(b))

	for i, a := range orig {
		assert.Equal(t, a, b[i], "File content mismatch at position %d", i)
	}
}

func TestReadMetadata(t *testing.T) {

	// Test MobyDick
	fp := "test_data/MobyDick.epub"
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

	assert.Equal(t, mdataWant, epub.Metadata)

	// Test The Stone Age
	fp = "test_data/TheStoneAgeInNorthAmericaVol2.epub"
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
	assert.Equal(t, mdataWant, epub.Metadata)

	// Test The Brothers Karamazov
	fp = "test_data/TheBrothersKaramazov.epub"
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
			{name: "Constance Garnett", role: "trl"},
		},
		Description: "",
		Uid:         "http://www.gutenberg.org/28054",
	}

	assert.Equal(t, mdataWant, epub.Metadata)

	// Test The Stones of Venice
	fp = "test_data/TheStonesOfVeniceVol2.epub"
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
			{name: "calibre (7.12.0) [https://calibre-ebook.com]", role: "bkp"},
		},
		Description: "",
		Uid:         "http://www.gutenberg.org/30755",
	}

	assert.Equal(t, mdataWant, epub.Metadata)
}

func TestWriteMetadata(t *testing.T) {
	tmpFp := "test_data/TestEpub.epub"

	og, err := os.ReadFile("test_data/TheBrothersKaramazov.epub")
	assert.Nil(t, err)

	err = os.WriteFile(tmpFp, og, os.ModePerm)
	assert.Nil(t, err)
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
			{name: "Bob Ross", role: "art"},
		},
		Description: "fakeMetadata",
		Uid:         "notauid",
	}

	epub.Metadata = newMetadata

	err = epub.WriteChanges()
	assert.Nil(t, err)

	epub.Close()

	epub, err = OpenEpub(tmpFp)
	if err != nil {
		t.Fatal(err)
	}
	defer epub.Close()

	assert.Equal(t, newMetadata, epub.Metadata)
}

func TestWriteCoverImage(t *testing.T) {
	tmpFp := "test_data/TestEpub.epub"

	og, err := os.ReadFile("test_data/TheBrothersKaramazov.epub")
	assert.Nil(t, err)

	err = os.WriteFile(tmpFp, og, os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(tmpFp)

	epub, err := OpenEpub(tmpFp)
	if err != nil {
		t.Fatal(err)
	}

	newImage, err := os.ReadFile("test_data/coverImg.png")
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
