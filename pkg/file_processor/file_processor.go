package file_processor

import (
	"2019_2_Covenant/tools/time_parser"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"path/filepath"
)

var MetadataError = fmt.Errorf("failed to get meta data from context")

type FileProcessor struct {
	rootDir string
	db      *sql.DB
}

func NewFileProcessor(rootDir string, db *sql.DB) *FileProcessor {
	return &FileProcessor{
		rootDir: rootDir,
		db:      db,
	}
}

func (fr *FileProcessor) ProcessAudio(stream Files_ProcessAudioServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return MetadataError
	}

	src, err := readFile(stream)
	if err != nil {
		return err
	}

	stream.SendAndClose(&empty.Empty{})

	destName := uuid.New().String()
	destPath := filepath.Join(fr.rootDir, "/music/", destName)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}

	defer dest.Close()
	if _, err = io.Copy(dest, src); err != nil {
		return err
	}

	duration, err := time_parser.TrackDuration(destPath)
	if err != nil {
		return err
	}

	id := md.Get("id")[0]
	_, err = fr.db.Exec("UPDATE tracks SET path='/resources/music/'||$1, duration=$2 WHERE id=$3", destName, duration, id)
	return err
}

func (fr *FileProcessor) ProcessImage(stream Files_ProcessImageServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return MetadataError
	}

	src, err := readFile(stream)
	if err != nil {
		return err
	}

	var query string
	var folder string
	switch md.Get("type")[0] {
	case "albumPhoto":
		folder = "/photos/albums/"
		query = "UPDATE albums SET photo = '/resources/photos/albums/'||$1 WHERE id = $2"
	case "artistPhoto":
		folder = "/photos/artists/"
		query = "UPDATE artists SET photo = '/resources/photos/artists/'||$1 WHERE id = $2"
	case "avatar":
		folder = "/avatars/"
		query = "UPDATE users SET avatar = '/resources/avatars/'||$1 WHERE id = $2"
	}

	destName := uuid.New().String()
	dest, err := os.Create(filepath.Join(fr.rootDir, folder,destName))
	if err != nil {
		return err
	}

	defer dest.Close()
	if _, err = io.Copy(dest, src); err != nil {
		return err
	}

	id := md.Get("id")[0]

	_, err = fr.db.Exec(query, destName, id)
	return err
}

type ChunkedStream interface {
	Recv() (*Chunk, error)
	SendAndClose(*empty.Empty) error
}

func readFile(stream ChunkedStream) (io.Reader, error) {
	file := &bytes.Buffer{}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		file.Write(chunk.Data)
	}

	err := stream.SendAndClose(&empty.Empty{})
	if err != nil {
		return nil, err
	}

	return file, nil
}
