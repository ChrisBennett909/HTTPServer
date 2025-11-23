import sqlite3

conn = sqlite3.connect("credDB.db")
cur = conn.cursor()


cur.execute("""
        CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        isAdmin INTEGER NOT NULL,
        sessionToken TEXT
        );
        """)

conn.commit()
conn.close()

print("Users table created")
