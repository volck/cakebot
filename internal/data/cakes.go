package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Cake struct {
	User_ID      string  `db:"USER_ID"`
	When         string  `db:"CAKEDAY"`
	CreatedAt    string  `db:"CREATEDAT" json:"-"`
	Excluded     int     `db:"EXCLUDED" json:"-"`
	Firstname    string  `db:"FIRSTNAME"`
	Lastname     string  `db:"LASTNAME"`
	Ldapdata     []uint8 `db:"LDAPDATA" json:"-"`
	Notification int     `db:"NOTIFICATION" json:"-"`
	Notified     int     `db:"NOTIFIED" json:"-"`
	NewDate      string  `json:"-"`
}

type CakeModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the cakedays table.
func (c CakeModel) NewDraw() (sql.Result, error) {

	query := `insert into cakedays (user_id, cakeday, createdat) select user_id, friday, sysdate from (
	with chefs as (select
	user_id
	from users
	order by dbms_random.random),
	cdays as (select
	friday
	from fridays where friday > (select nvl(max(cakeday), sysdate) from cakedays) order by friday)
	select
	u.user_id,
	c.friday
	from (select user_id, rownum as line from chefs) u
	join (select friday, rownum as line from cdays) c on (u.line = c.line))
`

	result, err := c.DB.Exec(query)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return result, nil
}

func (c CakeModel) Insert(cake *Cake) error {

	query := `INSERT INTO CAKEDAYS(CAKEDAYS.user_id, CAKEDAYS.cakeday, CAKEDAYS.createdat, NOTIFICATION) VALUES(:USERID, TO_DATE(:WHEN, 'DD.MM.YY'), SYSDATE, 0)`

	args := []any{cake.User_ID, cake.When}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&cake.User_ID, &cake.When, &cake.Notification)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrRecordNotFound
	default:
		return nil
	}
}

func (c CakeModel) GetCurrent() (*Cake, error) {

	query := `select
cakedays.USER_ID, cakedays.CAKEDAY, users.firstname, users.lastname
from cakedays
INNER JOIN users
ON users.user_id= cakedays.user_id where
to_char(cakeday, 'YYYY-IW') = to_char(sysdate, 'YYYY-IW')`

	// Declare a Movie struct to hold the data returned by the query.
	var cake Cake
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := c.DB.QueryRow(query).Scan(
		&cake.User_ID,
		&cake.When,
		&cake.Firstname,
		&cake.Lastname,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Otherwise, return a pointer to the cake struct.
	return &cake, nil
}

func (c CakeModel) GetByDate(theDate string) (*Cake, error) {
	query := `select cakedays.USER_ID, cakedays.CAKEDAY, users.firstname, users.lastname from cakedays INNER JOIN users ON users.user_id= cakedays.user_id WHERE cakedays.cakeday = TO_DATE(:theDate, 'DD.MM.YY')`

	var cake Cake
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := c.DB.QueryRow(query, theDate).Scan(
		&cake.User_ID,
		&cake.When,
		&cake.Firstname,
		&cake.Lastname,
	)

	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Otherwise, return a pointer to the Movie struct.
	return &cake, nil

}

func (c CakeModel) GetUid(id string) ([]*Cake, error) {

	query := "select cakedays.USER_ID, cakedays.CAKEDAY, users.firstname, users.lastname from cakedays INNER JOIN users ON users.user_id= cakedays.user_id WHERE users.user_id = :ID collate binary_ci ORDER BY cakeday ASC"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := c.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	cakes := []*Cake{}
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var cake Cake
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.

		err := rows.Scan(
			&cake.User_ID,
			&cake.When,
			&cake.Firstname,
			&cake.Lastname,
		)
		if err != nil {
			return nil, err
		}
		cakes = append(cakes, &cake)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of movies.
	return cakes, nil

}

// Add a placeholder method for fetching a specific record from the cakedays table.
func (c CakeModel) GetAll() ([]*Cake, error) {

	// Create a context with a 3-second timeout.
	query := `select
cakedays.USER_ID, cakedays.CAKEDAY, cakedays.CREATEDAT, users.firstname, users.lastname
from cakedays
INNER JOIN users
ON users.user_id= cakedays.user_id
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	cakes := []*Cake{}
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var cake Cake
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.

		err := rows.Scan(
			&cake.User_ID,
			&cake.When,
			&cake.CreatedAt,
			&cake.Firstname,
			&cake.Lastname,
		)
		if err != nil {
			return nil, err
		}
		cakes = append(cakes, &cake)
	}
	// Add the Movie struct to the slice.

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of movies.
	return cakes, nil
}

// Add a placeholder method for updating a specific record in the cakedays table.
func (c CakeModel) Update(cake *Cake) error {

	query := `
UPDATE cakedays
SET cakedays.USER_ID = :USERID, cakedays.CREATEDAT = SYSDATE, cakedays.cakeday = TO_DATE(:NEWDATE, 'DD.MM.YYYY')
WHERE cakedays.cakeday = TO_DATE(:WHEN, 'DD.MM.YYYY')`

	fmt.Printf("the cake:\nuid: %s\nnewdate: %s\noriginal:%s\n\n\n", cake.User_ID, cake.NewDate, cake.When)

	args := []any{
		cake.User_ID,
		cake.NewDate,
		cake.When,
	}
	err := c.DB.QueryRow(query, args...).Scan(&cake.NewDate)
	if err != nil {
		fmt.Println("queryrow err", err)
	}
	return nil
}

// Add a placeholder method for deleting a specific record from the cakedays table.
func (c CakeModel) Unset(theDate string) error {
	query := `UPDATE cakedays
SET cakedays.USER_ID = NULL, CREATEDAT = SYSDATE
WHERE cakedays.cakeday = TO_DATE(:theDate, 'DD.MM.YYYY')`
	returnDate := ""
	err := c.DB.QueryRow(query, theDate).Scan(&returnDate)
	if err != nil {
		fmt.Println("queryrow err", err)
	}

	return nil

}

func (c CakeModel) SetNotified(cake *Cake) error {
	query := `UPDATE cakedays
set notified = 1
WHERE user_id = :USER_ID AND cakeday = to_date(:CAKEDAY, 'DD.MM.YYYY')`
	ctx := context.Background()
	result, err := c.DB.ExecContext(ctx, query, cake.User_ID, cake.When)
	if err != nil {
		log.Fatal("notified", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ExpectedSingleRowAffected
	}
	return nil
}
