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

	// gtest.C(t, func(t *gtest.T) {
	// 	_, err := db.TableFields(ctx, "t_user t_user2")
	// 	gtest.AssertNE(err, nil)
	// })
}

// func TestFilteredLink(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		s := dblink.FilteredLink()
// 		gtest.AssertEQ(s, "oracle:xxx@127.0.0.1:1521/XE")
// 	})
// }

func Test_DB_Query(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableName := "t_inf_user"
		createTable(tableName)

		resOne, err := db.Query(ctx, fmt.Sprintf("SELECT * from %s", tableName))
		t.AssertNil(err)
		g.Dump(resOne)

		resTwo := make([]User, 0)
		err = db.Schema(TestDbName).Model(tableName).Scan(&resTwo)
		t.AssertNil(err)

		resThree := make([]User, 0)
		model := db.Model(tableName)
		// model.Where("id", g.Slice{401877392097280})
		// model.Where("account_name like ?", "%"+"xzh"+"%")
		model.Where("deleted", 1).Order("created_time desc")

		total, err := model.Count()
		t.AssertNil(err)
		g.Dump(total)

		err = model.Scan(&resThree)
		// err = model.Page(pageNo, pageSize).Scan(&result)
		t.AssertNil(err)
	})
}

func TestDoInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		createTable("C_insert")
		// defer dropTable("DoInsert")

		i := 1122
		data := g.Map{
			"ID":           i,
			"ACCOUNT_NAME": fmt.Sprintf(`A%d`, i),
			"PWD_RESET":    1,
			"CREATED_TIME": gtime.Now().String(),
		}
		_, err := db.Schema(TestDbName).Insert(ctx, "DoInsert", data)
		gtest.Assert(err, nil)
	})

	// gtest.C(t, func(t *gtest.T) {
	// 	createTable("DoInsert")
	// 	// defer dropTable("DoInsert")

	// 	i := 22
	// 	data := g.Map{
	// 		"ID":           i,
	// 		"ACCOUNT_NAME": fmt.Sprintf(`t%d`, i),
	// 		"PWD_RESET":    2,
	// 		"CREATED_TIME": gtime.Now().String(),
	// 	}
	// 	_, err := db.Schema(TestDbName).Save(ctx, "DoInsert", data, 10)
	// 	gtest.Assert(err, nil)

	// 	_, err = db.Schema(TestDbName).Replace(ctx, "DoInsert", data, 10)
	// 	gtest.Assert(err, nil)
	// })
}

// func Test_DB_Exec(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		_, err := db.Exec(ctx, "SELECT ? from dual", 1)
// 		t.AssertNil(err)

// 		_, err = db.Exec(ctx, "ERROR")
// 		t.AssertNE(err, nil)
// 	})
// }

// func Test_DB_Insert(t *testing.T) {
// 	// table := createTable()
// 	// defer dropTable(table)

// 	gtest.C(t, func(t *gtest.T) {
// 		_, err := db.Schema("DP").Insert(ctx, "DP.t_inf_user", g.Map{
// 			"ID":           122233,
// 			"ACCOUNT_NAME": "t1",
// 			"USER_NAME":    "25d55ad283aa400af464c76d713c07ad",
// 			"SALT":         "T1",
// 			"CREATED_TIME": gtime.Now().String(),
// 		})
// 		t.AssertNil(err)

// 		// normal map
// 		result, err := db.Schema("DP").Insert(ctx, "DP.t_inf_user", g.Map{
// 			"ID":           "233",
// 			"ACCOUNT_NAME": "t2",
// 			"USER_NAME":    "25d55ad283aa400af464c76d713c07ad",
// 			"SALT":         "name_2",
// 			"CREATED_TIME": gtime.Now().String(),
// 		})
// 		t.AssertNil(err)
// 		n, _ := result.RowsAffected()
// 		t.Assert(n, 1)

