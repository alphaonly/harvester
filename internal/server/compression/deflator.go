package compression

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Deflator struct {
	Level           int    //flate.BestCompression
	ContentEncoding string //deflate

}

func (d Deflator) CompressionHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//validation for compression
		if strings.Contains(r.Header.Get("Content-Encoding"), "deflate") {
			//deflate compressed response

			//read body
			byteData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//Compression logic
			d.ContentEncoding = r.Header.Get("Content-Encoding")
			compressedByteData, err := d.Compress(byteData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//compressed response to next handler
			if next != nil {
				r.Body = io.NopCloser(bytes.NewReader(compressedByteData))
				next.ServeHTTP(w, r)
				return
			}
			//compressed response to requester
			w.Write(compressedByteData)
			if err != nil {
				log.Println("compressed data writing error")
				http.Error(w, "response writing error", http.StatusInternalServerError)
				return
			}
		} else {
			//uncompressed response to next handler
			if next != nil {
				// r.Body = io.NopCloser(bytes.NewReader(byteData))
				next.ServeHTTP(w, r)
				return
			}
		}
		//uncompressed response

	})
}

func (d Deflator) WriteResponseBodyHandler(next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}
		w.Write(byteData)
		if err != nil {
			log.Println("byteData writing error")
			http.Error(w, "byteData writing error", http.StatusInternalServerError)
			return
		}
	})
}
func (d Deflator) DeCompressionHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//validation for decompression
		if strings.Contains(r.Header.Get("Content-Encoding"), "deflate") {

			//always read body
			byteData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			//Decompression logic
			d.ContentEncoding = r.Header.Get("Content-Encoding")
			decompressedByteData, err := d.Decompress(byteData)
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
			w.Write(decompressedByteData)
			if err != nil {
				log.Println("decompressed writing error")
				http.Error(w, "decompressed writing error", http.StatusInternalServerError)
				return
			}
		} else {

			//compressed response
			//initial response to next handler
			if next != nil {

				// r.Body = io.NopCloser(bytes.NewReader(byteData))
				next.ServeHTTP(w, r)
				return
			}
		}

	})
}

func (d Deflator) Compress(data []byte) ([]byte, error) {
	if d.ContentEncoding != "deflate" {
		return data, nil
	}

	var b bytes.Buffer
	w, err := flate.NewWriter(&b, d.Level)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

func (d Deflator) Decompress(data []byte) ([]byte, error) {
	if d.ContentEncoding != "deflate" {
		return data, nil
	}
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()

	var b bytes.Buffer

	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}

// type gzipWriter struct {
// 	http.ResponseWriter
// 	Writer io.Writer
// }

// func (w gzipWriter) Write(b []byte) (int, error) {
// 	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
// 	return w.Writer.Write(b)
// }

// func gzipandle(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
// 		if err != nil {
// 			io.WriteString(w, err.Error())
// 			return
// 		}
// 		defer gz.Close()

// 		w.Header().Set("Content-Encoding", "gzip")

// 		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
// 	})
// }
