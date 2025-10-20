import asyncio
import websockets
import json

url = "https://local"

async def listen():
    
    async with websockets.connect(url) as socket:

        while True:
            received = socket.recv()
            processReception(received)


asyncio.get_event_loop().run_until_complete(listen)    


def processReception(received):
    package = json.dumps(received)


    pass
    