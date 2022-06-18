package internal

import (
	"testing"

	"fmt"

	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities/creates"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func TestGorm(t *testing.T) {
	gormDB, err := gorm.Open(sqlite.Open("entities_test.tdb"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	gormDB.AutoMigrate(&entities.Cluster{}, &entities.Node{}, &entities.NodeStatus{})

	clusterRepository := repository.New[entities.Cluster](gormDB)
	nodeRepository := repository.New[entities.Node](gormDB)
	nodeStatusRepository := repository.New[entities.NodeStatus](gormDB)

	cluster1, err := clusterRepository.ExecuteOne(
		creates.CreateCluster("clust_123"),
	)
	fmt.Println(cluster1, err)

	node1, err := nodeRepository.ExecuteOne(
		creates.CreateNode("node_1", "clust_123"),
	)
	fmt.Println(node1, err)

	cpuUsage, _ := decimal.NewFromString("11.22")
	nodeStatus1, err := nodeStatusRepository.ExecuteOne(
		creates.CreateNodeStatus("node_1", time.Now(), cpuUsage, uint64(1), uint64(2)),
	)
	fmt.Println(nodeStatus1, err)
}
