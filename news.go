package crawler

import (
	"github.com/go-some/txtanalyzer"
)

type News struct {
	Title, Body, Url, Time, Origin, ImgUrl string
	BodySum                                string
	HasGraphImg                            bool
	EntitiesInTitle                        []txtanalyzer.Entity
	PersonList                             []string
	OrgList                                []string
	ProdList                               []string
}
