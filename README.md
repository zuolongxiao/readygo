# ReadyGo

## 简介

Golang API 快速开发脚手架, 适合小型后台管理系统的快速开发

## 特性

- 基于 Gin + GORM 构建
- 遵循 RESTful API 设计规范
- JWT 认证
- RBAC 权限控制

## 开发环境部署

1. 获取代码

```bash
git clone https://github.com/zuolongxiao/readygo.git
```

1. 下载依赖

```bash
# 设置国内代理，可选操作
go env -w GOPROXY=https://goproxy.cn,direct

cd readygo
go mod tidy
```

1. 修改配置文件（或者通过环境变量设置）

```bash
cp config.sample.yaml config.yaml
vi config.yaml
# 修改以下配置项:
# JWT.Secret: "your secret"

# Database.Host: 127.0.0.1
# Database.Port: 3306
# Database.Password: your password
# Database.Name = readygo
```

1. 创建数据库

```bash
mysql -h127.0.0.1 -P3306 -uroot -p -e "create database readygo"
```

1. 系统初始设置

```bash
# 迁移数据库
go run main.go admin migrate

# 导入权限
go run main.go admin permission

# 创建超级用户,不指定-p则生成随机密码
go run main.go admin create -u root -p yourpassword
```

1. 启动服务

```bash
go run main.go serve
```

1. 测试

```bash
# 获取 JWT token, 后续请求要携带这个 token
curl --request POST 'http://127.0.0.1:8000/api/auth' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "root",
    "password": "password"
}'

# 获取用户列表
curl --request GET 'http://127.0.0.1:8000/api/v1/admins' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 完整 API 列表请查看 routing/routes/v1/sys.go 文件
```

## 开发流程(以 tag 的 CRUD 为例)

1. 创建 model

```go
// 完整代码请查看: models/tag.go

// 定义数据库字段和 go 的数据类型映射，和用于迁移的字段标签
type Tag struct {
    Base

    Name  string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
    State string `gorm:"type:char(1);default:N;not null"`
}

// API 返回的 JSON 需要显示的字段
type TagView struct {
    BaseView

    Name  string `json:"name"`
    State string `json:"state"`
}

// 创建 tag 时可提交的字段以及验证规则，额外的字段会被过滤掉
// 使用 github.com/go-playground/validator/v10 作为验证引擎
type TagCreate struct {
    Name  string `json:"name" binding:"required,max=50"`
    State string `json:"state" binding:"omitempty,oneof=N Y"`
}

// 更新 tag 时可提交的字段以及验证规则，额外的字段会被过滤掉
type TagUpdate struct {
    Name  string `json:"name" binding:"max=50"`
    State string `json:"state" binding:"omitempty,oneof=N Y"`
}

// 保存到数据库前执行的钩子，用于做入库前的最后校验
// 例如：验证 tag 名称是否已存在，如果已存在，则返回一个错误，并终止执行
// 否则返回 nil
// 其他钩子请查看 gorm 文档
func (m *Tag) BeforeSave(tx *gorm.DB) error {
    var count int64
    if err := tx.Model(m).Where("id <> ? AND name = ?", m.ID, m.Name).Limit(1).Count(&count).Error; err != nil {
        return errs.DBError(err.Error())
    }

    if count > 0 {
        return errs.DuplicatedError("tag.name")
    }

    return nil
}

// 根据URL的查询字符串来构建SQL查询语句
// URL查询字符串还支持以下公共参数：
// size：获取的记录数，默认20，例如：size=50
// sort：排序规则，默认 sort=+id，根据id升序排序，可选 sort=-id，根据id降序排序
// offset：偏移量，用于遍历查询结果，默认0，当返回的结果中offset为0时，说明已经完成遍历，否则在下次请求中带上offset
// 例如URL: GET /api/v1/tags?size=5&sort=-id&offset=5&name=tag&state=N
// 将产生SQL语句：SELECT `g_tag`.`id`,`g_tag`.`created_at`,`g_tag`.`updated_at`,`g_tag`.`name`,`g_tag`.`state` FROM `g_tag` WHERE id < 5 AND name LIKE 'tag%' AND state = 'N' ORDER BY id DESC LIMIT 5
func (*Tag) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
    if name := c.Query("name"); name != "" {
        db = db.Where("name LIKE ?", name+"%")
    }

    if state := c.Query("state"); state != "" {
        db = db.Where("state = ?", state)
    }

    return db
}
```

1. 创建 API

