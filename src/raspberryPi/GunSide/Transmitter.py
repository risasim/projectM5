from rpi_hardware_pwm import HardwarePWM
import time


def shoot():
    IR_CARRIER_FREQ = 38000
    DUTY_CYCLE = 50

    try:
        pwm = HardwarePWM(pwm_channel=2, hz=IR_CARRIER_FREQ, chip=0)

        print("start sending...")
        pwm.start(DUTY_CYCLE)
        time.sleep(1)

    except Exception as e:
        print(f"an error occurred: \n {e}")
    finally:
        if 'pwm' in locals():
            pwm.stop()
            print("stopped sending")

