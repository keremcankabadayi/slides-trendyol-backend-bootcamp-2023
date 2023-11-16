package persistence

import (
	"be-bootcamp-2023/application/domain/cqrs/command"
	"be-bootcamp-2023/application/domain/cqrs/query"
	"be-bootcamp-2023/pkg/couchbase"
	"context"
	"errors"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/google/martian/log"
	"time"
)

type CouchbaseRepository struct {
	collection *couchbase.Collection
	cluster    *couchbase.Cluster
	bucketName string
}

func NewCouchbaseRepository(cluster *couchbase.Cluster, bucketName string) (*CouchbaseRepository, error) {
	bucket := cluster.Bucket(bucketName)
	// We wait until the bucket is definitely connected and setup.
	err := bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return nil, err
	}
	return &CouchbaseRepository{
		collection: bucket.DefaultCollection(),
		cluster:    cluster,
		bucketName: bucketName,
	}, nil
}

func (repository *CouchbaseRepository) FindById(ctx context.Context, documentId string) (*query.Content, error) {
	docOut, err := repository.collection.Get(documentId, &gocb.GetOptions{
		Context: ctx,
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	var content = &query.Content{}
	err = docOut.Content(content)
	if err != nil {
		log.Errorf(fmt.Sprintf("Content deserialize error, content id: %s, err: %v", documentId, err))
		return nil, err
	}
	return content, nil
}

func (repository *CouchbaseRepository) Insert(documentId string, content *command.Content) error {
	_, err := repository.collection.Insert(documentId, content, nil)
	if err != nil {
		log.Errorf(fmt.Sprintf("Content Insert error, id: %s, err: %v", documentId, err))
		return err
	}
	return nil
}

func (repository *CouchbaseRepository) Upsert(documentId string, cas gocb.Cas, content *command.Content) error {
	_, err := repository.collection.Replace(documentId, content, &gocb.ReplaceOptions{Cas: cas})
	_ = fmt.Sprintf("Upsert: Replacing document id %s if document not found it tries to insert", documentId)

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return repository.Insert(documentId, content)
		}
		if errors.Is(err, gocb.ErrCasMismatch) {
			fmt.Printf("CasMismatch for operation %s", err.Error())
			return err
		}
		return err
	}

	log.Infof(fmt.Sprintf("CouchbaseRepository is Succeeded for Upsert. ContentId: %s", documentId))
	return nil
}

func (repository *CouchbaseRepository) MutateIn(documentId string, path string, subDocument interface{}) error {
	mops := []gocb.MutateInSpec{
		gocb.UpsertSpec(path, subDocument, &gocb.UpsertSpecOptions{}),
	}
	_, err := repository.collection.MutateIn(documentId, mops, &gocb.MutateInOptions{
		Timeout: 2 * time.Second,
	})
	if err != nil {
		return err
	}
	return err
}

func (repository *CouchbaseRepository) Get(documentId string) (*query.Content, error) {
	result, err := repository.collection.Get(documentId, &gocb.GetOptions{})

	if err != nil {
		return nil, err
	}

	var data *query.Content
	if err = result.Content(&data); err != nil {
		return nil, err
	}

	return data, err
}

func (repository *CouchbaseRepository) GetSpec(documentId string, fieldName string) (*string, error) {
	specs := []gocb.LookupInSpec{gocb.GetSpec(fieldName, &gocb.GetSpecOptions{})}
	result, err := repository.collection.LookupIn(documentId, specs, &gocb.LookupInOptions{})
	_ = result
	var data *string
	if err = result.ContentAt(0, &data); err != nil {
		return nil, err
	}

	return data, err
}

func (repository *CouchbaseRepository) GetCas(documentId string) (gocb.Cas, error) {
	result, err := repository.collection.Exists(documentId, nil)
	if err != nil {
		return 0, err
	}

	cas := result.Cas()
	return cas, err
}
