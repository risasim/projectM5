#Used for running async threads and such
import asyncio
#Used for keeping a stable connection to the server for live updates
import websockets
#Used to convert between objects and json
import json
#Used to make a http request at the start of a connection.
import requests
import os
from web import web
from threading import Thread,Lock,Event
from queue import Queue,Empty
from datetime import datetime
import mediaplayer

from gpiozero import Button
import time
from GunSide.Gunled import changecolor
from signal import pause
from GunSide.Transmitter import shoot




class PlayerState:

    def __init__(self):
        self.health = 1
        self.timeAlive = 0
        self.lock = Lock()





class WebClient:

    def __init__(self,url, piSn, apikey):

        self.http_url = f"http{url}"
        self.ws_url = f"ws{url}"
        self.api_key = apikey
        self.pi_sn = piSn
        self.player = PlayerState()
        self.running = Event()
        self.queue = Queue()
        self.token = None
        self.auth_header = None
        self.web_sock = None
        print("init done")


        
    
    def authenticate(self):
        try:
            headers = {"Content-Type": "application/json"}
            payload = {
            "apiKey":self.api_key,
            "piSn": self.pi_sn
            }
            res = requests.post(self.http_url + "/piAuth",json=payload,headers=headers, timeout=10
                #what do we do with the response of this request?
            )
            res.raise_for_status()

            response = res.json()
            self.token = response.get("token")
            self.auth_header = {"Authorization": f"Bearer {self.token}"}
            print(f"authentication done token is:{self.token}")
        except requests.exceptions.RequestException:
            print("requeset.exceptions.requestException happend")
        except Exception as e:
            print(f"Error:\n{e}")

    async def handler(self):
        
        if not self.auth_header:
            raise ValueError("authenticate first")
        
        ws_endpoint = f"{self.ws_url}/api/wsPis"

        try: 
            async with websockets.connect(
                ws_endpoint,
                additional_headers=self.auth_header,
                ping_interval=20,
                ping_timeout=10
                ) as webSocket:
                    self.web_sock = webSocket
                    receive_task = asyncio.create_task(self._receive_message(webSocket))
                    send_task = asyncio.create_task(self._send_periodic_pings(webSocket))

                    await asyncio.gather(receive_task, send_task)
        except websockets.exceptions.WebSocketException as e:
                print(f"websocket error:\n{e}")
        except Exception as e:
                print(f"error:\n{e}")

    async def _receive_message(self, websocket):

        try:
            async for message in websocket:
                try:
                    object = json.loads(message)
                    global alive
                    global in_game
                    match object.msgtype:
                        case "Start":
                            if object.Data.active:
                                alive = True
                            else:
                                alive = False

                            in_game = True

                        case "HitResponseMsg":
                            if object.Data.playSound:
                                mediaplayer.playSound(object.Data.soundName)

                            if object.Data.dead:
                                alive = False
                            
                            if object.Data.revive:
                                def revive(object):
                                    global alive
                                    time.sleep(object.Data.reviveIn)
                                    alive = True
                                Thread(target=revive, args=(object,))
                        case "End":
                            in_game = False
                except json.JSONDecodeError as e:
                    print(f"json error:\n{e}")
        except Exception as e:
            print(f"error:\n{e}")


    async def _send_periodic_pings(self, websocket):

        try:
            while True:
                await asyncio.sleep(10)  # Send every 10 seconds
                ping_message = {
                    "type": "ping",
                    "piSn": self.pi_sn,
                    "timestamp": datetime.now().isoformat()
                }
                await websocket.send(json.dumps(ping_message))
        except websockets.exceptions.ConnectionClosed as e:
            print(f"websocket error:\m{e}")

    
    async def send_hit_data(self, websocket):
        obj =  {}
        obj.msgType = "HitDataMsg"
        obj.Data = {}
        obj.Data.victim = self.pi_sn
        obj.Data.timestamp = datetime.now().isoformat()
        await websocket.send(json.dumps(obj))






#game relevant vars

alive = True
in_game = True
counter = 0

async def main():
    base_url = "://116.203.97.62:8080"
    pi_Sn = "ae616eb0e54290a6"
    api_key = "123e4567-e89b-12d3-a456-426614174000"

    client = WebClient(base_url,pi_Sn, api_key)
    try:
        client.authenticate()
    except Exception as e:
        print(f"authentication error:\n{e}")
        return
    
    print("init sensors")
    ir_sensor = Button(17, pull_up=True, bounce_time=0.5)
    button = Button(15, pull_up=True)


    def buttonpress():
        global counter
        global alive
        if not alive:
            changecolor("NONE")
            return
        changecolor("RED")
        shoot() #maybe do this async
        
        counter += 1
        print("counter: " + str((counter % 6)))
        if counter % 6 == 0:
            time.sleep(2)
        
        changecolor("GREEN")


    def signal_received(client):
        print("signal received RECEIVED")
        client.send_hit_data(client.web_sock)

    print("setting functions")
    ir_sensor.when_pressed = lambda: signal_received(client)
    button.when_pressed = buttonpress
    
    print("awaiting handler")
    await client.handler()


if __name__ == "__main__":
    asyncio.run(main())






