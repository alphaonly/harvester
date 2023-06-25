package common

import (
	"context"
	"errors"
	"net/http"

	"github.com/alphaonly/harvester/internal/common/logging"
	"github.com/alphaonly/harvester/internal/schema"
)

type RWDataComposite struct {
	R *http.Request
	W http.ResponseWriter
}

func NewRWDataComposite(r *http.Request, w http.ResponseWriter) *RWDataComposite {

	return &RWDataComposite{r, w}
}

// RunNextHandler - runs next handler if it exists
func RunNextHandler(rwData *RWDataComposite, handler http.Handler, byteBody []byte) error {
	if handler == nil {
		logging.LogPrintln(errors.New("no next handler, possible error"))
		return nil
	}
	//write handled body for further handle
	ctxWithValue := context.WithValue(rwData.R.Context(), schema.PKey1, schema.PreviousBytes(byteBody))
	//call further handler with context parameters
	handler.ServeHTTP(rwData.W, rwData.R.WithContext(ctxWithValue))

	return nil
}
