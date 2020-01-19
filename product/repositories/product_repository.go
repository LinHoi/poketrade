package repositories

import (
	"database/sql"
	"product/datamodels"
	"product/common"
	"strconv"
)

//1. define interface
//2. implement interface

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll()([]*datamodels.Product, error)
	SubProductNum(productID int64) error
}

type ProductManager struct {
	table string
	mysqlConn *sql.DB
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		conn ,err := common.NewMysqlConn()
		if err != nil {
			return err
		}

		p.mysqlConn = conn
	}
	if p.table == ""{
		p.table = "product"
	}

	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productId int64,err  error) {
	if err = p.Conn(); err != nil{
		return
	}
	sql :="insert into product set productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(product.ProductName,product.ProductNum,product.ProductImage,product.ProductUrl)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err!= nil {
		return false
	}
	sql := "delete from product where ID = ?"
	stmt,err  := p.mysqlConn.Prepare(sql)
	if err != nil {
		 return  false
	}
	_, err = stmt.Exec(productID)
	if err != nil {
		 return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "update product " +
		"set productName = ?," +
		"productNum = ?," +
		"productImage = ?," +
		"productUrl= ? " +
		"where ID = "+strconv.FormatInt(product.ID,10)

	stmt ,err := p.mysqlConn.Prepare(sql)
	if err != nil { return err}

	_, err = stmt.Exec(product.ProductName,product.ProductNum,product.ProductImage,product.ProductUrl)
	if err != nil {return err}
	return nil
}

func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "select * from "+p.table+" where ID = "+strconv.FormatInt(productID,10)

	row ,err := p.mysqlConn.Query(sql)
	if err != nil { return &datamodels.Product{},err}
	defer  row.Close()


	result := common.GetResultRow(row)
	if len(result) == 0 { return &datamodels.Product{}, err}

	//传值还是传指针？注意判别
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)

	return productResult, nil
}

func (p *ProductManager) SelectAll() (productList []*datamodels.Product,err error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "select * from "+p.table


	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil { return nil,err}

	result := common.GetResultRows(rows)
	if len(result) == 0 { return nil, nil}


	for _, r := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(r, product)
		productList = append(productList, product)
	}
	return
}

func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update "+p.table+" set"+" productNum = productNum-1 where ID = "+strconv.FormatInt(productID,10)
	stmt ,err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
func NewProductManager(table string, db *sql.DB) IProduct{
	return &ProductManager{table, db}
}

