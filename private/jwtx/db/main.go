package main

import (
	"flag"

	"github.com/5-say/go-tools/tools/db"
	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Config struct {
	DB db.Config
}

var configFile = flag.String("f", "../rpc/etc/config.yaml", "the config file")

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath: "./dao/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode

		FieldNullable:     true, // 数据库中的字段可为空，则生成 struct 字段为指针类型
		FieldCoverable:    true, // 如果数据库中字段有默认值，则生成指针类型的字段，以避免零值（zero-value）问题，见：https://gorm.io/docs/create.html#Default-Values
		FieldSignable:     true, // 检测整数字段的无符号类型，调整生成的数据类型
		FieldWithIndexTag: true, // 为结构体生成 gorm index tag，如：gorm:"index:idx_name"
		FieldWithTypeTag:  true, // 为结构体生成 gorm type tag，如：gorm:"type:varchar(12)"
	})

	// 数据类型映射 ..
	g.WithDataTypeMap(map[string]func(detailType gorm.ColumnType) (dataType string){
		"int": func(detailType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"tinyint": func(detailType gorm.ColumnType) (dataType string) {
			return "int8"
		},
	})

	// 设置数据库连接
	var config Config
	conf.MustLoad(*configFile, &config)
	config.DB.LogLevel = 2
	g.UseDB(db.Open(config.DB, db.LogWriter()))

	// 读取所有表
	g.ApplyBasic(g.GenerateAllTable()...)

	// // Generate basic type-safe DAO API for struct `model.User` following conventions
	// g.ApplyBasic(model.User{})

	// // Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	// 生成文件
	g.Execute()
}
