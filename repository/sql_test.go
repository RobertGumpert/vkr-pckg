package repository

import (
	"github.com/RobertGumpert/vkr-pckg/dataModel"
	"github.com/RobertGumpert/vkr-pckg/runtimeinfo"
	"testing"
)

var storageProvider = SQLCreateConnection(
	TypeStoragePostgres,
	DSNPostgres,
	nil,
	"postgres",
	"toster123",
	"vkr-db",
	"5432",
	"disable",
)

func connect() IRepository {
	sqlRepository := NewSQLRepository(
		storageProvider,
	)
	return sqlRepository
}

func TestTruncate(t *testing.T) {
	_ = connect()
	storageProvider.SqlDB.Exec("TRUNCATE TABLE repositories CASCADE")
	storageProvider.SqlDB.Exec("TRUNCATE TABLE issues CASCADE")
}

func TestMigration(t *testing.T) {
	storageProvider.SqlDB.Exec("drop table repository_models cascade")
	storageProvider.SqlDB.Exec("drop table issue_models cascade")
	storageProvider.SqlDB.Exec("drop table nearest_issues_models cascade")
	storageProvider.SqlDB.Exec("drop table nearest_repositories_models cascade")
	storageProvider.SqlDB.Exec("drop table repositories_key_words_models cascade")
	storageProvider.SqlDB.Exec("drop table number_issue_intersections_models cascade")
	_ = connect()
}

func TestAddFlow(t *testing.T) {
	db := connect()
	//
	repository := dataModel.RepositoryModel{
		URL:         "a",
		Name:        "a",
		Owner:       "a",
		Topics:      []string{"a", "a", "a"},
		Description: "a",
	}
	repositories := []dataModel.RepositoryModel{
		{
			URL:         "b",
			Name:        "b",
			Owner:       "b",
			Topics:      []string{"b", "b", "b"},
			Description: "b",
		},
		{
			URL:         "c",
			Name:        "c",
			Owner:       "c",
			Topics:      []string{"c", "c", "c"},
			Description: "c",
		},
	}
	//
	err := db.AddRepository(&repository)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	//
	err = db.AddRepositories(repositories)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	//
	issues1 := []dataModel.IssueModel{
		{
			RepositoryID:       repository.ID,
			Number:             1,
			URL:                "a",
			Title:              "a",
			State:              "a",
			Body:               "a",
			TitleDictionary:    []string{"a"},
			TitleFrequencyJSON: []byte{23},
		},
	}
	err = db.AddIssues(issues1)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	//
	issues2 := make([]dataModel.IssueModel, 0)
	for index, repo := range repositories {
		issues2 = append(issues2, dataModel.IssueModel{
			RepositoryID:       repo.ID,
			Number:             index,
			URL:                repo.URL,
			Title:              repo.URL,
			State:              repo.URL,
			Body:               repo.URL,
			TitleDictionary:    []string{"b"},
			TitleFrequencyJSON: []byte{33},
		})
	}
	err = db.AddIssues(issues2)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	//
	list, err := db.GetIssueRepository(repository.ID)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	runtimeinfo.LogInfo(list)
	//
	err = db.AddNearestIssues(dataModel.NearestIssuesModel{
		RepositoryID:             repository.ID,
		IssueID:                  issues1[0].ID,
		RepositoryIDNearestIssue: repositories[0].ID,
		NearestIssueID:           issues2[0].ID,
		CosineDistance:           75.9,
		Intersections:            []string{"b"},
	})
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	//
	list, err = db.GetIssuesBesidesGroupRepositories(repositories[0].ID, repositories[1].ID)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	for _, elem := range list {
		runtimeinfo.LogInfo("ISSUE REPO: ", elem.RepositoryID)
		if elem.RepositoryID != repository.ID {
			t.Fatal("NOT WORKING")
		}
	}
	//
	list, err = db.GetIssuesOnlyGroupRepositories(repositories[0].ID, repositories[1].ID)
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	for _, elem := range list {
		runtimeinfo.LogInfo("ISSUE REPO: ", elem.RepositoryID)
		if elem.RepositoryID != repositories[0].ID && elem.RepositoryID != repositories[1].ID {
			t.Fatal("NOT WORKING")
		}
	}
	//
	err = db.CloseConnection()
	if err != nil {
		runtimeinfo.LogError(err)
		t.Fatal()
	}
	runtimeinfo.LogInfo("Ok")
}

//func BenchmarkAddKeyWords(b *testing.B) {
//	b.ReportAllocs()
//	db := connect()
//	for i := 0; i < b.N; i++ {
//		inc := strconv.Itoa(i)
//		key := "Key" + inc
//		_, _ = db.AddKeyWord(
//			key,
//			dataModel.RepositoriesIncludeKeyWordsJSON{
//				Repositories: []dataModel.RepositoryModel{
//					{
//						Name: key,
//					},
//				},
//			},
//		)
//	}
//}
//
//func BenchmarkReadKeyWords(b *testing.B) {
//	b.ReportAllocs()
//	db := connect()
//	for i := 0; i < b.N; i++ {
//		inc := strconv.Itoa(i)
//		key := "Key" + inc
//		_, _ = db.GetKeyWord(
//			key,
//		)
//	}
//}
//
//func BenchmarkUpdateKeyWords(b *testing.B) {
//	b.ReportAllocs()
//	db := connect()
//	for i := 0; i < b.N; i++ {
//		inc := strconv.Itoa(i)
//		key := "Update Key" + inc
//		_, _ = db.UpdateKeyWord(
//			key,
//			dataModel.RepositoriesIncludeKeyWordsJSON{
//				Repositories: []dataModel.RepositoryModel{
//					{
//						Name: key,
//					},
//				},
//			},
//		)
//	}
//}
