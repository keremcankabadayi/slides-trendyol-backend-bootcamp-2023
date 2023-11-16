package main

import (
	"be-bootcamp-2023/application/domain/cqrs/command"
	"be-bootcamp-2023/application/infrastructure/persistence"
	"be-bootcamp-2023/application/presentation"
	"be-bootcamp-2023/pkg/config"
	"be-bootcamp-2023/pkg/couchbase"
	"be-bootcamp-2023/pkg/util"
	"context"
	"errors"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := config.New()
	if err != nil {
		panic("application configuration could not read, error " + err.Error())
	}

	couchbaseConfig := &application.Couchbase
	cluster, err := couchbase.ConnectClusterWithConfig(couchbaseConfig)

	couchbaseRepository, err := persistence.NewCouchbaseRepository(cluster, couchbaseConfig.BucketName)
	if err != nil {
		fmt.Printf("Could not create %s repository, err: %v", couchbaseConfig.BucketName, err)
		return
	}

	service := presentation.NewContentService(couchbaseRepository)

	contentId := int64(123)
	culture := "tr-TR"

	// Case 1: Create data
	//createData(service, contentId, culture)

	// Case 2: Update data
	//cas, _ := service.GetCas(getDocumentId(contentId, culture))
	//updateData(service, contentId, *cas, culture)

	// Case 3: Mutate data
	//mutateData(service, contentId, culture)

	// Case 4: Mutate data with status model
	//mutateDataWithNewStatusModel(service, contentId, culture)

	// Case 5: Get data
	//getData(service, contentId, culture)

	// Case 6: Get lookup data
	//getLookupData(service, contentId, culture)

	// Case 7: Exceptions: Not Found
	//getData(service, 1234, culture)

	//Case 8: Exceptions: CASMismatchException

	cas, _ := service.GetCas(getDocumentId(contentId, culture))
	fmt.Printf("first transaction cas: %d \n", cas)
	var wg sync.WaitGroup
	wg.Add(1)
	go updateRoutine(&wg, service, *cas, contentId, culture)
	wg.Wait()

	updateData(service, contentId, *cas, culture)

}

func createData(service *presentation.ContentService, id int64, culture string) {
	content := &command.Content{
		Id:               id,
		Name:             "Oversize Kapşonlu Sweatshirt",
		Description:      "Büyük beden ve bol kesim",
		Status:           command.ContentStatus{Id: 0},
		Culture:          culture,
		CreationDate:     util.GetFormattedNow(),
		LastModifiedDate: util.GetFormattedNow(),
		CreatedBy:        "bootcamp-trainer@trendyol.com",
		LastModifiedBy:   "bootcamp-trainer@trendyol.com",
	}

	docId := getDocumentId(id, culture)
	err := service.Insert(docId, content)

	if err != nil {
		return
	}
}

func updateData(service *presentation.ContentService, id int64, cas gocb.Cas, culture string) {
	content := &command.Content{
		Id:               id,
		Name:             "Oversize Kapşonlu Sweatshirt",
		Description:      "Büyük beden ve bol kesim",
		Status:           command.ContentStatus{Id: 3},
		Culture:          culture,
		CreationDate:     util.GetFormattedNow(),
		LastModifiedDate: util.GetFormattedNow(),
		CreatedBy:        "bootcamp-trainer@trendyol.com",
		LastModifiedBy:   "bootcamp-trainer@trendyol.com",
	}

	docId := getDocumentId(id, culture)
	err := service.Upsert(docId, cas, content)

	if err != nil {
		return
	}

}

func mutateData(service *presentation.ContentService, id int64, culture string) {
	contentStatus := &command.ContentStatus{
		Id: 5,
	}
	docId := getDocumentId(id, culture)
	err := service.MutateIn(docId, "status", contentStatus)

	if err != nil {
		return
	}
}

func mutateDataWithNewStatusModel(service *presentation.ContentService, id int64, culture string) {
	contentStatus := &command.ContentStatus{
		Id:               6,
		LastModifiedBy:   "bootcamp-speaker@trendyol.com",
		LastModifiedDate: util.GetFormattedNow(),
	}
	docId := getDocumentId(id, culture)
	err := service.MutateIn(docId, "status", contentStatus)

	if err != nil {
		return
	}
}

func getData(service *presentation.ContentService, id int64, culture string) {
	docId := getDocumentId(id, culture)
	content, err := service.Get(docId)

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return
		}
		return
	}

	fmt.Println(content)
}

func getLookupData(service *presentation.ContentService, id int64, culture string) {
	docId := getDocumentId(id, culture)
	field, err := service.GetSpec(docId, "description")

	if err != nil {
		return
	}

	_ = field
}

func updateRoutine(wg *sync.WaitGroup, service *presentation.ContentService, cas gocb.Cas, contentId int64, culture string) {
	defer wg.Done()
	newCas, _ := service.GetCas(getDocumentId(contentId, culture))
	fmt.Printf("new cas: %d \n", newCas)
	updateData(service, contentId, cas, culture)
}

func getDocumentId(id int64, culture string) string {
	return fmt.Sprintf("%d_%s", id, culture)
}
