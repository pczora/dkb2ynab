from bs4 import BeautifulSoup
import mechanicalsoup
import getpass
import csv

    
def main():

    username = input("Username: ")
    password = getpass.getpass()

    browser = mechanicalsoup.StatefulBrowser()

    browser.open("https://banking.zinspilot.de/login")
    browser.select_form('form[name="f"]')
    browser["j_username"] = username
    browser["j_password"] = password
    response = browser.submit_selected()
    browser.open("https://banking.zinspilot.de/benutzer/anlagen")
    page = browser.page
    investment_rows = page.find_all("li", class_="badges-item") 

    investments = []
    for ir in investment_rows:
        bank_name = ir["data-bankname"]
        value = ir.find(class_="badge-prize-value").text
        link = ir.find(class_="btn-info")
        browser.follow_link(link["href"])
        table = browser.page.find("table", class_="zebra-striped")
        table_rows = table.find_all("tr")
        with open(f"{bank_name}.csv", "wt+", newline="") as f:
            writer = csv.writer(f)
            for tr in table_rows:
                csv_row = []
                for header_cell in tr.find_all("th"):
                    csv_row.append(header_cell.get_text().strip())

                for cell in tr.find_all("td", class_=""):
                        csv_row.append(cell.get_text().replace("\n", " ").strip())
                for cell in tr.find_all("td", class_="booking-explanation-box"):
                        csv_row.append(cell.get_text().replace("\n", " ").strip())
                for cell in tr.find_all("td", class_="rtl"):
                        h5 = cell.find("h5")
                        if h5:
                            csv_row.append(h5.text.strip())
                        else:
                            csv_row.append(cell.text.strip())
                writer.writerow(csv_row)

if __name__ == "__main__":
    main()
