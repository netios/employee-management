// ready : mysql db
/* 
--
-- Table structure for table `employee`
--
CREATE DATABASE geuntaedb;
CREATE TABLE `employee` (
	`idx` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'index',
	`userid` INT(11) NOT NULL COMMENT '직원CODE',
	`employee_name` VARCHAR(50) NOT NULL COMMENT '직원이름',
	`employee_action` ENUM('출근','퇴근') NOT NULL COMMENT '발생이벤트. 1=출근, 2=퇴근, 3=야근, 4=휴일출근,5=휴일퇴근',
	`employee_ip` VARCHAR(15) NOT NULL COMMENT '접속한IP',
	`regdate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '발생날짜',
	PRIMARY KEY (`idx`),
	INDEX `userid` (`userid`)
)
COMMENT='근태관리 테이블'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
*/

// ready : go get
/*
go get github.com/go-sql-driver/mysql
go get github.com/labstack/echo
go get github.com/labstack/echo/middleware
go get github.com/JonathanMH/goClacks/echo
*/

package main

import (
  "database/sql"
  "fmt"
  "net/http"

  _ "github.com/go-sql-driver/mysql"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"

)

func main() {
  // Echo instance
  e := echo.New()


  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
  }))
   
  // Struct
  type Employee struct {
			Id      string `json:"idx"`
			UserId  string `json:"userid"`
			Name    string `json:"employee_name"`
			ActCode string `json: "employee_action"`
			Ip      string `json : "employee_ip"`
			RegDate string `json : "regdate"`
  }

  type Userinfo struct {
			UserId			string
			Username 		string
			Stats 			string
			Office_in		string
  }
  
  type Userinfomodel struct {
		Userinfo_value	Userinfo
		Userinfo_array []Userinfo
  }
  
  const (
    dbName = "geuntaedb"
	dbUser = "root"
    dbPass = "password"
    dbHost = "127.0.0.1"
    dbPort = "3306"
  )
  
  
   dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbName)
   
   db, err := sql.Open("mysql", dbSource)
   
	if err != nil {
				fmt.Println(err.Error())
	} else {
				fmt.Println("db is connected")
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
				fmt.Println(err.Error())
	}
  
  // Route => handler
  e.GET("/", func(c echo.Context) error {

    return c.JSON(http.StatusOK, "Hi!")
  })

  e.GET("/id/:id", func(c echo.Context) error {
    requested_id := c.Param("id")
    fmt.Println(requested_id);
    return c.JSON(http.StatusOK, requested_id)
  })
  
  //------------------------------------------------------------------------ 전체 리스트 보기
  e.GET("/employee/:id", func(c echo.Context) error {
	requested_id := c.Param("id")
	fmt.Println(requested_id)
	var userid string
	var name string
	var actcode string
	var ip string
	var regdate string
 
	err = db.QueryRow("SELECT userid, employee_name, employee_action, employee_ip, regdate FROM employee WHERE userid = ? AND regdate> curdate() ORDER BY regdate LIMIT 1", requested_id).Scan(&userid, &name, &actcode, &ip, &regdate)
 
	if err != nil {
		fmt.Println(err)
	}
 
	response := Employee{UserId: userid, Name: name, ActCode: actcode, Ip: ip, RegDate: regdate}
	return c.JSON(http.StatusOK, response)
	
  })
  
  //------------------------------------------------------------------------ 출근을 했는지 여부 조사
  e.GET("/employee_cnt/:id", func(c echo.Context) error {
	requested_id := c.Param("id")
	fmt.Println(requested_id)
	var tfflag string
	
	err = db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM employee WHERE userid = ? AND regdate> curdate() LIMIT 1", requested_id).Scan(&tfflag)
	if err != nil {
		fmt.Println(err)
	}
	return c.String(http.StatusOK, tfflag)
  })
  
  //------------------------------------------------------------------------ 출근 기록 시간
  e.GET("/employee_time_in/:id", func(c echo.Context) error {
	requested_id := c.Param("id")
	fmt.Println(requested_id)
	var regdate string
 
	err = db.QueryRow("SELECT regdate FROM employee WHERE userid = ? AND regdate> curdate() ORDER BY regdate LIMIT 1", requested_id).Scan(&regdate)
	
	if err != nil {
		fmt.Println(err)
	}

	return c.String(http.StatusOK, regdate)
  })
 
   //------------------------------------------------------------------------ 퇴근 기록 시간
  e.GET("/employee_time_out/:id", func(c echo.Context) error {
	requested_id := c.Param("id")
	fmt.Println(requested_id)
	var regdate string
 
	err = db.QueryRow("SELECT regdate FROM employee WHERE userid = ? AND regdate> curdate() AND employee_action = '퇴근' ORDER BY regdate DESC LIMIT 1", requested_id).Scan(&regdate)
	
	if err != nil {
		fmt.Println(err)
	}

	return c.String(http.StatusOK, regdate)
  })

  //------------------------------------------------------------------------ 현재 DB 시간
  e.GET("/employee_today", func(c echo.Context) error {

	var u Userinfomodel
 
	sql := "select * from view_today"
	
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	
    for rows.Next() {
		err = rows.Scan(&u.Userinfo_value.UserId, &u.Userinfo_value.Username, &u.Userinfo_value.Stats, &u.Userinfo_value.Office_in)
		if err != nil {
			fmt.Println(err)
		}
		u.Userinfo_array = append(u.Userinfo_array, u.Userinfo_value)
    }
	//** 중요 DB에 null 값이 있으면 에러가 남 그래서 반드시 view로 null를 예외처리하여 string화 한다.
	// null을 처리하려면 따로 구현해야함.	
	return c.JSON(http.StatusOK, u)
  })
  
  //------------------------------------------------------------------------ 값 입력하기
  e.POST("/employee", func(c echo.Context) error {
	emp := new(Employee)
	if err := c.Bind(emp); err != nil {
				return err
	}
	//
	sql := "INSERT INTO employee(userid, employee_name, employee_action, employee_ip) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(sql)

	if err != nil {
				fmt.Print(err.Error())
	}
	defer stmt.Close()
	result, err2 := stmt.Exec(emp.UserId, emp.Name, emp.ActCode, emp.Ip)

	// Exit if we get an error
	if err2 != nil {
				panic(err2)
	}
	fmt.Println(result.LastInsertId())

	return c.JSON(http.StatusCreated, emp)
  })

  e.Logger.Fatal(e.Start(":4000"))
}
