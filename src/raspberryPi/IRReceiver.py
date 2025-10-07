import lirc
import time
import sys

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
    except lirc.error.LircError:
        print("error: cant make connection to lircd")
    except Exception as e:
        print(f"something went wrong {e}")
    finally:
        if "sockid" in locals():
            lirc.deinit()

