import mysql.connector
from mysql.connector import Error
from secrets import host, database_name, user, password

class database:
    def __init__(self):
        self._conn = mysql.connector.connect(
            host=host,
            database=database_name,
            user=user,
            passwd=password
        )
        self._cursor = self._conn.cursor()
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tab):
        self.close()

    @property
    def connection(self):
        return self._conn
    
    @property
    def cursor(self):
        return self._cursor
    
    def commit(self):
        self.connection.commit()

    def close(self, commit=True):
        if commit:
            self.commit()
        self.connection.close()

    def execute(self, sql, params=None):
        self.cursor.execute(sql, params or ())

    def fetchall(self):
        return self.cursor.fetchall()

'''
db = DataBase()
sql = "INSERT INTO power_cuts ( il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date, is_happened) VALUES ( %s, %s, %s, %s, %s, %s, %s, %s, %s)"
val = ["q", "w", "e", "r", "t", "y", "u", "i", False]
db.execute(sql, val)
db.commit()

db = DataBase()
sql = "INSERT INTO addresses ( il, ilce, mah, sok) VALUES ( %s, %s, %s, %s)"
val = ["İL", "İLÇE", "FALANCA MAHALLESİ", "VUMARIN OĞLU TİMAR Sokak"]
db.execute(sql, val)
db.commit()
'''