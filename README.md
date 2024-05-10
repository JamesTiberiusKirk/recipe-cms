# Recipe CMS
This is for personal use, but feel free to fork and use it yourself.

---

Made with go, echo, templ, tailwindcss and PostgresSQL.
To help wit developent, ait, hot-reloader-proxy and migrator are used.

## Developent
1. Start dev db `docker-compose up` 
2. Setup .env and install some dependencies `make setup` 
3. Start dev server with `make -j run` 
    - `-j` is so it can run the different watch jobs concurently

Hot reloader proxy will by default start the proxy server on port 5001

# Roadmap
- [ ] Figure out what todo with the homepage
- [ ] On every recipe add a recipe exporter page/button to import to my fitness pal
    - [More info](https://support.myfitnesspal.com/hc/en-us/articles/360032271592-How-does-the-Recipe-Importer-on-the-website-work)
- [ ] Add nutrients info table support
- [ ] Add total weight/portions support
- [ ] On ingredients/seasonings add a way to ignore amount
     - I.E. for seasonings I don't care about the amount, just want the list
- [ ] Maybe think of a way to reorder items in a recipe?
