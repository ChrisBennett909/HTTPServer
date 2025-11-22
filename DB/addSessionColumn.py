import sqlite3

conn = sqlite3.connect("credDB.db")
cursor = conn.cursor()
table = "users"

cursor.execute("""
ALTER TABLE users ADD COLUMN sessionToken TEXT;
""")

conn.commit()
conn.close()

print("Column Added")
