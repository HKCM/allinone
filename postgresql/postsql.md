# postgresql

## 安装PostgreSQL

```shell
sudo apt-get update
sudo apt-get install postgresql postgresql-client
```

PostgreSQL 安装完成后默认是已经启动的,但是也可以通过下面的方式来手动启动服务。
```shell
sudo /etc/init.d/postgresql start   # 开启
sudo /etc/init.d/postgresql stop    # 关闭
sudo /etc/init.d/postgresql restart # 重启
```

连接PostgreSQL
```shell
psql -U postgres -h localhost -p 5432
```

## 用户相关

创建数据库用户dbuser,并设置密码。
```sql
CREATE USER dbuser WITH PASSWORD 'password';
```

创建用户数据库,这里为exampledb,并指定所有者为dbuser。
```sql
CREATE DATABASE exampledb OWNER dbuser;
```

将exampledb数据库的所有权限都赋予dbuser,否则dbuser只能登录控制台,没有任何数据库操作权限。
```sql
GRANT ALL PRIVILEGES ON DATABASE exampledb to dbuser;
```

## 操作

创建数据库
```sql
CREATE DATABASE dbname;
```
用`\l`查看数据库

选择数据库
使用 `\c + database` 来进入数据库

创建表
```sql
CREATE TABLE COMPANY(
   ID INT PRIMARY KEY     NOT NULL,
   NAME           TEXT    NOT NULL,
   AGE            INT     NOT NULL,
   ADDRESS        CHAR(50),
   SALARY         REAL,
   JOIN_DATE      DATE
);
```

用`\d`查看数据表
用`\d tablename` 查看表格信息

插入数据
```sql
INSERT INTO COMPANY (ID,NAME,AGE,ADDRESS,SALARY,JOIN_DATE) VALUES (1, 'Paul', 32, 'California', 20000.00,'2001-07-13');

-- 使用默认值
INSERT INTO COMPANY (ID,NAME,AGE,ADDRESS,SALARY,JOIN_DATE) VALUES (3, 'Teddy', 23, 'Norway', 20000.00, DEFAULT );

-- 插入多行数据
INSERT INTO COMPANY (ID,NAME,AGE,ADDRESS,SALARY,JOIN_DATE) VALUES (4, 'Mark', 25, 'Rich-Mond ', 65000.00, '2007-12-13' ), (5, 'David', 27, 'Texas', 85000.00, '2007-12-13');

```


查询数据
```sql
SELECT * FROM company;

SELECT ID,NAME FROM company;

--找出 AGE 以 2 开头的数据
SELECT * FROM COMPANY WHERE AGE::text LIKE '2%';

--找出 address 字段中含有 - 字符的数据
SELECT * FROM COMPANY WHERE ADDRESS  LIKE '%-%';

--来查询SALARY=10000的数据
SELECT * FROM COMPANY WHERE SALARY = 10000;

--找出 AGE(年龄) 字段大于等于 25,并且 SALARY(薪资) 字段大于等于 65000 的数据
SELECT * FROM COMPANY WHERE AGE >= 25 AND SALARY >= 65000;

--找出 AGE(年龄) 字段大于等于 25,或者 SALARY(薪资) 字段大于等于 65000 的数据
SELECT * FROM COMPANY WHERE AGE >= 25 OR SALARY >= 65000;

--在公司表中找出 AGE(年龄) 字段不为空的记录
SELECT * FROM COMPANY WHERE AGE IS NOT NULL;

--列出了 AGE(年龄) 字段为 25 或 27 的数据
SELECT * FROM COMPANY WHERE AGE IN ( 25, 27 );
SELECT * FROM COMPANY WHERE AGE **NOT** IN ( 25, 27 );

--列出了 AGE(年龄) 字段在 25 到 27 的数据
SELECT * FROM COMPANY WHERE AGE BETWEEN 25 AND 27;

--子查询语句中读取 SALARY(薪资) 字段大于 65000 的数据,然后通过 EXISTS 运算符判断它是否返回行,如果有返回行则读取所有的 AGE(年龄) 字段。
SELECT AGE FROM COMPANY
        WHERE EXISTS (SELECT AGE FROM COMPANY WHERE SALARY > 65000);
```

