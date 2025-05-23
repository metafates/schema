# Validators

This table features all available validators.

| Name | Description |
| ---- | ----------- |
| `Any[T]` | Any accepts any value of T. |
| `Zero[T]` | Zero accepts all zero values.<br/><br/>The zero value is:<br/>- 0 for numeric types,<br/>- false for the boolean type, and<br/>- "" (the empty string) for strings.<br/><br/>See [NonZero]. |
| `NonZero[T]` | NonZero accepts all non-zero values.<br/><br/>The zero value is:<br/>- 0 for numeric types,<br/>- false for the boolean type, and<br/>- "" (the empty string) for strings.<br/><br/>See [Zero]. |
| `Positive[T]` | Positive accepts all positive real numbers excluding zero.<br/><br/>See [Positive0] for zero including variant. |
| `Negative[T]` | Negative accepts all negative real numbers excluding zero.<br/><br/>See [Negative0] for zero including variant. |
| `Positive0[T]` | Positive0 accepts all positive real numbers including zero.<br/><br/>See [Positive] for zero excluding variant. |
| `Negative0[T]` | Negative0 accepts all negative real numbers including zero.<br/><br/>See [Negative] for zero excluding variant. |
| `Even[T]` | Even accepts integers divisible by two. |
| `Odd[T]` | Odd accepts integers not divisible by two. |
| `Email[T]` | Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>". |
| `URL[T]` | URL accepts a single url.<br/>The url may be relative (a path, without a host) or absolute (starting with a scheme).<br/><br/>See also [HTTPURL]. |
| `HTTPURL[T]` | HTTPURL accepts a single http(s) url.<br/><br/>See also [URL]. |
| `IP[T]` | IP accepts an IP address.<br/>The address can be in dotted decimal ("192.0.2.1"),<br/>IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18"). |
| `IPV4[T]` | IPV4 accepts an IP V4 address (e.g. "192.0.2.1"). |
| `IPV6[T]` | IPV6 accepts an IP V6 address, including IPv4-mapped IPv6 addresses.<br/>The address can be regular IPv6 ("2001:db8::68"), or IPv6 with<br/>a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18"). |
| `MAC[T]` | MAC accepts an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address. |
| `CIDR[T]` | CIDR accepts CIDR notation IP address and prefix length,<br/>like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291. |
| `Base64[T]` | Base64 accepts valid base64 encoded strings. |
| `Charset0[T, F]` | Charset0 accepts (possibly empty) text which contains only runes acceptable by filter.<br/>See [Charset] for a non-empty variant. |
| `Charset[T, F]` | Charset accepts non-empty text which contains only runes acceptable by filter.<br/>See also [Charset0]. |
| `Latitude[T]` | Latitude accepts any number in the range [-90; 90].<br/><br/>See also [Longitude]. |
| `Longitude[T]` | Longitude accepts any number in the range [-180; 180].<br/><br/>See also [Latitude]. |
| `InPast[T]` | InFuture accepts any time after current timestamp.<br/><br/>See also [InPast]. |
| `InFuture[T]` | InFuture accepts any time after current timestamp.<br/><br/>See also [InPast]. |
| `Unique[S, T]` | Unique accepts a slice-like of unique values.<br/><br/>See [UniqueSlice] for a slice shortcut. |
| `UniqueSlice[T]` | Unique accepts a slice of unique values.<br/><br/>See [Unique] for a more generic version. |
| `NonEmpty[S, T]` | NonEmpty accepts a non-empty slice-like (len > 0).<br/><br/>See [NonEmptySlice] for a slice shortcut. |
| `NonEmptySlice[T]` | NonEmptySlice accepts a non-empty slice (len > 0).<br/><br/>See [NonEmpty] for a more generic version. |
| `MIME[T]` | MIME accepts RFC 1521 mime type string. |
| `UUID[T]` | UUID accepts a properly formatted UUID in one of the following formats:<br/>  - xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx<br/>  - urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx<br/>  - xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx<br/>  - {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx} |
| `JSON[T]` | JSON accepts valid json encoded text. |
| `CountryAlpha2[T]` | CountryAlpha2 accepts case-insensitive ISO 3166 2-letter country code. |
| `CountryAlpha3[T]` | CountryAlpha3 accepts case-insensitive ISO 3166 3-letter country code. |
| `CountryAlpha[T]` | CountryAlpha accepts either [CountryAlpha2] or [CountryAlpha3]. |
| `CurrencyAlpha[T]` | CurrencyAlpha accepts case-insensitive ISO 4217 alphabetic currency code. |
| `LangAlpha2[T]` | LangAlpha2 accepts case-insensitive ISO 639 2-letter language code. |
| `LangAlpha3[T]` | LangAlpha3 accepts case-insensitive ISO 639 3-letter language code. |
| `LangAlpha[T]` | LangAlpha accepts either [LangAlpha2] or [LangAlpha3]. |
| `And[T, A, B]` | And is a meta validator that combines other validators with AND operator.<br/>Validators are called in the same order as specified by type parameters.<br/><br/>See also [Or], [Not]. |
| `Or[T, A, B]` | Or is a meta validator that combines other validators with OR operator.<br/>Validators are called in the same order as type parameters.<br/><br/>See also [And], [Not]. |
| `Not[T, V]` | Not is a meta validator that inverts given validator.<br/><br/>See also [And], [Or]. |
