from gpiozero import LED
import time

green_led = LED(23)
red_led = LED(24)


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