|实例|描述|
|-----|-----|
|WHERE SALARY::text LIKE '200%'	|找出 SALARY 字段中以 200 开头的数据。|
|WHERE SALARY::text LIKE '%200%'|找出 SALARY 字段中含有 200 字符的数据。|
|WHERE SALARY::text LIKE '_00%'	|找出 SALARY 字段中在第二和第三个位置上有 00 的数据。|
|WHERE SALARY::text LIKE '2 % %'|找出 SALARY 字段中以 2 开头的字符长度大于 3 的数据。|
|WHERE SALARY::text LIKE '%2'	|找出 SALARY 字段中以 2 结尾的数据|
|WHERE SALARY::text LIKE '_2%3'	|找出 SALARY 字段中 2 在第二个位置上并且以 3 结尾的数据|
|WHERE SALARY::text LIKE '2___3'|找出 SALARY 字段中以 2 开头,3 结尾并且是 5 位数的数据|

更新数据
```sql
UPDATE COMPANY SET SALARY = 15000 WHERE ID = 3;

UPDATE COMPANY SET ADDRESS = 'Texas', SALARY=20000;
```

删除记录
```sql
DELETE FROM COMPANY WHERE ID = 2;

# 删除整张表的数据
DELETE FROM COMPANY;
```

修改表结构

用 ALTER TABLE 在一张已存在的表上添加列的语法如下:
```sql
ALTER TABLE table_name ADD column_name datatype;
```

在一张已存在的表上 DROP COLUMN(删除列),语法如下:
```sql
ALTER TABLE table_name DROP COLUMN column_name;
```

修改表中某列的 DATA TYPE(数据类型),语法如下:
```sql
ALTER TABLE table_name ALTER COLUMN column_name TYPE datatype;
```

数据类型:https://www.runoob.com/postgresql/postgresql-data-type.html

给表中某列添加 NOT NULL 约束,语法如下:
```sql
ALTER TABLE table_name MODIFY column_name datatype NOT NULL;
```

给表中某列 ADD UNIQUE CONSTRAINT( 添加 UNIQUE 约束),语法如下:
```sql
ALTER TABLE table_name
ADD CONSTRAINT MyUniqueConstraint UNIQUE(column1, column2...);
```

给表中 ADD CHECK CONSTRAINT(添加 CHECK 约束),语法如下:
```sql
ALTER TABLE table_name
ADD CONSTRAINT MyUniqueConstraint CHECK (CONDITION);
```

给表 ADD PRIMARY KEY(添加主键),语法如下:
```sql
ALTER TABLE table_name
ADD CONSTRAINT MyPrimaryKey PRIMARY KEY (column1, column2...);
```

DROP CONSTRAINT (删除约束),语法如下:
```sql
ALTER TABLE table_name
DROP CONSTRAINT MyUniqueConstraint;
```

DROP PRIMARY KEY (删除主键),语法如下:
```sql
ALTER TABLE table_name
DROP CONSTRAINT MyPrimaryKey;
```

### 11.删除表
```shell
drop table department;
drop table department, company;
```
用`\d`查看数据表

### 删除数据库

```shell
DROP DATABASE IF EXISTS db_name
```

### 控制台命令

```
\q: 退出
\h:查看SQL命令的解释,比如\h select。
\?:查看psql命令列表。
\l:列出所有数据库。
\c [database_name]:连接其他数据库。
\d:列出当前数据库的所有表格。
\d [table_name]:列出某一张表格的结构。
\du:列出所有用户。
\e:打开文本编辑器。
\password: 设置密码
\conninfo:列出当前数据库和连接的信息。
```

### 查询数据库大小

```sql
SELECT pg_size_pretty(pg_database_size('rcdb'));
```