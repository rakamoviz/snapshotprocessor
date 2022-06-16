package db

import (
	"testing"

	"fmt"

	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/creates"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/repository"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func TestGorm(t *testing.T) {
	gormDB, err := gorm.Open(sqlite.Open("test.tdb"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	gormDB.AutoMigrate(&models.Cluster{}, &models.Node{}, &models.NodeStatus{})

	clusterRepository := repository.MakeRepository[models.Cluster](gormDB)
	nodeRepository := repository.MakeRepository[models.Node](gormDB)
	nodeStatusRepository := repository.MakeRepository[models.NodeStatus](gormDB)

	cluster1, err := clusterRepository.Execute(
		creates.CreateCluster("clust_123"),
	)
	fmt.Println(cluster1, err)

	node1, err := nodeRepository.Execute(
		creates.CreateNode("node_1", "clust_123"),
	)
	fmt.Println(node1, err)

	cpuUsage, _ := decimal.NewFromString("11.22")
	nodeStatus1, err := nodeStatusRepository.Execute(
		creates.CreateNodeStatus("node_1", time.Now(), cpuUsage, uint64(1), uint64(2)),
	)
	fmt.Println(nodeStatus1, err)
}
