from secrets import sender_mail, sender_pass, dest_mail
from smtplib import SMTP_SSL as SMTP
from email.mime.text import MIMEText

def mail(content):
    SMTPserver = 'smtp.gmail.com'
    sender =     sender_mail
    destination = [dest_mail]

    USERNAME = sender_mail
    PASSWORD = sender_pass

    text_subtype = 'plain' #or html, xml

    subject="Planlı Kesinti Bildirisi"

    try:
        msg = MIMEText(content, text_subtype)
        msg['Subject']=       subject
        msg['From']   = sender

        conn = SMTP(SMTPserver)
        conn.set_debuglevel(False)
        conn.login(USERNAME, PASSWORD)
        try:
            conn.sendmail(sender, destination, msg.as_string())
        finally:
            conn.quit()

    except:
        sys.exit( "mail failed; %s" % "CUSTOM_ERROR" )

def mail_cut(il, ilce, mah, sok, plan_start, plan_finish, ann_type, ann_date):
    try:
        con="""\
            Selam! Bir elektrik kesintisi sizi bekliyor :(
            
            %s %s %s %s'da %s ile %s süresi arasında %s nedeniyle elektrik kesintisi planlanıyor.

            Duyuru eklenme tarihi: %s

            Lütfen kesintiye karşı hazırlıklarınızı yapın.
            """ % (mah, sok, ilce, il, plan_start, plan_finish, ann_type, ann_date)
        mail(con)
    finally:
        return 0
    return 1

if __name__ == "__main__":
    con="""\
        Test message!
        """
    mail(con)