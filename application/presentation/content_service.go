package presentation

import (
	"be-bootcamp-2023/application/domain/cqrs/command"
	"be-bootcamp-2023/application/domain/cqrs/query"
	"be-bootcamp-2023/application/infrastructure/persistence"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"time"
)

type ContentService struct {
	repository *persistence.CouchbaseRepository
}

func NewContentService(repository *persistence.CouchbaseRepository) *ContentService {
	return &ContentService{
		repository: repository,
	}
}

func (c *ContentService) Insert(documentId string, content *command.Content) error {
	err := c.repository.Insert(documentId, content)

	if err != nil {
		return err
	}
	return err
}

func (c *ContentService) Upsert(documentId string, cas gocb.Cas, content *command.Content) error {
	fmt.Printf("trying to update cas: %d \n", cas)
	err := c.repository.Upsert(documentId, cas, content)
	time.Sleep(1 * time.Second)

	if err != nil {
		return err
	}
	return err
}

func (c *ContentService) MutateIn(documentId string, path string, contentStatus *command.ContentStatus) error {
	err := c.repository.MutateIn(documentId, path, contentStatus)

	if err != nil {
		return err
	}
	return err
}

func (c *ContentService) Get(documentId string) (*query.Content, error) {
	content, err := c.repository.Get(documentId)

	if err != nil {
		return nil, err
	}
	return content, err
}

func (c *ContentService) GetSpec(documentId string, fieldName string) (*string, error) {
	field, err := c.repository.GetSpec(documentId, fieldName)

	if err != nil {
		return nil, err
	}
	return field, err
}

func (c *ContentService) GetCas(documentId string) (*gocb.Cas, error) {
	cas, err := c.repository.GetCas(documentId)

	if err != nil {
		return nil, err
	}
	return &cas, err
}
