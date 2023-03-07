package database

import (
	"database/sql"
	"testing"

	"github.com/gointensivo2/internal/entity"
	"github.com/stretchr/testify/suite"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, error := sql.Open("sqlite3", ":memory:")
	suite.NoError(error)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price FLOAT NOT NULL, tax FLOAT NOT NULL, final_price FLOAT NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestSavingOrder() {
	order, error := entity.NewOrder("123", 10.0, 1.0)
	suite.NoError(error)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	error = repo.Save(order)
	suite.NoError(error)

	var orderResult entity.Order
	suite.NoError(suite.Db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", order.ID).Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice))

	suite.NoError(error)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}
