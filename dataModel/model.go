package dataModel

import (
"github.com/lib/pq"
"gorm.io/gorm"
)

type RepositoriesIncludeKeyWordsJSON struct {
	Repositories []RepositoryModel `json:"repositories"`
}

type NearestRepositoriesJSON struct {
	Repositories map[uint]float64 `json:"repositories"`
}

type TitleFrequencyJSON struct {
	Dictionary map[string]float64 `json:"dictionary"`
}

//
//
//

type RepositoryModel struct {
	gorm.Model
	URL                                    string                     `gorm:"not null;"`
	Name                                   string                     `gorm:"not null; index:repository_name,unique;"`
	Owner                                  string                     `gorm:"not null;"`
	Topics                                 pq.StringArray             `gorm:"not null; type:text[];"`
	Description                            string                     `gorm:"not null;"`
	Issues                                 []IssueModel               `gorm:"foreignKey:RepositoryID; constraint:OnDelete:CASCADE;"`
	NearestIssuesOtherRepositories         []NearestIssuesModel       `gorm:"foreignKey:RepositoryID; constraint:OnDelete:CASCADE;"`
	IssuesNearestToIssuesOtherRepositories []NearestIssuesModel       `gorm:"foreignKey:RepositoryIDNearestIssue; constraint:OnDelete:CASCADE;"`
	NearestIRepositories                   []NearestRepositoriesModel `gorm:"foreignKey:RepositoryID; constraint:OnDelete:CASCADE;"`
	//
	IntersectionsIssues             []NumberIssueIntersectionsModel `gorm:"foreignKey:RepositoryID; constraint:OnDelete:CASCADE;"`
	AsComparableIntersectionsIssues []NumberIssueIntersectionsModel `gorm:"foreignKey:ComparableRepositoryID; constraint:OnDelete:CASCADE;"`
}

type IssueModel struct {
	gorm.Model
	RepositoryID                     uint
	Number                           int                  `gorm:"not null;"`
	URL                              string               `gorm:"not null;"`
	Title                            string               `gorm:"not null;"`
	State                            string               `gorm:"not null;"`
	Body                             string               `gorm:"not null;size:1000000"`
	TitleDictionary                  pq.StringArray       `gorm:"not null; type:text[];"`
	TitleFrequencyJSON               []byte               `gorm:"not null;"`
	NearestIssuesOtherRepositories   []NearestIssuesModel `gorm:"foreignKey:IssueID; constraint:OnDelete:CASCADE;"`
	NearestToIssuesOtherRepositories []NearestIssuesModel `gorm:"foreignKey:NearestIssueID; constraint:OnDelete:CASCADE;"`
}

type NearestIssuesModel struct {
	gorm.Model
	RepositoryID             uint
	IssueID                  uint
	NearestIssueID           uint
	RepositoryIDNearestIssue uint
	//
	TitleNumberIntersections float64
	TitleCosineDistance      float64
	BodyNumberIntersections  float64
	BodyCosineDistance       float64
	Rank                     float64
	//
	CosineDistance float64        `gorm:"not null;"`
	Intersections  pq.StringArray `gorm:"not null; type:text[];"`
}

type NearestRepositoriesModel struct {
	gorm.Model
	RepositoryID uint
	Repositories []byte `gorm:"not null;"`
}

type RepositoriesKeyWordsModel struct {
	gorm.Model
	KeyWord      string `gorm:"not null; index:key_word,unique;"`
	Position     int64
	Repositories []byte `gorm:"not null;"`
}

type NumberIssueIntersectionsModel struct {
	gorm.Model
	RepositoryID           uint
	ComparableRepositoryID uint
	NumberIntersections    float64 `gorm:"not null;"`
	//
	RepositoryCountIssues int64 `gorm:"not null;"`
	CountNearestPairs     int64 `gorm:"not null;"`
}
