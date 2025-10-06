from gpiozero import Button
import time
import IRTransmitter as IRT

GPIO_PIN = 15 # enter the pin that will be used

sensor = Button(GPIO_PIN, pull_up=True)



def checkValidPress():
    ## check if it is not spammed if so change color to orange / red to and do not shoot
    return False

def buttonpress():

    if (checkValidPress):
        # send signal to to send IRTransmitter
        IRT.shootwithinfo()





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