package repository_tests

import (
	"github.com/tryffel/fusio/storage/models"
	"testing"
)

func TestCreateQueryPipeline(t *testing.T) {
	db := getDatabaseFromArgs()
	defer db.RemoveAllRecords()
	if db == nil {
		t.Error("No connection to database")
		return
	}

	user := db.createUser()
	pipeline := &models.Pipeline{
		OwnerId: user.ID,
		Enabled: true,
		Name:    "test",
	}

	// Just store interfaces and test they are stored correctly
	i1 := models.PipelineItem{Type: "user", Data: models.User{Name: "user_a"}}
	i2 := models.PipelineItem{Type: "device", Data: models.Device{Name: "d_1"}}

	pipeline.Blocks.Items = append(pipeline.Blocks.Items, i1)
	pipeline.Blocks.Items = append(pipeline.Blocks.Items, i2)

	err := db.Pipeline.Create(pipeline)
	if err != nil {
		t.Error(err)
	}

	p, err := db.Pipeline.FindbyId(pipeline.Id)
	if err != nil {
		t.Error(err)
	}

	if p.Blocks.Items[0].Type != "user" || p.Blocks.Items[1].Type != "device" {
		t.Error("block type doesn't match")
	}

	if p.Blocks.Items[0].Data.(map[string]interface{})["Name"].(string) != "user_a" {
		t.Error("block data doesn't match")
	}
}
