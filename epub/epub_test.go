package epub

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const f64MIN = -179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.0000000000000000

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
		title:      "Moby Dick; Or, The Whale",
		titleSort:  "",
		author:     "Herman Melville",
		authorSort: "Melville, Herman",
		language:   "en",
		series:     "",
		seriesNum:  f64MIN,
		subjects: []string{
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
		isbn:         "",
		publisher:    "",
		pubDate:      "2001-07-01",
		rights:       "Public domain in the USA.",
		contributors: []Contributor{},
		description:  "",
		uid:          "http://www.gutenberg.org/2701",
	}

	// NaN != NaN, so check it first then sub in f64 min
	assert.NotEqual(t, math.NaN(), epub.Metadata.seriesNum)
	epub.Metadata.seriesNum = f64MIN

	assert.Equal(t, mdataWant, epub.Metadata)

	// Test The Stone Age
	fp = "test_data/TheStoneAgeInNorthAmericaVol2.epub"
	epub, err = OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}
	mdataWant = &Metadata{
		title:        "The stone age in North America, vol. II",
		titleSort:    "",
		author:       "Warren K. Moorehead",
		authorSort:   "Moorehead, Warren K. (Warren King)",
		language:     "en",
		series:       "The Stone Age In North America",
		seriesNum:    2,
		subjects:     []string{},
		isbn:         "",
		publisher:    "",
		pubDate:      "2024-09-07",
		rights:       "Public domain in the USA.",
		contributors: []Contributor{},
		description:  "",
		uid:          "http://www.gutenberg.org/74390",
	}
	assert.Equal(t, mdataWant, epub.Metadata)

	// Test The Brothers Karamazov
	fp = "test_data/TheBrothersKaramazov.epub"
	epub, err = OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	mdataWant = &Metadata{
		title:      "The Brothers Karamazov",
		titleSort:  "Brothers Karamazov, The",
		author:     "Fyodor Dostoyevsky",
		authorSort: "Dostoyevsky, Fyodor",
		language:   "en",
		series:     "",
		seriesNum:  f64MIN,
		subjects: []string{
			"Didactic fiction",
			"Fathers and sons -- Fiction",
			"Russia -- Social life and customs -- 1533-1917 -- Fiction",
			"Brothers -- Fiction",
		},
		isbn:      "0374528373",
		publisher: "",
		pubDate:   "2009-02-12",
		rights:    "Public domain in the USA.",
		contributors: []Contributor{
			{name: "Constance Garnett", role: "trl"},
		},
		description: "",
		uid:         "http://www.gutenberg.org/28054",
	}

	// NaN != NaN, so check it first then sub in f64 min
	assert.NotEqual(t, math.NaN(), epub.Metadata.seriesNum)
	epub.Metadata.seriesNum = f64MIN

	assert.Equal(t, mdataWant, epub.Metadata)

	// Test The Stones of Venice
	fp = "test_data/TheStonesOfVeniceVol2.epub"
	epub, err = OpenEpub(fp)
	if err != nil {
		t.Fatal(err)
	}

	mdataWant = &Metadata{
		title:      "The Stones of Venice",
		titleSort:  "Stones of Venice, The",
		author:     "John Ruskin",
		authorSort: "Ruskin, John",
		language:   "en",
		series:     "The Stones of Venice",
		seriesNum:  2,
		subjects: []string{
			"Architecture -- Italy -- Venice",
		},
		isbn:      "",
		publisher: "",
		pubDate:   "2009-12-31",
		rights:    "Public domain in the USA.",
		contributors: []Contributor{
			{name: "calibre (7.12.0) [https://calibre-ebook.com]", role: "bkp"},
		},
		description: "",
		uid:         "http://www.gutenberg.org/30755",
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
		title:      "newTitle",
		titleSort:  "titleNew",
		author:     "newAuthor",
		authorSort: "authorNew",
		language:   "klingon",
		series:     "newSeries",
		seriesNum:  42,
		subjects: []string{
			"subject1",
			"subject2",
		},
		isbn:      "8675309",
		publisher: "newPub",
		pubDate:   "1999-12-31",
		rights:    "",
		contributors: []Contributor{
			{name: "Bob Ross", role: "art"},
		},
		description: "fakeMetadata",
		uid:         "notauid",
		nubayrahId:  "test_book",
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