```go
// 完整代码请查看: api/v1/tag.go

// 获取列表
func ListTags(c *gin.Context) {
    w := utils.NewContextWrapper(c)
    s := services.New(&models.Tag{})

    var tags []models.TagView
    if err := s.Find(&tags, c); err != nil {
        w.Respond(err, nil)
        return
    }

    data := map[string]interface{}{
        "list":   tags,
        "offset": s.GetOffset(),
    }

    w.Respond(nil, data)
}

// 根据id获取指定记录
func GetTag(c *gin.Context) {
    w := utils.NewContextWrapper(c)
    s := services.New(&models.Tag{})

    var tag models.TagView
    if err := s.GetByID(&tag, c.Param("id")); err != nil {
        w.Respond(err, nil)
        return
    }

    w.Respond(nil, tag)
}

// 新建记录
func CreateTag(c *gin.Context) {
    w := utils.NewContextWrapper(c)

    // 将用户提交的 JSON 绑定到结构体并进行验证
    binding := models.TagCreate{}
    if err := w.Bind(&binding); err != nil {
        w.Respond(err, nil)
        return
    }

    // 用绑定后的结构体对 model 进行填充
    m := models.Tag{}
    s := services.New(&m)
    if err := s.Fill(&binding); err != nil {
        w.Respond(err, nil)
        return
    }

    // 对 model 设置额外信息，并创建一条记录
    m.CreatedBy = w.GetUsername()
    if err := s.Create(); err != nil {
        w.Respond(err, nil)
        return
    }

    data := map[string]interface{}{
        "id":         m.ID,
        "created_at": m.CreatedAt.Time,
    }

    w.Respond(nil, data)
}

// 更新记录
func UpdateTag(c *gin.Context) {
    w := utils.NewContextWrapper(c)

    // 将用户提交的 JSON 绑定到结构体并进行验证
    binding := models.TagUpdate{}
    if err := w.Bind(&binding); err != nil {
        w.Respond(err, nil)
        return
    }

    // 根据id从数据库获取记录并填充到 model
    m := models.Tag{}
    s := services.New(&m)
    if err := s.LoadByID(c.Param("id")); err != nil {
        w.Respond(err, nil)
        return
    }

    // 使用绑定后的结构体覆盖掉 model 中对应的字段，使 model 得到更新
    if err := s.Fill(&binding); err != nil {
        w.Respond(err, nil)
        return
    }

    // 对 model 设置额外信息，并更新到数据库
    m.UpdatedBy = w.GetUsername()
    if err := s.Save(); err != nil {
        w.Respond(err, nil)
        return
    }

    data := map[string]interface{}{
        "id":         m.ID,
        "updated_at": m.UpdatedAt.Time,
    }

    w.Respond(nil, data)
}

// 删除记录
func DeleteTag(c *gin.Context) {
    w := utils.NewContextWrapper(c)

    // 查询记录是否存在
    s := services.New(&models.Tag{})
    if err := s.LoadByID(c.Param("id")); err != nil {
        w.Respond(err, nil)
        return
    }

    // 执行删除操作
    if err := s.Delete(); err != nil {
        w.Respond(err, nil)
        return
    }

    w.Respond(nil, nil)
}
```

1. 添加路由

```go
// 完整代码请查看: routing/routes/v1/tag.go

// 定义路由规则
var tag = []routes.Route{
    {
        // 每条路由对应一个权限，Handler名称即为权限名称，
        // 关于权限请看后面纤细解释
        Method:  "GET", // HTTP方法
        Pattern: "/tags/:id", // 路由匹配模式
        Handler: apiv1.GetTag, // 处理函数，对应上一步创建的API函数
        Flag:    "Y", // 导入标识，""或"-":不导入数据库，"Y":导入数据库并启用权限，"N":导入数据库并禁用权限
        Desc:    "Get tag", // 描述
    },
    {
        Method:  "GET",
        Pattern: "/tags",
        Handler: apiv1.ListTags,
        Flag:    "Y",
        Desc:    "List tags",
    },
    {
        Method:  "POST",
        Pattern: "/tags",
        Handler: apiv1.CreateTag,
        Flag:    "Y",
        Desc:    "Create tag",
    },
    {
        Method:  "PUT",
        Pattern: "/tags/:id",
        Handler: apiv1.UpdateTag,
        Flag:    "Y",
        Desc:    "Update tag",
    },
    {
        Method:  "DELETE",
        Pattern: "/tags/:id",
        Handler: apiv1.DeleteTag,
        Flag:    "Y",
        Desc:    "Delete tag",
    },
}

// 添加路由
func init() {
    Routes = append(Routes, tag...)
}
```

1. 迁移表

