# PokeTCGScraper
Scraper to find best seller in CardMarket (best price &amp; most amount of cards)

## TODO
- [ ] Functions to print struct as json
- [ ] Change multiple browser for multiple pages and only one browser
- [ ] Add concurency for round() and look for other places to add it
- [ ] Loop wanted list to start flow for each instead of having the url hardcoded
- [ ] Import wantedList from external sources (e.g. csv, CardMarket API (probably not))
- [ ] TBD...

## Flow
### Via Card

 - Get card seller
 - Get seller
 - Follow "Via Seller" Flow
 - Save data from every card(in wanted list) for every seller and sort it

### Via Seller
 - Loop through other cards in wanted list from the seller
