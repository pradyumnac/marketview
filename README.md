# Marketview
Lists end of market data for Indian Equity Markets
Based on publicly available data from NSE & other servers

## Goal
View end of day market data in a easily reabale format in the terminal and web browser

1. server/ contains files that serves files in a Github Pages like static server
	Accesible at [Marketview](pradyumnac.github.io/marketview)
2. cli/ contains code that can bebuilt as a terminal application
3. fetcher/ contains code that fetches 
4. data/ contains files fetched from publicly available data that is used to show
	to the user
	1. eod - end of market data for scrips and indices
	2. mf - end of market data for mutual funds
	3. vix - end of day volatility data
	Refer to this []file](data/readme.md) to understand directory strycture

## Scope  
1. This project will only show EOD data. Nothing intraday for trading purposes


## FAQ
1. Why does the script `nse_cd` return no data intermittently?
Listing of companies shown in fzf list are from both BSE and NSE. Some companies
are solely listed in BSE but not in NSE. This mght be one of those cases

If you feel, you have encountered a bug, feel free to report
However, include the run log for debug purposes 
