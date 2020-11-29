package customers

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"context"
	"time"
	"errors"
	"log"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type Service struct{
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Customer struct {
	ID 		int64		`json:"id"`
	Name	string		`json:"name"`
	Phone	string		`json:"phone"`
	Active	bool		`json:"active"`
	Created	time.Time	`json:"created"`	
}

func (s *Service) ByID(ctx context.Context, id int64)(*Customer,error){
	item := &Customer{}

	err := s.db.QueryRowContext(ctx,`
	SELECT id,name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID,&item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	
	return item, nil
}

func (s *Service) All(ctx context.Context)([]*Customer,error){
	items := make([]*Customer,0)
	rows,err := s.db.QueryContext(ctx,`
	SELECT id,name, phone, active, created FROM customers ORDER BY id
	`)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	for rows.Next() {
		item := &Customer{}
		rows.Scan(&item.ID,&item.Name, &item.Phone, &item.Active, &item.Created)
		items = append(items,item)
	}
	return items,nil
}

func (s *Service) AllActive(ctx context.Context)([]*Customer,error){
	items := make([]*Customer,0)
	rows,err := s.db.QueryContext(ctx,`
	SELECT id,name, phone, active, created FROM customers WHERE active= true ORDER BY id;
	`)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	for rows.Next() {
		item := &Customer{}
		rows.Scan(&item.ID,&item.Name, &item.Phone, &item.Active, &item.Created)
		items = append(items,item)
	}
	return items,nil
}

func (s *Service) Create(ctx context.Context,name string, phone string)(*Customer,error){
	customer:= &Customer{
		Name: name,
		Phone: phone,
		Active: true,
	}

	err := s.db.QueryRowContext(ctx,`
	INSERT INTO customers(name,phone) VALUES ($1,$2) ON CONFLICT (phone) DO UPDATE SET name= excluded.name RETURNING id,created;
	`,name,phone).Scan(&customer.ID,&customer.Created)
	if err != nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return customer,nil
}

func (s *Service) Update(ctx context.Context,id int64,name string, phone string)(*Customer,error){
	customer:= &Customer{
		ID: id,
		Name: name,
		Phone: phone,
	}

	err := s.db.QueryRowContext(ctx,`
	UPDATE customers SET name =$1,phone=$2 WHERE id =$3 RETURNING active,created
	`,name,phone,id).Scan(&customer.Active,&customer.Created)
	if err != nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return customer,nil
}


func (s *Service) RemoveByID(ctx context.Context,id int64)(*Customer,error){
	customer:= &Customer{}
	err := s.db.QueryRowContext(ctx,`
	DELETE FROM customers WHERE id= $1 RETURNING id,name,phone,active,created
	`,id).Scan(&customer.ID,&customer.Name,&customer.Phone,&customer.Active,&customer.Created)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return customer,nil
}

func (s *Service) BlockByID(ctx context.Context,id int64)(*Customer,error){
	customer:= &Customer{}
	err := s.db.QueryRowContext(ctx,`
	UPDATE customers SET active= false WHERE id= $1 RETURNING id,name,phone,active,created
	`,id).Scan(&customer.ID,&customer.Name,&customer.Phone,&customer.Active,&customer.Created)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return customer,nil
}

func (s *Service) UnBlockByID(ctx context.Context,id int64)(*Customer,error){
	customer:= &Customer{}
	err := s.db.QueryRowContext(ctx,`
	UPDATE customers SET active= true WHERE id= $1 RETURNING id,name,phone,active,created
	`,id).Scan(&customer.ID,&customer.Name,&customer.Phone,&customer.Active,&customer.Created)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return customer,nil
}