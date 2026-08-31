package client

type StorageClient struct {
}

func NewStorageClient() *StorageClient {
	return &StorageClient{}
}

// func (s *StorageClient) GetClient() *StorageClient {
// 	return s
// }

func (s *StorageClient) Upload() error {
	return nil
}

func (s *StorageClient) GetPresignedUrl() error {
	return nil
}
