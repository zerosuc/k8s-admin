package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"go-admin/internal/model"
)

func newApiCache() *gotest.Cache {
	record1 := &model.Api{}
	record1.ID = 1
	record2 := &model.Api{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewApiCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_apiCache_Set(t *testing.T) {
	c := newApiCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Api)
	err := c.ICache.(ApiCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ApiCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_apiCache_Get(t *testing.T) {
	c := newApiCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Api)
	err := c.ICache.(ApiCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ApiCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ApiCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_apiCache_MultiGet(t *testing.T) {
	c := newApiCache()
	defer c.Close()

	var testData []*model.Api
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Api))
	}

	err := c.ICache.(ApiCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ApiCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Api))
	}
}

func Test_apiCache_MultiSet(t *testing.T) {
	c := newApiCache()
	defer c.Close()

	var testData []*model.Api
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Api))
	}

	err := c.ICache.(ApiCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_apiCache_Del(t *testing.T) {
	c := newApiCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Api)
	err := c.ICache.(ApiCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_apiCache_SetCacheWithNotFound(t *testing.T) {
	c := newApiCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Api)
	err := c.ICache.(ApiCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewApiCache(t *testing.T) {
	c := NewApiCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewApiCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewApiCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
