package storage

import (
	"io"
	"user-service/config"

	"github.com/labstack/gommon/log"
	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseI interface {
	UploadFile(path string, file io.Reader) (string, error)
}

type supabaseStruct struct {
	cfg *config.Config
}

func New(cfg *config.Config) SupabaseI {
	return &supabaseStruct{
		cfg: cfg,
	}
}

func (s *supabaseStruct) UploadFile(path string, file io.Reader) (string, error) {
	client := storage_go.NewClient(s.cfg.Storage.URL, s.cfg.Storage.Key, map[string]string{
		"Content-Type": "image/png",
	})

	_, err := client.UploadFile(s.cfg.Storage.Bucket, path, file)
	if err != nil {
		log.Errorf("[supabaseStruct-1] UploadFile: %v", err)
		return "", err
	}

	result := client.GetPublicUrl(s.cfg.Storage.Bucket, path)

	return result.SignedURL, nil
}
