import requests
import mysql.connector
from mysql.connector import Error
from database import database
from mail import mail_cut

def is_on_adresses(address_list, addr):
    for i in address_list:
        if i == addr:
            return True


db = database()
sql = "SELECT * FROM addresses"
db.execute(sql)
addresses = db.fetchall()

sql = "SELECT * FROM power_cuts"
db.execute(sql)
cuts = db.fetchall()

iller = []
ilceler = []
mahalleler = []
sokaklar = []
for address in addresses:
    iller.append(address[0])
    ilceler.append(address[1])
    mahalleler.append(address[2])
    sokaklar.append(address[3])



url = 'https://cc.meramedas.com.tr/services/publicdata.ashx?m=mrm_gb1&il=42&ilce=1827&mahalle=&ay=&yil='

resp = requests.get(url)
data = resp.json()
#print(data)
print(len(data))

for i in range(len(data)):
    #print(data[i]['MahalleKoyAdi'])
    if is_on_adresses(sokaklar, data[i]['CaddeSokakAdi']):
        il = data[i]['IlAdi']
        ilce = data[i]['IlceAdi']
        mah = data[i]['MahalleKoyAdi']
        sok = data[i]['CaddeSokakAdi']
        plan_start = data[i]['PlanlananBaslangic']
        plan_finish = data[i]['PlanlananBitis']
        ann_type = data[i]['IlanTipi']
        ann_date = data[i]['YayinTarihi']
        print(i)
        print(il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date)
        #print(mail_cut(il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date))
        sql = "INSERT INTO power_cuts ( il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date, is_happened, is_mailed) VALUES ( %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
        val = [il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date, False, False]
        db.execute(sql, val)
        db.commit()

sql = "SELECT * FROM power_cuts"
db.execute(sql)
cuts = db.fetchall()

for cut in cuts:
    if not cut[9]: #if it's not mailed yet...
        if not mail_cut(il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date):
            #cursor.execute ("UPDATE tblTableName SET Year=%s, Month=%s, Day=%s, Hour=%s, Minute=%s WHERE Server='%s' " % (Year, Month, Day, Hour, Minute, ServerID))
            sql = "UPDATE power_cuts SET is_mailed = 1 WHERE sok = %s"
            val = [cut[3]] #make is_mailed variable 1 if its has same sok value
            db.execute(sql, val)
            db.commit()