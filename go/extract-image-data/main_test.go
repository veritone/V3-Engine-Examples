package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

/*
TestExifDataFromJpeg verifies that we can accurately extract data from a JPG file
*/
func TestExifDataFromJpeg(t *testing.T) {
	file, err := os.Open("testdata/animal.jpg")
	if err != nil {
		t.Fatalf("Cannot open testdata/animal.jpg: %s", err.Error())
	}

	vtnStandard := getExifDataAsVtnStandard(file)
	//printJson(vtnStandard)

	if len(vtnStandard.Series) != 1 {
		t.Fatalf("Should have 1 series, but has %d", len(vtnStandard.Series))
	}
	if vtnStandard.Series[0].StartTimeMs != 0 {
		t.Fatalf("Start time should be 0, but is %d", vtnStandard.Series[0].StartTimeMs)
	}
	if vtnStandard.Series[0].StopTimeMs != 0 {
		t.Fatalf("Stop time should be 0, but is %d", vtnStandard.Series[0].StopTimeMs)
	}
	if vtnStandard.Series[0].Vendor.ExifError != "" {
		t.Fatalf("Exif error: %s", vtnStandard.Series[0].Vendor.ExifError)
	}

	// spot check a couple of values
	assertThatExifContains(t, vtnStandard.Series[0].Vendor.Exif, exif.DateTime, `"2008:07:31 10:38:11"`)
	assertThatExifContains(t, vtnStandard.Series[0].Vendor.Exif, exif.XResolution, `"72/1"`)
}

/*
TestExifDataFromTiff verifies we can accurately extract data from a TIFF file
*/
func TestExifDataFromTiff(t *testing.T) {
	file, err := os.Open("testdata/DudleyLeavittUtah.tiff")
	if err != nil {
		t.Fatalf("Cannot open testdata/DudleyLeavittUtah.tiff: %s", err.Error())
	}

	vtnStandard := getExifDataAsVtnStandard(file)
	//printJson(vtnStandard)

	if len(vtnStandard.Series) != 1 {
		t.Fatalf("Should have 1 series, but has %d", len(vtnStandard.Series))
	}
	if vtnStandard.Series[0].StartTimeMs != 0 {
		t.Fatalf("Start time should be 0, but is %d", vtnStandard.Series[0].StartTimeMs)
	}
	if vtnStandard.Series[0].StopTimeMs != 0 {
		t.Fatalf("Stop time should be 0, but is %d", vtnStandard.Series[0].StopTimeMs)
	}
	if vtnStandard.Series[0].Vendor.ExifError != "" {
		t.Fatalf("Exif error: %s", vtnStandard.Series[0].Vendor.ExifError)
	}

	// spot check a couple of values
	assertThatExifContains(t, vtnStandard.Series[0].Vendor.Exif, exif.Orientation, "1")
	assertThatExifContains(t, vtnStandard.Series[0].Vendor.Exif, exif.ImageWidth, "196")
	assertThatExifContains(t, vtnStandard.Series[0].Vendor.Exif, exif.ImageLength, "257")
	assertThatExifContains(t, vtnStandard.Series[0].Vendor.Exif, exif.DateTime, `"2009:09:26 01:11:52"`)
	assertThatExifDoesNotContain(t, vtnStandard.Series[0].Vendor.Exif, exif.GPSDateStamp)
}

/*
TestExifDataFromGif verifies that unsupported files (animated GIF) generates the expected error in the
vtn-standard file, but does not fail to process
 */
func TestExifDataFromGif(t *testing.T) {
	file, err := os.Open("testdata/vulture.gif")
	if err != nil {
		t.Fatalf("Cannot open testdata/vulture.gif: %s", err.Error())
	}

	vtnStandard := getExifDataAsVtnStandard(file)
	//printJson(vtnStandard)

	if len(vtnStandard.Series) != 1 {
		t.Fatalf("Should have 1 series, but has %d", len(vtnStandard.Series))
	}
	if vtnStandard.Series[0].StartTimeMs != 0 {
		t.Fatalf("Start time should be 0, but is %d", vtnStandard.Series[0].StartTimeMs)
	}
	if vtnStandard.Series[0].StopTimeMs != 0 {
		t.Fatalf("Stop time should be 0, but is %d", vtnStandard.Series[0].StopTimeMs)
	}
	if vtnStandard.Series[0].Vendor.ExifError != "exif: failed to find exif intro marker" {
		t.Fatalf("EXIF extraction should have failed, but did not")
	}
}

