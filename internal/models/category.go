package models

import (
	"database/sql"
	"fmt"
	"time"
	// "forum/internal/db"
)

// Category represents a forum category
type Category struct {
	Id 			int
	PostID 		int
	UserID 		int
	Status 		string
	Content 	string
	CreatedAt 	time.Time
}


// GetAllDistinctCategories returns all unique category statuses from the category table.
func GetalldistCat(db *sql.DB) ([]string, error) {
    rows, err := db.Query("SELECT DISTINCT status FROM category")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var category string
        if err := rows.Scan(&category); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
	fmt.Printf("Fetched categories: %v\n", categories)
    return categories, nil
}