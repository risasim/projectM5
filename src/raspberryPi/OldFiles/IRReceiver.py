import lirc
import time
import sys
import web
import json
import datetime

def receiveIR():
    try:
        sockid = lirc.connect("rpi_receiver", blocking=False)

        while True:
            codes = lirc.nextcode()

            if codes:

                data = codes[0] # [code, Name, Remote, Count]

                hex = data[0] # should be the hex of the of pre_data + code of 64 bit / 16 hex so 
                              # 0x1111111122222222
                              # 0x{pre_dat}{code}
                message = formatHit(hex)
                web.sendMessage(message) # don't know if it is correctly typed
    except lirc.error.LircError:
        print("error: cant make connection to lircd")
    except Exception as e:
        print(f"something went wrong {e}")
    finally:
        if "sockid" in locals():
            lirc.deinit()


def formatHit(hex):
    x = {
        "msgtype": "HitDataMsg",
        "Data": {
            "victim": 1, # me
            "shooter": (hex >> 32),
            "timestamp":  datetime.datetime.now().isoformat()
        }
    }