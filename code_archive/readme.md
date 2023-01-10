# Company Annoucements Tracker


## Methodology 
Currently , I am using [digdeeper.one](https://digdeeper.one/announcements) as data source.
Browser based scrapping using [browserless.io](https://www.browserless.io/)
Plan to switch to BSE/NSE server in future.

## Usage
1. Create Python virtual environment and activate
2. Run `pip install -r requirements.txt`
3. Get an API key from [browserless.io](https://www.browserless.io/)
4. rename file `.env.sample` to `.env` and paste the API key into this file
    1. `.env` is developer secret and should not be checked into git. Approprioate `.gitignore` entry is set
5. Run "python tasks.py"

# License
MIT Licence. Do not run this in production.
No infringement of digdeeper's hard work is intended. 
I merely want to get an mvp up and running and then switch data source to BSE/NSE
