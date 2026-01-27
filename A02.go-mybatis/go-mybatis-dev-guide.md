# Go-MyBatis Comprehensive Developer Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Installation and Setup](#installation-and-setup)
3. [Core Concepts](#core-concepts)
4. [Quick Start Tutorial](#quick-start-tutorial)
5. [Database Configuration](#database-configuration)
6. [XML Mapper Files](#xml-mapper-files)
7. [Dynamic SQL](#dynamic-sql)
8. [CRUD Operations](#crud-operations)
9. [Template Tags](#template-tags)
10. [Optimistic Locking](#optimistic-locking)
11. [Logical Deletion](#logical-deletion)
12. [Transactions and AOP](#transactions-and-aop)
13. [Dynamic Data Sources](#dynamic-data-sources)
14. [Connection Pooling](#connection-pooling)
15. [Logging Configuration](#logging-configuration)
16. [XML Generation Tool](#xml-generation-tool)
17. [Best Practices](#best-practices)
18. [Production Deployment](#production-deployment)
19. [Troubleshooting](#troubleshooting)

---

## Introduction

Go-MyBatis is a SQL mapper ORM framework for Golang that brings MyBatis-like functionality to the Go ecosystem. It's a fork of zhuxiujia/GoMybatis and provides powerful features including:

- **High Performance**: Claims up to 456,621 TPS (transactions per second) in benchmark tests
- **Dynamic SQL**: Full support for MyBatis-style dynamic queries with 15+ XML tags
- **Transaction Management**: Declarative transactions with 8 propagation behaviors
- **Optimistic Locking**: Built-in support for version-based concurrency control
- **Logical Deletion**: Soft delete functionality for data recovery
- **Intelligent Expression Processing**: Complex expressions like `#{foo.Bar}`, `#{arg+1}`
- **Database Agnostic**: Support for MySQL, PostgreSQL, SQLite, Oracle, TiDB, CockroachDB, and more

### When to Use Go-MyBatis

**Use Go-MyBatis when:**
- You need complex dynamic SQL queries
- You're migrating Java/MyBatis projects to Go
- You want fine-grained control over SQL
- You need transaction propagation behaviors
- You prefer SQL-first approach over ORM abstractions

**Consider alternatives when:**
- You want a pure Go-idiomatic ORM (consider GORM, sqlx)
- You need simpler CRUD operations without complex queries
- You want compile-time SQL validation

---

## Installation and Setup

### Prerequisites

- Go 1.11+ (with Go modules support)
- Database driver for your target database

### Installation Steps

**1. Install Go-MyBatis**

```bash
go get github.com/zhuxiujia/GoMybatis
```

> **Note**: The repository `jba/go-mybatis` is a fork. Use `zhuxiujia/GoMybatis` for the most active development.

**2. Install Database Driver**

For MySQL:
```bash
go get github.com/go-sql-driver/mysql
```

For PostgreSQL:
```bash
go get github.com/lib/pq
```

For SQLite:
```bash
go get github.com/mattn/go-sqlite3
```

**3. Project Structure**

```
myproject/
├── main.go
├── mapper/
│   ├── user_mapper.xml
│   ├── order_mapper.xml
│   └── product_mapper.xml
├── model/
│   ├── user.go
│   ├── order.go
│   └── product.go
├── repository/
│   └── repository.go
└── config/
    └── database.go
```

---

## Core Concepts

### 1. Engine

The `GoMybatisEngine` is the core component that:
- Manages database connections
- Parses XML mapper files
- Generates proxy implementations for mapper interfaces
- Handles transaction management

### 2. Mapper

A mapper is a Go struct with function fields that correspond to SQL operations defined in XML files.

```go
type UserMapper struct {
    SelectById   func(id int64) (*User, error)
    Insert       func(user *User) (int64, error)
    Update       func(user *User) (int64, error)
    Delete       func(id int64) (int64, error)
}
```

### 3. XML Mapper Files

XML files contain SQL statements with dynamic elements:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <select id="SelectById">
        SELECT * FROM users WHERE id = #{id}
    </select>
</mapper>
```

### 4. ResultMap

ResultMaps define how database columns map to Go struct fields, especially important for:
- Different naming conventions (snake_case vs camelCase)
- Complex object relationships
- Custom type conversions

---

## Quick Start Tutorial

### Step 1: Define Your Model

**model/user.go**

```go
package model

import "time"

type User struct {
    Id         int64     `json:"id"`
    Username   string    `json:"username"`
    Email      string    `json:"email"`
    Password   string    `json:"password"`
    Status     int       `json:"status"`
    Version    int       `json:"version"`        // For optimistic locking
    CreateTime time.Time `json:"create_time"`
    UpdateTime time.Time `json:"update_time"`
    DeleteFlag int       `json:"delete_flag"`    // For logical deletion
}
```

### Step 2: Create XML Mapper

**mapper/user_mapper.xml**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap" tables="users">
        <id column="id" property="id"/>
        <result column="username" property="username" langType="string"/>
        <result column="email" property="email" langType="string"/>
        <result column="password" property="password" langType="string"/>
        <result column="status" property="status" langType="int"/>
        <result column="version" property="version" langType="int" version_enable="true"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="update_time" property="updateTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int" 
                logic_enable="true" logic_undelete="1" logic_deleted="0"/>
    </resultMap>
    
    <select id="SelectById" resultType="User">
        SELECT * FROM users WHERE id = #{id} AND delete_flag = 1
    </select>
    
    <select id="SelectByUsername" resultType="User">
        SELECT * FROM users WHERE username = #{username} AND delete_flag = 1
    </select>
    
    <insert id="Insert">
        INSERT INTO users(username, email, password, status, version, 
                         create_time, update_time, delete_flag)
        VALUES(#{username}, #{email}, #{password}, #{status}, 1, 
               NOW(), NOW(), 1)
    </insert>
    
    <update id="Update">
        UPDATE users 
        SET username = #{username},
            email = #{email},
            version = version + 1,
            update_time = NOW()
        WHERE id = #{id} AND version = #{version}
    </update>
    
    <delete id="Delete">
        UPDATE users SET delete_flag = 0 WHERE id = #{id}
    </delete>
</mapper>
```

### Step 3: Define Mapper Interface

**repository/user_repository.go**

```go
package repository

import (
    "myproject/model"
)

type UserMapper struct {
    SelectById       func(id int64) (*model.User, error)
    SelectByUsername func(username string) (*model.User, error)
    Insert           func(user *model.User) (int64, error)
    Update           func(user *model.User) (int64, error)
    Delete           func(id int64) (int64, error)
}
```

### Step 4: Initialize Engine and Load Mappers

**main.go**

```go
package main

import (
    "fmt"
    "io/ioutil"
    _ "github.com/go-sql-driver/mysql"
    "github.com/zhuxiujia/GoMybatis"
    "myproject/repository"
    "myproject/model"
)

func main() {
    // Initialize engine
    engine := GoMybatis.GoMybatisEngine{}.New()
    
    // Database connection string
    // Format: username:password@tcp(host:port)/database?params
    dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    
    err := engine.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    
    // Load XML mapper file
    xmlBytes, err := ioutil.ReadFile("mapper/user_mapper.xml")
    if err != nil {
        panic(err)
    }
    
    // Create mapper instance
    var userMapper repository.UserMapper
    
    // Bind XML to mapper
    engine.WriteMapperPtr(&userMapper, xmlBytes)
    
    // Use the mapper
    user := &model.User{
        Username: "john_doe",
        Email:    "john@example.com",
        Password: "hashed_password",
        Status:   1,
    }
    
    // Insert
    id, err := userMapper.Insert(user)
    if err != nil {
        fmt.Println("Insert error:", err)
        return
    }
    fmt.Println("Inserted user with ID:", id)
    
    // Select
    retrievedUser, err := userMapper.SelectById(id)
    if err != nil {
        fmt.Println("Select error:", err)
        return
    }
    fmt.Printf("Retrieved user: %+v\n", retrievedUser)
}
```

### Step 5: Run Your Application

```bash
go run main.go
```

---

## Database Configuration

### Connection String Formats

**MySQL**
```go
dsn := "username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local"
```

**PostgreSQL**
```go
dsn := "host=localhost port=5432 user=username password=password dbname=database sslmode=disable"
```

**SQLite**
```go
dsn := "./database.db"
```

**SQL Server**
```go
dsn := "sqlserver://username:password@host:port?database=dbname"
```

### Connection Parameters

**Important MySQL Parameters:**
- `charset=utf8mb4` - Support full Unicode including emojis
- `parseTime=True` - Parse DATE and DATETIME to time.Time
- `loc=Local` - Use local timezone
- `timeout=30s` - Connection timeout
- `readTimeout=30s` - Read timeout
- `writeTimeout=30s` - Write timeout

### Multiple Database Support

```go
// Database configuration
type DatabaseConfig struct {
    DriverName string
    DataSource string
    MaxOpen    int
    MaxIdle    int
}

var configs = map[string]DatabaseConfig{
    "mysql": {
        DriverName: "mysql",
        DataSource: "root:pass@tcp(localhost:3306)/db1?charset=utf8mb4",
        MaxOpen:    100,
        MaxIdle:    20,
    },
    "postgres": {
        DriverName: "postgres",
        DataSource: "host=localhost port=5432 user=user password=pass dbname=db2",
        MaxOpen:    50,
        MaxIdle:    10,
    },
}
```

---

## XML Mapper Files

### DTD Declaration

Always start your XML mapper with the DTD declaration:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
```

### ResultMap Configuration

**Basic ResultMap**

```xml
<resultMap id="UserResultMap" tables="users">
    <id column="id" property="id"/>
    <result column="username" property="username" langType="string"/>
    <result column="email" property="email" langType="string"/>
    <result column="created_at" property="createdAt" langType="time.Time"/>
</resultMap>
```

**ResultMap with Special Features**

```xml
<resultMap id="AdvancedResultMap" tables="products">
    <!-- Primary key -->
    <id column="id" property="id"/>
    
    <!-- Regular fields -->
    <result column="name" property="name" langType="string"/>
    <result column="price" property="price" langType="float64"/>
    
    <!-- Optimistic lock version -->
    <result column="version" property="version" langType="int" version_enable="true"/>
    
    <!-- Logical deletion -->
    <result column="delete_flag" property="deleteFlag" langType="int" 
            logic_enable="true" logic_undelete="1" logic_deleted="0"/>
    
    <!-- Timestamps -->
    <result column="create_time" property="createTime" langType="time.Time"/>
    <result column="update_time" property="updateTime" langType="time.Time"/>
</resultMap>
```

### Supported langType Values

| Go Type | langType Value |
|---------|---------------|
| string | string |
| int, int8, int16, int32, int64 | int |
| uint, uint8, uint16, uint32, uint64 | uint |
| float32, float64 | float64 |
| bool | bool |
| time.Time | time.Time |
| []byte | []byte |

### Parameter Binding

**Single Parameter**

```xml
<select id="SelectById">
    SELECT * FROM users WHERE id = #{id}
</select>
```

```go
SelectById func(id int64) (*User, error)
```

**Multiple Parameters with mapperParams Tag**

```xml
<select id="SelectByUsernameAndEmail">
    SELECT * FROM users 
    WHERE username = #{username} AND email = #{email}
</select>
```

```go
SelectByUsernameAndEmail func(username string, email string) (*User, error) 
    `mapperParams:"username,email"`
```

**Struct Parameter**

```xml
<insert id="Insert">
    INSERT INTO users(username, email, password)
    VALUES(#{username}, #{email}, #{password})
</insert>
```

```go
Insert func(user *User) (int64, error)
```

**Nested Properties**

```xml
<select id="SelectByAddress">
    SELECT * FROM users 
    WHERE city = #{address.city} AND country = #{address.country}
</select>
```

```go
type Address struct {
    City    string
    Country string
}

type User struct {
    Id      int64
    Address Address
}

SelectByAddress func(user *User) ([]User, error)
```

---

## Dynamic SQL

Go-MyBatis supports 15+ dynamic SQL tags inspired by MyBatis:

### 1. if - Conditional Logic

```xml
<select id="SearchUsers" resultType="User">
    SELECT * FROM users
    WHERE delete_flag = 1
    <if test="username != null and username != ''">
        AND username LIKE #{username + '%'}
    </if>
    <if test="email != null">
        AND email = #{email}
    </if>
    <if test="status != null">
        AND status = #{status}
    </if>
</select>
```

```go
type SearchCriteria struct {
    Username string
    Email    string
    Status   *int
}

SearchUsers func(criteria SearchCriteria) ([]User, error)
```

### 2. choose, when, otherwise - Switch Case

```xml
<select id="SearchByType" resultType="User">
    SELECT * FROM users
    WHERE delete_flag = 1
    <choose>
        <when test="searchType == 'username'">
            AND username = #{searchValue}
        </when>
        <when test="searchType == 'email'">
            AND email = #{searchValue}
        </when>
        <otherwise>
            AND (username LIKE #{searchValue + '%'} OR email LIKE #{searchValue + '%'})
        </otherwise>
    </choose>
</select>
```

```go
SearchByType func(searchType string, searchValue string) ([]User, error) 
    `mapperParams:"searchType,searchValue"`
```

### 3. where - Dynamic WHERE Clause

The `<where>` tag automatically:
- Removes leading AND/OR
- Only adds WHERE if content is not empty

```xml
<select id="DynamicSearch" resultType="User">
    SELECT * FROM users
    <where>
        delete_flag = 1
        <if test="username != null">
            AND username LIKE #{username + '%'}
        </if>
        <if test="minAge != null">
            AND age >= #{minAge}
        </if>
        <if test="maxAge != null">
            AND age <= #{maxAge}
        </if>
    </where>
</select>
```

### 4. set - Dynamic UPDATE Clause

The `<set>` tag automatically:
- Removes trailing commas
- Only adds SET if content is not empty

```xml
<update id="UpdateSelective">
    UPDATE users
    <set>
        <if test="username != null">
            username = #{username},
        </if>
        <if test="email != null">
            email = #{email},
        </if>
        <if test="status != null">
            status = #{status},
        </if>
        update_time = NOW()
    </set>
    WHERE id = #{id}
</update>
```

### 5. foreach - Batch Operations

**IN Clause**

```xml
<select id="SelectByIds" resultType="User">
    SELECT * FROM users
    WHERE id IN
    <foreach collection="ids" item="item" open="(" close=")" separator=",">
        #{item}
    </foreach>
</select>
```

```go
SelectByIds func(ids []int64) ([]User, error) `mapperParams:"ids"`
```

**Batch Insert**

```xml
<insert id="BatchInsert">
    INSERT INTO users(username, email, status)
    VALUES
    <foreach collection="users" item="user" separator=",">
        (#{user.username}, #{user.email}, #{user.status})
    </foreach>
</insert>
```

```go
BatchInsert func(users []User) (int64, error) `mapperParams:"users"`
```

**Foreach Attributes:**
- `collection` - Parameter name (array or slice)
- `item` - Variable name for current iteration
- `index` - Index variable (optional)
- `open` - Opening string
- `close` - Closing string
- `separator` - Separator between iterations

### 6. trim - Custom Prefix/Suffix Handling

```xml
<select id="DynamicTrimSearch" resultType="User">
    SELECT * FROM users
    <trim prefix="WHERE" prefixOverrides="AND |OR ">
        <if test="username != null">
            AND username = #{username}
        </if>
        <if test="email != null">
            AND email = #{email}
        </if>
    </trim>
</select>
```

### 7. bind - Create Variables

```xml
<select id="SearchWithPattern" resultType="User">
    <bind name="pattern" value="'%' + keyword + '%'" />
    SELECT * FROM users
    WHERE username LIKE #{pattern} OR email LIKE #{pattern}
</select>
```

### 8. sql and include - Reusable SQL Fragments

```xml
<sql id="userColumns">
    id, username, email, status, create_time, update_time
</sql>

<sql id="activeUserWhere">
    WHERE delete_flag = 1 AND status = 1
</sql>

<select id="SelectActiveUsers" resultType="User">
    SELECT <include refid="userColumns"/>
    FROM users
    <include refid="activeUserWhere"/>
</select>

<select id="CountActiveUsers">
    SELECT COUNT(*)
    FROM users
    <include refid="activeUserWhere"/>
</select>
```

### Complex Dynamic SQL Example

```xml
<select id="AdvancedSearch" resultType="User">
    SELECT DISTINCT u.*
    FROM users u
    <if test="includeOrders == true">
        LEFT JOIN orders o ON u.id = o.user_id
    </if>
    <where>
        u.delete_flag = 1
        <if test="username != null and username != ''">
            AND u.username LIKE #{username + '%'}
        </if>
        <if test="email != null and email != ''">
            AND u.email = #{email}
        </if>
        <if test="statusList != null and statusList.length > 0">
            AND u.status IN
            <foreach collection="statusList" item="status" open="(" close=")" separator=",">
                #{status}
            </foreach>
        </if>
        <if test="minAge != null or maxAge != null">
            <trim prefix="AND (" suffix=")" prefixOverrides="AND |OR ">
                <if test="minAge != null">
                    AND u.age >= #{minAge}
                </if>
                <if test="maxAge != null">
                    AND u.age <= #{maxAge}
                </if>
            </trim>
        </if>
        <choose>
            <when test="sortBy == 'username'">
                ORDER BY u.username
            </when>
            <when test="sortBy == 'email'">
                ORDER BY u.email
            </when>
            <otherwise>
                ORDER BY u.create_time DESC
            </otherwise>
        </choose>
    </where>
    <if test="limit != null and limit > 0">
        LIMIT #{limit}
        <if test="offset != null">
            OFFSET #{offset}
        </if>
    </if>
</select>
```

---

## CRUD Operations

### Create (Insert)

**Single Insert**

```xml
<insert id="Insert" useGeneratedKeys="true" keyProperty="id">
    INSERT INTO users(username, email, password, status, create_time)
    VALUES(#{username}, #{email}, #{password}, #{status}, NOW())
</insert>
```

```go
Insert func(user *User) (int64, error)

// Usage
user := &User{
    Username: "alice",
    Email:    "alice@example.com",
    Password: "hashed_pw",
    Status:   1,
}
id, err := userMapper.Insert(user)
// user.Id will be populated with auto-generated ID
```

**Batch Insert**

```xml
<insert id="BatchInsert">
    INSERT INTO users(username, email, status, create_time)
    VALUES
    <foreach collection="users" item="user" separator=",">
        (#{user.username}, #{user.email}, #{user.status}, NOW())
    </foreach>
</insert>
```

```go
BatchInsert func(users []User) (int64, error) `mapperParams:"users"`

// Usage
users := []User{
    {Username: "user1", Email: "user1@example.com", Status: 1},
    {Username: "user2", Email: "user2@example.com", Status: 1},
    {Username: "user3", Email: "user3@example.com", Status: 1},
}
affected, err := userMapper.BatchInsert(users)
```

### Read (Select)

**Select Single Record**

```xml
<select id="SelectById" resultType="User">
    SELECT * FROM users WHERE id = #{id} AND delete_flag = 1
</select>
```

```go
SelectById func(id int64) (*User, error)

// Usage
user, err := userMapper.SelectById(123)
if err != nil {
    // Handle error
}
if user == nil {
    // Not found
}
```

**Select Multiple Records**

```xml
<select id="SelectAll" resultType="User">
    SELECT * FROM users WHERE delete_flag = 1 ORDER BY create_time DESC
</select>
```

```go
SelectAll func() ([]User, error)

// Usage
users, err := userMapper.SelectAll()
```

**Pagination**

```xml
<select id="SelectPage" resultType="User">
    SELECT * FROM users
    WHERE delete_flag = 1
    <if test="keyword != null">
        AND (username LIKE #{keyword + '%'} OR email LIKE #{keyword + '%'})
    </if>
    ORDER BY create_time DESC
    LIMIT #{limit} OFFSET #{offset}
</select>
```

```go
type PageRequest struct {
    Keyword string
    Limit   int
    Offset  int
}

SelectPage func(req PageRequest) ([]User, error)

// Usage
page := PageRequest{
    Keyword: "john",
    Limit:   20,
    Offset:  40, // Page 3 (20 items per page)
}
users, err := userMapper.SelectPage(page)
```

### Update

**Full Update**

```xml
<update id="Update">
    UPDATE users
    SET username = #{username},
        email = #{email},
        status = #{status},
        version = version + 1,
        update_time = NOW()
    WHERE id = #{id} AND version = #{version}
</update>
```

```go
Update func(user *User) (int64, error)

// Usage
user.Username = "new_username"
user.Email = "new_email@example.com"
affected, err := userMapper.Update(user)
if affected == 0 {
    // Optimistic lock conflict or record not found
}
```

**Selective Update**

```xml
<update id="UpdateSelective">
    UPDATE users
    <set>
        <if test="username != null">
            username = #{username},
        </if>
        <if test="email != null">
            email = #{email},
        </if>
        <if test="status != null">
            status = #{status},
        </if>
        version = version + 1,
        update_time = NOW()
    </set>
    WHERE id = #{id} AND version = #{version}
</update>
```

### Delete

**Physical Delete (Hard Delete)**

```xml
<delete id="DeletePhysical">
    DELETE FROM users WHERE id = #{id}
</delete>
```

**Logical Delete (Soft Delete)**

```xml
<delete id="Delete">
    UPDATE users 
    SET delete_flag = 0, update_time = NOW()
    WHERE id = #{id} AND delete_flag = 1
</delete>
```

```go
Delete func(id int64) (int64, error)

// Usage
affected, err := userMapper.Delete(123)
```

---

## Template Tags

Template tags provide automatic CRUD generation based on resultMap configuration. This is one of the most powerful features of Go-MyBatis.

### Prerequisites

Template tags **require** a properly configured `resultMap`:

```xml
<resultMap id="BaseResultMap" tables="users">
    <id column="id" property="id"/>
    <result column="username" property="username" langType="string"/>
    <result column="email" property="email" langType="string"/>
    <result column="status" property="status" langType="int"/>
    <result column="version" property="version" langType="int" version_enable="true"/>
    <result column="create_time" property="createTime" langType="time.Time"/>
    <result column="update_time" property="updateTime" langType="time.Time"/>
    <result column="delete_flag" property="deleteFlag" langType="int" 
            logic_enable="true" logic_undelete="1" logic_deleted="0"/>
</resultMap>
```

### insertTemplete

Auto-generates INSERT statements.

**XML Configuration:**

```xml
<!-- Basic insert template -->
<insertTemplete/>

<!-- With custom columns -->
<insertTemplete columns="username,email,password,status"/>
```

**Go Mapper:**

```go
type UserMapper struct {
    // Single insert - returns affected rows
    InsertTemplete func(user User) (int64, error)
    
    // Batch insert - note the slice parameter
    InsertTempleteBatch func(users []User) (int64, error) `mapperParams:"users"`
}
```

**Usage Example:**

```go
// Single insert
user := User{
    Username: "bob",
    Email:    "bob@example.com",
    Status:   1,
}
affected, err := userMapper.InsertTemplete(user)

// Batch insert
users := []User{
    {Username: "alice", Email: "alice@example.com", Status: 1},
    {Username: "charlie", Email: "charlie@example.com", Status: 1},
}
affected, err := userMapper.InsertTempleteBatch(users)
```

**Features:**
- Automatically sets `delete_flag` to `logic_undelete` value (1)
- Skips null fields with `test="field != null"`
- Supports batch insertion

### selectTemplete

Auto-generates SELECT statements with WHERE conditions.

**XML Configuration:**

```xml
<!-- Select all (with logical delete filter) -->
<selectTemplete/>

<!-- Select with dynamic WHERE conditions -->
<selectTemplete wheres="username?username=#{username}"/>

<!-- Multiple conditions (comma-separated) -->
<selectTemplete wheres="username?username=#{username},status?status=#{status}"/>

<!-- Using null-safe expression (*?*) -->
<selectTemplete wheres="username?*?*,email?email=#{email}"/>
```

**Condition Syntax:**
- Format: `fieldName?condition`
- `*?*` = null expression (no condition if field is null)
- Multiple conditions separated by commas

**Go Mapper:**

```go
SelectTemplete func(username string, status int) ([]User, error) 
    `mapperParams:"username,status"`
```

**Generated SQL (conceptual):**

```sql
SELECT * FROM users
WHERE delete_flag = 1
  AND username = ?
  AND status = ?
```

### updateTemplete

Auto-generates UPDATE statements with optimistic locking.

**XML Configuration:**

```xml
<!-- Update specific fields with WHERE conditions -->
<updateTemplete 
    sets="username?username=#{username},email?email=#{email},status?status=#{status}"
    wheres="id?id=#{id}"/>
```

**Go Mapper:**

```go
UpdateTemplete func(user User) (int64, error)
```

**Features:**
- Automatically increments `version` field
- WHERE clause includes `version` for optimistic locking
- Only updates specified fields in `sets`
- Sets `update_time` automatically

**Generated SQL (conceptual):**

```sql
UPDATE users
SET username = ?,
    email = ?,
    status = ?,
    version = version + 1,
    update_time = NOW()
WHERE id = ? AND version = ?
```

### deleteTemplete

Auto-generates logical DELETE (soft delete) statements.

**XML Configuration:**

```xml
<!-- Delete by conditions -->
<deleteTemplete wheres="id?id=#{id}"/>

<!-- Delete by multiple conditions -->
<deleteTemplete wheres="username?username=#{username},status?status=#{status}"/>
```

**Go Mapper:**

```go
DeleteTemplete func(id int64) (int64, error) `mapperParams:"id"`
```

**Features:**
- Sets `delete_flag` to `logic_deleted` value (0)
- Not a physical delete - record remains in database
- Enables data recovery

**Generated SQL (conceptual):**

```sql
UPDATE users
SET delete_flag = 0,
    update_time = NOW()
WHERE id = ?
```

### Complete Template Tags Example

**XML:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap" tables="products">
        <id column="id" property="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="description" property="description" langType="string"/>
        <result column="price" property="price" langType="float64"/>
        <result column="stock" property="stock" langType="int"/>
        <result column="version" property="version" langType="int" version_enable="true"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="update_time" property="updateTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int" 
                logic_enable="true" logic_undelete="1" logic_deleted="0"/>
    </resultMap>
    
    <!-- Auto-generated CRUD -->
    <insertTemplete/>
    <selectTemplete wheres="name?name LIKE #{name + '%'},price?price <= #{price}"/>
    <updateTemplete 
        sets="name?name=#{name},price?price=#{price},stock?stock=#{stock}" 
        wheres="id?id=#{id}"/>
    <deleteTemplete wheres="id?id=#{id}"/>
</mapper>
```

**Go Code:**

```go
type Product struct {
    Id          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    Version     int       `json:"version"`
    CreateTime  time.Time `json:"create_time"`
    UpdateTime  time.Time `json:"update_time"`
    DeleteFlag  int       `json:"delete_flag"`
}

type ProductMapper struct {
    InsertTemplete       func(product Product) (int64, error)
    InsertTempleteBatch  func(products []Product) (int64, error) `mapperParams:"products"`
    SelectTemplete       func(name string, price float64) ([]Product, error) `mapperParams:"name,price"`
    UpdateTemplete       func(product Product) (int64, error)
    DeleteTemplete       func(id int64) (int64, error) `mapperParams:"id"`
}

// Usage
func main() {
    // ... engine initialization ...
    
    // Insert
    product := Product{
        Name:        "Laptop",
        Description: "High-performance laptop",
        Price:       999.99,
        Stock:       50,
    }
    id, err := productMapper.InsertTemplete(product)
    
    // Select by name and max price
    products, err := productMapper.SelectTemplete("Lap", 1500.00)
    
    // Update
    product.Id = id
    product.Price = 899.99
    product.Stock = 45
    affected, err := productMapper.UpdateTemplete(product)
    
    // Logical delete
    affected, err = productMapper.DeleteTemplete(id)
}
```

---

## Optimistic Locking

Optimistic locking prevents lost updates in concurrent scenarios without database locks.

### How It Works

1. Each record has a `version` field (integer)
2. When reading, the version is retrieved
3. When updating, the WHERE clause includes the current version
4. If another transaction updated the record, the version won't match and update fails

### Configuration

**1. Add version field to your struct:**

```go
type Product struct {
    Id      int64
    Name    string
    Price   float64
    Version int  // Optimistic lock version
}
```

**2. Configure in resultMap:**

```xml
<resultMap id="ProductResultMap" tables="products">
    <id column="id" property="id"/>
    <result column="name" property="name" langType="string"/>
    <result column="price" property="price" langType="float64"/>
    <result column="version" property="version" langType="int" version_enable="true"/>
</resultMap>
```

**3. Create update statement:**

```xml
<update id="UpdatePrice">
    UPDATE products
    SET price = #{price},
        version = version + 1
    WHERE id = #{id} AND version = #{version}
</update>
```

Or use template tag:

```xml
<updateTemplete sets="price?price=#{price}" wheres="id?id=#{id}"/>
```

### Usage Example

```go
// Concurrent update scenario
func UpdateProductPrice(mapper ProductMapper, productId int64, newPrice float64) error {
    // 1. Read current state
    product, err := mapper.SelectById(productId)
    if err != nil {
        return err
    }
    
    // 2. Modify
    product.Price = newPrice
    
    // 3. Update with version check
    affected, err := mapper.Update(product)
    if err != nil {
        return err
    }
    
    // 4. Check if update succeeded
    if affected == 0 {
        return errors.New("optimistic lock conflict: record was modified by another transaction")
    }
    
    return nil
}
```

### Handling Conflicts

**Strategy 1: Retry**

```go
func UpdateWithRetry(mapper ProductMapper, productId int64, newPrice float64, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        product, err := mapper.SelectById(productId)
        if err != nil {
            return err
        }
        
        product.Price = newPrice
        affected, err := mapper.Update(product)
        if err != nil {
            return err
        }
        
        if affected > 0 {
            return nil // Success
        }
        
        // Conflict detected, retry
        time.Sleep(time.Millisecond * 100)
    }
    
    return errors.New("max retries exceeded")
}
```

**Strategy 2: Notify User**

```go
affected, err := mapper.Update(product)
if affected == 0 {
    return errors.New("the record was modified by another user, please refresh and try again")
}
```

### Best Practices

1. **Always read before update**: Don't construct updates without reading current state
2. **Handle conflicts gracefully**: Use retry logic or user notification
3. **Version field type**: Use `int`, `int32`, or `int64`
4. **Initial version**: Set to 1 or 0 when inserting
5. **Don't skip version**: Always include version in UPDATE WHERE clause

---

## Logical Deletion

Logical deletion (soft delete) marks records as deleted without physically removing them from the database, enabling:
- Data recovery
- Audit trails
- Referential integrity

### Configuration

**1. Add delete_flag field:**

```go
type User struct {
    Id         int64
    Username   string
    DeleteFlag int  // 1=active, 0=deleted
}
```

**2. Configure in resultMap:**

```xml
<resultMap id="UserResultMap" tables="users">
    <id column="id" property="id"/>
    <result column="username" property="username" langType="string"/>
    <result column="delete_flag" property="deleteFlag" langType="int" 
            logic_enable="true" 
            logic_undelete="1" 
            logic_deleted="0"/>
</resultMap>
```

**Attributes:**
- `logic_enable="true"` - Enable logical deletion
- `logic_undelete="1"` - Value for active records
- `logic_deleted="0"` - Value for deleted records

### Automatic Behaviors

When using template tags with logical deletion enabled:

**insertTemplete** - Automatically sets delete_flag=1:
```sql
INSERT INTO users(username, delete_flag) VALUES(?, 1)
```

**selectTemplete** - Automatically filters deleted records:
```sql
SELECT * FROM users WHERE delete_flag = 1 AND username = ?
```

**deleteTemplete** - Updates instead of deletes:
```sql
UPDATE users SET delete_flag = 0 WHERE id = ?
```

### Manual SQL with Logical Deletion

**Soft Delete:**

```xml
<delete id="Delete">
    UPDATE users 
    SET delete_flag = 0, 
        delete_time = NOW()
    WHERE id = #{id} AND delete_flag = 1
</delete>
```

**Select Active Only:**

```xml
<select id="SelectActive" resultType="User">
    SELECT * FROM users 
    WHERE delete_flag = 1
    ORDER BY create_time DESC
</select>
```

**Select Including Deleted:**

```xml
<select id="SelectAll" resultType="User">
    SELECT * FROM users
    ORDER BY create_time DESC
</select>
```

**Restore Deleted Record:**

```xml
<update id="Restore">
    UPDATE users 
    SET delete_flag = 1,
        restore_time = NOW()
    WHERE id = #{id} AND delete_flag = 0
</update>
```

**Permanent Delete:**

```xml
<delete id="PermanentDelete">
    DELETE FROM users WHERE id = #{id}
</delete>
```

### Usage Example

```go
type UserMapper struct {
    Insert          func(user User) (int64, error)
    SelectById      func(id int64) (*User, error)
    SelectAll       func() ([]User, error)
    SelectDeleted   func() ([]User, error)
    Delete          func(id int64) (int64, error)
    Restore         func(id int64) (int64, error)
    PermanentDelete func(id int64) (int64, error)
}

// Soft delete
affected, err := userMapper.Delete(userId)

// Verify it's gone from normal queries
user, err := userMapper.SelectById(userId) // Returns nil

// Restore
affected, err = userMapper.Restore(userId)

// Now it's back
user, err = userMapper.SelectById(userId) // Returns user

// Permanent delete (use with caution)
affected, err = userMapper.PermanentDelete(userId)
```

### Database Schema

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    delete_flag TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 0=deleted',
    delete_time DATETIME NULL COMMENT 'When deleted',
    restore_time DATETIME NULL COMMENT 'When restored',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_delete_flag (delete_flag),
    INDEX idx_username (username, delete_flag)
);
```

---

## Transactions and AOP

Go-MyBatis provides declarative transaction management with 8 transaction propagation behaviors, similar to Spring Framework.

### Transaction Propagation Types

| Type | Description |
|------|-------------|
| `PROPAGATION_REQUIRED` | Use existing transaction or create new one (default) |
| `PROPAGATION_SUPPORTS` | Use existing transaction or execute without transaction |
| `PROPAGATION_MANDATORY` | Use existing transaction or throw error |
| `PROPAGATION_REQUIRES_NEW` | Always create new transaction, suspend existing |
| `PROPAGATION_NOT_SUPPORTED` | Execute without transaction, suspend existing |
| `PROPAGATION_NEVER` | Execute without transaction or throw error if exists |
| `PROPAGATION_NESTED` | Execute within nested transaction if exists |
| `PROPAGATION_NOT_REQUIRED` | Create new transaction or throw error if exists |

### Basic Transaction Usage

**1. Define Service with Transaction Tags**

```go
type UserService struct {
    userMapper *UserMapper
    
    // Transaction method - rolls back on error or panic
    CreateUser func(username, email string) (*User, error) `tx:"" rollback:"error"`
    
    // Transaction with specific propagation
    UpdateUserEmail func(userId int64, email string) error `tx:"PROPAGATION_REQUIRED" rollback:"error"`
    
    // Nested transaction
    ProcessOrder func(userId int64, amount float64) error `tx:"PROPAGATION_NESTED" rollback:"error"`
}
```

**Tag Attributes:**
- `tx:""` or `tx:"PROPAGATION_REQUIRED"` - Transaction propagation behavior
- `rollback:"error"` - Rollback if method returns non-nil error
- Can also rollback on `panic`

**2. Implement Service Methods**

```go
func NewUserService(mapper *UserMapper) *UserService {
    return &UserService{
        userMapper: mapper,
        
        CreateUser: func(username, email string) (*User, error) {
            // Validate
            if username == "" || email == "" {
                return nil, errors.New("username and email required")
            }
            
            // Check if exists
            existing, _ := mapper.SelectByUsername(username)
            if existing != nil {
                return nil, errors.New("username already exists")
            }
            
            // Create user
            user := &User{
                Username: username,
                Email:    email,
                Status:   1,
            }
            
            _, err := mapper.Insert(user)
            if err != nil {
                return nil, err // Transaction will rollback
            }
            
            return user, nil
        },
        
        UpdateUserEmail: func(userId int64, email string) error {
            user, err := mapper.SelectById(userId)
            if err != nil {
                return err
            }
            if user == nil {
                return errors.New("user not found")
            }
            
            user.Email = email
            _, err = mapper.Update(user)
            return err
        },
    }
}
```

**3. Register Service with AOP Proxy**

```go
func main() {
    // Initialize engine
    engine := GoMybatis.GoMybatisEngine{}.New()
    engine.Open("mysql", dsn)
    
    // Load mapper
    var userMapper UserMapper
    engine.WriteMapperPtr(&userMapper, xmlBytes)
    
    // Create service
    userService := NewUserService(&userMapper)
    
    // Register with AOP proxy (enables transaction management)
    GoMybatis.AopProxyService(userService, &engine)
    
    // Use service - transactions are automatic
    user, err := userService.CreateUser("alice", "alice@example.com")
    if err != nil {
        log.Println("Transaction rolled back:", err)
    } else {
        log.Println("Transaction committed, user created:", user.Id)
    }
}
```

### Nested Transactions Example

```go
type OrderService struct {
    orderMapper   *OrderMapper
    productMapper *ProductMapper
    userMapper    *UserMapper
    
    // Parent transaction
    PlaceOrder func(userId int64, productId int64, quantity int) error `tx:"PROPAGATION_REQUIRED" rollback:"error"`
    
    // Child transactions
    DeductStock func(productId int64, quantity int) error `tx:"PROPAGATION_NESTED" rollback:"error"`
    UpdateUserCredit func(userId int64, amount float64) error `tx:"PROPAGATION_NESTED" rollback:"error"`
}

func NewOrderService(orderMapper *OrderMapper, productMapper *ProductMapper, userMapper *UserMapper) *OrderService {
    service := &OrderService{
        orderMapper:   orderMapper,
        productMapper: productMapper,
        userMapper:    userMapper,
    }
    
    service.PlaceOrder = func(userId int64, productId int64, quantity int) error {
        // This starts the parent transaction
        
        // Child transaction 1: Deduct stock
        err := service.DeductStock(productId, quantity)
        if err != nil {
            return err // Rolls back parent transaction
        }
        
        // Child transaction 2: Update user credit
        err = service.UpdateUserCredit(userId, 100.0)
        if err != nil {
            return err // Rolls back parent transaction
        }
        
        // Create order
        order := &Order{
            UserId:    userId,
            ProductId: productId,
            Quantity:  quantity,
        }
        _, err = orderMapper.Insert(order)
        if err != nil {
            return err // Rolls back everything
        }
        
        return nil // Commits all nested transactions
    }
    
    service.DeductStock = func(productId int64, quantity int) error {
        product, err := productMapper.SelectById(productId)
        if err != nil {
            return err
        }
        
        if product.Stock < quantity {
            return errors.New("insufficient stock")
        }
        
        product.Stock -= quantity
        _, err = productMapper.Update(product)
        return err
    }
    
    service.UpdateUserCredit = func(userId int64, amount float64) error {
        user, err := userMapper.SelectById(userId)
        if err != nil {
            return err
        }
        
        user.Credit += amount
        _, err = userMapper.Update(user)
        return err
    }
    
    return service
}
```

### Transaction Rollback Scenarios

**Rollback on Error:**
```go
ProcessPayment func(orderId int64) error `tx:"" rollback:"error"`

// Implementation
ProcessPayment: func(orderId int64) error {
    // ... some operations ...
    if paymentFailed {
        return errors.New("payment failed") // Triggers rollback
    }
    return nil // Commits transaction
}
```

**Rollback on Panic:**
```go
ProcessOrder func(orderId int64) error `tx:"" rollback:"error"`

// Implementation
ProcessOrder: func(orderId int64) error {
    // ... some operations ...
    if criticalError {
        panic("critical error") // Triggers rollback
    }
    return nil
}
```

### Best Practices

1. **Keep transactions short**: Long transactions lock resources
2. **Avoid external I/O in transactions**: Don't call external APIs within transactions
3. **Use appropriate propagation**: Choose based on business requirements
4. **Handle errors explicitly**: Return errors to trigger rollbacks
5. **Test rollback scenarios**: Ensure rollbacks work correctly
6. **Use nested transactions carefully**: Can impact performance

---

## Dynamic Data Sources

Dynamic data sources allow routing database operations to different databases at runtime based on custom logic.

### Configuration

**1. Open multiple database connections:**

```go
// Primary database
engine := GoMybatis.GoMybatisEngine{}.New()
err := engine.Open("mysql", "root:pass@tcp(primary:3306)/db1?charset=utf8mb4")

// Secondary databases
GoMybatis.Open("mysql", "root:pass@tcp(secondary:3306)/db2?charset=utf8mb4")
GoMybatis.Open("postgres", "host=tertiary port=5432 user=user password=pass dbname=db3")
```

**2. Create a router:**

```go
router := GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
    // Route based on mapper package name
    if strings.Contains(mapperName, "user.") {
        dsn := "root:pass@tcp(user-db:3306)/userdb?charset=utf8mb4"
        return &dsn
    }
    
    if strings.Contains(mapperName, "order.") {
        dsn := "root:pass@tcp(order-db:3306)/orderdb?charset=utf8mb4"
        return &dsn
    }
    
    if strings.Contains(mapperName, "analytics.") {
        dsn := "host=analytics-db port=5432 user=user dbname=analyticsdb"
        return &dsn
    }
    
    // Default: return nil to use primary database
    return nil
})

// Set the router
engine.SetDataSourceRouter(router)
```

### Routing Strategies

**Strategy 1: By Mapper Package**

```go
router := GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
    switch {
    case strings.HasPrefix(mapperName, "repository.user."):
        dsn := userDatabaseDSN
        return &dsn
    case strings.HasPrefix(mapperName, "repository.product."):
        dsn := productDatabaseDSN
        return &dsn
    default:
        return nil // Use default
    }
})
```

**Strategy 2: Read/Write Splitting**

```go
type ReadWriteRouter struct {
    writeDSN string
    readDSNs []string
    current  int
    mu       sync.Mutex
}

func (r *ReadWriteRouter) Route(mapperName string) *string {
    // Writes go to master
    if strings.Contains(mapperName, "Insert") || 
       strings.Contains(mapperName, "Update") || 
       strings.Contains(mapperName, "Delete") {
        return &r.writeDSN
    }
    
    // Reads go to replicas (round-robin)
    r.mu.Lock()
    defer r.mu.Unlock()
    dsn := r.readDSNs[r.current]
    r.current = (r.current + 1) % len(r.readDSNs)
    return &dsn
}

// Usage
rwRouter := &ReadWriteRouter{
    writeDSN: "root:pass@tcp(master:3306)/db",
    readDSNs: []string{
        "root:pass@tcp(replica1:3306)/db",
        "root:pass@tcp(replica2:3306)/db",
        "root:pass@tcp(replica3:3306)/db",
    },
}

router := GoMybatis.GoMybatisDataSourceRouter{}.New(rwRouter.Route)
```

**Strategy 3: Tenant-Based (Multi-Tenancy)**

```go
type TenantRouter struct {
    tenantDBMap map[string]string
    mu          sync.RWMutex
}

func (r *TenantRouter) SetTenant(tenantId string, dsn string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.tenantDBMap[tenantId] = dsn
}

func (r *TenantRouter) Route(mapperName string) *string {
    // Extract tenant from context or thread-local storage
    tenantId := GetCurrentTenantId() // Your implementation
    
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    if dsn, exists := r.tenantDBMap[tenantId]; exists {
        return &dsn
    }
    
    return nil // Default database
}

// Usage
tenantRouter := &TenantRouter{
    tenantDBMap: make(map[string]string),
}
tenantRouter.SetTenant("tenant1", "root:pass@tcp(tenant1-db:3306)/db")
tenantRouter.SetTenant("tenant2", "root:pass@tcp(tenant2-db:3306)/db")

router := GoMybatis.GoMybatisDataSourceRouter{}.New(tenantRouter.Route)
```

**Strategy 4: Time-Based (Sharding)**

```go
router := GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
    now := time.Now()
    year := now.Year()
    
    // Route to year-specific database
    dsn := fmt.Sprintf("root:pass@tcp(db-server:3306)/orders_%d", year)
    return &dsn
})
```

### Complete Example

```go
package main

import (
    "fmt"
    "strings"
    _ "github.com/go-sql-driver/mysql"
    "github.com/zhuxiujia/GoMybatis"
)

func main() {
    // Initialize engine with primary database
    engine := GoMybatis.GoMybatisEngine{}.New()
    primaryDSN := "root:password@tcp(primary:3306)/maindb?charset=utf8mb4&parseTime=True"
    err := engine.Open("mysql", primaryDSN)
    if err != nil {
        panic(err)
    }
    
    // Register additional databases
    userDSN := "root:password@tcp(user-server:3306)/userdb?charset=utf8mb4&parseTime=True"
    orderDSN := "root:password@tcp(order-server:3306)/orderdb?charset=utf8mb4&parseTime=True"
    
    GoMybatis.Open("mysql", userDSN)
    GoMybatis.Open("mysql", orderDSN)
    
    // Create router
    router := GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
        fmt.Printf("Routing mapper: %s\n", mapperName)
        
        if strings.Contains(mapperName, "UserMapper") {
            fmt.Println("-> Routing to user database")
            return &userDSN
        }
        
        if strings.Contains(mapperName, "OrderMapper") {
            fmt.Println("-> Routing to order database")
            return &orderDSN
        }
        
        fmt.Println("-> Using primary database")
        return nil
    })
    
    engine.SetDataSourceRouter(router)
    
    // Load mappers
    var userMapper UserMapper
    var orderMapper OrderMapper
    
    engine.WriteMapperPtr(&userMapper, userXML)
    engine.WriteMapperPtr(&orderMapper, orderXML)
    
    // Operations automatically route to correct databases
    user, _ := userMapper.SelectById(1)      // Routes to userdb
    order, _ := orderMapper.SelectById(100)  // Routes to orderdb
}
```

---

## Connection Pooling

Connection pooling is critical for performance in production. Go-MyBatis uses the standard `database/sql` connection pool.

### Configuration

```go
package main

import (
    "database/sql"
    "time"
    _ "github.com/go-sql-driver/mysql"
    "github.com/zhuxiujia/GoMybatis"
)

func main() {
    engine := GoMybatis.GoMybatisEngine{}.New()
    
    // Open database
    dsn := "root:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
    err := engine.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    
    // Get underlying *sql.DB
    db := engine.GormDB() // Or engine.GetDB() depending on version
    
    // Configure connection pool
    db.SetMaxOpenConns(100)        // Maximum open connections
    db.SetMaxIdleConns(20)         // Maximum idle connections
    db.SetConnMaxLifetime(time.Hour)      // Maximum connection lifetime
    db.SetConnMaxIdleTime(time.Minute * 10) // Maximum idle time
}
```

### Pool Parameters Explained

**SetMaxOpenConns(n int)**
- Maximum number of open connections to the database
- Includes both idle and in-use connections
- Default: unlimited (0)
- Recommended: 100-200 for typical applications

**SetMaxIdleConns(n int)**
- Maximum number of idle connections in the pool
- Should be less than MaxOpenConns
- Default: 2
- Recommended: 20-50 for high-traffic applications

**SetConnMaxLifetime(d time.Duration)**
- Maximum lifetime of a connection
- After this time, connection is closed and recreated
- Helps with database-side connection limits
- Recommended: 1 hour

**SetConnMaxIdleTime(d time.Duration)**
- Maximum idle time before closing connection
- Helps free up resources during low traffic
- Recommended: 10 minutes

### Tuning Guidelines

**Low Traffic Application:**
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour * 2)
db.SetConnMaxIdleTime(time.Minute * 15)
```

**Medium Traffic Application:**
```go
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(20)
db.SetConnMaxLifetime(time.Hour)
db.SetConnMaxIdleTime(time.Minute * 10)
```

**High Traffic Application:**
```go
db.SetMaxOpenConns(200)
db.SetMaxIdleConns(50)
db.SetConnMaxLifetime(time.Minute * 30)
db.SetConnMaxIdleTime(time.Minute * 5)
```

### Monitoring Connection Pool

```go
package main

import (
    "fmt"
    "time"
)

func monitorConnectionPool(db *sql.DB) {
    ticker := time.NewTicker(time.Second * 10)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := db.Stats()
        
        fmt.Printf("Connection Pool Stats:\n")
        fmt.Printf("  Open Connections: %d\n", stats.OpenConnections)
        fmt.Printf("  In Use: %d\n", stats.InUse)
        fmt.Printf("  Idle: %d\n", stats.Idle)
        fmt.Printf("  Wait Count: %d\n", stats.WaitCount)
        fmt.Printf("  Wait Duration: %s\n", stats.WaitDuration)
        fmt.Printf("  Max Idle Closed: %d\n", stats.MaxIdleClosed)
        fmt.Printf("  Max Lifetime Closed: %d\n", stats.MaxLifetimeClosed)
        fmt.Println("---")
    }
}

// Usage
go monitorConnectionPool(db)
```

### Best Practices

1. **Set MaxOpenConns**: Never use unlimited connections
2. **Monitor pool stats**: Track WaitCount and WaitDuration
3. **Balance with database limits**: Don't exceed database max_connections
4. **Consider application scaling**: MaxOpenConns * num_instances < database max_connections
5. **Tune based on metrics**: Adjust based on actual usage patterns

---

## Logging Configuration

### Enable Logging

```go
engine := GoMybatis.GoMybatisEngine{}.New()
engine.SetLogEnable(true)
```

### Custom Log Implementation

```go
import (
    "log"
    "os"
    "github.com/zhuxiujia/GoMybatis"
)

// Custom logger that implements GoMybatis.Log interface
type CustomLogger struct {
    logger *log.Logger
}

func (l *CustomLogger) Println(messages []byte) {
    l.logger.Println(string(messages))
}

func main() {
    engine := GoMybatis.GoMybatisEngine{}.New()
    
    // Create custom logger
    file, _ := os.OpenFile("sql.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    customLogger := &CustomLogger{
        logger: log.New(file, "[SQL] ", log.LstdFlags),
    }
    
    // Set custom logger
    engine.SetLogEnable(true)
    engine.SetLog(customLogger)
    
    // Or use standard logger
    engine.SetLog(&GoMybatis.LogStandard{
        PrintlnFunc: func(messages []byte) {
            log.Println(string(messages))
        },
    })
}
```

### Structured Logging Example

```go
import (
    "encoding/json"
    "time"
)

type StructuredLogger struct {
    file *os.File
}

type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    SQL       string `json:"sql"`
    Duration  string `json:"duration"`
}

func (l *StructuredLogger) Println(messages []byte) {
    entry := LogEntry{
        Timestamp: time.Now().Format(time.RFC3339),
        Level:     "INFO",
        SQL:       string(messages),
        Duration:  "0ms", // Parse from messages if available
    }
    
    jsonBytes, _ := json.Marshal(entry)
    l.file.Write(append(jsonBytes, '\n'))
}

// Usage
logFile, _ := os.OpenFile("sql_structured.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
structuredLogger := &StructuredLogger{file: logFile}
engine.SetLog(structuredLogger)
```

---

## XML Generation Tool

Automatically generate XML mapper files from Go structs.

### Step 1: Annotate Your Struct

```go
package model

import "time"

type UserAddress struct {
    Id            int64     `json:"id" gm:"id"`              // Primary key
    UserId        int64     `json:"user_id"`
    RealName      string    `json:"real_name"`
    Phone         string    `json:"phone"`
    AddressDetail string    `json:"address_detail"`
    Version       int       `json:"version" gm:"version"`    // Optimistic lock
    CreateTime    time.Time `json:"create_time"`
    DeleteFlag    int       `json:"delete_flag" gm:"logic"`  // Logical deletion
}
```

**Annotations:**
- `gm:"id"` - Marks as primary key
- `gm:"version"` - Marks as optimistic lock version field
- `gm:"logic"` - Marks as logical deletion flag

### Step 2: Create Generation Tool

**tools/xml_generator.go**

```go
package main

import (
    "fmt"
    "reflect"
    "github.com/zhuxiujia/GoMybatis"
    "myproject/model"
)

func main() {
    // Define your structs
    structs := []interface{}{
        model.UserAddress{},
        model.Order{},
        model.Product{},
    }
    
    // Generate XML for each struct
    for _, bean := range structs {
        typeName := reflect.TypeOf(bean).Name()
        tableName := "biz_" + GoMybatis.StructToSnakeString(bean)
        
        xml := GoMybatis.CreateXml(tableName, bean)
        filename := typeName + "Mapper.xml"
        
        err := GoMybatis.OutPutXml(filename, xml)
        if err != nil {
            fmt.Printf("Error generating %s: %v\n", filename, err)
        } else {
            fmt.Printf("Generated: %s\n", filename)
        }
    }
}
```

### Step 3: Run Generator

```bash
cd tools
go run xml_generator.go
```

### Generated XML Example

**UserAddressMapper.xml**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable Logical Delete Fields-->
    <!--logic_deleted Logically delete deleted fields-->
    <!--logic_undelete Logically Delete Undeleted Fields-->
    <!--version_enable Optimistic lock version field, support int, int8, int16, int32, Int64-->
    <resultMap id="BaseResultMap" tables="biz_user_address">
        <id column="id" property="id"/>
        <result column="id" property="id" langType="int64" />
        <result column="user_id" property="user_id" langType="int64" />
        <result column="real_name" property="real_name" langType="string" />
        <result column="phone" property="phone" langType="string" />
        <result column="address_detail" property="address_detail" langType="string" />
        <result column="version" property="version" langType="int" version_enable="true" />
        <result column="create_time" property="create_time" langType="time.Time" />
        <result column="delete_flag" property="delete_flag" langType="int" logic_enable="true" logic_undelete="1" logic_deleted="0" />
    </resultMap>
</mapper>
```

### Customization

After generation, you can add custom queries:

```xml
<mapper>
    <!-- Auto-generated resultMap -->
    <resultMap id="BaseResultMap" tables="biz_user_address">
        <!-- ... -->
    </resultMap>
    
    <!-- Add your custom queries -->
    <select id="SelectByUserId" resultType="UserAddress">
        SELECT * FROM biz_user_address 
        WHERE user_id = #{userId} AND delete_flag = 1
    </select>
    
    <select id="SearchByCity" resultType="UserAddress">
        SELECT * FROM biz_user_address
        WHERE address_detail LIKE #{city + '%'} AND delete_flag = 1
    </select>
    
    <!-- Use template tags -->
    <insertTemplete/>
    <selectTemplete wheres="user_id?user_id=#{userId}"/>
    <updateTemplete sets="real_name?real_name=#{realName},phone?phone=#{phone}" wheres="id?id=#{id}"/>
    <deleteTemplete wheres="id?id=#{id}"/>
</mapper>
```

---

## Best Practices

### 1. Project Structure

```
myproject/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── database.go
│   ├── model/
│   │   ├── user.go
│   │   ├── order.go
│   │   └── product.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── order_repository.go
│   │   └── product_repository.go
│   └── service/
│       ├── user_service.go
│       └── order_service.go
├── mapper/
│   ├── user_mapper.xml
│   ├── order_mapper.xml
│   └── product_mapper.xml
├── migrations/
│   ├── 001_create_users.sql
│   └── 002_create_orders.sql
└── go.mod
```

### 2. XML File Organization

**Development:** Store XML files in project directory for IDE support

```
mapper/
├── user_mapper.xml
├── order_mapper.xml
└── product_mapper.xml
```

**Production:** Embed XML files using `embed` package (Go 1.16+)

```go
package main

import (
    _ "embed"
    "github.com/zhuxiujia/GoMybatis"
)

//go:embed mapper/user_mapper.xml
var userMapperXML []byte

//go:embed mapper/order_mapper.xml
var orderMapperXML []byte

func main() {
    engine := GoMybatis.GoMybatisEngine{}.New()
    // ... setup ...
    
    engine.WriteMapperPtr(&userMapper, userMapperXML)
    engine.WriteMapperPtr(&orderMapper, orderMapperXML)
}
```

### 3. Error Handling

```go
// Check for specific errors
user, err := userMapper.SelectById(id)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    return nil, fmt.Errorf("database error: %w", err)
}

// Check affected rows
affected, err := userMapper.Delete(id)
if err != nil {
    return err
}
if affected == 0 {
    return errors.New("user not found or already deleted")
}
```

### 4. Repository Pattern

```go
package repository

type UserRepository interface {
    GetById(id int64) (*model.User, error)
    GetByUsername(username string) (*model.User, error)
    Create(user *model.User) error
    Update(user *model.User) error
    Delete(id int64) error
    List(criteria SearchCriteria) ([]model.User, error)
}

type userRepositoryImpl struct {
    mapper UserMapper
}

func NewUserRepository(mapper UserMapper) UserRepository {
    return &userRepositoryImpl{mapper: mapper}
}

func (r *userRepositoryImpl) GetById(id int64) (*model.User, error) {
    user, err := r.mapper.SelectById(id)
    if err == sql.ErrNoRows {
        return nil, nil // Not found
    }
    return user, err
}

// ... other methods ...
```

### 5. Use Context for Cancellation

```go
// Custom mapper with context
type UserMapper struct {
    SelectByIdWithContext func(ctx context.Context, id int64) (*User, error)
}

// Usage
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := userMapper.SelectByIdWithContext(ctx, userId)
```

### 6. Avoid N+1 Queries

**Bad:**
```go
users, _ := userMapper.SelectAll()
for _, user := range users {
    orders, _ := orderMapper.SelectByUserId(user.Id) // N queries
    user.Orders = orders
}
```

**Good:**
```xml
<select id="SelectUsersWithOrders" resultType="User">
    SELECT u.*, o.*
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    WHERE u.delete_flag = 1
</select>
```

### 7. Validate Input

```go
func (s *UserService) CreateUser(username, email string) error {
    // Validate
    if username == "" {
        return errors.New("username required")
    }
    if !isValidEmail(email) {
        return errors.New("invalid email format")
    }
    if len(username) < 3 || len(username) > 50 {
        return errors.New("username must be 3-50 characters")
    }
    
    // Proceed with creation
    // ...
}
```

### 8. Use Prepared Statements

Go-MyBatis automatically uses prepared statements with `#{}` syntax. Never use `${}` with user input:

**Safe:**
```xml
<select id="SelectByUsername">
    SELECT * FROM users WHERE username = #{username}
</select>
```

**Unsafe (SQL Injection):**
```xml
<select id="SelectByUsername">
    SELECT * FROM users WHERE username = '${username}'
</select>
```

### 9. Index Optimization

```sql
-- Add indexes for frequently queried columns
CREATE INDEX idx_username ON users(username);
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_status_delete_flag ON users(status, delete_flag);

-- Composite index for common queries
CREATE INDEX idx_user_search ON users(username, email, status) WHERE delete_flag = 1;
```

### 10. Testing

```go
package repository_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
    // Setup test database
    engine := setupTestEngine(t)
    defer engine.Close()
    
    // Create mapper
    var mapper UserMapper
    engine.WriteMapperPtr(&mapper, xmlBytes)
    
    // Test
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    err := mapper.Insert(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.Id)
    
    // Verify
    retrieved, err := mapper.SelectById(user.Id)
    assert.NoError(t, err)
    assert.Equal(t, user.Username, retrieved.Username)
}
```

---

## Production Deployment

### 1. Database Connection Configuration

```go
package config

import (
    "fmt"
    "os"
    "strconv"
)

type DatabaseConfig struct {
    Driver          string
    Host            string
    Port            int
    Database        string
    Username        string
    Password        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime int // seconds
}

func LoadDatabaseConfig() DatabaseConfig {
    port, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
    maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
    maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "20"))
    maxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "3600"))
    
    return DatabaseConfig{
        Driver:          getEnv("DB_DRIVER", "mysql"),
        Host:            getEnv("DB_HOST", "localhost"),
        Port:            port,
        Database:        getEnv("DB_NAME", "myapp"),
        Username:        getEnv("DB_USER", "root"),
        Password:        getEnv("DB_PASSWORD", ""),
        MaxOpenConns:    maxOpen,
        MaxIdleConns:    maxIdle,
        ConnMaxLifetime: maxLifetime,
    }
}

func (c DatabaseConfig) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        c.Username, c.Password, c.Host, c.Port, c.Database)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 2. Graceful Shutdown

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Initialize engine
    engine := GoMybatis.GoMybatisEngine{}.New()
    // ... setup ...
    
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Start server
    go startServer()
    
    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutting down...")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Close database connections
    if db := engine.GormDB(); db != nil {
        db.Close()
    }
    
    log.Println("Shutdown complete")
}
```

### 3. Health Checks

```go
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    // Check database connection
    if err := db.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "unhealthy",
            "error":  err.Error(),
        })
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
}
```

### 4. Monitoring and Metrics

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    dbQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_queries_total",
            Help: "Total number of database queries",
        },
        []string{"mapper", "operation"},
    )
    
    dbQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "db_query_duration_seconds",
            Help: "Duration of database queries",
        },
        []string{"mapper", "operation"},
    )
)

// Wrap mapper calls with metrics
func (r *userRepository) GetById(id int64) (*User, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        dbQueriesTotal.WithLabelValues("UserMapper", "SelectById").Inc()
        dbQueryDuration.WithLabelValues("UserMapper", "SelectById").Observe(duration)
    }()
    
    return r.mapper.SelectById(id)
}
```

### 5. Docker Deployment

**Dockerfile:**

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/mapper ./mapper

EXPOSE 8080
CMD ["./server"]
```

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=myapp
      - DB_USER=appuser
      - DB_PASSWORD=apppass
      - DB_MAX_OPEN_CONNS=100
      - DB_MAX_IDLE_CONNS=20
    depends_on:
      - mysql
    restart: unless-stopped
  
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=rootpass
      - MYSQL_DATABASE=myapp
      - MYSQL_USER=appuser
      - MYSQL_PASSWORD=apppass
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

### 6. Environment-Specific Configuration

**.env.development:**
```
DB_HOST=localhost
DB_PORT=3306
DB_NAME=myapp_dev
DB_USER=dev_user
DB_PASSWORD=dev_pass
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
LOG_LEVEL=debug
```

**.env.production:**
```
DB_HOST=prod-db-cluster.us-east-1.rds.amazonaws.com
DB_PORT=3306
DB_NAME=myapp_prod
DB_USER=prod_user
DB_PASSWORD=${VAULT_DB_PASSWORD}
DB_MAX_OPEN_CONNS=200
DB_MAX_IDLE_CONNS=50
LOG_LEVEL=info
```

---

## Troubleshooting

### Common Issues and Solutions

#### 1. XML Parsing Errors

**Problem:** `XML parse error: element not found`

**Solution:**
- Verify DTD declaration at top of XML file
- Check for unclosed tags
- Ensure proper XML structure
- Validate against MyBatis DTD

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!-- Your mappings -->
</mapper>
```

#### 2. Parameter Binding Errors

**Problem:** `parameter 'xxx' not found`

**Solution:**
Add `mapperParams` tag to function signature:

```go
// Wrong
SelectByUsername func(username string) (*User, error)

// Correct
SelectByUsername func(username string) (*User, error) `mapperParams:"username"`
```

#### 3. Connection Pool Exhaustion

**Problem:** `Error: too many connections`

**Solution:**
```go
// Increase pool size or reduce connection lifetime
db.SetMaxOpenConns(200)  // Increase max connections
db.SetConnMaxLifetime(time.Minute * 30)  // Reduce lifetime
```

#### 4. Optimistic Lock Conflicts

**Problem:** Update returns 0 affected rows

**Solution:**
```go
// Implement retry logic
for i := 0; i < 3; i++ {
    user, _ := mapper.SelectById(id)
    user.Name = newName
    affected, err := mapper.Update(user)
    if affected > 0 {
        break // Success
    }
    time.Sleep(time.Millisecond * 100)
}
```

#### 5. Type Conversion Errors

**Problem:** `cannot convert type`

**Solution:**
Ensure `langType` in resultMap matches Go struct field type:

```xml
<!-- Go field: CreateTime time.Time -->
<result column="create_time" property="createTime" langType="time.Time"/>

<!-- Go field: Status int -->
<result column="status" property="status" langType="int"/>
```

#### 6. Transaction Not Rolling Back

**Problem:** Transaction commits even on error

**Solution:**
- Ensure `rollback:"error"` tag is present
- Function must return error type
- Call `AopProxyService` to register service

```go
type UserService struct {
    CreateUser func(user *User) error `tx:"" rollback:"error"`
}

GoMybatis.AopProxyService(&userService, &engine)
```

#### 7. Nil Pointer Errors

**Problem:** `panic: runtime error: invalid memory address`

**Solution:**
- Check if `WriteMapperPtr` was called
- Verify XML file was loaded
- Ensure engine is initialized

```go
// Initialize engine
engine := GoMybatis.GoMybatisEngine{}.New()
err := engine.Open("mysql", dsn)

// Load XML
var mapper UserMapper
engine.WriteMapperPtr(&mapper, xmlBytes)

// Now safe to use
user, err := mapper.SelectById(1)
```

---

## Conclusion

Go-MyBatis brings the power and flexibility of MyBatis to the Go ecosystem. This guide covered:

- ✅ Installation and setup
- ✅ Core concepts and quick start
- ✅ XML mapper configuration
- ✅ Dynamic SQL with 15+ tags
- ✅ Template tags for auto-generated CRUD
- ✅ Optimistic locking and logical deletion
- ✅ Transaction management with 8 propagation behaviors
- ✅ Dynamic data sources and routing strategies
- ✅ Connection pooling and performance tuning
- ✅ Logging and monitoring
- ✅ XML generation tools
- ✅ Best practices and production deployment

### Next Steps

1. **Explore Examples**: Check [GoMybatis examples](https://github.com/zhuxiujia/GoMybatis/tree/master/example)
2. **Read Documentation**: Visit [official docs](https://zhuxiujia.github.io/gomybatis.io/)
3. **Join Community**: Contribute or ask questions on GitHub
4. **Benchmark**: Test performance in your specific use case

### Additional Resources

- **GitHub Repository**: https://github.com/zhuxiujia/GoMybatis
- **Documentation**: https://zhuxiujia.github.io/gomybatis.io/
- **MyBatis Reference**: https://mybatis.org/mybatis-3/
- **Go database/sql**: https://pkg.go.dev/database/sql

---

**Last Updated**: January 2026  
**Go-MyBatis Version**: Latest (check GitHub for current version)  
**License**: Apache-2.0
