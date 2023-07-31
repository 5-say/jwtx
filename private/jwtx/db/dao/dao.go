package dao

import (
	"github.com/5-say/go-tools/tools/db"
	"github.com/5-say/zero-auth/private/jwtx/db/dao/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// commonDB ..
var commonDB *gorm.DB

// SetCommon .. 设置公共连接
func SetCommon(config db.Config, logWriter logger.Writer) {
	commonDB = db.Open(config, logWriter)
}

// Common .. 调用公共连接
func Common() *query.Query {
	db := commonDB.Session(&gorm.Session{NewDB: true})
	return query.Use(db)
}

// New .. 调用新连接
func New(config db.Config, logWriter logger.Writer) *query.Query {
	db := db.Open(config, logWriter)
	return query.Use(db)
}

// LockUpdate .. 排他锁
func LockUpdate(tx *gorm.DB) *gorm.DB {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"})
}
