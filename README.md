# opg-infra-costs
OPG Infrastructure Costs: Managed by opg-org-infra &amp; Terraform
This project is to auto generate a spreadsheet with cost tables and charts using the AWS SDK in php. this is pslit into two commands, one to fetch the data, one to generate the spreadsheet.


## Requirements

Requires PHP 7.1+ and `composer` to be installed and the `Makefile` presumes usage of `aws-vault` for idenity management.


# Commands

## Download

`./aws-costs costs-to-file` uses the environment credentials with the SDK to download costs to a json file stored in `./tmp`.


# Generation

`./aws-costs spreadsheet` looks for and converts all json files within `./tmp` to data sources and generates a spreadsheet based on a series of dimensions and tags.


# Usage

There is a `Makefile` which includes all steps to make this easy, just run:

<pre>make</pre>

This will handle all the downloads and generation. If you only want to download data files run:

<pre>make get-costs</pre>

For just spreadsheet creation, run

<pre>make create-spreadsheet</pre>