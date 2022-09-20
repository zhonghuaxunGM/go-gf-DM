// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	// "github.com/gogf/gf/v2/frame/g"
	// "github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_DB_Ping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

func TestTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := []string{"A_tables4", "A_tables2"}

		for _, v := range tables {
			// createInitTable(v)
			createTable(v)
		}
		// TODO Question1
		// result, err := db.Tables(ctx)
		result, err := db.Tables(ctx, TestDbName)
		gtest.Assert(err, nil)

		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if strings.ToUpper(tables[i]) == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}

		result, err = dblink.Tables(ctx, TestDbName)
		gtest.Assert(err, nil)
		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if strings.ToUpper(tables[i]) == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}

		// for _, v := range tables {
		// dropTable(v)
		// }
	})
}

func TestTableFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := "b_tables"

		createTable(tables)
		// defer dropTable("Fields")
		var expect = map[string][]interface{}{
			"ID":           {"BIGINT", false},
			"ACCOUNT_NAME": {"VARCHAR", false},
			"PWD_RESET":    {"TINYINT", false},
			"DELETED":      {"INT", false},
			"CREATED_TIME": {"TIMESTAMP", false},
		}

		// _, err := dbErr.TableFields(ctx, "Fields")
		// gtest.AssertNE(err, nil)

		// res, err := dblink.TableFields(ctx, "Fields")
		res, err := db.TableFields(ctx, tables, TestDbName)
		gtest.Assert(err, nil)

		for k, v := range expect {
			_, ok := res[k]
			gtest.AssertEQ(ok, true)

			gtest.AssertEQ(res[k].Name, k)
			gtest.Assert(res[k].Type, v[0])
			gtest.Assert(res[k].Null, v[1])
			g.Dump(res)
		}

	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.TableFields(ctx, "t_user t_user2")
		gtest.AssertNE(err, nil)
	})
}

// func TestFilteredLink(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		s := dblink.FilteredLink()
// 		gtest.AssertEQ(s, "oracle:xxx@127.0.0.1:1521/XE")
// 	})
// }

func Test_DB_Query(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableName := "c_insert"
		// createTable(tableName)

		resOne, err := db.Query(ctx, fmt.Sprintf("SELECT * from %s", tableName))
		t.AssertNil(err)
		g.Dump("resOne", resOne)

		resTwo := make([]User, 0)
		err = db.Schema(TestDbName).Model(tableName).Scan(&resTwo)
		t.AssertNil(err)

		resThree := make([]User, 0)
		model := db.Model(tableName)
		// model.Where("id", g.Slice{401877392097280})
		// model.Where("account_name like ?", "%"+"xzh"+"%")
		model.Where("deleted", 0).Order("created_time desc")

		total, err := model.Count()
		t.AssertNil(err)
		g.Dump("total", total)

		err = model.Scan(&resThree)
		// err = model.Page(pageNo, pageSize).Scan(&result)
		t.AssertNil(err)
	})
}

func TestSave(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// createTable("DoInsert")
		// defer dropTable("DoInsert")
		data := []User{
			{
				ID:          5556,
				AccountName: "user_3",
				CreatedTime: time.Now(),
			},
			{
				ID:          22,
				AccountName: "user_3",
				CreatedTime: time.Now(),
			},
			{
				ID:          223,
				AccountName: "user_3",
				CreatedTime: time.Now(),
			},
		}
		_, err := db.Schema(TestDbName).Model("C_insert").Data(data).Save()
		gtest.Assert(err, nil)

		// _, err = db.Schema(TestDbName).Replace(ctx, "DoInsert", data, 10)
		// gtest.Assert(err, nil)
	})
}

func TestDoInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// createTable("C_insert")
		// defer dropTable("C_insert")
		i := 66
		data := User{
			ID:          int64(i),
			AccountName: fmt.Sprintf(`A%d222s`, i),
			PwdReset:    23,
			// CreatedTime: time.Now(),
		}
		_, err := db.Schema(TestDbName).Model("C_insert").OmitEmpty().Data(&data).Insert()
		gtest.Assert(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// createTable("C_insert")
		// defer dropTable("C_insert")
		i := 55
		data := g.Map{
			"ID":           i,
			"ACCOUNT_NAME": fmt.Sprintf(`A%d222s`, i),
			"PWD_RESET":    23,
			// "CREATED_TIME": gtime.Now().String(),
		}
		_, err := db.Schema(TestDbName).Insert(ctx, "C_insert", data)
		gtest.Assert(err, nil)
	})

}

// func Test_DB_Exec(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		_, err := db.Exec(ctx, "SELECT ? from dual", 1)
// 		t.AssertNil(err)

// 		_, err = db.Exec(ctx, "ERROR")
// 		t.AssertNE(err, nil)
// 	})
// }

