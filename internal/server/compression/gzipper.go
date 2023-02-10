package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func GZipCompressionHandler(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//validation for compression
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			log.Println("GZipCompressionHandler invoked")
			//gzip compressed response

			//read body
			byteData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//Compression logic
			compressedByteData, err := GzipCompress(byteData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//compressed response to next handler
			if next != nil {
				r.Body = io.NopCloser(bytes.NewReader(*compressedByteData))
				next.ServeHTTP(w, r)
				return
			}
			log.Fatal("Gzip comp handler requires next handler not nil")
		}
	}
}

func GZipWriteResponseBodyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			return
		}
		log.Println("GZipWriteResponseBodyHandler invoked")

		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}
		_, err = w.Write(byteData)
		if err != nil {
			log.Println("byteData writing error")
			http.Error(w, "byteData writing error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
	}
}
func GZipDeCompressionHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GZipDeCompressionHandler invoked")

		//validation for decompression
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {

			//read compressed body
			byteData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//Decompression logic

			decompressedByteData, err := GzipDecompress(byteData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//decompressed response to next handler
			if next != nil {
				r.Body = io.NopCloser(bytes.NewReader(decompressedByteData))
				next.ServeHTTP(w, r)
				return
			}
			//decompressed response to requester
			_, err = w.Write(decompressedByteData)
			if err != nil {
				log.Println("decompressed writing error")
				http.Error(w, "decompressed writing error", http.StatusInternalServerError)
				return
			}
		}

		//compressed response
		//initial response to next handler
		if next != nil {
			// r.Body = io.NopCloser(bytes.NewReader(byteData))
			next.ServeHTTP(w, r)
			return
		}
		log.Fatal("Gzip deComp handler requires next handler not nil")

	}
}

func GzipCompress(data []byte) (*[]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to GzipCompress temporary buffer: %v", err)
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed GzipCompress data: %v", err)
	}

	compressedBytes := b.Bytes()

	return &compressedBytes, nil
}

func GzipDecompress(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		log.Fatalf("Gzip decompress error:%v", err)
	}
	defer func() {
		err := gzipReader.Close()
		if err != nil {
			log.Fatal("Close error")
		}
	}()
	decompressedBytes, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed GzipDecompress data: %v", err)
	}

	return decompressedBytes, nil
}
