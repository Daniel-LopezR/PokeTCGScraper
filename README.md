# PokeTCGScraper
Scraper to find best seller in CardMarket (best price &amp; most amount of cards)

## Thoughts about refactoring to C
 - Main C program could be in charge of checking browser and generating a file with browser options. Therefore go code that runs rod checks that file and the argument sent from C to decide what browser to use.
 - Results are written from the go code to a FIFO file and read from C program(reading until some kind of FINISH keyword is found), then it calls a finish reading function in go that clears the file. The idea is to only have one program modify the file.
 - ....

## OLD_TODO
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
