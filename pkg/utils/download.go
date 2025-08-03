package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func Download(client *http.Client, _url, path string, headers map[string]string) error {
	parsedUrl, err := url.Parse(_url)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		errorMsg := string(errorBody)
		return fmt.Errorf("failed request with status code %d and body: %s", response.StatusCode, errorMsg)
	}

	filename := path

	if strings.TrimSpace(filename) == "" {
		segments := strings.Split(parsedUrl.Path, "/")
		filename = segments[len(segments)-1]
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer outfile.Close()
	defer response.Body.Close()

	contentLength := response.Header.Get("content-length")

	if contentLength == "" {
		fmt.Println("no content length do io.Copy")
		_, err = io.Copy(outfile, response.Body)
		if err != nil {
			return err
		}
		return nil
	}

	contentLengthInt, err := strconv.Atoi(contentLength)
	if err != nil {
		return err
	}

	fmt.Println("content-length:", contentLengthInt)

	bytesRead := 0
	buf := make([]byte, CHUNK_SIZE)
	lastBytesRead := 0

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	lastTs := time.Now()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				elapsed := time.Since(lastTs).Seconds()

				deltaBytes := bytesRead - lastBytesRead
				speed := float64(deltaBytes) / 1024 / 1024 / elapsed

				remaining := contentLengthInt - bytesRead
				eta := float64(remaining) / float64(deltaBytes) * float64(elapsed)

				fmt.Printf("%s Speed: %.2f MB/s (eta: %.2f seconds)\n", path, speed, eta)

				lastBytesRead = bytesRead
				lastTs = time.Now()
			}
		}
	}()

	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			_, writeErr := outfile.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			bytesRead += n

		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("\n%s Done\n", path)
				ticker.Stop()
				done <- true
				break
			} else {
				return err
			}
		}

		fmt.Printf("%s: progress: %.2f\r", path, float32(bytesRead)/float32(contentLengthInt))
	}

	return nil
}
