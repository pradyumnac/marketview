
import string
from time import sleep
import math

MAXRETRIES_SCROLL: int = 10
DIGDEEPER_ANNOUNCEMENTS_URL: string = "https://digdeeper.one/announcements"

async def get_current_rowcount(page):
  return await page.evaluate('''document.querySelectorAll(".list-dummy").length''');


async def triggerScroll(page, previous_row_count):

  '''
  Triggers an infinitescroll and waits for the content to get reloaded
  before return
  '''
  
  await page.evaluate('''window.scrollBy(0, 2000)''');

  # Wait for trigger based infinite scroll content update
  retry_count: int = 0
  current_row_count = previous_row_count;
  while(current_row_count <= previous_row_count and retry_count < MAXRETRIES_SCROLL):
    retry_count+=1
    sleep(1)    # wait 1 second before cheching again
    print(f"Waiting for reload: Iteration [{retry_count}], current row count: [{current_row_count}]")
    current_row_count = await get_current_rowcount(page);
  
  return current_row_count


async def get_n_digdeeper_cell_data(page, number_of_entries_to_load = 200):
  '''
  Get n number of data from digdeeper
  '''
  await page.goto(DIGDEEPER_ANNOUNCEMENTS_URL)


  row_count_to_load: int =  math.ceil(number_of_entries_to_load/3);
  current_row_count: int = await get_current_rowcount(page)
  
  while current_row_count < row_count_to_load:
    current_row_count = await triggerScroll(page, current_row_count)

  print(f"{current_row_count*3} cells loaded")
  

  # After N loads, the current date data to be grabbed should be there
  return await page.evaluate('''() => {
    
    let ret_vars = []
    document.querySelectorAll(".GroupItem").forEach(groupitem => {
      let ret_var = {}
      ret_var["link"]    = groupitem.querySelector(".bubble-r-box a")['href'];
      ret_var["name"]    = groupitem.querySelector(".bubble-element .Text .content").textContent;
      ret_var["desc"]    = groupitem.querySelectorAll(".bubble-element .Text")[1].textContent;
      ret_var["dt"]      = groupitem.querySelectorAll(".bubble-element .Text")[2].textContent;
      ret_var["tag"]     = groupitem.querySelectorAll(".bubble-element .Text")[3].textContent;
      ret_vars.push(ret_var)
    });

    return ret_vars;
  }''')






