from gpiozero import LED
import time

greenPin = 16
redPin = 17

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