// 		// struct
// 		type User struct {
// 			ID           int    `gconv:"ID"`
// 			ACCOUNT_NAME string `json:"ACCOUNT_NAME"`
// 			USER_NAME    string `gconv:"USER_NAME"`
// 			SALT         string `gconv:"SALT"`
// 			CREATED_TIME string `json:"CREATED_TIME"`
// 		}
// 		timeStr := gtime.Now().String()
// 		result, err = db.Schema("DP").Insert(ctx, "DP.t_inf_user", User{
// 			ID:           33,
// 			ACCOUNT_NAME: "user_3",
// 			USER_NAME:    "25d55ad283aa400af464c76d713c07ad",
// 			SALT:         "name_3",
// 			CREATED_TIME: timeStr,
// 		})
// 		t.AssertNil(err)
// 		n, _ = result.RowsAffected()
// 		t.Assert(n, 1)

// 		one, err := db.Schema("DP").Model("DP.t_inf_user").Where("ID", 33).One()
// 		t.AssertNil(err)
// 		fmt.Println(one)
// 		t.Assert(one["ID"].Int(), 33)
// 		t.Assert(one["ACCOUNT_NAME"].String(), "user_3")
// 		t.Assert(one["USER_NAME"].String(), "25d55ad283aa400af464c76d713c07ad")
// 		t.Assert(one["SALT"].String(), "name_3")
// 		t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)

// 		// *struct
// 		timeStr = gtime.Now().String()
// 		result, err = db.Schema("DP").Insert(ctx, "DP.t_inf_user", &User{
// 			ID:           4,
// 			ACCOUNT_NAME: "t4",
// 			USER_NAME:    "25d55ad283aa400af464c76d713c07ad",
// 			SALT:         "name_4",
// 			CREATED_TIME: timeStr,
// 		})
// 		t.AssertNil(err)
// 		n, _ = result.RowsAffected()
// 		t.Assert(n, 1)

// 		one, err = db.Schema("DP").Model("DP.t_inf_user").Where("ID", 4).One()
// 		t.AssertNil(err)
// 		t.Assert(one["ID"].Int(), 4)
// 		t.Assert(one["ACCOUNT_NAME"].String(), "t4")
// 		t.Assert(one["USER_NAME"].String(), "25d55ad283aa400af464c76d713c07ad")
// 		t.Assert(one["SALT"].String(), "name_4")
// 		t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)

// 		// batch with Insert
// 		timeStr = gtime.Now().String()
// 		r, err := db.Schema("DP").Insert(ctx, "DP.t_inf_user", g.Slice{
// 			g.Map{
// 				"ID":           200,
// 				"ACCOUNT_NAME": "t200",
// 				"USER_NAME":    "25d55ad283aa400af464c76d71qw07ad",
// 				"SALT":         "T200",
// 				"CREATED_TIME": timeStr,
// 			},
// 			g.Map{
// 				"ID":           300,
// 				"ACCOUNT_NAME": "t300",
// 				"USER_NAME":    "25d55ad283aa400af464c76d713c07ad",
// 				"SALT":         "T300",
// 				"CREATED_TIME": timeStr,
// 			},
// 		})
// 		t.AssertNil(err)
// 		n, _ = r.RowsAffected()
// 		t.Assert(n, 2)

// 		one, err = db.Schema("DP").Model("DP.t_inf_user").Where("ID", 200).One()
// 		t.AssertNil(err)
// 		t.Assert(one["ID"].Int(), 200)
// 		t.Assert(one["ACCOUNT_NAME"].String(), "t200")
// 		t.Assert(one["USER_NAME"].String(), "25d55ad283aa400af464c76d71qw07ad")
// 		t.Assert(one["SALT"].String(), "T200")
// 		t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)
// 	})
// }

