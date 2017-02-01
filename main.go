package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}

func _main() error {
	db, err := sql.Open("mysql", "root:@tcp(mysql:3306)/test")
	if err != nil {
		return err
	}

	for {
		if err := db.Ping(); err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	if err := readyTables(db); err != nil {
		return err
	}
	if err := create(db); err != nil {
		return err
	}
	if err := readList(db); err != nil {
		return err
	}
	// if err := union(conn); err != nil {
	// 	return err
	// }

	return nil
}

func readyTables(db *sql.DB) error {
	for _, q := range []string{
		`
			create table users (
				id int not null auto_increment,
				name varchar(255),
				gender varchar(255),
				age int,
				primary key (id)
			)
		`,
		`
			create table posts (
				id int not null auto_increment,
				user_id int not null,
				body varchar(1000),
				primary key (id),
				index idx_user_id (user_id),
				foreign key (user_id) references users (id)
			)
		`,
	} {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

type User struct {
	ID     sql.NullInt64
	Name   sql.NullString
	Gender sql.NullString
	Age    sql.NullInt64
}

func (u User) InsertSQL() (string, sql.NullInt64, sql.NullString, sql.NullString, sql.NullInt64) {
	return "INSERT INTO users (id, name, gender, age) VALUES (?, ?, ?, ?)", u.ID, u.Name, u.Gender, u.Age
}

type Users []User

func (us Users) SelectSQL() string {
	return "SELECT * FROM users"
}

func (us *Users) ReadRows(rows *sql.Rows) error {
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	fmt.Println(cols)
	for i := 0; rows.Next(); i++ {
		u := User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Gender, &u.Age); err != nil {
			return err
		}
		*us = append(*us, u)
	}
	return nil
}

type Post struct {
	ID     sql.NullInt64
	Body   sql.NullString
	UserID sql.NullInt64 `db:"user_id"`
	User   User          `db:"-"`
}

func (u Post) InsertSQL() (string, sql.NullInt64, sql.NullString, sql.NullInt64) {
	return "INSERT INTO posts (id, body, user_id) VALUES (?, ?, ?)", u.ID, u.Body, u.UserID
}

type Posts []Post

func (ps Posts) SelectSQL() string {
	return "SELECT posts.id as `posts-id`, posts.body, posts.user_id, users.id, users.name, users.gender, users.age FROM posts LEFT JOIN users ON posts.user_id = users.id"
}

func (ps *Posts) ReadRows(rows *sql.Rows) error {
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	fmt.Println(cols)
	for i := 0; rows.Next(); i++ {
		p := Post{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.Body, &p.User.ID, &p.User.Name, &p.User.Gender, &p.User.Age); err != nil {
			return err
		}
		*ps = append(*ps, p)
	}
	return nil
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{Valid: true, String: s}
}

func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Valid: true, Int64: i}
}

func NewNullBool(b bool) sql.NullBool {
	return sql.NullBool{Valid: true, Bool: b}
}

func NewNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Valid: true, Float64: f}
}

func create(db *sql.DB) error {
	for _, u := range []User{
		{
			Name:   NewNullString("Foo"),
			Gender: NewNullString("male"),
			Age:    NewNullInt64(29),
		},
		{
			Name:   NewNullString("Bar"),
			Gender: NewNullString("female"),
			Age:    NewNullInt64(17),
		},
		{
			Name:   NewNullString("Baz"),
			Gender: NewNullString("male"),
			Age:    NewNullInt64(41),
		},
		{
			Name:   NewNullString("Qux"),
			Gender: NewNullString("female"),
			Age:    NewNullInt64(32),
		},
		{
			Name:   NewNullString("Hoge"),
			Gender: NewNullString("male"),
			Age:    NewNullInt64(11),
		},
		{
			Name:   NewNullString("Fuga"),
			Gender: NewNullString("female"),
			Age:    NewNullInt64(51),
		},
	} {
		res, err := db.Exec(u.InsertSQL())
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", res)
	}

	for _, p := range []Post{
		{
			UserID: NewNullInt64(1),
			Body:   NewNullString("AAAAAAAAAAA"),
		},
		{
			UserID: NewNullInt64(1),
			Body:   NewNullString("BBBBBBBBBBBBB"),
		},
		{
			UserID: NewNullInt64(2),
			Body:   NewNullString("CCCCCC"),
		},
	} {
		res, err := db.Exec(p.InsertSQL())
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", res)
	}

	return nil
}

func readList(db *sql.DB) error {
	var (
		// us Users
		ps Posts
	)

	rows, err := db.Query(ps.SelectSQL())
	if err != nil {
		return err
	}
	if err := ps.ReadRows(rows); err != nil {
		return nil
	}

	// rows, err := db.Query(us.SelectSQL())
	// if err != nil {
	// 	return err
	// }
	// us.ReadRows(rows)
	// users, err := json.MarshalIndent(us, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("users = %v\n", string(users))

	// umap := map[int64]User{}
	// for _, u := range us {
	// 	if u.ID.Valid {
	// 		umap[u.ID.Int64] = u
	// 	}
	// }
	// for i, p := range ps {
	// 	if p.ID.Valid {
	// 		p.User = umap[p.UserID.Int64]
	// 		ps[i] = p
	// 	}
	// }

	posts, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("posts = %v\n", string(posts))

	return nil
}

// func uniq(ns []sql.NullInt64) []int64 {
// 	m := map[int64]bool{}
// 	for _, n := range ns {
// 		if n.Valid {
// 			m[n.Int64] = true
// 		}
// 	}
// 	var s []int64
// 	for i, _ := range m {
// 		s = append(s, i)
// 	}
// 	return s
// }

// func union(conn *sql.Connection) error {
// 	var us []User
// 	sess := conn.NewSession(nil)
// 	if _, err := sess.Select("*").From(
// 		sql.Union(
// 			sql.Select("*").From("users").Where(
// 				sql.And(
// 					sql.Gt("age", 20),
// 					sql.Eq("gender", "male"),
// 				),
// 			),
// 			sql.Select("*").From("users").Where(
// 				sql.And(
// 					sql.Lt("age", 50),
// 					sql.Eq("gender", "female"),
// 				),
// 			),
// 		).As("uni"),
// 	).Load(&us); err != nil {
// 		return err
// 	}
//
// 	// users, err := json.MarshalIndent(us, "", "  ")
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Printf("union users = %v\n", string(users))
//
// 	return nil
// }
