package main

import (
	"database/sql" 
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/denisenkom/go-mssqldb"
)

type ProductionRequest struct {
	ProductId string `json:"productid"` 
}

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func getProduct(c echo.Context) error {
	// เชื่อมต่อกับฐานข้อมูล MSSQL
	db, err := sql.Open("mssql", "server=yourserver;user id=yourusername;password=yourpassword;port=1433;database=yourdatabase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ดึงค่า ProductId จาก URL parameter
	productId := c.QueryParam("productid")

	// Query เพื่อค้นหาข้อมูลสินค้าจากฐานข้อมูล
	rows, err := db.Query("exec dbo.getProdTypeMaster = ?", productId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// อ่านข้อมูลจาก rows และสร้าง struct Product
	var product Product
	for rows.Next() {
		err := rows.Scan(&product.Name, &product.Price)
		if err != nil {
			log.Fatal(err)
		}
	}

	// ส่งข้อมูลสินค้ากลับไปยังผู้ใช้
	return c.JSON(http.StatusOK, product)
}

func main() {
	e := echo.New()

	// เรียกใช้งาน getProduct เมื่อมี HTTP GET มาที่ /product
	e.GET("/product", getProduct)

	// เริ่มต้นเซิร์ฟเวอร์ที่พอร์ต 8080
	e.Start(":9092")
}
