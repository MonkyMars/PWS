package services

import (
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/types"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) GetFileByID(fileID string) (*types.File, error) {
	query := Query().SetOperation("select").SetTable("files").SetLimit(1)
	query.Where["file_id"] = fileID

	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return nil, nil
	}

	return data.Single, nil
}

func (fs *FileService) GetFilesBySubjectID(subjectID string) ([]types.File, error) {
	query := Query().SetOperation("select").SetTable("files")
	query.Where["subject_id"] = subjectID

	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
