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
url = "http://116.203.97.62:8080"
websockurl = "ws://116.203.97.62:8080"
socket = None
piNumb = "ae616eb0e54290a6"

#game relevant vars

alive = False
gamealive = False



class PlayerState:

    def __init__(self):
        self.health = 1
        self.timeAlive = 0
        self.lock = Lock()





class WebClient:

    def __init__(self,url):

        self.serverUrl = url
        self.websockurl = websockurl
        self.player = PlayerState()
        self.running = Event()
        self.queue = Queue()
        # self.receiverThread = ReceiverThread(self)
        # self.transmitterThread = TransmitterThread(self)

        #Handle the initial get request to obtain a token for the websocket. 
        token = None
        self.auth_header = {
            "Authorization": f"Bearer {token}"
        }
        try:
            complete = self.serverUrl + "/piAuth"
            print(complete)
            headers = {
                "Content-Type": "application/json"
                    }
            payload = {
            "apiKey":"123e4567-e89b-12d3-a456-426614174000",
            "piSn": piNumb
            }
            r = requests.post(complete,json=payload,headers=headers
                #what do we do with the response of this request?
            )
            if r.status_code == 200:
                received = r.json()
                print(received)
                print(r.text)
                token = received.get(token) 
            else: 
                print("not successedesdfesafd")
                print(r.status_code)
                print(r.text)
        except requests.exceptions.RequestException:
            print("requeset.exceptions.requestException happend")

    async def handler(self):
        try: 
            async with websockets.connect(self.websockurl + "/api/wsPis",additional_headers=self.auth_header) as webSocket:

                print("Creation of Websocket Completed.")

                while not self.running.is_set():

                    try:
                        message = self.queue.get_nowait()
                        await webSocket.send(json.dumps(message))
                    except Empty:
                        pass

                    try:
                        message = await asyncio.wait_for(webSocket.recv,0.1 )
                        self.processReception(self,json.loads(message))
                    except asyncio.TimeoutError:
                        pass
                    except websockets.exceptions.ConnectionClosedOK:
                        print("Connection Closed!")
                        break
                    await asyncio.sleep(0.1)
        except Exception as e:
            print(f"Websocket Error!\n{e}")

    #Handle the received message
    
    def processReception(self,object):
            global alive
            global gamealive
            match object.msgtype:
                case "Start":
                    if object.Data.active:
                        alive = True
                    else:
                        alive = False

                    gamealive = True

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
                    gamealive = False



    def start(self):
        self.running.clear()
        # self.transmitterThread.start()
        # self.receiverThread.start()

        asyncio.run(self.handler())

        # self.transmitterThread.join()
        # self.receiverThread.join()

    def stop(self):
        self.running.set()
            

#
#Class representing the thread which runs the receiver logic
#
# class ReceiverThread(Thread):
    
#     def __init__(self,webc):
#         self.connection = webc
#         super().__init__()
#         print("Started receiver thread.")
#         pass
#     def run(self):
#         while not self.connection.running.is_set():
#             print('in receiverthread running')
            

        

#
#Class representing the thread which runs the transmitter logic.
#Edit run() to alter the functionality of the thread. __init__ can be edited to add properties or instantiate other classes.
# #
# class TransmitterThread(Thread):

#     def __init__(self,webc):
#         self.connection = webc
#         super().__init__()
#         print("Started transmitter thread.")
#     def run(self):
#         while not self.connection.running.is_set():
#             print('in traansmittertrhead running')
            
#Instantiates the first class WebClient, passing the given url as a paramater. 









try:

    client = WebClient(url)
    GPIO_PIN = 15 # enter the pin that will be used
    counter = 0

    PIN = 17
    ir_sensor = None
        
    ir_sensor = Button(PIN, pull_up=True, bounce_time=0.5)



    sensor = Button(GPIO_PIN, pull_up=True)
    original_handler = None   


    def buttonpress():
        global counter
        global original_handler
        if original_handler is None:
            original_handler = sensor.when_pressed
        
        
        sensor.when_pressed = None
        try:
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
        finally:
            sensor.when_pressed = buttonpress



        
    def signal_received():
        global i
        print(f"signal received i got hit for the : {i} time")
        i += 1

        obj =  {}
        obj.msgType = "HitDataMsg"
        obj.Data = {}
        obj.Data.victim = id
        obj.Data.timestamp = datetime.now().isoformat()
        client.queue.put(json.dumps(obj))
        
        
        
        


    ir_sensor.when_pressed = signal_received
    sensor.when_pressed = buttonpress




        
    try:
        client.start()
    except KeyboardInterrupt:
        print("Interruption Occured.")
    finally:
        client.stop()
except Exception as e:
    print(f"e = {e}")
finally:
    sensor.close()
    ir_sensor.close()






