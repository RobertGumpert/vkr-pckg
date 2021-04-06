package repository

import "github.com/RobertGumpert/vkr-pckg/dataModel"

type IStorage interface {
	HasEntities() error
	CreateEntities() error
	Migration() error
	CloseConnection() error
}

type IRepository interface {
	IStorage
	//
	// REPOSITORY:
	//
	AddRepository(repository *dataModel.RepositoryModel) error
	AddRepositories(repositories []dataModel.RepositoryModel) error
	AddNearestRepositories(repositoryId uint, nearest dataModel.NearestRepositoriesJSON) error
	UpdateNearestRepositories(repositoryId uint, nearest dataModel.NearestRepositoriesJSON) error
	GetRepositoryByName(name string) (dataModel.RepositoryModel, error)
	GetRepositoryByID(repositoryId uint) (dataModel.RepositoryModel, error)
	GetNearestRepositories(repositoryId uint) (dataModel.NearestRepositoriesJSON, error)
	GetAllRepositories() ([]dataModel.RepositoryModel, error)
	RewriteAllNearestRepositories(repositoryId []uint, models []dataModel.NearestRepositoriesJSON) error
	//
	// ISSUE
	//
	AddIssue(issue *dataModel.IssueModel) error
	AddIssues(issues []dataModel.IssueModel) error
	AddNearestIssues(nearest dataModel.NearestIssuesModel) error
	GetIssuesBesidesGroupRepositories(repositoryId ...uint) ([]dataModel.IssueModel, error)
	GetIssueByID(issueId uint) (dataModel.IssueModel, error)
	GetIssueRepository(repositoryId uint) ([]dataModel.IssueModel, error)
	GetNearestIssuesForIssue(issueId uint) ([]dataModel.NearestIssuesModel, error)
	GetNearestIssuesForRepository(repositoryId uint) ([]dataModel.NearestIssuesModel, error)
	//
	// KEYWORDS:
	//
	AddKeyWord(keyWord string, position int64, repositories dataModel.RepositoriesIncludeKeyWordsJSON) (dataModel.RepositoriesKeyWordsModel, error)
	UpdateKeyWord(keyWord string, position int64,  repositories dataModel.RepositoriesIncludeKeyWordsJSON) (dataModel.RepositoriesKeyWordsModel, error)
	GetKeyWord(keyWord string) (dataModel.RepositoriesKeyWordsModel, error)
	GetAllKeyWords() ([]dataModel.RepositoriesKeyWordsModel, error)
	RewriteAllKeyWords(models []dataModel.RepositoriesKeyWordsModel) error
}

