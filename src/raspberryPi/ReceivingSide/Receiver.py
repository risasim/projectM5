from gpiozero import Button
from signal import pause
import time





def formatHit(message):
    formatted = message
    # format

    
    sendToServer(formatted)
    pass


def sendToServer(message):
    sendMessage(message)

def receiveInfo():
    pass

def receive():
    PIN = 17
    ir_sensor = Button(PIN, pull_up=True, bounce_time=0.3)

    def signal_received():
        print("signal received i got hit")
        pass


    ir_sensor.when_pressed = signal_received

    pause()






