package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

//UserSignup: 通过用户名和密码完成用户注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to insert, err:%s" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:%s" + err.Error())
		return false
	}
	//判断表中是否有username重复
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

//UserSignIn: 判断用户名和密码是否正确
func UserSignIn(username string, enc_passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == enc_passwd {
		return true
	}
	return false
}

//UpdateToken: 刷新用户登录token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	//执行sql查询
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
