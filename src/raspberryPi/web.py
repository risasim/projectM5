import asyncio 
import json
from websockets.asyncio.client import connect
conn = None 

async def createConnection(ip):
    async with connect(ip) as connection:
        conn = connection
        async for message in connection:
            await processOnReception(message)

    


async def processOnReception(message):
    received = json.load(message)
    
    print(message)


async def sendMessage(message):
    await conn.send(message)