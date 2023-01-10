import asyncio
import pyppeteer
from digdeeper import get_n_digdeeper_cell_data
from environs import Env
import json

env = Env()
env.read_env()
api_token = env("BROWSERLESS_API_KEY")

async def main():
  '''
  The main async io even loop
  '''
  browser = await pyppeteer.launcher.connect(
    browserWSEndpoint='wss://chrome.browserless.io?token='+api_token
    )
  page = await browser.newPage()
  await page.setViewport({"width": 1280, "height": 926});
  await page.setUserAgent(
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36'
  );
  
  batch_values:list[dict] = await get_n_digdeeper_cell_data(page, 200)
  
  with open("ddo-sample-data-{len(batch_values)}.json","w") as outfile:
    outfile.write(json.dumps(batch_values))
  await browser.close()

asyncio.run(main())