/*
TestProcessHandler tests the process handler by simulating passing in a request file like the Engine Toolkit would
*/
func TestProcessHandler(t *testing.T) {
	request := createTestProcessRequestForFile(t, "testdata/animal.jpg")
	response := httptest.NewRecorder()
	srv := newServer()
	srv.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status: %d: %s", response.Code, response.Body.String())
	}
	jsonMap := make(map[string](interface{}))
	if err := json.Unmarshal(response.Body.Bytes(), &jsonMap); err != nil {
		t.Fatalf("%s", err)
	}
	//printJson(jsonMap)

	// check the timestamps
	series := jsonMap["series"].([]interface{})
	series1 := series[0].(map[string]interface{})
	if series1["startTimeMs"].(float64) != 1000 {
		t.Error()
	}
	if series1["stopTimeMs"].(float64) != 2000 {
		t.Error()
	}

	// spot-check the exif values
	vendor := series1["vendor"].(map[string]interface{})
	exif := vendor["exif"].(map[string]interface{})
	if exif["DateTime"] != "2008:07:31 10:38:11" {
		t.Error()
	}
	if exif["PixelYDimension"].([]interface{})[0].(float64) != 68 {
		t.Error()
	}
	if exif["PixelXDimension"].([]interface{})[0].(float64) != 100 {
		t.Error()
	}
}

/* printJson is a helper function that just prints a pretty-printed version of the object
 * to the console
 */
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

func createTestProcessRequestForFile(t *testing.T, filePath string) *http.Request {
	var buf bytes.Buffer
	m := multipart.NewWriter(&buf)
	_ = m.WriteField("startOffsetMS", "1000")
	_ = m.WriteField("endOffsetMS", "2000")
	f, err := m.CreateFormFile("chunk", filepath.Base(filePath))
	if err != nil {
		t.Fatalf("%s", err)
	}
	src, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if _, err := io.Copy(f, src); err != nil {
		t.Fatalf("%s", err)
	}
	if err := m.Close(); err != nil {
		t.Fatalf("%s", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/process", &buf)
	request.Header.Set("Content-Type", m.FormDataContentType())

	return request
}

const expectedOutput = `{
	"series": [{
		"startTimeMs": 1000,
		"stopTimeMs": 2000,
		"vendor": {
			"exif": {
				"ApertureValue": ["368640/65536"],
				"ColorSpace": [1],
				"ComponentsConfiguration": "",
				"CustomRendered": [0],
				"DateTime": "2008:07:31 10:38:11",
				"DateTimeDigitized": "2008:05:30 15:56:01",
				"DateTimeOriginal": "2008:05:30 15:56:01",
				"ExifIFDPointer": [214],
				"ExifVersion": "0221",
				"ExposureBiasValue": ["0/1"],
				"ExposureMode": [1],
				"ExposureProgram": [1],
				"ExposureTime": ["1/160"],
				"FNumber": ["71/10"],
				"Flash": [9],
				"FlashpixVersion": "0100",
				"FocalLength": ["135/1"],
				"FocalPlaneResolutionUnit": [2],
				"FocalPlaneXResolution": ["3888000/876"],
				"FocalPlaneYResolution": ["2592000/583"],
				"GPSInfoIFDPointer": [978],
				"GPSVersionID": [2, 2, 0, 0],
				"ISOSpeedRatings": [100],
				"InteroperabilityIFDPointer": [948],
				"InteroperabilityIndex": "R98",
				"Make": "Canon",
				"MeteringMode": [5],
				"Model": "Canon EOS 40D",
				"Orientation": [1],
				"PixelXDimension": [100],
				"PixelYDimension": [68],
				"ResolutionUnit": [2],
				"SceneCaptureType": [0],
				"ShutterSpeedValue": ["483328/65536"],
				"Software": "GIMP 2.4.5",
				"SubSecTime": "00",
				"SubSecTimeDigitized": "00",
				"SubSecTimeOriginal": "00",
				"ThumbJPEGInterchangeFormat": [1090],
				"ThumbJPEGInterchangeFormatLength": [1378],
				"UserComment": "",
				"WhiteBalance": [0],
				"XResolution": ["72/1"],
				"YCbCrPositioning": [2],
				"YResolution": ["72/1"]
			}
		}
	}]
}`
