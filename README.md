# OPG Infra Costs

Managed by `opg-org-infra` Terraform.

## Purpose

This `go` version has been updated to improve performance over the `php` original and move to the next step of integration with the upcoming OPG metrics tool.

It provides a series of sub commands (using `flag`) with `detail` being core and others adding wrappers for ease and reformatting of that data for transmission / storage. The commands are listed below

## Commands

All commands require relevant `aws` credentials for `breakglass` or `billing` roles and therefore must be run using our normal prefix of aws vault:

```bash
aws-vault exec identity -- $COMMAND
```

### List

- [Detail](#detail)
- [Month to Date](#month-to-date)
- [Year to Date](#year-to-date)
- [Cost Increases](#cost-increases)
- [Excel](#excel)
- [Metrics](#metrics)


### Detail

`detail` provides the most customisation and filtering abilities of all the commands allowing lookups against all of our AWS accounts.

Full list of options for this command:

```bash
opg-infra-costs detail
    -account string
    	Limit the accounts list to only this *name* only - eg Sirius
    -env string
    	Limit the accounts list to only this environment
    -granularity string
    	Grouping for the cost data to be either DAILY or MONTHLY (default "DAILY")
    -service string
    	Limit the aws services to only this service name (uppercase)
    -end-date string
    	End date following 2006-01-02T15:04:05Z07:00 format, defaults to today at midnight
    -start-date string
    	Start date following 2006-01-02T15:04:05Z07:00 format, defaults to yesterday at midnight
    -data-columns
            Display these column - eg Account.Name,Account.Environment,Cost - needs to align with data-group-by
    -data-group-by
            Group the data by columns within cost - eg Account.Name,Account.Environment - would merge cost data to that level
    -data-headers
            Header names for columns - eg AccountName,Environment,Cost - needs to align with data-group-by

```

NOTE: All output is shown as a table and is not in a set order (limitation of `go range`).

#### Listing all costs for the last 24 hours

This is the default setup, so you can simply run

```bash
opg-infra-costs detail
```

#### Listing all costs for the last 24 hours for a specific AWS service

This is the default setup, so you can simply run

```bash
opg-infra-costs detail \
-service="EC2"
```
Limit the cost data to only those relating to EC2 - the `-service` uses `contains` to match.


#### List all costs for a specific account and environment between dates

By using `-account` and `-environment` you can filter the cost data based on those matched and by utilising `-start-date` and `-end-date` you can adjust the time period of the data retrived.

```bash
opg-infra-costs detail \
-account="sirius" \
-env="dev" \
-start-date="2020-09-01T00:00:00Z" \
-end-date="2020-10-01T00:00:00Z" \
```

This command will fetch a full listing of cost data from AWS for the `development sirius` account between the start and end of September for each day.


#### A months total for each service for a specific account and environment

The previous example lists the costs for every day, that can be great for keeping a close eye on spend, but by changing the `-granuality` you can widen that scope to reduce the data being shown for easier viewing

```bash
opg-infra-costs detail \
-account="sirius" \
-env="dev" \
-start-date="2020-09-01T00:00:00Z" \
-end-date="2020-10-01T00:00:00Z" \
-granularity="MONTHLY"
```

This changes the cost data to be grouped by the month in total, rather than each day


#### Grouping and summation by fields

For more custom queries, such as seeing the monthly cost of AWS Config for two months over all sirius related environments you can using the `-data` prefixed arguments to adjust the filtering and grouping of the data and accumulating the costs.

* `-data-columns` filters what columns are shown in the tabular output.
* `-data-group-by` changes how the cost data is group and therefore accumulated
* `-data-headers` sets the tabular header strings

These fields relate to the `struct` being used in the project - so may give unusal results.

These are advanced usages and typically common ones will have a related command to handle it more cleanly.

An exmaple command that does find two months of `tax` costs for sirius, grouped by month and environment would look like this:

```bash
opg-infra-costs detail \
-account=sirius \
-service=tax \
-granularity=MONTHLY \
-data-group-by="Account.Name,Account.Environment,Service,Date" \
-data-columns="Account.Name,Account.Environment,Service,Date,Cost" \
-data-headers="Account,Env,Service,Date,Cost" \
-start-date=2020-07-01T00:00:00Z \
-end-date=2020-09-01T00:00:00Z
```

### Month to Date

Month to date provides a simple top line number by default for the current months costs:

```bash
opg-infra-costs mtd
```

It does provide filtering for `-account`, `-env` and `-service` filtering and a flag to display tabular breakfown of the costs.

#### Viewing the total costs for sirius' dev environments

To get a total sum for the current month for a specific account and environment you can run:

```bash
opg-infra-costs mtd \
-account="sirius" \
-env="dev"
```

This command would display total cost for sirius' dev environment.

#### Show a small tabular breakdown

Instead of getting just top line, you may want to see the costs for each version of sirius individually as well (so seeing line cost for dev, pre-prod & prod):

```bash
opg-infra-costs mtd \
-account="sirius" \
-breakdown
```

### Year to Date

Year to date provides a simple top line number by default for the current years total costs:

```bash
opg-infra-costs ytd
```

Provides the same filtering as `mtd` account


### Cost Increases

Command `increases` compares two sets of monthly cost data and flags significant increases in value between services.

```bash
opg-infra-costs increases
```

#### Selecting months

Change the two months being compared over all accounts

```bash
opg-infra-costs increases \
-a="2020-09" \
-b="2020-10"
```

#### Filtering

You can limit the account and environment being searched for in the same was as `detail` command:

```bash
opg-infra-costs increases \
-account="sirius" \
-env="prod"
```

#### Change percentage trigger

Adjust the percentage change to flag as an increase cost:

```bash
opg-infra-costs increases \
-percentage-change=50
```

#### Change ignore cost barrier

As we are likely to get large percentage changes on low value items we set a base cost level to detection. Any cost below this is not included in the cost comparison.

```bash
opg-infra-costs increases \
-baseline-cost=100
```

This would exclude any line item cost below $100 in value from the comparision table


### Excel

Generate a year to date spreadsheet of costs grouped in various ways with sparkline trends to a file called `costs.xlsx`

```bash
opg-infra-costs excel
```

This has arguments to filter by `-account` and `-env`, but not `-service`


### Metrics

*WIP* - Command to send data to the mentrics api end point

```bash
opg-infra-costs metrics
```
