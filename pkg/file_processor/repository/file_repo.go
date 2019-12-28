package repository

import (
	"2019_2_Covenant/pkg/file_processor"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
	"io"
)

type FileRepository struct {
	client file_processor.FilesClient
}

func NewFileRepository(client file_processor.FilesClient) *FileRepository {
	return &FileRepository{
		client: client,
	}
}

func (fr *FileRepository) ProcessTrack(ctx context.Context, file io.Reader, id uint64) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "id", fmt.Sprintf("%d", id))
	stream, err := fr.client.ProcessAudio(ctx)
	if err != nil {
		return err
	}

	return sendFile(file, stream)
}

func (fr *FileRepository) ProcessAlbumPhoto(ctx context.Context, file io.Reader, id uint64) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "type", "albumPhoto", "id", fmt.Sprintf("%d", id))
	stream, err := fr.client.ProcessImage(ctx)
	if err != nil {
		return err
	}

	return sendFile(file, stream)
}

func (fr *FileRepository) ProcessArtistPhoto(ctx context.Context, file io.Reader, id uint64) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "type", "artistPhoto", "id", fmt.Sprintf("%d", id))
	stream, err := fr.client.ProcessImage(ctx)
	if err != nil {
		return err
	}

	return sendFile(file, stream)
}

func (fr *FileRepository) ProcessAvatar(ctx context.Context, file io.Reader, id uint64) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "type", "avatar", "id", fmt.Sprintf("%d", id))
	stream, err := fr.client.ProcessImage(ctx)
	if err != nil {
		return err
	}

	return sendFile(file, stream)
}

type chunkedStream interface {
	Send(*file_processor.Chunk) error
	CloseAndRecv() (*empty.Empty, error)
}

func sendFile(file io.Reader, stream chunkedStream) error {
	buff := make([]byte, 400096)

	n, err := file.Read(buff)
	if err != nil {
		return err
	}

	for {
		err = stream.Send(&file_processor.Chunk{Data: buff[:n]})
		if err != nil {
			return err
		}

		n, err = file.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return nil
}
