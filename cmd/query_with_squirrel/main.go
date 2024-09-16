package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"

	desc "github.com/Ivanrumanchev/auth/pkg/user_v1"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу users
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(gofakeit.Name(), gofakeit.Email(), 1).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int
	err = pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	// Делаем запрос на выборку записей из таблицы users
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}

	var id int
	var name, email string
	var role desc.Role
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan users: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n", id, name, email, role, createdAt, updatedAt)
	}

	// Делаем запрос на обновление записи в таблице users
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("role", 2).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update users: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы users
	builderSelectOne := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n", id, name, email, role, createdAt, updatedAt)
}
