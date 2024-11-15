package chain

import (
	"context"
	"encoding/json"
	"net/http"

	"pkg/errors"
)

func DefaultResponseEncoder(ctx context.Context, w http.ResponseWriter, response any) error {
	_, span := tracer.Start(ctx, "DefaultResponseEncoder")
	defer span.End()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errors.InternalServer.Wrap(err)
	}
	return nil
}
