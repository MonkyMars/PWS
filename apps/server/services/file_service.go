package services

import (
	"fmt"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) GetFileByID(fileID string) (*types.File, error) {
	query := Query().SetOperation("select").SetTable("files").SetLimit(1).SetSelect(database.PrefixQuery(lib.TableFiles, []string{
		"id", "subject_id", "name", "created_at", "uploaded_by", "mime_type", "file_id", "folder_id", "updated_at", "url",
	}))
	query.Where[fmt.Sprintf("public.%s.file_id", lib.TableFiles)] = fileID
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
	query := Query().SetOperation("select").SetTable("files").SetSelect(database.PrefixQuery(lib.TableFiles, []string{
		"id", "subject_id", "name", "created_at", "uploaded_by", "mime_type", "file_id", "folder_id", "updated_at", "url",
	}))
	query.Where[fmt.Sprintf("public.%s.subject_id", lib.TableFiles)] = subjectID

	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
