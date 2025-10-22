from gpiozero import Button
import time
import Gunled
import Transmitter
from signal import pause
from web import alive




GPIO_PIN = 15 # enter the pin that will be used


timeList = []

sensor = Button(GPIO_PIN, pull_up=True)

minimum = 1000000000

def checkValidPress():
    
    if not alive:
        return "RED"
    
    deltaList = []
    for i in range(0, len(timeList)):
        j = i + 1
        if (j >= len(timeList)):
            break
        print("-------------------")
        print(timeList[j])
        deltaList.append(timeList[j] - timeList[i])
    

    av = (sum(deltaList)/ len(deltaList)) if not len(deltaList) == 0  else minimum + 1
    if ( av < minimum):
        Gunled.changecolor("RED")
        return "RED"
    else:
        Gunled.changecolor("GREEN")
        return "GREEN"
    ## check if it is not spammed if so change color to orange / red to and do not shoot

def buttonpress():
    if (len(timeList) >= 5):
        timeList.pop(0)
    timeList.append(time.time_ns())
    color = checkValidPress()
    if not (color == "RED"):
        # send signal to to send IRTransmitter
        Transmitter.shootWithInfo()
        Gunled.changecolor("NONE")
        time.sleep(0.05)
        Gunled.changecolor(color)






sensor.when_pressed = buttonpress

pause()

## KY-004 Button
## GND (right) connect to ground on rpi.
## VCC (left) connect to 3.3V pin (3V3) on rpi.
## OUT (mid) connect to GPIO pin on rpi. 
## important note: activate pull-up on the GPIO pin for reliability 
