package mysql

import (
	"crawler/internal/repository"
	"crawler/pkg/config"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewDB 创建并初始化数据库连接
func NewDB(dbConfig config.MySQLConfig) (*sqlx.DB, error) {
	if err := ensureDatabaseExists(dbConfig); err != nil {
		return nil, fmt.Errorf("确保数据库存在失败: %w", err)
	}

	dbConn, err := sqlx.Connect("mysql", buildConnectionString(dbConfig))
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	configureConnectionPool(dbConn, dbConfig)

	// 初始化数据库表
	if err := initializeDatabaseTables(dbConn); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("初始化数据库表失败: %w", err)
	}

	return dbConn, nil
}

// ensureDatabaseExists 确保数据库存在，不存在则创建
func ensureDatabaseExists(dbConfig config.MySQLConfig) error {
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
	)

	rootConn, err := sqlx.Connect("mysql", rootDSN)
	if err != nil {
		return fmt.Errorf("连接MySQL根用户失败: %w", err)
	}
	defer rootConn.Close()

	createDBQuery := fmt.Sprintf(repository.CreateDatabaseSQL, dbConfig.Database)
	if _, err := rootConn.Exec(createDBQuery); err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}

	return nil
}

// initializeDatabaseTables 初始化数据库表结构
func initializeDatabaseTables(dbConn *sqlx.DB) error {
	if _, err := dbConn.Exec(repository.CreateArticlesTableSQL); err != nil {
		return fmt.Errorf("创建文章表失败: %w", err)
	}
	return nil
}

// buildConnectionString 构建数据库连接字符串
func buildConnectionString(dbConfig config.MySQLConfig) string {
	connParams := make([]string, 0, 8) // 预分配合适的容量

	// appendParam 添加连接参数
	appendParam := func(paramName string, paramValue interface{}) {
		if paramValue != nil {
			connParams = append(connParams, fmt.Sprintf("%s=%v", paramName, paramValue))
		}
	}

	// 添加必要的连接参数
	appendParam("charset", dbConfig.Charset)
	appendParam("parseTime", dbConfig.ParseTime)
	appendParam("loc", "Local")
	appendParam("timeout", dbConfig.Timeout)
	appendParam("readTimeout", dbConfig.ReadTimeout)
	appendParam("writeTimeout", dbConfig.WriteTimeout)
	appendParam("collation", dbConfig.Collation)

	// 构建完整的连接字符串
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		strings.Join(connParams, "&"),
	)
}

// configureConnectionPool 配置数据库连接池
func configureConnectionPool(dbConn *sqlx.DB, dbConfig config.MySQLConfig) {
	dbConn.SetMaxOpenConns(dbConfig.MaxOpenConns)
	dbConn.SetMaxIdleConns(dbConfig.MaxIdleConns)
	dbConn.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	dbConn.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)
}
