import csv

PKG = "iso"
COUNTRIES = "ISO-3166.csv"
CURRENCIES = "ISO-4217.csv"


def make_p(f): return lambda *s: print(*s, file=f)


def gen_countries():
    OUTPUT = "countries.go"

    data = []

    with open(COUNTRIES, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',', quotechar='"')
        header = next(reader)


        for row in reader:
            row = dict(zip(header, row))

            data.append(row)


    with open(OUTPUT, "w") as f:
        p = make_p(f)

        p(f"package {PKG}")
        p()
        p("var CountryAlpha2 = map[string]struct{}{")
        for row in data:
            alpha2 = row["alpha-2"].lower()
            p(f'\t"{alpha2}": {{}},')
        p("}")
        p()
        p("var CountryAlpha3 = map[string]struct{}{")
        for row in data:
            alpha3 = row["alpha-3"].lower()
            p(f'\t"{alpha3}": {{}},')
        p("}")


def gen_currencies():
    OUTPUT = "currencies.go"

    data = []

    with open(CURRENCIES, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',', quotechar='"')
        header = next(reader)


        for row in reader:
            row = dict(zip(header, row))

            data.append(row)


    with open(OUTPUT, "w") as f:
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
