package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var (
	readyStatus    int = 0
	verboseLogging     = true
)

func main() {
	if err := http.ListenAndServe("0.0.0.0:8080", newServer()); err != nil {
		fmt.Fprintf(os.Stderr, "exif: %s", err)
		os.Exit(1)
	}
}

func newServer() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("/ready", handleReady)
	s.HandleFunc("/process", handleProcess)
	return s
}

// handleReady tells the engine toolkit whether we're ready to process content yet. It can
// return http.StatusUnavailable (503) during initialization. Once ready to process content, it
// must return http.StatusOK (200). If initialization fails, it can return
// http.StatusInternalServerError (500) to terminate the engine toolkit. Normally, an engine
// this simple would just return http.StatusOK all the time, but for illustrative purposes we
// will initially return http.StatusUnavailable to simulate some start-up time for this engine
// (you can see the Ready webhook test change status when running in test mode - see README)
func handleReady(w http.ResponseWriter, r *http.Request) {
	if readyStatus == 0 {
		// initialize to unavailable (simulate some start-up operation)
		readyStatus = http.StatusServiceUnavailable

		// switch to OK after 5 seconds
		go func() {
			time.Sleep(5 * time.Second)
			readyStatus = http.StatusOK
		}()
	}

	w.WriteHeader(readyStatus)
}

// handleProcess runs the incoming chunk through an exif decoder and writes the results to the
// response as a JSON vtn-standard file.
func handleProcess(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(512 * 1024 * 1024)

	printRequest(request)

	// get (only) the request fields we are going to need for processing or output

	// get start offset, which we need to copy to the output
	startMs, err := getRequestInteger(request, "startOffsetMS")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// get the stop offset, which we need to copy to the output
	stopMs, err := getRequestInteger(request, "endOffsetMS")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// verify the media type of the input file
	mediaType, err := getRequestString(request, "chunkMimeType")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/tiff" {
		http.Error(responseWriter, fmt.Sprintf("Chunk has media type '%s'. Supported types are ['image/jpeg', 'image/tiff']", mediaType),
			http.StatusBadRequest)
		return
	}

	// get the URL of the file we need to process
	cacheURI, err := getRequestString(request, "cacheURI")
	if err != nil {
		http.Error(responseWriter,
			fmt.Sprintf("Field cacheURI could not be read so there is nothing to process: %v", err.Error()),
			http.StatusBadRequest)
		return
	}

	// get a reader for the input file
	cacheURIResponse, err := http.Get(cacheURI)
	if err != nil {
		http.Error(responseWriter,
			fmt.Sprintf("Cannot retrieve file '%s' to process: %v", cacheURI, err.Error()),
			http.StatusInternalServerError)
		return
	}
	defer cacheURIResponse.Body.Close()

	// get the exif data from the file in a vtn-standard structure
	log.Printf("Getting EXIF data for %s", cacheURI)
	vtnStandard := getExifDataAsVtnStandard(cacheURIResponse.Body)
	vtnStandard.Series[0].StartTimeMs = startMs
	vtnStandard.Series[0].StopTimeMs = stopMs

	// return the response in JSON format
	if err := json.NewEncoder(responseWriter).Encode(vtnStandard); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}

// printRequest will dump all the request information to the logs. This is for debugging or
// exploring the information that is sent to the engine. Only produces output if verboseLogging
// is set to true
func printRequest(request *http.Request) {
	if !verboseLogging {
		return
	}

	// generate report
	log.Println("  Header fields:")
	for key, values := range request.Header {
		if len(values) == 1 {
			log.Printf("    %s -> %v\n", key, values[0])
		} else {
			log.Printf("    %s -> %v\n", key, values)
		}
	}

	log.Println("  Form fields:")
	for key, values := range request.Form {
		if len(values) == 1 {
			log.Printf("    %s -> %v\n", key, values[0])
		} else {
			log.Printf("    %s -> %v\n", key, values)
		}
	}
}

// getRequestString is a convenience method for extracting a string from a request form
func getRequestString(r *http.Request, fieldName string) (string, error) {
	value := r.FormValue(fieldName)
	if value == "" {
		return "", fmt.Errorf("field '%s' could not be found in the request", fieldName)
	}

	return value, nil
}

// getRequestInteger is a convenience method for extracting an integer from a request form
func getRequestInteger(r *http.Request, fieldName string) (int, error) {
	value, err := getRequestString(r, fieldName)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("field %s ('%s') could not be converted to an integer: %v", fieldName, value, err.Error())
	}

	return intValue, nil
}

// getExifDataAsVtnStandard extracts the EXIF data from a file and returns it wrapped in a
// vtn-standard structure. Any errors extracting the EXIF data will be encoded into the
// structure but not returned explicitly. file is not closed by this function
func getExifDataAsVtnStandard(file io.Reader) vtnStandard {
	var vtnStandard vtnStandard
	vtnStandard.Series = make([]vtnSeries, 1)

	exifData, err := exif.Decode(file)
	vtnStandard.Series[0].Vendor.Exif = exifData
	if err != nil {
		vtnStandard.Series[0].Vendor.ExifError = err.Error()
	}
	return vtnStandard
}

// vtnStandard defines a skeleton for a response compatible with the vtn-standard (see
// https://docs.veritone.com/#/developer/engines/standards/engine-output/?id=engine-output-standard-vtn-standard)
// This does not define the entire vtn-standard, just barely enough to represent the exif data
// as a custom vendor structure
type vtnStandard struct {
	Series []vtnSeries `json:"series"`
}

type vtnSeries struct {
	StartTimeMs int       `json:"startTimeMs"`
	StopTimeMs  int       `json:"stopTimeMs"`
	Vendor      vtnVendor `json:"vendor"`
}

type vtnVendor struct {
	Exif      *exif.Exif `json:"exif,omitempty"`
	ExifError string     `json:"exifError,omitempty"`
}