```go
// 完整代码请查看: commands/migration.go

// 将 &models.Tag{}, 添加到最后一个参数
func RunMigration() {
    models.DB.AutoMigrate(
        &models.Admin{},
        &models.Authorization{},
        &models.Permission{},
        &models.Role{},
        &models.Tag{},
    )
}

// 然后命令行执行: go run main.go migration:run
```

1. 导入权限,可选操作

```bash
# 命令行执行
go run main.go permission:load
```

1. 测试

> 文件修改后，需重启HTTP服务，使修改生效

```bash
# 新建记录
curl --request POST 'http://127.0.0.1:8000/api/v1/tags' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "test",
    "state": "Y"
}'

# 获取列表
curl --request GET 'http://127.0.0.1:8000/api/v1/tags' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 获取指定记录
curl --request GET 'http://127.0.0.1:8000/api/v1/tags/1' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 更新记录
curl --request PUT 'http://127.0.0.1:8000/api/v1/tags/1' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "test-updated",
    "state": "N"
}'

# 删除记录
curl --request DELETE 'http://127.0.0.1:8000/api/v1/tags/1' \
--header 'Authorization: Bearer <JWT TOKEN>'
```

## 权限说明

- 每个API对应一个权限，API的函数名即为权限名，权限需要导入到数据库并且`is_enabled`设置为`Y`才生效，否则所有角色拥有该API权限。
- 权限可以通过命令行导入或者通过API来管理。

### 用户权限管理流程

1. 新建角色
2. 将角色和权限关联
3. 设置用户角色，则用户拥有该角色的权限，ID为1的用户为超级管理员，超级管理员拥有所有权限，不用设置角色

### 权限管理相关API

1. 权限管理

```bash
# 获取权限列表
curl --request GET 'http://127.0.0.1:8000/api/v1/permissions' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 新建权限，通常使通过命令行导入
# name: 权限名, 需要和API函数名对应，否则不起作用
# title: 权限标题，用于显示
# group: 权限分组，便于查询
# note: 备注
# is_enabled: 是否启用
curl --request POST 'http://127.0.0.1:8000/api/v1/permissions' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "ListTags",
    "title": "List tags",
    "group": "",
    "note": "",
    "is_enabled": "Y"
}'

# 修改权限
curl --request PUT 'http://127.0.0.1:8000/api/v1/permissions/<permissionID>' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "ListTags",
    "title": "List tags",
    "group": "tag",
    "note": "list tags note",
    "is_enabled": "N"
}'

# 删除权限
curl --request DELETE 'http://127.0.0.1:8000/api/v1/permissions/<permissionID>' \
--header 'Authorization: Bearer <JWT TOKEN>'
```

1. 角色管理

```bash
# 获取角色列表
curl --request GET 'http://127.0.0.1:8000/api/v1/roles' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 新建角色
# name: 角色名称
curl --request POST 'http://127.0.0.1:8000/api/v1/roles' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "admin"
}'

# 修改角色
curl --request PUT 'http://127.0.0.1:8000/api/v1/roles/<roleID>' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Administrators"
}'

# 删除角色
curl --request DELETE 'http://127.0.0.1:8000/api/v1/roles/<roleID>' \
--header 'Authorization: Bearer <JWT TOKEN>'
```

1. 用户管理

```bash
# 获取用户列表
curl --request GET 'http://127.0.0.1:8000/api/v1/admins' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 新建用户
curl --request POST 'http://127.0.0.1:8000/api/v1/admins' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "role_id": 1,
    "username": "admin",
    "password": "admin",
    "is_locked": "N"
}'

# 修改用户
curl --request PUT 'http://127.0.0.1:8000/api/v1/admins/<adminID>' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "role_id": 1,
    "username": "admin",
    "password": "admin",
    "is_locked": "N"
}'

# 删除用户
curl --request DELETE 'http://127.0.0.1:8000/api/v1/admins/<adminID>' \
--header 'Authorization: Bearer <JWT TOKEN>'
```

1. 角色与权限关联

```bash
# 获取指定角色已关联的权限
curl --request GET 'http://127.0.0.1:8000/api/v1/roles/<roleID>/permissions' \
--header 'Authorization: Bearer <JWT TOKEN>'

# 为指定角色关联一个权限
curl --request POST 'http://127.0.0.1:8000/api/v1/roles/<roleID>/permissions' \
--header 'Authorization: Bearer <JWT TOKEN>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "permission_id": <permissionID>
}'

# 移除指定角色下的一个权限
curl --request DELETE 'http://127.0.0.1:8000/api/v1/roles/<roleID>/permissions/<permissionID>' \
--header 'Authorization: Bearer <JWT TOKEN>'
```
