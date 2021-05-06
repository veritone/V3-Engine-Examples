package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
)

// TestExifDataFromJpeg verifies that we can accurately extract data from a JPG file
func TestExifDataFromJpeg(t *testing.T) {
	file, err := os.Open("testdata/animal.jpg")
	if err != nil {
		t.Fatalf("Cannot open testdata/animal.jpg: %s", err.Error())
	}

	vtnStandard := getExifDataAsVtnStandard(file)
	//printJson(vtnStandard)

	if len(vtnStandard.Object) != 1 {
		t.Fatalf("Should have 1 series, but has %d", len(vtnStandard.Object))
	}
	if vtnStandard.Object[0].Vendor.ExifError != "" {
		t.Fatalf("Exif error: %s", vtnStandard.Object[0].Vendor.ExifError)
	}

	// spot check a couple of values
	assertThatExifContains(t, vtnStandard.Object[0].Vendor.Exif, exif.DateTime, `"2008:07:31 10:38:11"`)
	assertThatExifContains(t, vtnStandard.Object[0].Vendor.Exif, exif.XResolution, `"72/1"`)
}

// TestExifDataFromTiff verifies we can accurately extract data from a TIFF file
func TestExifDataFromTiff(t *testing.T) {
	file, err := os.Open("testdata/DudleyLeavittUtah.tiff")
	if err != nil {
		t.Fatalf("Cannot open testdata/DudleyLeavittUtah.tiff: %s", err.Error())
	}

	vtnStandard := getExifDataAsVtnStandard(file)
	//printJson(vtnStandard)

	if len(vtnStandard.Object) != 1 {
		t.Fatalf("Should have 1 series, but has %d", len(vtnStandard.Object))
	}
	if vtnStandard.Object[0].Vendor.ExifError != "" {
		t.Fatalf("Exif error: %s", vtnStandard.Object[0].Vendor.ExifError)
	}

	// spot check a couple of values
	assertThatExifContains(t, vtnStandard.Object[0].Vendor.Exif, exif.Orientation, "1")
	assertThatExifContains(t, vtnStandard.Object[0].Vendor.Exif, exif.ImageWidth, "196")
	assertThatExifContains(t, vtnStandard.Object[0].Vendor.Exif, exif.ImageLength, "257")
	assertThatExifContains(t, vtnStandard.Object[0].Vendor.Exif, exif.DateTime, `"2009:09:26 01:11:52"`)
	assertThatExifDoesNotContain(t, vtnStandard.Object[0].Vendor.Exif, exif.GPSDateStamp)
}

// TestExifDataFromGif verifies that unsupported files (animated GIF) generates the expected
// error in the vtn-standard file, but does not fail to process
func TestExifDataFromGif(t *testing.T) {
	file, err := os.Open("testdata/vulture.gif")
	if err != nil {
		t.Fatalf("Cannot open testdata/vulture.gif: %s", err.Error())
	}

	vtnStandard := getExifDataAsVtnStandard(file)
	//printJson(vtnStandard)

	if len(vtnStandard.Object) != 1 {
		t.Fatalf("Should have 1 series, but has %d", len(vtnStandard.Object))
	}
	if vtnStandard.Object[0].Vendor.ExifError != "exif: failed to find exif intro marker" {
		t.Fatalf("EXIF extraction should have failed, but did not")
	}
}

// TestProcessHandler tests the process handler by simulating passing in a request file like the
// Engine Toolkit would
func TestProcessHandlerOnJpeg(t *testing.T) {
	// send a request to our processor
	imageUri := "https://github.com/veritone/V3-Engine-Examples/blob/master/go/extract-image-data/testdata/animal.jpg?raw=true"
	request := createTestProcessRequestForFile(t, imageUri, "image/jpeg")
	response := httptest.NewRecorder()
	srv := newServer()
	srv.ServeHTTP(response, request)

	// extract results
	if response.Code != http.StatusOK {
		t.Fatalf("got status: %d: %s", response.Code, response.Body.String())
	}
	jsonMap := make(map[string](interface{}))
	if err := json.Unmarshal(response.Body.Bytes(), &jsonMap); err != nil {
		t.Fatalf("%s", err)
	}
	// printJson(jsonMap)

	// check the timestamps
	obj := jsonMap["object"].([]interface{})
	obj1 := obj[0].(map[string]interface{})

	// spot-check the exif values
	vendor := obj1["vendor"].(map[string]interface{})
	exif := vendor["exif"].(map[string]interface{})
	if exif["DateTime"] != "2008:07:31 10:38:11" {
		t.Errorf("DateTime: expected '2008:07:31 10:38:11' but got '%v'", exif["DateTime"])
	}
	if exif["PixelYDimension"].([]interface{})[0].(float64) != 68 {
		t.Errorf("PixelYDimension: expected '68' but got '%v'", exif["PixelYDimension"].([]interface{})[0])
	}
	if exif["PixelXDimension"].([]interface{})[0].(float64) != 100 {
		t.Errorf("PixelXDimension: expected '100' but got '%v'", exif["PixelXDimension"].([]interface{})[0])
	}
}

// TestProcessHandler tests the process handler by simulating passing in a request for an
// unsupported file type
func TestProcessHandlerOnUnsupportedGif(t *testing.T) {
	// send a request to our processor
	imageUri := "https://github.com/veritone/V3-Engine-Examples/blob/master/go/extract-image-data/testdata/vulture.gif?raw=true"
	request := createTestProcessRequestForFile(t, imageUri, "image/gif")
	response := httptest.NewRecorder()
	srv := newServer()
	srv.ServeHTTP(response, request)

	// extract results
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status '400', but got '%v'", response.Code)
	}
	if !strings.Contains(response.Body.String(), "media type 'image/gif'") {
		t.Fatalf("expected error to be that the media type was incorrect, but was '%v'", response.Body.String())
	}
}

func createTestProcessRequestForFile(t *testing.T, fileUrl string, mediaType string) *http.Request {
	formData := url.Values{}
	formData.Set("startOffsetMS", "1000")
	formData.Set("endOffsetMS", "2000")
	formData.Set("chunkMimeType", mediaType)
	formData.Set("cacheURI", fileUrl)
	// formData.Set("payload", "{\"verbose\":\"false\"}")

	request, err := http.NewRequest(http.MethodPost, "/process", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("%s", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(formData.Encode())))

	return request
}

// printJson is a helper function that just prints a pretty-printed version of the object to the
// console
func printJson(obj interface{}) {
	jsonString, _ := json.MarshalIndent(obj, "", "  ")
	_, _ = os.Stdout.Write(jsonString)
	fmt.Println()
}

/* assertThatExifContains is a convenience function to assert that an exif structure contains
 * an expected value
 */
func assertThatExifContains(t *testing.T, exif *exif.Exif, fieldName exif.FieldName, expectedValue string) {
	tag, err := exif.Get(fieldName)
	if err != nil {
		t.Fatalf("Field '%s' has no value in the EXIF data, was expecting '%s'",
			fieldName, expectedValue)
		return
	}
	if tag.String() != expectedValue {
		t.Fatalf("Expected EXIF to contain '%s' for field %s, but it contained %s",
			expectedValue, fieldName, tag.String())
	}
}

/* assertThatExifDoesNotContain verifies that the EXIF structure does not contain a value for the field
 */
func assertThatExifDoesNotContain(t *testing.T, exif *exif.Exif, fieldName exif.FieldName) {
	tag, err := exif.Get(fieldName)
	if err == nil {
		t.Fatalf("Field '%s' should not exist, but it does (with value '%s')",
			fieldName, tag.String())
	}
}
