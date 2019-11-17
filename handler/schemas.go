package handler

var schemas = []string{
	`CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		parent INTEGER,
		FOREIGN KEY(parent) REFERENCES categories(id));`,
	`CREATE TABLE IF NOT EXISTS parts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pn TEXT,
		manufacturer TEXT,
		category INTEGER,
		value DECIMAL(6, 18), -- pico (-12) to mega (6)
		unit TEXT,
		package TEXT,
		description TEXT,
		location TEXT,
		inventory INTEGER,
		FOREIGN KEY(category) REFERENCES categories(id));`,
	`CREATE TABLE IF NOT EXISTS distributors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		part INTEGER,
		name TEXT,
		url TEXT,
		FOREIGN KEY(part) REFERENCES parts(id));`,
}

func (h *Handler) createTables() error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	for _, s := range schemas {
		_, err = tx.Exec(s)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
