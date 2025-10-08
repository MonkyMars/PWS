package services

import (
	"fmt"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
)

type ContentService struct{}

func NewContentService() *ContentService {
	return &ContentService{}
}

func (cs *ContentService) GetFileByID(fileID string) (*types.File, error) {
	query := Query().SetOperation("select").SetTable("files").SetLimit(1).SetSelect(database.PrefixQuery(lib.TableFiles, []string{
		"id", "subject_id", "name", "created_at", "uploaded_by", "mime_type", "file_id", "folder_id", "updated_at", "url",
	}))
	query.Where[fmt.Sprintf("public.%s.file_id", lib.TableFiles)] = fileID
	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return &types.File{}, nil
	}

	return data.Single, nil
}

func (cs *ContentService) GetFilesBySubjectID(subjectId string, folderId string) ([]types.File, error) {
	query := Query().SetOperation("select").SetTable("files").SetSelect(database.PrefixQuery(lib.TableFiles, []string{
		"id", "subject_id", "name", "created_at", "uploaded_by", "mime_type", "file_id", "folder_id", "updated_at", "url",
	}))
	query.Where[fmt.Sprintf("public.%s.subject_id", lib.TableFiles)] = subjectId
	query.Where[fmt.Sprintf("public.%s.folder_id", lib.TableFiles)] = folderId

	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return []types.File{}, nil
	}

	return data.Data, nil
}

func (cs *ContentService) GetFoldersByParentID(subjectId string, parentId string) ([]types.Folder, error) {
	query := Query().SetOperation("select").SetTable("folders").SetSelect(database.PrefixQuery(lib.TableFolders, []string{
		"id", "subject_id", "name", "created_at", "parent_id",
	}))
	query.Where[fmt.Sprintf("public.%s.subject_id", lib.TableFolders)] = subjectId
	query.Where[fmt.Sprintf("public.%s.parent_id", lib.TableFolders)] = parentId

	data, err := database.ExecuteQuery[types.Folder](query)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return []types.Folder{}, nil
	}

	return data.Data, nil
}
