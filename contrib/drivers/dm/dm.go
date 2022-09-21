package dm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "gitee.com/chunanyong/dm"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type DriverDM struct {
	*gdb.Core
}

var (
	tableFieldsMap = gmap.New(true)
)

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"dm"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

func New() gdb.Driver {
	return &DriverDM{}
}

func (d *DriverDM) New(core *gdb.Core, node gdb.ConfigNode) (gdb.DB, error) {
	return &DriverDM{
		Core: core,
	}, nil
}

func (d *DriverDM) Open(config gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "dm"
	)
	// dm://userName:password@ip:port/dbname
	if config.Link != "" {
		source = config.Link
		// Custom changing the schema in runtime.
		// if config.Name != "" {
		// source, _ = gregex.ReplaceString(`/([\w\.\-]+)+`, "/"+config.Name, source)
		// }
	} else {
		source = fmt.Sprintf(
			"dm://%s:%s@%s:%s/%s?charset=%s",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
		// if config.Timezone != "" {
		// source = fmt.Sprintf("%s&loc=%s", source, url.QueryEscape(config.Timezone))
		// }
	}
	g.Dump("DriverDM.Open()::source", source)
	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`dm.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

func (d *DriverDM) GetChars() (charLeft string, charRight string) {
	return `"`, `"`
}

func (d *DriverDM) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	// TODO support multiple schema
	if len(schema) == 0 {
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Schema is empty`)
		// schema = []string{"DP"}
	}
	// select * from all_tables where owner = 'DP';
	result, err = d.DoSelect(ctx, link, fmt.Sprintf(`SELECT * FROM ALL_TABLES WHERE OWNER IN ('%s')`, schema[0]))
	if err != nil {
		return
	}

	for _, m := range result {
		if v, ok := m["IOT_NAME"]; ok {
			tables = append(tables, v.String())
		}
	}
	g.Dump("DriverDM.Tables()::tables", tables)
	return
}

