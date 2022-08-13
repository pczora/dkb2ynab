from dkb_robo import DKBRobo
from pprint import pprint
from datetime import datetime, timedelta
import getpass
import csv

def main() -> int:
    dkb_user = input("user: ")
    dkb_password = getpass.getpass()
    tan_insert = False
    debug = True
    today = datetime.today()
    thirty_days_ago = today - timedelta(days=30)

    with DKBRobo(dkb_user, dkb_password, tan_insert, debug) as dkb:
        print(dkb.last_login)
        for account in dkb.account_dic.values(): 
            transactions = dkb.get_transactions(account['transactions'], account['type'], thirty_days_ago.strftime("%d.%m.%Y"), today.strftime("%d.%m.%Y"))
            pprint(transactions)
            with open(account['account'] + '.csv', 'w', newline='') as csvfile:
                headers = transactions[0].keys()
                csvwriter = csv.DictWriter(csvfile, fieldnames=headers)
                csvwriter.writeheader()
                csvwriter.writerows(transactions)

if __name__ == "__main__":
    main()
