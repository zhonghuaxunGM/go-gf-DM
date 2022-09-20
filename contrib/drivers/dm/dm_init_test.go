// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	_ "gitee.com/chunanyong/dm"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	db     gdb.DB
	dblink gdb.DB
	dbErr  gdb.DB
	ctx    context.Context
)

const (
	TableSize = 10

// TableName       = "inf_group"
// TableNamePrefix = "t_"
// TestSchema = "SYSDBADP"
)

const (
	TestDbIP    = "127.0.0.1"
	TestDbPort  = "5236"
	TestDbUser  = "SYSDBADP"
	TestDbPass  = "SYSDBADP"
	TestDbName  = "SYSDBADP"
	TestDbType  = "dm"
	TestCharset = "utf8"
)

func init() {
	node := gdb.ConfigNode{
		Host:             TestDbIP,
		Port:             TestDbPort,
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             TestDbName,
		Type:             TestDbType,
		Role:             "master",
		Charset:          TestCharset,
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
	}

	nodeLink := gdb.ConfigNode{
		Type: TestDbType,
		Name: TestDbName,
		Link: fmt.Sprintf(
			"dm://%s:%s@%s:%s/%s?charset=%s",
			TestDbUser, TestDbPass, TestDbIP, TestDbPort, TestDbName, TestCharset,
		),
	}

	nodeErr := gdb.ConfigNode{
		Host:    TestDbIP,
		Port:    TestDbPort,
		User:    TestDbUser,
		Pass:    "1234",
		Name:    TestDbName,
		Type:    TestDbType,
		Role:    "master",
		Charset: TestCharset,
		Weight:  1,
	}

	gdb.AddConfigNode(gdb.DefaultGroupName, node)
	if r, err := gdb.New(node); err != nil {
		gtest.Fatal(err)
	} else {
		db = r
	}

	gdb.AddConfigNode("dblink", nodeLink)
	if r, err := gdb.New(nodeLink); err != nil {
		gtest.Fatal(err)
	} else {
		dblink = r
	}

	gdb.AddConfigNode("dbErr", nodeErr)
	if r, err := gdb.New(nodeErr); err != nil {
		gtest.Fatal(err)
	} else {
		dbErr = r
	}

	ctx = context.Background()
}

func createTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("random_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE "%s"
(
"ID" BIGINT NOT NULL,
"ACCOUNT_NAME" VARCHAR(128) DEFAULT '' NOT NULL,
"PWD_RESET" TINYINT DEFAULT 0 NOT NULL,
"DESCRIPTION" VARCHAR(256) DEFAULT '' NOT NULL,
"PHOTO" VARCHAR(1024) DEFAULT '' NOT NULL,
"ENABLED" INT DEFAULT 1 NOT NULL,
"DELETED" INT DEFAULT 0 NOT NULL,
"CREATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"CREATED_TIME" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
"UPDATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"UPDATED_TIME" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
NOT CLUSTER PRIMARY KEY("ID")) STORAGE(ON "MAIN", CLUSTERBTR) ;
	`, name)); err != nil {
		gtest.Fatal(err)
	}

	return
}

type User struct {
	ID          int64     `gorm:"column:id" db:"id" json:"id" form:"id"`
	AccountName string    `gorm:"column:account_name" db:"account_name" json:"account_name" form:"account_name"`
	PwdReset    int64     `gorm:"column:pwd_reset" db:"pwd_reset" json:"pwd_reset" form:"pwd_reset"`
	Phone       string    `gorm:"column:phone" db:"phone" json:"phone" form:"phone"`
	Description string    `gorm:"column:description" db:"description" json:"description" form:"description"`
	Photo       string    `gorm:"column:photo" db:"photo" json:"photo" form:"photo"`
	Enabled     int64     `gorm:"column:enabled" db:"enabled" json:"enabled" form:"enabled"`
	Deleted     int64     `gorm:"column:deleted" db:"deleted" json:"deleted" form:"deleted"`
	CreatedBy   string    `gorm:"column:created_by" db:"created_by" json:"created_by" form:"created_by"`
	CreatedTime time.Time `gorm:"column:created_time" db:"created_time" json:"created_time" form:"created_time"`
	UpdatedBy   string    `gorm:"column:updated_by" db:"updated_by" json:"updated_by" form:"updated_by"`
	UpdatedTime time.Time `gorm:"column:updated_time" db:"updated_time" json:"updated_time" form:"updated_time"`
}

func createInitTable(table ...string) (name string) {
	name = createTable(table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":           i,
			"account_name": fmt.Sprintf(`name_%d`, i),
			"pwd_reset":    0,
			"create_time":  gtime.Now().String(),
		})
	}
	result, err := db.Insert(context.Background(), name, array.Slice())
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func dropTable(table string) {
	count, err := db.GetCount(ctx, "SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = ?", strings.ToUpper(table))
	if err != nil {
		gtest.Fatal(err)
	}

	if count == 0 {
		return
	}
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE %s", table)); err != nil {
		gtest.Fatal(err)
	}
}