// func Test_DB_BatchInsert(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		table := createTable()
// 		defer dropTable(table)
// 		r, err := db.Insert(ctx, table, g.List{
// 			{
// 				"ID":          2,
// 				"PASSPORT":    "t2",
// 				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
// 				"NICKNAME":    "name_2",
// 				"CREATE_TIME": gtime.Now().String(),
// 			},
// 			{
// 				"ID":          3,
// 				"PASSPORT":    "user_3",
// 				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
// 				"NICKNAME":    "name_3",
// 				"CREATE_TIME": gtime.Now().String(),
// 			},
// 		}, 1)
// 		t.AssertNil(err)
// 		n, _ := r.RowsAffected()
// 		t.Assert(n, 2)

// 	})

// 	gtest.C(t, func(t *gtest.T) {
// 		table := createTable()
// 		defer dropTable(table)
// 		// []interface{}
// 		r, err := db.Insert(ctx, table, g.Slice{
// 			g.Map{
// 				"ID":          2,
// 				"PASSPORT":    "t2",
// 				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
// 				"NICKNAME":    "name_2",
// 				"CREATE_TIME": gtime.Now().String(),
// 			},
// 			g.Map{
// 				"ID":          3,
// 				"PASSPORT":    "user_3",
// 				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
// 				"NICKNAME":    "name_3",
// 				"CREATE_TIME": gtime.Now().String(),
// 			},
// 		}, 1)
// 		t.AssertNil(err)
// 		n, _ := r.RowsAffected()
// 		t.Assert(n, 2)
// 	})

// 	// batch insert map
// 	gtest.C(t, func(t *gtest.T) {
// 		table := createTable()
// 		defer dropTable(table)
// 		result, err := db.Insert(ctx, table, g.Map{
// 			"ID":          1,
// 			"PASSPORT":    "t1",
// 			"PASSWORD":    "p1",
// 			"NICKNAME":    "T1",
// 			"CREATE_TIME": gtime.Now().String(),
// 		})
// 		t.AssertNil(err)
// 		n, _ := result.RowsAffected()
// 		t.Assert(n, 1)
// 	})
// }

// func Test_DB_BatchInsert_Struct(t *testing.T) {
// 	// batch insert struct
// 	gtest.C(t, func(t *gtest.T) {
// 		table := createTable()
// 		defer dropTable(table)

// 		type User struct {
// 			Id         int         `c:"ID"`
// 			Passport   string      `c:"PASSPORT"`
// 			Password   string      `c:"PASSWORD"`
// 			NickName   string      `c:"NICKNAME"`
// 			CreateTime *gtime.Time `c:"CREATE_TIME"`
// 		}
// 		user := &User{
// 			Id:         1,
// 			Passport:   "t1",
// 			Password:   "p1",
// 			NickName:   "T1",
// 			CreateTime: gtime.Now(),
// 		}
// 		result, err := db.Insert(ctx, table, user)
// 		t.AssertNil(err)
// 		n, _ := result.RowsAffected()
// 		t.Assert(n, 1)
// 	})
// }

// func Test_DB_Update(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)

// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.Update(ctx, table, "password='987654321'", "id=3")
// 		t.AssertNil(err)
// 		n, _ := result.RowsAffected()
// 		t.Assert(n, 1)

// 		one, err := db.Model(table).Where("ID", 3).One()
// 		t.AssertNil(err)
// 		t.Assert(one["ID"].Int(), 3)
// 		t.Assert(one["PASSPORT"].String(), "user_3")
// 		t.Assert(strings.TrimSpace(one["PASSWORD"].String()), "987654321")
// 		t.Assert(one["NICKNAME"].String(), "name_3")
// 	})
// }

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

// func Test_DB_Delete(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.Delete(ctx, table, "1=1")
// 		t.AssertNil(err)
// 		n, _ := result.RowsAffected()
// 		t.Assert(n, TableSize)
// 	})
// }

// func Test_Empty_Slice_Argument(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
// 	gtest.C(t, func(t *gtest.T) {
// 		result, err := db.GetAll(ctx, fmt.Sprintf(`select * from %s where id in(?)`, table), g.Slice{})
// 		t.AssertNil(err)
// 		t.Assert(len(result), 0)
// 	})
// }
