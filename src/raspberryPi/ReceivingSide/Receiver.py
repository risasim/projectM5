from gpiozero import Button
from signal import pause
import time



def formatHit(message):
    formatted = message
    # format



def receiveInfo():
    pass
    
def signal_received():
    global i
    print(f"signal received i got hit for the : {i} time")
    i += 1

def receive():
    PIN = 17
    ir_sensor = None
    
    try:
    
        ir_sensor = Button(PIN, pull_up=True)
        ir_sensor.when_pressed = signal_received

        pause()
    except Exception as e:
        print(f"an error occurred: \n {e}")
    finally:
        if ir_sensor:
            ir_sensor.close()
            
            
if __name__ == "__main__":
    global i
    i = 0
    print("doing main; recieve")
    receive()



