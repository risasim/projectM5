#Used for running async threads and such
import asyncio
#Used for keeping a stable connection to the server for live updates
import websockets
#Used to convert between objects and json
import json
#Used to make a http request at the start of a connection.
import requests
import os
from threading import Thread,Lock,Event
from queue import Queue,Empty
from datetime import datetime
import mediaplayer
import mimetypes


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
        self.sound_path = None
        self.websocket_loop = None
        print("init done")
        


    def download_file(self):
        try:
            print(self.auth_header)
            response = requests.get(url=f"{self.http_url}/api/sound", headers=self.auth_header) #, timeout=10#, stream = True
            response.raise_for_status()
            content_type = response.headers.get("Content-Type")
            extension = ""

            if content_type:
                extension = mimetypes.guess_extension(content_type, strict=False)

            if not extension:
                extension = ".bin"

            script_dir = os.path.dirname(__file__)
            download_dir = "downloads"
            download_dir = os.path.join(script_dir, download_dir)
            filename = f"downloaded_sound{extension}"            
            full_path = os.path.join(download_dir, filename)
            self.sound_path = full_path
            os.makedirs(download_dir, exist_ok=True)

            with open(full_path, 'wb') as f:
                # if streams true
                # for chunk in response.iter_content(chunk_size=8192):
                #     f.write(chunk)

                # if not streams true
                f.write(response.content)
        except requests.exceptions.RequestException as e:
            print(f"an error occured during requests :\n{e}")
        except IOError as e:
            print(f"an error occured while writing file:\n{e}")
    
    def authenticate(self):
        try:
            headers = {"Content-Type": "application/json"}
            payload = {
            "apiKey":self.api_key,
            "piSn": self.pi_sn
            }
            res = requests.post(self.http_url + "/api/piAuth",json=payload,headers=headers, timeout=10
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
                    self.websocket_loop = asyncio.get_running_loop()
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
                print("message in websock")
                try:
                    object = json.loads(message)
                    print(object)
                    global alive
                    global in_game
                    global infected
                    match object["msgtype"]:
                        case "Start":
                            if object["Data"]["active"]:
                                alive = True
                                infected = False
                            else:
                                alive = False
                                infected = True

                            in_game = True
                            
                            self.download_file()

                        case "HitResponseMsg":
                            if object["Data"]["playSound"]:
                                mediaplayer.playSound(self.sound_path) #object.Data.soundName

                            if object["Data"]["dead"]:
                                alive = False
                            
                            if object["Data"]["revive"]:
                                arg = object["Data"]["revivein"]
                                def revive(revive_in):
                                    global alive
                                    time.sleep(revive_in)
                                    alive = True
                                    changecolor("GREEN")
                                    
                                Thread(target=revive, args=(arg,))
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
            print(f"websocket error:\n{e}")

    
    def send_hit_data(self):
        print("inside send_hit data")

        obj =  {}
        obj["msgType"] = "HitDataMsg"
        obj["Data"] = {}
        obj["Data"]["victim"] = self.pi_sn
        obj["Data"]["timestamp"] = datetime.now().isoformat()
        message = json.dumps(obj)

        async def send_coroutine():
            print("in send coroutine")
            try:
                await self.web_sock.send(message)
            except Exception as e:
                print(f"an error happend:\n{e}")

        def schedule_send():
            print("in schedule send ")
            try:
                asyncio.run_coroutine_threadsafe(send_coroutine(), self.websocket_loop)
            except Exception as e:
                print(f"Error scheduling hit data send: {e}")
            
        send_thread = Thread(target=schedule_send)
        send_thread.start()








#game relevant vars

alive = False
in_game = False
counter = 0
infected = False
shoot_thread = None
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
    
    client.download_file()
    print(type(client.sound_path))
    mediaplayer.playSound(client.sound_path)
    print("init sensors")
    ir_sensor = Button(17, pull_up=True, bounce_time=0.5)
    ir_sensor_back = Button(27, pull_up=True, bounce_time=0.5)
    button = Button(15, pull_up=True)

    def buttonpress_helper():
        print("in helper function")
        global counter
        global alive
        if not alive:
            changecolor("NONE")
            return
        changecolor("RED")
        shoot()
        
        counter += 1
        print("counter: " + str((counter % 6)))
        if counter % 6 == 0:
            changecolor("NONE")
            time.sleep(2)
        
        changecolor("GREEN")


    def buttonpress():
        global shoot_thread
        print("in buttonpress")
        if shoot_thread is not None and shoot_thread.is_alive():
            return
        shoot_thread = Thread(target=buttonpress_helper)
        shoot_thread.start()
        
        

    def signal_received(client: WebClient, back, alive, infected):
        if not alive and not infected:
            return 
        if back: print("this is cause of the back sensor")
        print("signal received RECEIVED")
        print("going into send_hit_data")
        client.send_hit_data()

    print("setting functions")
    ir_sensor.when_pressed = lambda: signal_received(client, False, alive, infected)
    ir_sensor_back.when_pressed = lambda: signal_received(client, True, alive, infected)
    button.when_pressed = buttonpress
    
    print("awaiting handler")
    await client.handler()


if __name__ == "__main__":
    asyncio.run(main())






