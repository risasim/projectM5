from gpiozero import Button
import time
import Gunled
from signal import pause
from Transmitter import shoot



GPIO_PIN = 15 # enter the pin that will be used

counter = 0


sensor = Button(GPIO_PIN, pull_up=True)


original_handler = None    

def buttonpress():
    global counter
    global original_handler
    if original_handler is None:
        original_handler = sensor.when_pressed
    
    
    sensor.when_pressed = None
    try:
        Gunled.changecolor("RED")
        shoot() #maybe do this async
        
        counter += 1
        print("counter: " + str((counter % 6)))
        if counter % 6 == 0:
            time.sleep(2)
        
        Gunled.changecolor("GREEN")
    finally:
        sensor.when_pressed = buttonpress






sensor.when_pressed = buttonpress

pause()

## KY-004 Button
## GND (right) connect to ground on rpi.
## VCC (left) connect to 3.3V pin (3V3) on rpi.
## OUT (mid) connect to GPIO pin on rpi. 
## important note: activate pull-up on the GPIO pin for reliability 
