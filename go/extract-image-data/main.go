package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

/* handleReady tells the engine toolkit whether we're ready to process content yet.
 * It can return http.StatusUnavailable (503) during initialization. Once ready to process content, it must
 * return http.StatusOK (200). If initialization fails, it can return http.StatusInternalServerError (500) to
 * terminate the engine toolkit.
 * Normally, an engine this simple would just return http.StatusOK all the time, but for illustrative purposes
 * we will initially return http.StatusUnavailable to simulate some start-up time for this engine (you can
 * see the Ready webhook test change status when running in test mode - see README)
 */
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

/* handleProcess runs the incoming chunk through an exif decoder and writes the results to the response as JSON.
 */
func handleProcess(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(512 * 1024 * 1024)

	printRequest(request)

	// get (only) the request fields we are going to need for processing or output
	startMs, err := getRequestInteger(request, "startOffsetMS")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	stopMs, err := getRequestInteger(request, "endOffsetMS")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// get the chunk file
	chunkFile, fileHeader, err := request.FormFile("chunk")
	if err != nil {
		http.Error(responseWriter, "Unable to retrieve the chunk from the request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer chunkFile.Close()
	log.Printf("Getting EXIF data for file %s", fileHeader.Filename)

	// get the exif data from the file in a vtn-standard structure
	vtnStandard := getExifDataAsVtnStandard(chunkFile)
	vtnStandard.Series[0].StartTimeMs = startMs
	vtnStandard.Series[0].StopTimeMs = stopMs

	// return the response in JSON format
	if err := json.NewEncoder(responseWriter).Encode(vtnStandard); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}

/*
printRequest will dump all the request information to the logs. Only produces output if verboseLogging is
set to true
*/
func printRequest(request *http.Request) {
	if !verboseLogging {
		return
	}

	// get a file name from the form to label the report with
	aFileName := "a file"
outerLoop:
	for _, headers := range request.MultipartForm.File {
		for _, header := range headers {
			aFileName = "'" + header.Filename + "'"
			break outerLoop
		}
	}
	log.Println("Handling /process request for", aFileName)

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

	log.Println("  Form files:")
	for name, headers := range request.MultipartForm.File {
		for index, header := range headers {
			log.Printf("    %s %d -> %s (%d bytes)\n", name, index, header.Filename, header.Size)
		}
	}
}

/* getRequestInteger is a convenience method for extracting an integer from a request form */
func getRequestInteger(r *http.Request, fieldName string) (int, error) {
	value := r.FormValue(fieldName)
	if value == "" {
		return 0, errors.New(fmt.Sprintf("Field '%s' could not be found in the request", fieldName))
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(
			fmt.Sprintf("Field %s ('%s') could not be converted to an integer: %v",
				fieldName, value, err.Error()),
		)
	}

	return intValue, nil
}

/* getExifDataAsVtnStandard extracts the EXIF data from a file and returns it wrapped in a vtn-standard
 * structure. Any errors extracting the EXIF data will be encoded into the structure but not returned explicitly.
 * file is not closed by this function
 */
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

/* vtnStandard defines a skeleton for a response compatible with the vtn-standard (see
 * https://docs.veritone.com/#/developer/engines/standards/engine-output/?id=engine-output-standard-vtn-standard)
 * This does not define the entire vtn-standard, just barely enough to represent the exif data as a custom
 * vendor structure
 */
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