func Test_DB_Insert(t *testing.T) {
	// table := createTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Schema(TestDbName).Insert(ctx, "c_insert", g.Map{
			"ID":           1,
			"ACCOUNT_NAME": "t1",
			"CREATED_TIME": gtime.Now().String(),
		})
		t.AssertNil(err)

		// normal map
		result, err := db.Schema(TestDbName).Insert(ctx, "c_insert", g.Map{
			"ID":           "2",
			"ACCOUNT_NAME": "t2",
			"CREATED_TIME": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// struct
		timeStr := time.Now()
		result, err = db.Schema(TestDbName).Insert(ctx, "c_insert", User{
			ID:          3,
			AccountName: "user_3",
			CreatedTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Schema(TestDbName).Model("c_insert").Where("ID", 3).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 3)
		t.Assert(one["ACCOUNT_NAME"].String(), "user_3")
		t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)

		// *struct
		timeStr = time.Now()
		result, err = db.Schema(TestDbName).Insert(ctx, "c_insert", &User{
			ID:          4,
			AccountName: "t4",
			CreatedTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err = db.Schema(TestDbName).Model("c_insert").Where("ID", 4).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 4)
		t.Assert(one["ACCOUNT_NAME"].String(), "t4")
		t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)

		// batch with Insert
		timeStr = time.Now()
		r, err := db.Schema(TestDbName).Insert(ctx, "c_insert", g.Slice{
			g.Map{
				"ID":           200,
				"ACCOUNT_NAME": "t200",
				"CREATED_TIME": timeStr,
			},
			g.Map{
				"ID":           300,
				"ACCOUNT_NAME": "t300",
				"CREATED_TIME": timeStr,
			},
		})
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		one, err = db.Schema(TestDbName).Model("c_insert").Where("ID", 200).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 200)
		t.Assert(one["ACCOUNT_NAME"].String(), "t200")
		t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)
	})
}

