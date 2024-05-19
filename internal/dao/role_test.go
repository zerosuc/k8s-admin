package dao

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"go-admin/internal/cache"
	"go-admin/internal/model"
)

func newRoleDao() *gotest.Dao {
	testData := &model.Role{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	//c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // to test mysql, disable caching
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewRoleCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = NewRoleDao(d.DB, c.ICache.(cache.RoleCache))

	return d
}

func Test_roleDao_Create(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(RoleDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleDao_DeleteByID(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	testData.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: false,
	}

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(RoleDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(RoleDao).DeleteByID(d.Ctx, 0)
	assert.Error(t, err)
}

func Test_roleDao_DeleteByIDs(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	testData.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: false,
	}

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(RoleDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(RoleDao).DeleteByIDs(d.Ctx, []uint64{0})
	assert.Error(t, err)
}

func Test_roleDao_UpdateByID(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(RoleDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(RoleDao).UpdateByID(d.Ctx, &model.Role{})
	assert.Error(t, err)

}

func Test_roleDao_GetByID(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(RoleDao).GetByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// notfound error
	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(2).
		WillReturnRows(rows)
	_, err = d.IDao.(RoleDao).GetByID(d.Ctx, 2)
	assert.Error(t, err)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(3, 4).
		WillReturnRows(rows)
	_, err = d.IDao.(RoleDao).GetByID(d.Ctx, 4)
	assert.Error(t, err)
}

func Test_roleDao_GetByCondition(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(RoleDao).GetByCondition(d.Ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: testData.ID,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// notfound error
	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(2).
		WillReturnRows(rows)
	_, err = d.IDao.(RoleDao).GetByCondition(d.Ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: 2,
			},
		},
	})
	assert.Error(t, err)
}

func Test_roleDao_GetByIDs(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(RoleDao).GetByIDs(d.Ctx, []uint64{testData.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.IDao.(RoleDao).GetByIDs(d.Ctx, []uint64{111})
	assert.Error(t, err)

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleDao_GetByLastID(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, err := d.IDao.(RoleDao).GetByLastID(d.Ctx, 0, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, err = d.IDao.(RoleDao).GetByLastID(d.Ctx, 0, 10, "unknown-column")
	assert.Error(t, err)
}

func Test_roleDao_GetByColumns(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(RoleDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Sort: "ignore count", // ignore test count(*)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, _, err = d.IDao.(RoleDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Columns: []query.Column{
			{
				Name:  "id",
				Exp:   "<",
				Value: 0,
			},
		},
	})
	assert.Error(t, err)

	// error test
	dao := &roleDao{}
	_, _, err = dao.GetByColumns(context.Background(), &query.Params{Columns: []query.Column{{}}})
	t.Log(err)
}

func Test_roleDao_CreateByTx(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	_, err := d.IDao.(RoleDao).CreateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleDao_DeleteByTx(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	testData.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: false,
	}

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(RoleDao).DeleteByTx(d.Ctx, d.DB, testData.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleDao_UpdateByTx(t *testing.T) {
	d := newRoleDao()
	defer d.Close()
	testData := d.TestData.(*model.Role)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(RoleDao).UpdateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}
