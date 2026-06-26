// Package media defines the port for storing incident images. The service
// depends only on this abstraction; the concrete object storage adapter lives
// in the infrastructure layer.
package media

import "context"

// Image is an uploaded image together with the metadata needed to store it.
type Image struct {
	Data        []byte
	ContentType string
	Ext         string
}

// Storage persists an incident image under the given key and returns a public
// URL pointing at it.
type Storage interface {
	Upload(ctx context.Context, key string, img Image) (string, error)
}
