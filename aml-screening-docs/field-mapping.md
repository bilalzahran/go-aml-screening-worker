# Field Mapping — World Check API → Prompt Template

## Subject Fields (KYC / CRM source)

| Go Struct Field | Prompt Placeholder | API / DB Source | Format |
|---|---|---|---|
| `Subject.Name` | `{{subject_name}}` | Customer record — full name | Native script preferred (e.g. 王俊明) |
| `Subject.DOB` | `{{subject_dob}}` | Customer record — date of birth | `YYYY-MM-DD` |
| `Subject.Nationality` | `{{subject_nationality}}` | Customer record — nationality | ISO 3166-1 alpha-3 (`CHN`, `USA`) |
| `Subject.Gender` | `{{subject_gender}}` | Customer record — gender | `MALE` / `FEMALE` |
| `Subject.EntityType` | `{{subject_entity_type}}` | Customer record — entity type | `INDIVIDUAL` / `ENTITY` |

## Hit Fields (World Check API response)

| Go Struct Field | Prompt Placeholder | API Field Path | Notes |
|---|---|---|---|
| `Hit.ReferenceID` | `{{hit_reference_id}}` | `referenceId` or `entityId` | Unique watchlist record ID |
| `Hit.PrimaryName` | `{{hit_primary_name}}` | `primaryName` / `entityName` | Main name on record |
| `Hit.MatchedTerm` | `{{hit_matched_term}}` | `matchedTerm` / `matchedName` | Alias that triggered the match |
| `Hit.Provider` | `{{hit_provider}}` | `providerType` | `WATCHLIST` / `PEP` / `SANCTION` |
| `Hit.Category` | `{{hit_category}}` | `category` | e.g. `POST CONVICTION`, `PEP` |
| `Hit.Gender` | `{{hit_gender}}` | `gender` | From hit biographic details |
| `Hit.Nationality` | `{{hit_nationality}}` | `nationality` / `countryOfNationality` | ISO alpha-3 |
| `Hit.Location` | `{{hit_location}}` | `addresses[0]` / `countryLocation` | Province/city, country |

## Aliases (`Hit.Aliases[]`)

| Go Struct Field | API Field Path | Values |
|---|---|---|
| `Alias.Name` | `names[].fullName` | The alias string |
| `Alias.Type` | `names[].nameType` | `PRIMARY` / `AKA` / `NATIVE_AKA` / `DBA` / `FORMERLY_KNOWN_AS` |
| `Alias.Language` | `names[].languageName` | Optional. e.g. `Chinese`, `Arabic` |

## Comparison Fields (`Hit.Comparisons[]`)

| Go Struct Field | API Field Path | Values |
|---|---|---|
| `ComparisonField.Field` | `secondaryFieldResults[].field` | e.g. `NATIVE_AKA`, `DATE_OF_BIRTH`, `GENDER` |
| `ComparisonField.Submitted` | `secondaryFieldResults[].submittedValue` | Value from your screening request |
| `ComparisonField.Matched` | `secondaryFieldResults[].matchedValue` | Value from the watchlist record |
| `ComparisonField.Result` | `secondaryFieldResults[].matchResult` | `EXACT` / `CLOSE` / `NOT_MATCHED` / `N/A` |

## Adverse Information (`Hit.AdverseInfo`)

| Go Struct Field | API Field Path | Notes |
|---|---|---|
| `AdverseInfo.Categories` | `categories[]` | e.g. `Special Interest Categories`, `Other Bodies` |
| `AdverseInfo.Keywords` | `keywords[]` | e.g. `Narcotics Trafficking (M:1SK)` |
| `AdverseInfo.Reports` | `furtherInformation[]` or `reports[]` | Free-text incident descriptions |

## Associates (`Hit.Associates[]`)

| Go Struct Field | API Field Path | Notes |
|---|---|---|
| `Associate.Name` | `associates[].primaryName` | Associated individual's name |
| `Associate.Category` | `associates[].category` | e.g. `POST CONVICTION` |
| `Associate.EntityType` | `associates[].entityType` | `INDIVIDUAL` / `ENTITY` |

## Sources (`Hit.Sources[]`)

| Go Struct Field | API Field Path | Notes |
|---|---|---|
| `string` | `sources[].url` or `webLinks[]` | Reference URLs for the adverse info |
