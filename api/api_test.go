package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"nubayrah/api/book"
	"nubayrah/api/router"
	"nubayrah/epub"
	"nubayrah/sqlite"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Create config and dirs to run nubary out of "./testHome"
func makeTestConfig() error {

	// todo: gorm doesn't free up the database file until the process
	// terminates so there can be an existing ./testHome even after cleanup
	os.RemoveAll("./testHome")

	homeDir := filepath.Join("./testHome", ".nubayrah")
	err := os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		return err
	}

	libraryRoot := filepath.Join(homeDir, "library")
	err = os.MkdirAll(libraryRoot, os.ModePerm)
	if err != nil {
		return err
	}

	viper.SetDefault("library_path", libraryRoot)
	viper.SetDefault("config_path", filepath.Join(homeDir, "config.yaml"))
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 5050)
	viper.SetDefault("db_path", filepath.Join(libraryRoot, "nubayrah.db"))

	// tells Viper to look for `dataRoot/config.yaml``
	viper.AddConfigPath(homeDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return nil
}

func startTestServer(t *testing.T) (*gorm.DB, error) {
	var srv *http.Server
	var DB *gorm.DB

	t.Cleanup(func() {
		if srv != nil {
			srv.Shutdown(context.Background())
		}

		if DB != nil {
			db, err := DB.DB()
			if err != nil {
				db.Close()
			}
		}
		os.RemoveAll("./testHome")
	})

	DB, err := sqlite.OpenDatabase()
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

	srv = &http.Server{Addr: addr, Handler: router.New(DB)}
	go func() {
		srv.ListenAndServe()
	}()

	return DB, nil
}

func makePOSTBody(path string) (*bytes.Buffer, string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("epub", filepath.Base(path))
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func uploadFile(path string) (*http.Response, error) {
	body, ct, err := makePOSTBody(path)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("http://%s:%d/books", viper.GetString("host"), viper.GetInt("port"))
	req, err := http.NewRequest("POST", addr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", ct)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func TestStartServer(t *testing.T) {
	err := makeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	_, err = startTestServer(t)
	if err != nil {
		t.Fatal(err)
	}

	addr := fmt.Sprintf("http://%s:%d/books", viper.GetString("host"), viper.GetInt("port"))
	r, err := http.Get(addr)
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != 200 {
		t.Fatal(fmt.Errorf("Unexpected response status: %d", r.StatusCode))
	}
}

func TestImportBook(t *testing.T) {

	err := makeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	DB, err := startTestServer(t)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := uploadFile("../test_data/MobyDick.epub")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 201 {
		t.Fatal(fmt.Errorf("Unexpected status code %d", resp.StatusCode))
	}

	// Verify that files exist
	_, err = os.Stat(filepath.Join("./testHome", ".nubayrah", "library", "Herman Melville", "Moby Dick; Or, The Whale.epub"))
	if err != nil {
		t.Fatal(err)
	}

	books := make([]*book.Book, 0)
	tx := DB.Find(&books)
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}

	if len(books) != 1 {
		t.Fatal(fmt.Errorf("Expected 1 entry in database, got %d", len(books)))
	}

	if books[0].Author != "Herman Melville" {
		t.Fatal(fmt.Errorf("Incorrect Author field on import. Wanted: `Herman Melville`, Have: `%s`", books[0].Author))
	}
}

func TestGetBooks(t *testing.T) {
	err := makeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	_, err = startTestServer(t)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := uploadFile("../test_data/MobyDick.epub")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 201 {
		t.Fatal(fmt.Errorf("Unexpected status code %d", resp.StatusCode))
	}

	addr := fmt.Sprintf("http://%s:%d/books", viper.GetString("host"), viper.GetInt("port"))

	resp, err = http.Get(addr)
	if err != nil {
		t.Fatal(err)
	}

	// b, err := io.ReadAll(resp.Body)
	// fmt.Println(string(b))

	var mdata []epub.Metadata
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&mdata)
	if err != nil {
		t.Fatal(err)
	}

	if len(mdata) != 1 {
		t.Fatal(fmt.Errorf("Incorrect number of entries in json. Want: 1 Have: %d", len(mdata)))
	}

	if mdata[0].Author != "Herman Melville" {
		t.Fatal(fmt.Errorf("Metadata content mismatch. Want author: Herman Melville Have: %s", mdata[0].Author))

	}

}
