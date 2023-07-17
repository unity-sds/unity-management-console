package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestStoreConfig(t *testing.T) {
	Convey("Given a list of core configurations to store", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		So(err, ShouldBeNil)
		DB = gormDB

		configsToStore := []models.CoreConfig{
			{ID: 1, Key: "key1", Value: "value1"},
			{ID: 2, Key: "key2", Value: "value2"},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "core_configs" \("key","value","id"\) VALUES \(\$1,\$2,\$3\),\(\$4,\$5,\$6\) ON CONFLICT \("key"\) DO UPDATE SET "value"="excluded"."value" RETURNING "id"$`).
			WithArgs("key1", "value1", 1, "key2", "value2", 2).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
		mock.ExpectCommit()

		Convey("When StoreConfig is called", func() {
			storedConfigs, err := StoreConfig(configsToStore)

			Convey("Then it should return stored configurations and no error", func() {
				So(err, ShouldBeNil)
				So(storedConfigs, ShouldResemble, configsToStore)
			})
		})
	})
}
func TestFetchConfig(t *testing.T) {
	Convey("Given a database with core configurations", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		So(err, ShouldBeNil)
		DB = gormDB

		rows := sqlmock.NewRows([]string{"ID", "Key", "Value"}).
			AddRow(1, "key1", "value1").
			AddRow(2, "key2", "value2")

		mock.ExpectQuery("^SELECT (.+) FROM \"core_configs\"").WillReturnRows(rows)

		Convey("When FetchConfig is called", func() {
			configs, err := FetchConfig()

			Convey("Then it should return all configurations and no error", func() {
				So(err, ShouldBeNil)
				So(configs, ShouldResemble, []models.CoreConfig{
					{ID: 1, Key: "key1", Value: "value1"},
					{ID: 2, Key: "key2", Value: "value2"},
				})
			})
		})
	})
}
