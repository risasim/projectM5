from gpiozero import Button
import time
import IRTransmitter as IRT
import Gunled





GPIO_PIN = 15 # enter the pin that will be used


timeList = []

sensor = Button(GPIO_PIN, pull_up=True)



def checkValidPress():
    deltaList = []
    for i in range(0, len(timeList)):
        j = i + 1
        if (j >= len(timeList)):
            break
        deltaList.append(timeList[j] - timeList[i])
    

    av = (sum(deltaList)/ len(deltaList))
    if ( av < 500):
        Gunled.changecolor("RED")
        return "RED"
    elif (av < 1000):
        Gunled.changecolor("ORANGE")
        return "ORANGE"
    else:
        Gunled.changecolor("GREEN")
        return "GREEN"
    ## check if it is not spammed if so change color to orange / red to and do not shoot

def buttonpress():
    if (len(timeList) >= 5):
        timeList.pop(0)
    timeList.append(time.perf_counter_ns)
    color = checkValidPress()
    if not (color == "RED"):
        # send signal to to send IRTransmitter
        IRT.shootwithinfo()
        Gunled.changecolor("NONE")
        time.sleep(0.05)
        Gunled.changecolor(color)






sensor.when_pressed = buttonpress

def main():
    try: 
        while(True):
            time.sleep(1)

    except KeyboardInterrupt:
        print("inturrupted")


## KY-004 Button
## GND (right) connect to ground on rpi.
## VCC (left) connect to 3.3V pin (3V3) on rpi.
## OUT (mid) connect to GPIO pin on rpi. 
## important note: activate pull-up on the GPIO pin for reliability 