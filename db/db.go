package db

/*
使用例子：
	conn := db.Connect("root", "123456", "127.0.0.1", 3306, "game",50,20)
	M := db.NewModel(conn, "t_user")
	res := M.Order("id desc").Where("id = 100000000001").Limit(1).Get()

*/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"game_server/core/logger"

	_ "github.com/go-sql-driver/mysql"
)

var (
	key        string
	value      string
	conditions string
	str        string
)

type Model struct {
	link      *sql.DB  //存储连接对象
	tableName string   //存储表名
	field     string   //存储字段
	allFields []string //存储当前表所有字段
	where     string   //存储where条件
	order     string   //存储order条件
	limit     string   //存储limit条件

}

func init() {

}

/**
 * 初始化连接数据库操作
 */
func Connect(user_name string, password string, ip string, port string, database string, maxConn int) *sql.DB {
	defer logger.Flush()
	//1.连接数据库
	connstr := user_name + ":" + password + "@tcp(" + ip + ":" + port + ")/" + database + "?charset=utf8"
	db, err := sql.Open("mysql", connstr)
	//2.判断连接
	if err != nil {
		logger.Errorf("connect mysql fail !", err)
		return nil
	}

	db.SetMaxOpenConns(maxConn)
	logger.Info("connect mysql success ! ]")
	return db
}

//构造方法
func NewModel(conn *sql.DB, table string) Model {
	var this Model
	this.field = "*"
	//1.存储操作的表名
	this.tableName = table
	this.link = conn
	//2.获得当前表的所有字段
	this.getFields()
	return this
}

/**
 * 获取当前表的所有字段
 */
func (this *Model) getFields() {

	//查看表结构
	sql := "DESC " + this.tableName
	//执行并发送SQL
	result, err := this.link.Query(sql)

	if err != nil {
		fmt.Printf("sql fail ! [%s]", err)
	}

	this.allFields = make([]string, 0)

	for result.Next() {
		var field string
		var Type interface{}
		var Null string
		var Key string
		var Default interface{}
		var Extra string
		err := result.Scan(&field, &Type, &Null, &Key, &Default, &Extra)
		if err != nil {
			fmt.Print("scan fail ! ", err.Error())
		}
		this.allFields = append(this.allFields, field)
	}

}

/**
 * 执行并发送SQL(查询)
 * @param string $sql  要查询的SQL语句
 * @return array 返回查询出来的二维数组
 */
func (this *Model) query(sql string) (interface{}, error) {

	rows2, err := this.link.Query(sql)

	//查询数据，取所有字段
	if err != nil {
		return nil, err
	}

	//返回所有列
	cols, err := rows2.Columns()

	if err != nil {
		return nil, err
	}

	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))

	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))

	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make(map[int]map[string]string)

	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...) //将slic地址传入
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result, nil

}

/**
 * 查询数据
 */
func (this *Model) Get() (interface{}, error) {
	sql := `select * from ` + this.tableName + ` ` + this.where + ` ` + this.order + ` ` + this.limit
	//执行并发送SQL
	result, err := this.query(sql)
	return result, err
}

/**
 * 设置要查询的字段信息
 * @param string $field  要查询的字段
 * @return object 返回自己，保证连贯操作
 */
func (this *Model) Field(field string) *Model {
	this.field = field
	return this
}

/**
 * order排序条件
 * @param string  $order  以此为基准进行排序
 * @return $this  返回自己，保证连贯操作
 */
func (this *Model) Order(order string) *Model {
	this.order = `order by ` + order
	return this
}

/**
 * limit条件
 * @param string $limit 输入的limit条件
 * @return $this 返回自己，保证连贯操作
 */
func (this *Model) Limit(limit int) *Model {
	this.limit = "limit " + strconv.Itoa(limit)
	return this
}

/**
 * where条件
 * @param string $where 输入的where条件
 * @return $this 返回自己，保证连贯操作
 */
func (this *Model) Where(where string) *Model {
	this.where = `where ` + where
	return this
}

/**
 * 统计总条数
 * @return int 返回总数
 */
func (this *Model) Count() (interface{}, error) {
	//准备SQL语句
	sql := `select count(*) as total from ` + this.tableName + ` limit 1`
	result, err := this.query(sql)
	return result, err
}

/**
 * 执行并发送SQL语句(增删改)
 * @param string $sql 要执行的SQL语句
 * @return bool|int|string 添加成功则返回上一次操作id,删除修改操作则返回true,失败则返回false
 */
func (this *Model) exec(sql string) (interface{}, error) {

	res, err := this.link.Exec(sql)

	if err != nil {
		return res, err
	}

	result, err := res.LastInsertId()
	if err != nil {
		return res, err
	}
	return result, nil
}

/**
 * 添加操作
 * @param array  $data 要添加的数组
 * @return bool|int|string 添加成功则返回上一次操作的id,失败则返回false
 */
func (this *Model) Add(data map[string]interface{}) (interface{}, error) {

	//过滤非法字段
	for k, v := range data {
		if res := in_array(k, this.allFields); res != true {
			delete(data, k)
		} else {
			key += `,` + k
			value += `,` + `'` + v.(string) + `'`
		}
	}

	//将map中取出的键转为字符串拼接
	key = strings.TrimLeft(key, ",")
	//将map中的值转化为字符串拼接
	value = strings.TrimLeft(value, ",")
	//准备SQL语句
	sql := `insert into ` + this.tableName + ` (` + key + `) values (` + value + `)`
	// //执行并发送SQL
	result, err := this.exec(sql)

	return result, err

}

/**
 * 删除操作
 * @param string $id 要删除的id
 * @return bool  删除成功则返回true,失败则返回false
 */
func (this *Model) Delete() (interface{}, error) {

	sql := `delete from ` + this.tableName + ` ` + this.where

	//执行并发送
	result, err := this.exec(sql)

	return result, err
}

/**
 * 修改操作
 * @param  array $data  要修改的数组
 * @return bool 修改成功返回true，失败返回false
 */
func (this *Model) Update(data map[string]interface{}) (interface{}, error) {

	//过滤非法字段
	for k, v := range data {
		if res := in_array(k, this.allFields); res != true {
			delete(data, k)
		} else {
			str += k + ` = '` + v.(string) + `',`
		}
	}

	//去掉最右侧的逗号
	str = strings.TrimRight(str, ",")

	//判断是否有条件
	if this.where == "" {
		fmt.Println("没有条件")
	}

	sql := `update ` + this.tableName + ` set ` + str + ` ` + this.where

	result, err := this.exec(sql)

	return result, err
}

//是否存在数组内
func in_array(need interface{}, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

//返回json
func returnRes(errCode int, res interface{}, msg interface{}) string {
	result := make(map[string]interface{})
	result["errCode"] = errCode
	result["result"] = res
	result["msg"] = msg
	data, _ := json.Marshal(result)
	return string(data)
}
