## Tree Schema

Below is a sample schema to represent an entire tree structure in a relational database:

```sql
CREATE TABLE trees (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  name VARCHAR(255) NOT NULL,
  start_time TIMESTAMP NULL,
  end_time TIMESTAMP NULL,
  result JSON NOT NULL,
  structure JSON NOT NULL, -- stores the Node tree structure
  INDEX idx_trees_created_at (created_at),
  INDEX idx_trees_active (active)
);
```

#### Benefits

* **Fast reads** for the entire tree.
* **No recursive joins** or closure tables required.
* **Simple caching** in Redis or the application layer.
* Ideal for cases where trees are **constructed once** and **read frequently**.

#### Trade-offs

* Querying specific nodes with SQL is **not efficient**.
* **Full JSON updates** are required when a node changes.

In this design, trees are **deleted and re-inserted** on updates, so updating the JSON blob is not an issue.
Direct node querying is unnecessary since all important data and states are available in the `trees` table.
This approach provides **excellent read performance** and **cacheability**, which aligns with the main use case.


## Usage Example

```go
import (
	"database/sql"
	"log"

	acc "dynaspec/accessor/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	accessor := acc.New(db)
	// Use accessor to load, store, or update trees
}
```