func (d *DriverDM) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	// Format dm table
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			"function TableFields supports only single table operations",
		)
	}

	// SET schema
	useSchema := d.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}

	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`dm_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result gdb.Result
				link   gdb.Link
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			// g.Dump("s:", strings.ToUpper(d.QuoteWord(table)))
			result, err = d.DoSelect(
				ctx, link,
				// select * from all_tab_columns where owner='DP' and Table_Name='T_SYS_LOG'
				fmt.Sprintf(`SELECT * FROM ALL_TAB_COLUMNS WHERE OWNER='%s' AND Table_Name= '%s'`, useSchema, strings.ToUpper(table)),
			)
			if err != nil {
				return nil
			}
			fields = make(map[string]*gdb.TableField)
			for _, m := range result {
				// m[NULLABLE] returns "N" "Y"
				// "N" means not null
				// "Y" means could be  null
				var nullable bool
				if m["NULLABLE"].String() != "N" {
					nullable = true
				}
				fields[m["COLUMN_NAME"].String()] = &gdb.TableField{
					Index: m["COLUMN_ID"].Int(),
					Name:  m["COLUMN_NAME"].String(),
					Type:  m["DATA_TYPE"].String(),
					Null:  nullable,
					// Key:     m["Key"].String(),
					Default: m["DATA_DEFAULT"].Val(),
					// Extra:   m["Extra"].String(),
					// Comment: m["Comment"].String(),
				}
			}
			// g.Dump("DriverDM.TableFields()::fields", fields)
			return fields
		},
	)
	if v != nil {
		fields = v.(map[string]*gdb.TableField)
	}
	return
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *DriverDM) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	defer func() {
		newSql, newArgs, err = d.Core.DoFilter(ctx, link, newSql, newArgs)
	}()
	// var index int
	// Convert placeholder char '?' to string "@px".
	// str, _ := gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
	// index++
	// return fmt.Sprintf("@p%d", index)
	// })
	// g.Dump("sql:", sql)
	str, _ := gregex.ReplaceString("\"", "", sql)
	str, _ = gregex.ReplaceString("\n", "", str)
	str, _ = gregex.ReplaceString("\t", "", str)
	// g.Dump("str:", str)

	newSql = strings.ToUpper(str)
	g.Dump("DriverDM.DoFilter()::newSql", newSql)
	newArgs = args
	g.Dump("DriverDM.DoFilter()::newArgs", newArgs)

	return
}

func (d *DriverDM) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	// Save option only on duplicate key : ID
	case gdb.InsertOptionSave:
		g.Dump("===========================list===========================", list)
		listLength := len(list)
		if listLength == 0 {
			return nil, gerror.NewCode(gcode.CodeNotSupported, `Save operation list is empty by dm driver`)
		}
		var (
			keys           []string
			keysWithTable  []string
			keysWithAssign []string
			selvalues      []string
			values         []string
		)
		charL, charR := d.GetChars()
		valuecharL, valuecharR := "'", "'"
		for k := range list[0] {
			keys = append(keys, charL+k+charR)
			keysWithTable = append(keysWithTable, "T2."+charL+k+charR)
			keysWithAssign = append(keysWithAssign, fmt.Sprintf(`T1.%s = T2.%s`, charL+k+charR, charL+k+charR))
		}
		for _, column := range keys {
			fmt.Println("===========================")
			fmt.Println(list[0])
			g.Dump(list[0])
			m := list[0]
			fmt.Println("column:", column)
			fmt.Println("list[0][column]:", m[column])
			if m[column] == nil {
				fmt.Println("list[0][column]:", list[0][column])
				continue
			}
			va := reflect.ValueOf(list[0][column])
			ty := reflect.TypeOf(list[0][column])
			d := ""
			switch ty.Kind() {
			case reflect.String:
				d = va.String()
			case reflect.Int:
				d = strconv.FormatInt(va.Int(), 10)
			case reflect.Int64:
				d = strconv.FormatInt(va.Int(), 10)
			default:
				fmt.Println("default")
			}
			selvalues = append(selvalues, fmt.Sprintf(valuecharL+"%s"+valuecharR+" AS "+charL+"%s"+charR, d, column))
		}
		fmt.Println(selvalues)
		for _, mapper := range list[1:] {
			var element []string
			for _, column := range keys {
				if mapper[column] == nil {
					continue
				}
				va := reflect.ValueOf(mapper[column])
				ty := reflect.TypeOf(mapper[column])
				switch ty.Kind() {
				case reflect.String:
					element = append(element, valuecharL+va.String()+valuecharR)
				case reflect.Int:
					element = append(element, strconv.FormatInt(va.Int(), 10))
				case reflect.Int64:
					element = append(element, strconv.FormatInt(va.Int(), 10))
				}
			}
			values = append(values, fmt.Sprintf(`UNION ALL SELECT %s FROM DUAL`, strings.Join(element, ",")))
		}

		var (
			batchResult      = new(gdb.SqlResult)
			selectValues     = strings.Join(selvalues, ",")
			sqlValues        = strings.Join(values, " ")
			keyStr           = strings.Join(keys, ",")
			keyStrWithTable  = strings.Join(keysWithTable, ",")
			keyStrWithAssign = strings.Join(keysWithAssign, ",")
		)
		sqlStr := fmt.Sprintf(`
MERGE INTO %s T1 USING (SELECT %s FROM DUAL %s) T2 ON (T1.ID = T2.ID) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s) WHEN MATCHED THEN UPDATE SET %s; 
COMMIT;
`, table, selectValues, sqlValues, keyStr, keyStrWithTable, keyStrWithAssign)
		g.Dump("===========================sqlStr===========================", sqlStr)
		r, err := d.DoExec(ctx, link, sqlStr)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected += n
		}
		return batchResult, nil
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}
