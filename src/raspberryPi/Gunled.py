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


def main():
    try:
        while (True):
            # red
            red_led.on()
            green_led.off()
            time.sleep(5)
            # green
            red_led.off()
            green_led.on()
            time.sleep(5)
            # orange
            red_led.on()
            green_led.on()
            time.sleep(5)
    except KeyboardInterrupt:
        print("interrupted")



