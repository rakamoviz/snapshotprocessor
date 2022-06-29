package internal

import (
	"context"
	"testing"

	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/internal/entities/createfirsts"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
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
	//nodeStatusRepository := repository.New[entities.NodeStatus](gormDB)

	ctx := context.Background()

	cluster := entities.Cluster{Code: "cluster_123"}
	cluster1, err := clusterRepository.ExecuteOne(ctx, createfirsts.Cluster(cluster))
	fmt.Println(cluster1, err)

	node := entities.Node{Code: "node_1", Cluster: *cluster1}
	node1, err := nodeRepository.ExecuteOne(ctx, createfirsts.Node(node))
	fmt.Println(node1, err)

	/*
		cpuUsage, _ := decimal.NewFromString("11.22")
		nodeStatus1, err := nodeStatusRepository.ExecuteOne(
			creates.CreateNodeStatus("node_1", time.Now(), cpuUsage, uint64(1), uint64(2)),
		)
		fmt.Println(nodeStatus1, err)
	*/
}
