import csv

# go package name
PKG = "iso"

# https://github.com/datasets/country-codes/blob/main/data/country-codes.csv
COUNTRIES = "countries.csv"
COUNTRIES_OUT = "countries.go"

# https://github.com/datasets/currency-codes/blob/main/data/codes-all.csv
CURRENCIES = "currencies.csv"
CURRENCIES_OUT = "currencies.go"


def make_p(f): return lambda *s: print(*s, file=f)


def gen_countries():
    data = []

    with open(COUNTRIES, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',', quotechar='"')
        header = next(reader)


        for row in reader:
            row = dict(zip(header, row))

            data.append(row)


    with open(COUNTRIES_OUT, "w") as f:
        p = make_p(f)

        p(f"package {PKG}")
        p()
        p("var CountryAlpha2 = map[string]struct{}{")
        for row in data:
            alpha2 = row["ISO3166-1-Alpha-2"].lower()
            p(f'\t"{alpha2}": {{}},')
        p("}")
        p()
        p("var CountryAlpha3 = map[string]struct{}{")
        for row in data:
            alpha3 = row["ISO3166-1-Alpha-3"].lower()
            p(f'\t"{alpha3}": {{}},')
        p("}")


def gen_currencies():
    data = []

    with open(CURRENCIES, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',', quotechar='"')
        header = next(reader)


        for row in reader:
            row = dict(zip(header, row))

            data.append(row)


    with open(CURRENCIES_OUT, "w") as f:
        p = make_p(f)

        visited = set()

        p(f"package {PKG}")
        p()
        p("var CurrencyAlpha = map[string]struct{}{")
        for row in data:
            code = row["AlphabeticCode"].lower()

            if code not in visited and code != "":
                p(f'\t"{code}": {{}},')
                visited.add(code)
        p("}")


def main():
    gen_countries()
    gen_currencies()


if __name__ == "__main__":
    main()
