import asyncio
import websockets
import json
from web import web
from threading import Thread,Lock,Event
from queue import Queue,Empty
url = "https://local"
socket = None

webInstance = web(url)


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
        self.receiverThread = ReceiverThread(self)
        self.transmitterThread = TransmitterThread(self)


    async def handler(self):
        try: 
            async with websockets.connect(self.serverUrl) as webSocket:

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
        except Exception:
            print("Websocket Error!")
    
        def processReception(object):

            match object.msgtype:

                case "start":
                    pass
                case "dead":
                    pass


        def start():
            self.running.clear()
            self.transmitterThread.start()
            self.receiverThread.start()

            asyncio.run(self.handler)

            self.transmitterThread.join()
            self.receiverThread.join()

        def stop():
            self.running.set()


class ReceiverThread(Thread):
    
    def __init__(self,webc):
        pass
    def run(self):
        pass

#
#Class representing the thread which runs the transmitter logic.
#Edit run() to alter the functionality of the thread. __init__ can be edited to add properties or instantiate other classes.
#
class TransmitterThread(Thread):
    def __init__(self,webc):
        pass
    def run(self):
        pass

#Instantiates the first class WebClient, passing the given url as a paramater. 
url = "whatdahellyurl.com"
client = WebClient(url)


try:
    client.start()
except KeyboardInterrupt:
    print("Interruption Occured.")
finally:
    client.stop()






