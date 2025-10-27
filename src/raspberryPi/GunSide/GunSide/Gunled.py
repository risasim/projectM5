from gpiozero import LED
import time

greenPin = 23
redPin = 24

green_led = LED(greenPin)
red_led = LED(redPin)


def changecolor(color):
    if (color == "RED"):
        red_led.on()
        green_led.off()
    if (color == "GREEN"):
        red_led.off()
        green_led.on()
    if (color == "ORANGE"):
        red_led.on()
        green_led.on()
    if (color == "NONE"):
        red_led.off()
        green_led.off()



def glmain():
    while(True):
        changecolor("GREEN")
        time.sleep(2)
        changecolor("ORANGE")
        time.sleep(2)
        changecolor("RED")
        time.sleep(2)
        
