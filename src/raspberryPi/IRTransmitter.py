# need pigpio lib for this (should be on the rpi automatically)
# make sure pigpio deamon is running
    # sudo apt-get update
    # sudo apt-get install pigpio python3-pigpio

    # sudo systemctl enable pigpiod
    # sudo systemctl start pigpiod

# also use a pwm pin (GPIO18)

import pigpio
import time


def shoot():

    PWM_GPIO = 18

    FREQ = 38000

    DUTY_CYCLE_MILLION = 500000

    pi = pigpio.pi()

    if not pi.connected:
        print("cant make connection to pigpio deamon check if it is running")
        exit()

    try:
        pi.hardware_PWM(PWM_GPIO, FREQ, DUTY_CYCLE_MILLION)
        print("sending carrier wave")

        time.sleep(5)

        pi.hardware_PWM(PWM_GPIO, 0, 0)
        print("stopped sending carrier wave")
    except Exception as e:
        print(f"something went wrong: {e}")

    finally:
        pi.stop()
        print("pigpio connection closed")



## for this we need LIRC
    # sudo apt update
    # sudo apt install lirc
## change the config
    # sudo nano /boot/firmware/config.txt
    # add : dtoverlay=gpio-ir-tx,gpio_pin=18 
    # with gpio_pin being the pin we are using 
    # reboot the pi 
    # after reboot there should be a /dev/lirc0 

## now we need define what we send 
    # first we stop the lirc service so we can edit it
        # sudo systemctl stop lircd

## make a new config file (sender.lircd.conf) replace the contents
## of /etc/lirc/lircd.conf with this basisconfig
    # sudo nano /etc/lirc/lircd.conf

# # Remote name is 'SENDER_A'
# # keyname is 'ID_SIGNAL'

# begin remote

#   name  SENDER_A
#   flags NEC
#   eps            30
#   aeps          100

#   header       9000  4500
#   one           560  1690
#   zero          560   560
#   ptrail        560
#   repeat       9000  2250
#   pre_data_bits   32
#   pre_data     0xAAAAAAAA
#   gap          108000
#   toggle_bit_mask 0x0

#       begin codes
#           ID_SIGNAL     0x12345678 # the unique code that is send
#       end codes

# end remote


## SENDER_A is the name of the remote (source-identifier)
## ID_SIGNAL is the name of the key (signal-identifier)
## 0x12345678 is the hex that is send. you can change the pre_data and ID_SIGNAL to make a unique 64 bit identifier 

## now start lirc again
    # sudo systemctl start lircd

## lastly make sure the lirc lib for pyton is installed
    # pip3 install python-lirc

import lirc
import time

def shootwithinfo():

    REMOTE = "SENDER_A"
    KEY = "ID_SIGNAL"

    try:
        
        client = lirc.client()
        
        client.send_once(REMOTE, KEY)
    except Exception as e:
        print(f"something went worng with sending the lirc commant {e}")
        print("make sure lircd service is running and the configuration is correct")
        print("check this with: sudo systemctl status lircd")
    finally:
        client.close()
        pass
