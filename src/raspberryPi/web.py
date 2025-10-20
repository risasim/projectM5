import asyncio 
import json
import websockets
import time
import mediaplayer as mp




class web:  
    async def __init__(self,url):
        async with websockets.connect(url) as sock:
            self.socket = sock
            while True:
                received = sock.recv()
                processReception(received)

    async def speak(self,msg):
        await self.socket.send(msg)
   
    async def processOnReception(message):
        received = json.load(message)
        type = received["msgtype"]
        match type:
            case "HitResponseMsg":
                data = received["Data"]
                if data["playSound"]:
                    if not mp.playSound(data["soundName"]):
                        # not downloaded so download it
                        speak("") # dont know yet endpoint for downloading message 
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


asyncio.get_event_loop().run_until_complete(listen) 
