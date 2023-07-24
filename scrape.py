import asyncio
import csv
import aiohttp


from lxml import html

header_values = {
    'name': 'Michael Foord',
    'location': 'Northampton',
    'language': 'English',
    'User-Agent': 'Mozilla 4/0',
    'Accept-Encoding': 'gzip',
    'Accept-Language': 'en-US,en;q=0.9,es;q=0.8',
    'Upgrade-Insecure-Requests': '0',
    'Referrer': 'https://www.google.com/'
}

def get_links():
    links = []
    base = "https://www.buzzfeed.com/archive/" # + y/m/d
    start = 2006
    end = 2023
    for year in range(start, end+1):
        for month in range(1, 13):
            for day in range(1, 32):
                links.append(base + f"{year}/{month}/{day}")
    return links

async def get_content_from_link(session, link):
    def texts_from_html_elements(elements):
        return [x.strip() for x in elements]

    try:
  
        async with session.get(link) as response:
            
            tree = html.fromstring(await response.text())

            title_path = '//div[2]/div/h2/a/text()[normalize-space()]'
            link_path = '//div[2]/div/h2/a/@href'
            desc_path = '//div[2]/div/p/text()[normalize-space()]'
            author_path = '//div[3]/div/div/a/span/text()[normalize-space()]'

            titles = tree.xpath(title_path)
            links = tree.xpath(link_path)
            descs = texts_from_html_elements(tree.xpath(desc_path))
            authors = texts_from_html_elements(tree.xpath(author_path))
            link_comp = link.split("/")
            date = link_comp[-3] + "/" + link_comp[-2] + "/" + link_comp[-1]
            print(date)
            return list(zip([date]*len(titles), range(0, len(titles)), titles, links, descs, authors))
    except:
        print("unable to get ", link.split(".com"[-1]))
        return []

async def get_content_from_links(links):
    contents = []
    async with aiohttp.ClientSession() as session:
        contents = await asyncio.gather(*[get_content_from_link(session, link) for link in links])
    if contents is not None:
        return [item for row in contents for item in row]
    else:
        return []
def main():
    links = get_links()
    x = asyncio.get_event_loop().run_until_complete(get_content_from_links(links))
    with open('./csv_file.csv', 'w', encoding="utf-8") as f:
        writer = csv.writer(f)
        writer.writerow(["date", "index", "titles", "links", "descs", "authors"])
        writer.writerows(x)
    print(x)


if __name__ == "__main__":
    main()