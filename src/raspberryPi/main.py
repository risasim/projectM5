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
url = "https://local"
socket = None


piNumb = "ae616eb0e54290a6"

class PlayerState:

    def __init__(self):
        self.health = 1
        self.timeAlive = 0
        self.lock = Lock()





class WebClient:

    def __init__(self,url):

        self.serverUrl = url
        self.player = PlayerState()
        self.running = Event()
        self.queue = Queue()
        # self.receiverThread = ReceiverThread(self)
        # self.transmitterThread = TransmitterThread(self)

        #Handle the initial get request to obtain a token for the websocket. 

        try:
            
            requests.get(self.serverUrl + "/piAuth",json={"apiKey": 1234, "piSn": piNumb},headers={"Accept": "application/json"}
                
            )
        except requests.exceptions.RequestException:
            pass
        

    async def handler(self):
        try: 
            async with websockets.connect(self.serverUrl + "/api/wsPis") as webSocket:

                print("Creation of Websocket Completed.")

                while not self.running.is_set():

                    try:
                        message = self.queue.get_nowait()
                        await webSocket.send(json.dumps(message))
                    except Empty:
                        pass

                    try:
                        message = await asyncio.wait_for(webSocket.recv,0.1 )
                        self.processReception(json.loads(message))
                    except asyncio.TimeoutError:
                        pass
                    except websockets.exceptions.ConnectionClosedOK:
                        print("Connection Closed!")
                        break
                    await asyncio.sleep(0.1)
        except Exception:
            print("Websocket Error!")
    
    def processReception(object):

            match object.msgtype:

                case "start":
                    pass
                case "dead":
                    pass


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




from gpiozero import Button
import time
import Gunside.Gunled
from signal import pause
from Gunside.Transmitter import shoot


if __name__ == "__main__":
    
    try:
        url = "http://116.203.97.62:8080"
        client = WebClient(url)
        global alive = False
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
                    Gunled.changecolor("NONE")
                    return
                Gunled.changecolor("RED")
                shoot() #maybe do this async
                
                counter += 1
                print("counter: " + str((counter % 6)))
                if counter % 6 == 0:
                    time.sleep(2)
                
                Gunled.changecolor("GREEN")
            finally:
                sensor.when_pressed = buttonpress



            
        def signal_received():
            global i
            print(f"signal received i got hit for the : {i} time")
            i += 1
            
            
            
            


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






