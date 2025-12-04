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
	if fileID == "" {
		return nil, fmt.Errorf("fileID parameter is required")
	}

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

func (cs *ContentService) GetFilesBySubjectID(subjectId, folderId string, hasPrivileges bool) ([]types.File, error) {
	query := Query().SetOperation("select").SetTable("files").SetSelect(database.PrefixQuery(lib.TableFiles, []string{
		"id", "subject_id", "name", "created_at", "uploaded_by", "mime_type", "file_id", "folder_id", "updated_at", "url",
	}))
	query.Where[fmt.Sprintf("public.%s.subject_id", lib.TableFiles)] = subjectId
	query.Where[fmt.Sprintf("public.%s.folder_id", lib.TableFiles)] = folderId
	if !hasPrivileges {
		query.Where[fmt.Sprintf("public.%s.public", lib.TableFiles)] = true
	}
	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return []types.File{}, nil
	}

	return data.Data, nil
}

func (cs *ContentService) GetFoldersByParentID(subjectId string, parentId string, hasPrivileges bool) ([]types.Folder, error) {
	query := Query().SetOperation("select").SetTable("folders").SetSelect(database.PrefixQuery(lib.TableFolders, []string{
		"id", "subject_id", "name", "created_at", "parent_id",
	}))
	query.Where[fmt.Sprintf("public.%s.subject_id", lib.TableFolders)] = subjectId
	query.Where[fmt.Sprintf("public.%s.parent_id", lib.TableFolders)] = parentId
	query.Where[fmt.Sprintf("public.%s.public", lib.TableFolders)] = !hasPrivileges

	data, err := database.ExecuteQuery[types.Folder](query)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return []types.Folder{}, nil
	}

	return data.Data, nil
}

// CreateFile creates a new file record in the database
func (cs *ContentService) CreateFile(fileData map[string]any) (*types.File, error) {
	query := Query().SetOperation("insert").SetTable("files")
	query.Data = fileData

	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		return nil, err
	}

	return data.Single, nil
}

// CreateFolder creates a new folder record in the database
func (cs *ContentService) CreateFolder(folderData map[string]any) (*types.Folder, error) {
	query := Query().SetOperation("insert").SetTable("folders")
	query.Data = folderData

	data, err := database.ExecuteQuery[types.Folder](query)
	if err != nil {
		return nil, err
	}

	return data.Single, nil
}

// CreateMultipleFiles creates multiple file records in the database using batch insert
func (cs *ContentService) CreateMultipleFiles(filesData []map[string]any) ([]*types.File, error) {
	if len(filesData) == 0 {
		return []*types.File{}, nil
	}

	// For batch operations, we need to use the database actions
	files := make([]*types.File, 0, len(filesData))

	for _, fileData := range filesData {
		file, err := cs.CreateFile(fileData)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// ContentServiceInterface defines the methods that any content service implementation must provide.
type ContentServiceInterface interface {
	// File operations
	GetFileByID(fileID string) (*types.File, error)
	GetFilesBySubjectID(subjectID, folderID string, hasPrivileges bool) ([]types.File, error)
	CreateFile(fileData map[string]any) (*types.File, error)

	// Folder operations
	GetFoldersByParentID(subjectID, parentID string, hasPrivileges bool) ([]types.Folder, error)
	CreateFolder(folderData map[string]any) (*types.Folder, error)

	// Batch operations for multiple file uploads
	CreateMultipleFiles(filesData []map[string]any) ([]*types.File, error)
}
