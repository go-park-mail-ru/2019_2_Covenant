package file_processor

import (
	"context"
	"io"
)

type Repository interface {
	ProcessTrack(ctx context.Context, file io.Reader, id uint64) error
	ProcessAlbumPhoto(ctx context.Context, file io.Reader, id uint64) error
	ProcessArtistPhoto(ctx context.Context, file io.Reader, id uint64) error
	ProcessAvatar(ctx context.Context, file io.Reader, id uint64) error
}