func Test_DB_BatchInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "c_insert"
		// table := createTable()
		// defer dropTable(table)
		r, err := db.Schema(TestDbName).Insert(ctx, table, g.List{
			{
				"ID":           22,
				"ACCOUNT_NAME": "t2",
				"CREATE_TIME":  gtime.Now().String(),
			},
			{
				"ID":           23,
				"ACCOUNT_NAME": "user_3",
				"CREATE_TIME":  gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		table := "c_insert"
		// table := createTable()
		// defer dropTable(table)
		// []interface{}
		r, err := db.Schema(TestDbName).Insert(ctx, table, g.Slice{
			g.Map{
				"ID":           32,
				"ACCOUNT_NAME": "32t2",
				"CREATE_TIME":  gtime.Now().String(),
			},
			g.Map{
				"ID":           33,
				"ACCOUNT_NAME": "33user_3",
				"CREATE_TIME":  gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
	})

	// batch insert map
	gtest.C(t, func(t *gtest.T) {
		table := "c_insert"
		// table := createTable()
		// defer dropTable(table)
		result, err := db.Schema(TestDbName).Insert(ctx, table, g.Map{
			"ID":           41,
			"ACCOUNT_NAME": "41user41",
			"CREATE_TIME":  gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_BatchInsert_Struct(t *testing.T) {
	// batch insert struct
	gtest.C(t, func(t *gtest.T) {
		table := "c_insert"
		// table := createTable()
		// defer dropTable(table)
		user := &User{
			ID:          5556,
			AccountName: "t1",
			// CreatedTime: time.Now(),
		}
		result, err := db.Schema(TestDbName).Model(table).OmitEmpty().Insert(user)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_Update(t *testing.T) {
	table := "c_insert"
	// table := createInitTable()
	// defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Update(ctx, table, "pwd_reset=6", "id=66")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("ID", 66).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 66)
		t.Assert(one["ACCOUNT_NAME"].String(), "A66222s")
	})
}

// func Test_DB_GetAll(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)

// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
// 		t.AssertNil(err)
// 		t.Assert(len(result), 1)
// 		t.Assert(result[0]["ID"].Int(), 1)
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
// 		t.AssertNil(err)
// 		t.Assert(len(result), 1)
// 		t.Assert(result[0]["ID"].Int(), 1)
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?)", table), g.Slice{1, 2, 3})
// 		t.AssertNil(err)
// 		t.Assert(len(result), 3)
// 		t.Assert(result[0]["ID"].Int(), 1)
// 		t.Assert(result[1]["ID"].Int(), 2)
// 		t.Assert(result[2]["ID"].Int(), 3)
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3})
// 		t.AssertNil(err)
// 		t.Assert(len(result), 3)
// 		t.Assert(result[0]["ID"].Int(), 1)
// 		t.Assert(result[1]["ID"].Int(), 2)
// 		t.Assert(result[2]["ID"].Int(), 3)
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3}...)
// 		t.AssertNil(err)
// 		t.Assert(len(result), 3)
// 		t.Assert(result[0]["ID"].Int(), 1)
// 		t.Assert(result[1]["ID"].Int(), 2)
// 		t.Assert(result[2]["ID"].Int(), 3)
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=?", table), g.Slice{1, 3})
// 		t.AssertNil(err)
// 		t.Assert(len(result), 3)
// 		t.Assert(result[0]["ID"].Int(), 1)
// 		t.Assert(result[1]["ID"].Int(), 2)
// 		t.Assert(result[2]["ID"].Int(), 3)
// 	})
// }

// func Test_DB_GetOne(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		record, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_1")
// 		t.AssertNil(err)
// 		t.Assert(record["NICKNAME"].String(), "name_1")
// 	})
// }

// func Test_DB_GetValue(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "user_3")
// 		t.AssertNil(err)
// 		t.Assert(value.Int(), 3)
// 	})
// }

// func Test_DB_GetCount(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		count, err := db.GetCount(ctx, fmt.Sprintf("SELECT * FROM %s", table))
// 		t.AssertNil(err)
// 		t.Assert(count, TableSize)
// 	})
// }

// func Test_DB_GetStruct(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime gtime.Time
// 		}
// 		user := new(User)
// 		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
// 		t.AssertNil(err)
// 		t.Assert(user.NickName, "name_3")
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime *gtime.Time
// 		}
// 		user := new(User)
// 		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
// 		t.AssertNil(err)
// 		t.Assert(user.NickName, "name_3")
// 	})
// }

// func Test_DB_GetStructs(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime gtime.Time
// 		}
// 		var users []User
// 		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
// 		t.AssertNil(err)
// 		t.Assert(len(users), TableSize-1)
// 		t.Assert(users[0].Id, 2)
// 		t.Assert(users[1].Id, 3)
// 		t.Assert(users[2].Id, 4)
// 		t.Assert(users[0].NickName, "name_2")
// 		t.Assert(users[1].NickName, "name_3")
// 		t.Assert(users[2].NickName, "name_4")
// 	})

// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime *gtime.Time
// 		}
// 		var users []User
// 		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
// 		t.AssertNil(err)
// 		t.Assert(len(users), TableSize-1)
// 		t.Assert(users[0].Id, 2)
// 		t.Assert(users[1].Id, 3)
// 		t.Assert(users[2].Id, 4)
// 		t.Assert(users[0].NickName, "name_2")
// 		t.Assert(users[1].NickName, "name_3")
// 		t.Assert(users[2].NickName, "name_4")
// 	})
// }

// func Test_DB_GetScan(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime gtime.Time
// 		}
// 		user := new(User)
// 		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
// 		t.AssertNil(err)
// 		t.Assert(user.NickName, "name_3")
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime gtime.Time
// 		}
// 		var user *User
// 		err := db.GetScan(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
// 		t.AssertNil(err)
// 		t.Assert(user.NickName, "name_3")
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime *gtime.Time
// 		}
// 		user := new(User)
// 		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
// 		t.AssertNil(err)
// 		t.Assert(user.NickName, "name_3")
// 	})

// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime gtime.Time
// 		}
// 		var users []User
// 		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
// 		t.AssertNil(err)
// 		t.Assert(len(users), TableSize-1)
// 		t.Assert(users[0].Id, 2)
// 		t.Assert(users[1].Id, 3)
// 		t.Assert(users[2].Id, 4)
// 		t.Assert(users[0].NickName, "name_2")
// 		t.Assert(users[1].NickName, "name_3")
// 		t.Assert(users[2].NickName, "name_4")
// 	})

// 	gtest.C(t, func(t *gtest.T) {
// 		type User struct {
// 			Id         int
// 			Passport   string
// 			Password   string
// 			NickName   string
// 			CreateTime *gtime.Time
// 		}
// 		var users []User
// 		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
// 		t.AssertNil(err)
// 		t.Assert(len(users), TableSize-1)
// 		t.Assert(users[0].Id, 2)
// 		t.Assert(users[1].Id, 3)
// 		t.Assert(users[2].Id, 4)
// 		t.Assert(users[0].NickName, "name_2")
// 		t.Assert(users[1].NickName, "name_3")
// 		t.Assert(users[2].NickName, "name_4")
// 	})
// }

func Test_DB_Delete(t *testing.T) {
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Delete(ctx, "c_insert", "1=1")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model("c_insert").Where("id", 23).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Empty_Slice_Argument(t *testing.T) {
	table := "c_insert"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf(`select * from %s where id in(?)`, table), g.Slice{})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}
