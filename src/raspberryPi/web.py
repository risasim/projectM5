import asyncio 
import json
from websockets.asyncio.client import connect
import time
import mediaplayer as mp



conn = None 
alive = False
inGame = False

async def createConnection(ip):
    async with connect(ip) as connection:
        conn = connection
        async for message in connection:
            await processOnReception(message)

    


async def processOnReception(message):
    received = json.load(message)
    type = received["msgtype"]
    match type:
        case "HitResponseMsg":
            data = received["Data"]
            if data["playSound"]:
                if not mp.playSound(data["soundName"]):
                    # not downloaded so download it
                    sendMessage("") # dont know yet endpoint for downloading message 
                # play sound
            if data["dead"]:
                alive = False
            if data["revive"]:
                time.sleep(data["reviveIn"] / 1000)
                alive = True

        case "Start":
            inGame = True
            if data["gamemode"] == "infected " or data["gamemode"] == "hotPatato":
                alive = False
            else :
                alive = True
        case "End":
            alive = False
            inGame = False


    
    print(message)


async def sendMessage(message):
    await conn.send(message)