package expenses

import "os"

var DbUrl = os.Getenv("DATABASE_URL")
