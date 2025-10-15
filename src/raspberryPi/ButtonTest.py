import time
from gpiozero import Button
from mediaplayer import mmain
from Gunled import changecolor
from threading import Thread


GPIO_PIN = 15 # enter the pin that will be used


sensor = Button(GPIO_PIN, pull_up=True)



def buttonpress():
	print("the button was pressed")
	mmain()
	changecolor("GREEN")
	
	
	

sensor.when_pressed = buttonpress


def main():
    try: 
        while(True):
            time.sleep(1)

    except KeyboardInterrupt:
        print("inturrupted")
        changecolor("RED")
        
main()
