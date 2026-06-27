
基于 https://github.com/talen8/FlecBlog.git 二次开发。

爱来自爱吃鱼的猫：https://blog.talen.top

### 1. 环境要求

| 依赖 | 版本要求 |
|------|---------|
| Go | >= 1.25 |
| Node.js | >= 22 |
| PostgreSQL | >= 15 |
| pnpm | 推荐 |

### 2. 配置环境变量

**根目录 `.env`**

```env
DB_PASSWORD=your_db_password
JWT_SECRET=your_jwt_secret
API_URL=http://localhost:8080/api/v1
```

**`server/.env`**

```env
SERVER_PORT=8080
SERVER_ALLOW_ORIGINS=*
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name  #数据库名
DB_USER=your_db_user  #数据库用户名
DB_PASSWORD=your_db_password  #数据库密码
JWT_SECRET=your_jwt_secret   #改为安全的长字符
```

**`blog/.env`**

```env
API_URL=http://127.0.0.1:8080/api/v1
```

**`admin/.env`**

```env
VITE_API_URL=http://127.0.0.1:8080/api/v1
```

### 3. 启动后端

```bash
cd server
go run cmd/main.go
```

后端启动后会自动执行数据库迁移，首次启动会初始化所有表结构。

### 4. 启动前台博客

```bash
cd blog
pnpm install
pnpm dev
```

默认运行在 `http://localhost:3000`。

### 5. 启动管理后台

```bash
cd admin
pnpm install
pnpm dev
```

默认运行在 `http://localhost:4000`。

---

## 默认管理员账号

| 字段 | 值 |
|------|-----|
| 邮箱 | `admin@example.com` |
| 密码 | `123456` |

# Markdown语法以及插件化语法
请见：https://blog.seln.cn/posts/hc1wfyr9
